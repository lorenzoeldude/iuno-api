//
// apple_jws.go
// iuno-api
//

package storekit

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strings"
)

// =====================================================
// JWS HEADER
// =====================================================

type appleJWSHeader struct {
	Alg string   `json:"alg"`
	X5C []string `json:"x5c"`
}

// =====================================================
// JWS ENVIRONMENT
// =====================================================
//
// We decode this before certificate verification so we can
// select the correct trust model.
//
// Xcode StoreKit Testing uses a self-signed StoreKit test
// certificate.
//
// Sandbox / Production use Apple's certificate chain.
//

type appleJWSEnvironment struct {
	Environment string `json:"environment"`
}

// =====================================================
// VERIFY APPLE JWS
// =====================================================
//
// Verifies:
//
// 1. JWS structure
// 2. ES256 algorithm
// 3. Correct certificate trust model
// 4. Certificate validity
// 5. JWS signature
//
// Trust models:
//
// Xcode:
//   StoreKitTestCertificate.cer
//
// Sandbox / Production:
//   Apple root certificate + certificate chain
//
// IMPORTANT:
// The certificate supplied by the JWS is NEVER trusted merely
// because Apple/Xcode supplied it.
//

func VerifyAppleJWS(
	signedPayload string,
) ([]byte, error) {

	// =====================================================
	// BASIC JWS STRUCTURE
	// =====================================================

	if signedPayload == "" {
		return nil, fmt.Errorf(
			"signed payload is empty",
		)
	}

	parts := strings.Split(
		signedPayload,
		".",
	)

	if len(parts) != 3 {
		return nil, fmt.Errorf(
			"invalid JWS format: expected 3 parts, got %d",
			len(parts),
		)
	}

	headerPart := parts[0]
	payloadPart := parts[1]
	signaturePart := parts[2]

	if headerPart == "" {
		return nil, fmt.Errorf(
			"JWS header is empty",
		)
	}

	if payloadPart == "" {
		return nil, fmt.Errorf(
			"JWS payload is empty",
		)
	}

	if signaturePart == "" {
		return nil, fmt.Errorf(
			"JWS signature is empty",
		)
	}

	// =====================================================
	// DECODE HEADER
	// =====================================================

	headerBytes, err :=
		base64.RawURLEncoding.DecodeString(
			headerPart,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode JWS header: %w",
			err,
		)
	}

	var header appleJWSHeader

	if err :=
		json.Unmarshal(
			headerBytes,
			&header,
		); err != nil {

		return nil, fmt.Errorf(
			"failed to decode JWS header JSON: %w",
			err,
		)
	}

	// =====================================================
	// VERIFY ALGORITHM
	// =====================================================

	if header.Alg != "ES256" {
		return nil, fmt.Errorf(
			"unsupported Apple JWS algorithm: %s",
			header.Alg,
		)
	}

	// =====================================================
	// DECODE PAYLOAD
	// =====================================================
	//
	// We need the environment before selecting the
	// certificate trust model.
	//
	// The environment itself is NOT trusted yet.
	// It becomes trusted only after the JWS signature
	// is successfully verified against the appropriate
	// trusted certificate.
	//

	payloadBytes, err :=
		base64.RawURLEncoding.DecodeString(
			payloadPart,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode JWS payload: %w",
			err,
		)
	}

	var environmentPayload appleJWSEnvironment

	if err :=
		json.Unmarshal(
			payloadBytes,
			&environmentPayload,
		); err != nil {

		return nil, fmt.Errorf(
			"failed to decode Apple JWS environment: %w",
			err,
		)
	}

	// =====================================================
	// SELECT TRUST MODEL
	// =====================================================

	switch environmentPayload.Environment {

	case "Xcode":

		return verifyXcodeJWS(
			header,
			headerPart,
			payloadPart,
			signaturePart,
			payloadBytes,
		)

	case "Sandbox", "Production":

		return verifyAppleJWS(
			header,
			headerPart,
			payloadPart,
			signaturePart,
			payloadBytes,
		)

	default:

		return nil, fmt.Errorf(
			"unknown Apple JWS environment: %s",
			environmentPayload.Environment,
		)
	}
}

// =====================================================
// VERIFY XCODE JWS
// =====================================================
//
// Xcode StoreKit Testing uses a self-signed StoreKit
// test certificate.
//
// The certificate exported from Xcode is a trusted root
// for this local test environment.
//
// Apple's documentation confirms that the Xcode test
// certificate is a root certificate and does not have
// a certificate chain.
//

