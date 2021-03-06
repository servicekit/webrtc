package network

import (
	"fmt"
	"sync"

	"github.com/pions/pkg/stun"
	"github.com/servicekit/webrtc/internal/datachannel"
	"github.com/servicekit/webrtc/internal/dtls"
	"github.com/servicekit/webrtc/internal/sctp"
	"github.com/servicekit/webrtc/internal/srtp"
	"github.com/servicekit/webrtc/pkg/ice"
	"github.com/servicekit/webrtc/pkg/rtp"
	"github.com/pkg/errors"
)

// Manager contains all network state (DTLS, SRTP) that is shared between ports
// It is also used to perform operations that involve multiple ports
type Manager struct {
	icePwd      []byte
	iceNotifier ICENotifier

	dtlsState *dtls.State

	certPairLock sync.RWMutex
	certPair     *dtls.CertPair

	bufferTransportGenerator BufferTransportGenerator
	bufferTransports         map[uint32]chan<- *rtp.Packet

	// https://tools.ietf.org/html/rfc3711#section-3.2.3
	// A cryptographic context SHALL be uniquely identified by the triplet
	//  <SSRC, destination network address, destination transport port number>
	// contexts are keyed by IP:PORT:SSRC
	srtpContextsLock sync.RWMutex
	srtpContexts     map[string]*srtp.Context

	sctpAssociation *sctp.Association

	portsLock sync.RWMutex
	ports     []*port
}

// NewManager creates a new network.Manager
func NewManager(icePwd []byte, bufferTransportGenerator BufferTransportGenerator, dataChannelEventHandler DataChannelEventHandler, iceNotifier ICENotifier) (m *Manager, err error) {
	m = &Manager{
		icePwd:                   icePwd,
		iceNotifier:              iceNotifier,
		bufferTransports:         make(map[uint32]chan<- *rtp.Packet),
		srtpContexts:             make(map[string]*srtp.Context),
		bufferTransportGenerator: bufferTransportGenerator,
	}
	m.dtlsState, err = dtls.NewState(true)
	if err != nil {
		return nil, err
	}

	m.sctpAssociation = sctp.NewAssocation(func(raw []byte) {
		m.portsLock.Lock()
		defer m.portsLock.Unlock()

		for _, p := range m.ports {
			if p.iceState == ice.ConnectionStateCompleted {
				p.sendSCTP(raw)
				return
			}
		}
	}, func(data []byte, streamIdentifier uint16, payloadType sctp.PayloadProtocolIdentifier) {
		switch payloadType {
		case sctp.PayloadTypeWebRTCDCEP:
			msg, err := datachannel.Parse(data)
			if err != nil {
				fmt.Println(errors.Wrap(err, "Failed to parse DataChannel packet"))
				return
			}
			switch msg := msg.(type) {
			case *datachannel.ChannelOpen:
				// Cannot return err
				ack := datachannel.ChannelAck{}
				ackMsg, err := ack.Marshal()
				if err != nil {
					fmt.Println("Error Marshaling ChannelOpen ACK", err)
					return
				}
				if err = m.sctpAssociation.HandleOutbound(ackMsg, streamIdentifier, sctp.PayloadTypeWebRTCDCEP); err != nil {
					fmt.Println("Error sending ChannelOpen ACK", err)
					return
				}
				dataChannelEventHandler(&DataChannelCreated{streamIdentifier: streamIdentifier, Label: string(msg.Label)})
			default:
				fmt.Println("Unhandled DataChannel message", m)
			}
		case sctp.PayloadTypeWebRTCString:
			fallthrough
		case sctp.PayloadTypeWebRTCBinary:
			fallthrough
		case sctp.PayloadTypeWebRTCStringEmpty:
			fallthrough
		case sctp.PayloadTypeWebRTCBinaryEmpty:
			dataChannelEventHandler(&DataChannelMessage{streamIdentifier: streamIdentifier, Body: data})
		default:
			fmt.Printf("Unhandled Payload Protocol Identifier %v \n", payloadType)
		}

	})

	return m, err
}

// Listen starts a new Port for this manager
func (m *Manager) Listen(address string) (boundAddress *stun.TransportAddr, err error) {
	p, err := newPort(address, m)
	if err != nil {
		return nil, err
	}

	m.ports = append(m.ports, p)
	return p.listeningAddr, nil
}

// Close cleans up all the allocated state
func (m *Manager) Close() {
	m.portsLock.Lock()
	defer m.portsLock.Unlock()

	err := m.sctpAssociation.Close()
	m.dtlsState.Close()
	for _, p := range m.ports {
		portError := p.close()
		if err != nil {
			err = errors.Wrapf(portError, " also: %s", err.Error())
		} else {
			err = portError
		}
	}
}

// DTLSFingerprint generates the fingerprint included in an SessionDescription
func (m *Manager) DTLSFingerprint() string {
	return m.dtlsState.Fingerprint()
}

// SendRTP finds a connected port and sends the passed RTP packet
func (m *Manager) SendRTP(packet *rtp.Packet) {
	m.portsLock.Lock()
	defer m.portsLock.Unlock()

	for _, p := range m.ports {
		if p.iceState == ice.ConnectionStateCompleted {
			p.sendRTP(packet)
			return
		}
	}
}

// SendDataChannelMessage sends a DataChannel message to a connected peer
func (m *Manager) SendDataChannelMessage(message []byte, streamIdentifier uint16) error {
	err := m.sctpAssociation.HandleOutbound(message, streamIdentifier, sctp.PayloadTypeWebRTCString)
	if err != nil {
		return errors.Wrap(err, "SCTP Association failed handling outbound packet")
	}

	return nil
}

func (m *Manager) iceHandler(p *port, oldState ice.ConnectionState) {
	// One port disconnected, scan the other ones
	if p.iceState == ice.ConnectionStateDisconnected {
		m.portsLock.Lock()
		defer m.portsLock.Unlock()

		for _, p := range m.ports {
			if p.iceState == ice.ConnectionStateCompleted {
				// Another peer is connected! We don't have to notify RTCPeerConnection
				break
			}
		}
		m.iceNotifier(p.iceState)
	} else {
		m.iceNotifier(p.iceState)
	}
}
