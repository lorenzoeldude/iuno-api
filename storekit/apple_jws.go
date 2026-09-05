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
// VERIFY APPLE JWS
// =====================================================
//
// Verifies:
//
// 1. JWS structure
// 2. ES256 algorithm
// 3. Apple certificate chain
// 4. Certificate validity
// 5. JWS signature
//
// The trusted Apple root certificate comes from our own
// configured certificate file.
//
// The root certificate supplied by the JWS is NEVER trusted
// merely because Apple supplied it.
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
	// CERTIFICATE CHAIN
	// =====================================================
	//
	// x5c is expected to contain Apple's certificate chain,
	// with the leaf certificate first.
	//
	// We require at least:
	//
	//   leaf
	//   intermediate
	//
	// The trusted root is loaded separately from disk.
	//

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

	// Support both:
	//
	//   DER certificate
	//
	// and:
	//
	//   PEM certificate
	//

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

	roots.AddCert(trustedRoot)

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

	if !ecdsa.VerifyASN1(
		publicKey,
		hash[:],
		signature,
	) {
		return nil, fmt.Errorf(
			"Apple JWS signature verification failed",
		)
	}

	// =====================================================
	// DECODE VERIFIED PAYLOAD
	// =====================================================

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
		pemDecode(
			certificateBytes,
		)

	if block == nil {
		return nil, fmt.Errorf(
			"failed to parse trusted Apple root certificate",
		)
	}

	certificate, err =
		x509.ParseCertificate(
			block,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse trusted Apple root certificate: %w",
			err,
		)
	}

	return certificate, nil
}

// =====================================================
// PEM DECODE
// =====================================================

func pemDecode(
	data []byte,
) ([]byte, error) {

	block, _ :=
		pem.Decode(
			data,
		)

	if block == nil {
		return nil, fmt.Errorf(
			"invalid PEM certificate",
		)
	}

	return block.Bytes, nil
}
