package webrtc

import (
	"errors"
	"fmt"
)

// Types of InvalidStateErrors
var (
	ErrConnectionClosed = errors.New("connection closed")
)

func rtcErrorString(obj interface{}, err error, subErr error) string {
	if subErr != nil {
		return fmt.Sprintf("%T: %v (%v)", obj, err, subErr)
	}

	return fmt.Sprintf("%T: %v", obj, err)
}

// InvalidStateError indicates the object is in an invalid state.
type InvalidStateError struct {
	Err    error
	SubErr error
}

func (e *InvalidStateError) Error() string {
	return rtcErrorString(e, e.Err, e.SubErr)
}

// Types of UnknownErrors
var (
	ErrNoConfig = errors.New("no configuration provided")
)

// UnknownError indicates the operation failed for an unknown transient reason
type UnknownError struct {
	Err    error
	SubErr error
}

func (e *UnknownError) Error() string {
	return rtcErrorString(e, e.Err, e.SubErr)
}

// Types of InvalidAccessErrors
var (
	ErrCertificateExpired = errors.New("certificate expired")
	ErrNoTurnCred         = errors.New("turn server credentials required")
	ErrTurnCred           = errors.New("invalid turn server credentials")
	ErrExistingTrack      = errors.New("track aready exists")
)

// InvalidAccessError indicates the object does not support the operation or argument.
type InvalidAccessError struct {
	Err    error
	SubErr error
}

func (e *InvalidAccessError) Error() string {
	return rtcErrorString(e, e.Err, e.SubErr)
}

// Types of NotSupportedErrors
var ()

// NotSupportedError indicates the operation is not supported.
type NotSupportedError struct {
	Err    error
	SubErr error
}

func (e *NotSupportedError) Error() string {
	return rtcErrorString(e, e.Err, e.SubErr)
}

// Types of InvalidModificationErrors
var (
	ErrModPeerIdentity         = errors.New("peer identity cannot be modified")
	ErrModCertificates         = errors.New("certificates cannot be modified")
	ErrModRtcpMuxPolicy        = errors.New("rtcp mux policy cannot be modified")
	ErrModICECandidatePoolSize = errors.New("ice candidate pool size cannot be modified")
)

// InvalidModificationError indicates the object can not be modified in this way.
type InvalidModificationError struct {
	Err    error
	SubErr error
}

func (e *InvalidModificationError) Error() string {
	return rtcErrorString(e, e.Err, e.SubErr)
}

// Types of SyntaxErrors
var (
	ErrURLSyntaxInvalid = errors.New("URL syntax is invalid")
)

// SyntaxError indicates the string did not match the expected pattern.
type SyntaxError struct {
	Err    error
	SubErr error
}

func (e *SyntaxError) Error() string {
	return rtcErrorString(e, e.Err, e.SubErr)
}

// NotReadableError indicates failure to read object properties properly
type NotReadableError struct {
	Err    error
	SubErr error
}

func (e *NotReadableError) Error() string {
	return rtcErrorString(e, e.Err, e.SubErr)
}

// Types of NotReadableErrors
var (
	ErrBadIdentityAssertion = errors.New("failed identity assertion")
)

// OperationError indicates general operation errors
type OperationError struct {
	Err    error
	SubErr error
}

func (e *OperationError) Error() string {
	return rtcErrorString(e, e.Err, e.SubErr)
}

// Types of NotReadableErrors
var (
	ErrSDPGenerationFailed = errors.New("sdp generation failed")
)
