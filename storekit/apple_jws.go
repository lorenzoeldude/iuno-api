//
//  apple_jws.go
//  iuno-api
//

package storekit

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
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
// The Apple root certificate is trusted from our own
// local certificate file, NOT from the incoming JWS.
//
// =====================================================

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

	parts :=
		strings.Split(
			signedPayload,
			".",
		)

	if len(parts) != 3 {
		return nil, fmt.Errorf(
			"invalid JWS format",
		)
	}

	headerPart := parts[0]
	payloadPart := parts[1]
	signaturePart := parts[2]

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
	// ALGORITHM
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

	if len(header.X5C) != 3 {
		return nil, fmt.Errorf(
			"invalid Apple certificate chain: expected 3 certificates, got %d",
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
				"failed to decode certificate %d: %w",
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
				"failed to parse certificate %d: %w",
				i,
				err,
			)
		}

		certificates[i] = certificate
	}

	leafCertificate :=
		certificates[0]

	intermediateCertificate :=
		certificates[1]

	rootCertificate :=
		certificates[2]

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
		x509.ParseCertificate(
			trustedRootBytes,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse trusted Apple root certificate: %w",
			err,
		)
	}

	// =====================================================
	// VERIFY ROOT MATCH
	// =====================================================

	if !rootCertificate.Equal(
		trustedRoot,
	) {
		return nil, fmt.Errorf(
			"Apple JWS root certificate is not trusted",
		)
	}

	// =====================================================
	// CERTIFICATE CHAIN
	// =====================================================

	roots :=
		x509.NewCertPool()

	roots.AddCert(
		trustedRoot,
	)

	intermediates :=
		x509.NewCertPool()

	intermediates.AddCert(
		intermediateCertificate,
	)

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
	// DECODE PAYLOAD
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