func verifyXcodeJWS(
	header appleJWSHeader,
	headerPart string,
	payloadPart string,
	signaturePart string,
	payloadBytes []byte,
) ([]byte, error) {

	// =====================================================
	// CERTIFICATE COUNT
	// =====================================================

	if len(header.X5C) != 1 {
		return nil, fmt.Errorf(
			"invalid Xcode StoreKit certificate chain: expected exactly 1 certificate, got %d",
			len(header.X5C),
		)
	}

	// =====================================================
	// LOAD XCODE TEST CERTIFICATE
	// =====================================================

	certificatePath :=
		os.Getenv(
			"APPLE_STOREKIT_TEST_CERT_PATH",
		)

	if certificatePath == "" {
		return nil, fmt.Errorf(
			"APPLE_STOREKIT_TEST_CERT_PATH is not configured",
		)
	}

	trustedCertificateBytes, err :=
		os.ReadFile(
			certificatePath,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to read Xcode StoreKit test certificate: %w",
			err,
		)
	}

	trustedCertificate, err :=
		parseAppleRootCertificate(
			trustedCertificateBytes,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse Xcode StoreKit test certificate: %w",
			err,
		)
	}

	// =====================================================
	// DECODE JWS CERTIFICATE
	// =====================================================

	certificateBytes, err :=
		base64.StdEncoding.DecodeString(
			header.X5C[0],
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode Xcode StoreKit certificate: %w",
			err,
		)
	}

	signingCertificate, err :=
		x509.ParseCertificate(
			certificateBytes,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse Xcode StoreKit certificate: %w",
			err,
		)
	}

	// =====================================================
	// ENSURE THE JWS CERTIFICATE IS OUR TRUSTED
	// XCODE CERTIFICATE
	// =====================================================
	//
	// Do NOT trust an arbitrary self-signed certificate.
	//
	// The certificate embedded in the JWS must be exactly
	// the certificate we explicitly exported from Xcode
	// and configured on the server.
	//

	fmt.Printf(
		"Apple StoreKit Xcode certificate fingerprints: JWS=%s TRUSTED=%s\n",
		certificateFingerprint(signingCertificate),
		certificateFingerprint(trustedCertificate),
	)

	if !signingCertificate.Equal(trustedCertificate) {
		return nil, fmt.Errorf(
			"Xcode StoreKit signing certificate does not match trusted StoreKit test certificate",
		)
	}

	// =====================================================
	// VERIFY CERTIFICATE VALIDITY
	// =====================================================

	if _, err :=
		signingCertificate.Verify(
			x509.VerifyOptions{
				Roots: func() *x509.CertPool {
					roots := x509.NewCertPool()
					roots.AddCert(trustedCertificate)
					return roots
				}(),

				KeyUsages: []x509.ExtKeyUsage{
					x509.ExtKeyUsageAny,
				},
			},
		); err != nil {

		return nil, fmt.Errorf(
			"Xcode StoreKit certificate verification failed: %w",
			err,
		)
	}

	// =====================================================
	// VERIFY JWS SIGNATURE
	// =====================================================

	signature, err :=
		base64.RawURLEncoding.DecodeString(
			signaturePart,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode Xcode JWS signature: %w",
			err,
		)
	}

	// ES256 uses:
	//
	// R = 32 bytes
	// S = 32 bytes
	//
	// JWS uses the JOSE format:
	//
	// R || S
	//

	if len(signature) != 64 {
		return nil, fmt.Errorf(
			"invalid ES256 JWS signature length: expected 64 bytes, got %d",
			len(signature),
		)
	}

	signingInput :=
		[]byte(
			headerPart + "." + payloadPart,
		)

	hash :=
		sha256.Sum256(
			signingInput,
		)

	publicKey, ok :=
		signingCertificate.PublicKey.(*ecdsa.PublicKey)

	if !ok {
		return nil, fmt.Errorf(
			"Xcode StoreKit signing certificate does not contain an ECDSA public key",
		)
	}

	r := new(big.Int).SetBytes(
		signature[:32],
	)

	s := new(big.Int).SetBytes(
		signature[32:],
	)

	if !ecdsa.Verify(
		publicKey,
		hash[:],
		r,
		s,
	) {
		return nil, fmt.Errorf(
			"Xcode StoreKit JWS signature verification failed",
		)
	}

	return payloadBytes, nil
}

// =====================================================
// VERIFY APPLE JWS
// =====================================================
//
// Sandbox / Production verification.
//
// This is your existing Apple certificate-chain
// verification and remains separate from Xcode.
//

