package webrtc

import (
	"testing"

	"github.com/servicekit/webrtc/internal/sdp"
)

func TestSetRemoteDescription(t *testing.T) {
	testCases := []struct {
		desc RTCSessionDescription
	}{
		{RTCSessionDescription{RTCSdpTypeOffer, sdp.NewJSEPSessionDescription("", false).Marshal()}},
	}

	for i, testCase := range testCases {
		peerConn, err := New(RTCConfiguration{})
		if err != nil {
			t.Errorf("Case %d: got error: %v", i, err)
		}
		err = peerConn.SetRemoteDescription(testCase.desc)
		if err != nil {
			t.Errorf("Case %d: got error: %v", i, err)
		}
	}
}
