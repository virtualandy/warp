package warp

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/fxamacker/cbor"
)

var goodP256Key *ecdsa.PrivateKey
var goodP256X cbor.RawMessage
var goodP256Y cbor.RawMessage
var goodP256COSE *COSEKey
var goodP256Raw cbor.RawMessage
var goodP384Key *ecdsa.PrivateKey
var goodP521Key *ecdsa.PrivateKey
var good1024Key *rsa.PrivateKey
var good2048Key *rsa.PrivateKey
var good4096Key *rsa.PrivateKey
var good25519Pub ed25519.PublicKey
var good25519Priv ed25519.PrivateKey
var mockCredentialID string
var mockRawCredentialID []byte
var mockCredential *testCred
var mockUser *testUser
var mockAttestedCredentialData AttestedCredentialData
var mockRawAttestedCredentialData []byte
var mockRawAuthData []byte
var mockAuthData AuthenticatorData
var mockRawAttestationObject cbor.RawMessage

func TestMain(m *testing.M) {
	var err error

	goodP256Key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Key gen error: %v", err)
	}

	goodP384Key, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		log.Fatalf("Key gen error: %v", err)
	}

	goodP521Key, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Fatalf("Key gen error: %v", err)
	}

	good1024Key, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("Key gen error: %v", err)
	}

	good2048Key, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Key gen error: %v", err)
	}

	good4096Key, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("Key gen error: %v", err)
	}

	good25519Pub, good25519Priv, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Key gen error: %v", err)
	}

	goodP256X, err = cbor.Marshal(goodP256Key.PublicKey.X.Bytes(), cbor.EncOptions{Sort: cbor.SortCTAP2})
	if err != nil {
		log.Fatalf("X marshal err: %v", err)
	}
	goodP256Y, err = cbor.Marshal(goodP256Key.PublicKey.X.Bytes(), cbor.EncOptions{Sort: cbor.SortCTAP2})
	if err != nil {
		log.Fatalf("Y marshal err: %v", err)
	}

	goodP256COSE = &COSEKey{
		Kty:       int(KeyTypeEC2),
		Alg:       int(AlgorithmES256),
		CrvOrNOrK: []byte{1},
		XOrE:      goodP256X,
		Y:         goodP256Y,
	}

	goodP256Raw, err = cbor.Marshal(goodP256COSE, cbor.EncOptions{Sort: cbor.SortCTAP2})
	if err != nil {
		log.Fatalf("COSEKey marshal err: %v", err)
	}

	mockCredentialID = "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"
	mockRawCredentialID, err = base64.RawURLEncoding.DecodeString(mockCredentialID)
	if err != nil {
		log.Fatalf("Credential ID decode err: %v", err)
	}

	mockCredential = &testCred{
		owner:     &testUser{},
		id:        mockRawCredentialID,
		publicKey: []byte(cborPublicKey),
		signCount: 0,
	}

	mockUser = &testUser{
		name: "jsmith",
		icon: "",
		id: []byte{
			0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14,
			0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24,
			0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c,
			0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55,
		},
		displayName: "John Smith",
		credentials: map[string]Credential{
			mockCredentialID: mockCredential,
		},
	}

	mockAttestedCredentialData = AttestedCredentialData{
		AAGUID: [16]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // authData.attestedCredentialData.aaguid
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // v
		},
		CredentialID: []byte{
			0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, // authData.attestedCredentialData.credentialID
			0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, // |
			0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, // |
			0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55, // v
		},
		CredentialPublicKey: COSEKey{
			Kty:       int(KeyTypeEC2),
			Alg:       int(AlgorithmES256),
			CrvOrNOrK: cbor.RawMessage{0x01},
			XOrE: cbor.RawMessage{
				0x58, 0x20,
				0x36, 0xc4, 0x85, 0xf8, 0x83, 0xda, 0xcf, 0xb3,
				0x63, 0xc8, 0xf6, 0x4d, 0x6a, 0x82, 0xe5, 0x65,
				0x3d, 0x7d, 0x36, 0x64, 0x2b, 0x3a, 0x10, 0x8b,
				0x51, 0x55, 0x5a, 0x8d, 0x33, 0x40, 0x7d, 0x5c,
			},
			Y: cbor.RawMessage{
				0x58, 0x20,
				0x69, 0xc9, 0x52, 0x21, 0x4f, 0xce, 0x43, 0xea,
				0x5f, 0x80, 0x43, 0x10, 0xbb, 0xe6, 0x3e, 0xd,
				0xee, 0xcb, 0xf1, 0xe9, 0xba, 0x69, 0x5d, 0xac,
				0x77, 0x53, 0xb1, 0x31, 0xbc, 0xbf, 0xf3, 0x98,
			},
		},
	}

	mockRawAttestedCredentialData = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // authData.attestedCredentialData.aaguid
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // v
		0x00, 0x20, // authData.attestedCredentialData.credentialIDLength = 32
		0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, // authData.attestedCredentialData.credentialID
		0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, // |
		0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, // |
		0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55, // v
		0xa5, // map of 5 items
		0x1,  // key 1 (Kty)
		0x2,  // 2 (EC2 key)
		0x3,  // key 3 (Alg)
		0x26, // -7
		0x20, // key -1
		0x1,  // 1 (P256 Curve)
		0x21, // key -2
		0x58, // byte string, >24 bytes
		0x20, // 32 bytes length
		0x36, 0xc4, 0x85, 0xf8, 0x83, 0xda, 0xcf, 0xb3,
		0x63, 0xc8, 0xf6, 0x4d, 0x6a, 0x82, 0xe5, 0x65,
		0x3d, 0x7d, 0x36, 0x64, 0x2b, 0x3a, 0x10, 0x8b,
		0x51, 0x55, 0x5a, 0x8d, 0x33, 0x40, 0x7d, 0x5c,
		0x22, // key -3
		0x58, // byte string, >24 bytes
		0x20, // 32 bytes length
		0x69, 0xc9, 0x52, 0x21, 0x4f, 0xce, 0x43, 0xea,
		0x5f, 0x80, 0x43, 0x10, 0xbb, 0xe6, 0x3e, 0xd,
		0xee, 0xcb, 0xf1, 0xe9, 0xba, 0x69, 0x5d, 0xac,
		0x77, 0x53, 0xb1, 0x31, 0xbc, 0xbf, 0xf3, 0x98,
	}

	mockRawAuthData = append([]byte{
		0xd8, 0x33, 0x51, 0x40, 0x80, 0xa0, 0xc7, 0x2b, //authdata.rpIDHash
		0x1e, 0xfa, 0x42, 0xb1, 0x8c, 0x96, 0xb9, 0x27, // |
		0x3e, 0x9f, 0x19, 0x3f, 0xa9, 0x80, 0xdb, 0x09, // |
		0xa0, 0x93, 0x33, 0x86, 0x5c, 0x2b, 0x32, 0xf3, // v
		0x41,                   // authData.Flags
		0x00, 0x00, 0x00, 0x01, // authData.SignCount
	}, mockRawAttestedCredentialData...)

	mockAuthData = AuthenticatorData{
		RPIDHash: [32]byte{
			0xd8, 0x33, 0x51, 0x40, 0x80, 0xa0, 0xc7, 0x2b, //authdata.rpIDHash
			0x1e, 0xfa, 0x42, 0xb1, 0x8c, 0x96, 0xb9, 0x27, // |
			0x3e, 0x9f, 0x19, 0x3f, 0xa9, 0x80, 0xdb, 0x09, // |
			0xa0, 0x93, 0x33, 0x86, 0x5c, 0x2b, 0x32, 0xf3, // v
		},
		UP:                     true,
		UV:                     false,
		AT:                     true,
		ED:                     false,
		SignCount:              1,
		AttestedCredentialData: mockAttestedCredentialData,
	}

	mockRawAttestationObject = append(cbor.RawMessage{
		0xa3,             // map, 3 items
		0x63,             // text string, 3 chars
		0x66, 0x6d, 0x74, // "fmt"
		0x64,                   // text string, 4 chars
		0x6e, 0x6f, 0x6e, 0x65, // "none"
		0x67,                                     // text string, 7 chars
		0x61, 0x74, 0x74, 0x53, 0x74, 0x6d, 0x74, // "attStmt"
		0xa0,                                           // null
		0x68,                                           // text string, 8 chars
		0x61, 0x75, 0x74, 0x68, 0x44, 0x61, 0x74, 0x61, // "authData"
		0x58, 0xa4, // byte string, 164 chars
	}, mockRawAuthData...)

	os.Exit(m.Run())
}

func TestSupportedKeyAlgorithms(t *testing.T) {
	algs := SupportedKeyAlgorithms()
	if !reflect.DeepEqual(algs, []COSEAlgorithmIdentifier{
		AlgorithmEdDSA,
		AlgorithmES512,
		AlgorithmES384,
		AlgorithmES256,
		AlgorithmPS512,
		AlgorithmPS384,
		AlgorithmPS256,
		AlgorithmRS512,
		AlgorithmRS384,
		AlgorithmRS256,
		AlgorithmRS1,
	}) {
		t.Fatal("Unexpected result")
	}
}

func TestSupportedAttestationStatementFormats(t *testing.T) {
	fmts := SupportedAttestationStatementFormats()
	if !reflect.DeepEqual(fmts, []AttestationStatementFormat{
		AttestationFormatNone,
	}) {
		t.Fatal("Unexpected result")
	}
}
