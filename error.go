package warp

import "fmt"

//Error represents an error in a WebAuthn relying party operation
type Error struct {
	err     string
	wrapped error
}

//Error implements the error interface
func (e Error) Error() string {
	return e.err
}

//Unwrap allows for error unwrapping
func (e Error) Unwrap() error {
	return e.wrapped
}

//Wrap returns a new error which contains the provided error wrapped with this
//error
func (e Error) Wrap(err error) Error {
	n := e
	n.wrapped = err
	return n
}

//Is establishes equality for error types
func (e Error) Is(target error) bool {
	return e.Error() == target.Error()
}

//NewError returns a new Error with a custom message
func NewError(fmStr string, els ...interface{}) Error {
	return Error{
		err: fmt.Sprintf(fmStr, els...),
	}
}

//Categorical top-level errors
var (
	ErrDecodeAttestedCredentialData = Error{err: "Error decoding attested credential data"}
	ErrDecodeAuthenticatorData      = Error{err: "Error decoding authenticator data"}
	ErrDecodeCOSEKey                = Error{err: "Error decoding raw public key"}
	ErrEncodeAttestedCredentialData = Error{err: "Error encoding attested credential data"}
	ErrEncodeAuthenticatorData      = Error{err: "Error encoding authenticator data"}
	ErrGenerateChallenge            = Error{err: "Error generating challenge"}
	ErrMarshalAttestationObject     = Error{err: "Error marshaling attestation object"}
	ErrOption                       = Error{err: "Option error"}
	ErrNotImplemented               = Error{err: "Not implemented"}
	ErrUnmarshalAttestationObject   = Error{err: "Error unmarshaling attestation object"}
	ErrVerifyAttestation            = Error{err: "Error verifying attestation"}
	ErrVerifyAuthentication         = Error{err: "Error verifying authentication"}
	ErrVerifyClientExtensionOutput  = Error{err: "Error verifying client extension output"}
	ErrVerifyRegistration           = Error{err: "Error verifying registration"}
	ErrVerifySignature              = Error{err: "Error verifying signature"}
)