func verifyAppleJWS(
	header appleJWSHeader,
	headerPart string,
	payloadPart string,
	signaturePart string,
	payloadBytes []byte,
) ([]byte, error) {

	// =====================================================
	// CERTIFICATE CHAIN
	// =====================================================

	if len(header.X5C) < 2 {
		return nil, fmt.Errorf(
			"invalid Apple certificate chain: expected at least 2 certificates, got %d",
			len(header.X5C),
		)
	}

	certificates :=
		make(
			[]*x509.Certificate,
			len(header.X5C),
		)

	for i, encodedCertificate := range header.X5C {

		certificateBytes, err :=
			base64.StdEncoding.DecodeString(
				encodedCertificate,
			)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to decode Apple certificate %d: %w",
				i,
				err,
			)
		}

		certificate, err :=
			x509.ParseCertificate(
				certificateBytes,
			)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse Apple certificate %d: %w",
				i,
				err,
			)
		}

		certificates[i] = certificate
	}

	leafCertificate :=
		certificates[0]

	// =====================================================
	// LOAD TRUSTED APPLE ROOT
	// =====================================================

	rootPath :=
		os.Getenv(
			"APPLE_ROOT_CA_PATH",
		)

	if rootPath == "" {
		return nil, fmt.Errorf(
			"APPLE_ROOT_CA_PATH is not configured",
		)
	}

	trustedRootBytes, err :=
		os.ReadFile(
			rootPath,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to read Apple root certificate: %w",
			err,
		)
	}

	trustedRoot, err :=
		parseAppleRootCertificate(
			trustedRootBytes,
		)

	if err != nil {
		return nil, err
	}

	// =====================================================
	// BUILD TRUSTED ROOT POOL
	// =====================================================

	roots :=
		x509.NewCertPool()

	roots.AddCert(
		trustedRoot,
	)

	// =====================================================
	// BUILD INTERMEDIATE POOL
	// =====================================================

	intermediates :=
		x509.NewCertPool()

	for i := 1; i < len(certificates); i++ {

		// Never treat an incoming certificate as a trusted
		// root. They are only intermediates.
		intermediates.AddCert(
			certificates[i],
		)
	}

	// =====================================================
	// VERIFY CERTIFICATE CHAIN
	// =====================================================

	_, err =
		leafCertificate.Verify(
			x509.VerifyOptions{
				Roots:         roots,
				Intermediates: intermediates,

				KeyUsages: []x509.ExtKeyUsage{
					x509.ExtKeyUsageAny,
				},
			},
		)

	if err != nil {
		return nil, fmt.Errorf(
			"Apple certificate chain verification failed: %w",
			err,
		)
	}

	// =====================================================
	// VERIFY JWS SIGNATURE
	// =====================================================

	signature, err :=
		base64.RawURLEncoding.DecodeString(
			signaturePart,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode JWS signature: %w",
			err,
		)
	}

	// ES256 uses the JOSE format:
	//
	// R (32 bytes) || S (32 bytes)
	//

	if len(signature) != 64 {
		return nil, fmt.Errorf(
			"invalid ES256 JWS signature length: expected 64 bytes, got %d",
			len(signature),
		)
	}

	signingInput :=
		[]byte(
			headerPart + "." + payloadPart,
		)

	hash :=
		sha256.Sum256(
			signingInput,
		)

	publicKey, ok :=
		leafCertificate.PublicKey.(*ecdsa.PublicKey)

	if !ok {
		return nil, fmt.Errorf(
			"Apple signing certificate does not contain an ECDSA public key",
		)
	}

	r := new(big.Int).SetBytes(
		signature[:32],
	)

	s := new(big.Int).SetBytes(
		signature[32:],
	)

	if !ecdsa.Verify(
		publicKey,
		hash[:],
		r,
		s,
	) {
		return nil, fmt.Errorf(
			"Apple JWS signature verification failed",
		)
	}

	return payloadBytes, nil
}

// =====================================================
// PARSE APPLE ROOT CERTIFICATE
// =====================================================
//
// Supports both DER and PEM encoded certificates.
//

func parseAppleRootCertificate(
	certificateBytes []byte,
) (*x509.Certificate, error) {

	// Try DER first.
	certificate, err :=
		x509.ParseCertificate(
			certificateBytes,
		)

	if err == nil {
		return certificate, nil
	}

	// Try PEM.
	block, _ :=
		pem.Decode(
			certificateBytes,
		)

	if block == nil {
		return nil, fmt.Errorf(
			"failed to parse trusted Apple root certificate",
		)
	}

	certificate, err =
		x509.ParseCertificate(
			block.Bytes,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse trusted Apple root certificate: %w",
			err,
		)
	}

	return certificate, nil
}

func certificateFingerprint(cert *x509.Certificate) string {
	hash := sha256.Sum256(cert.Raw)

	parts := make([]string, len(hash))

	for i, b := range hash {
		parts[i] = fmt.Sprintf("%02X", b)
	}

	return strings.Join(parts, ":")
}
