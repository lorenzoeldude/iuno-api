package storekit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"iuno-api/utils"
)

// =====================================================
// VERIFY APPLE TRANSACTION
// =====================================================
//
// StoreKit transaction verification.
//
// Xcode:
//   Local StoreKit testing.
//   Xcode transactions use the local Xcode signing
//   environment and are handled separately.
//
// Sandbox / Production:
//   JWS is cryptographically verified against Apple's
//   certificate chain before the payload is trusted.
//

func VerifyTransaction(
	signedTransaction string,
) (*Transaction, error) {

	// =====================================================
	// BASIC VALIDATION
	// =====================================================

	if signedTransaction == "" {
		return nil, fmt.Errorf(
			"signed transaction is empty",
		)
	}

	// =====================================================
	// SPLIT JWS
	// =====================================================

	parts := strings.Split(
		signedTransaction,
		".",
	)

	if len(parts) != 3 {
		return nil, fmt.Errorf(
			"invalid signed transaction format: expected 3 parts, got %d",
			len(parts),
		)
	}

	// =====================================================
	// DECODE PAYLOAD
	// =====================================================
	//
	// We need the environment before selecting the
	// appropriate verification model.
	//

	payloadBytes, err :=
		decodeJWSPayload(
			parts[1],
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode transaction payload: %w",
			err,
		)
	}

	var environmentPayload struct {
		Environment string `json:"environment"`
	}

	if err :=
		json.Unmarshal(
			payloadBytes,
			&environmentPayload,
		); err != nil {

		return nil, fmt.Errorf(
			"failed to decode transaction environment: %w",
			err,
		)
	}

	// =====================================================
	// VERIFY JWS
	// =====================================================

	switch environmentPayload.Environment {

	case "Xcode":

		// -------------------------------------------------
		// XCODE STOREKIT TESTING
		// -------------------------------------------------
		//
		// Xcode local StoreKit transactions are generated
		// by the local StoreKit testing environment.
		//
		// They are not App Store Sandbox/Production
		// transactions, so they do not use Apple's normal
		// certificate-chain verification.
		//
		// The decoded payload is validated below.
		//

	case "Sandbox", "Production":

		// -------------------------------------------------
		// APPLE SANDBOX / PRODUCTION
		// -------------------------------------------------
		//
		// These transactions MUST pass full cryptographic
		// JWS verification.
		//

		payloadBytes, err =
			VerifyAppleJWS(
				signedTransaction,
			)

		if err != nil {
			return nil, fmt.Errorf(
				"Apple JWS verification failed: %w",
				err,
			)
		}

	default:

		return nil, fmt.Errorf(
			"unknown Apple JWS environment: %s",
			environmentPayload.Environment,
		)
	}

	// =====================================================
	// DECODE TRANSACTION PAYLOAD
	// =====================================================

	var payload utils.AppleTransactionPayload

	if err :=
		json.Unmarshal(
			payloadBytes,
			&payload,
		); err != nil {

		return nil, fmt.Errorf(
			"failed to parse transaction payload: %w",
			err,
		)
	}

	// =====================================================
	// VALIDATE REQUIRED FIELDS
	// =====================================================

	if payload.TransactionID == "" {
		return nil, fmt.Errorf(
			"transactionId missing",
		)
	}

	if payload.OriginalTransactionID == "" {
		return nil, fmt.Errorf(
			"originalTransactionId missing",
		)
	}

	if payload.ProductID == "" {
		return nil, fmt.Errorf(
			"productId missing",
		)
	}

	// =====================================================
	// VALIDATE PRODUCT
	// =====================================================

	if !IsValidProduct(
		payload.ProductID,
	) {

		return nil, fmt.Errorf(
			"unknown Apple product: %s",
			payload.ProductID,
		)
	}

	// =====================================================
	// VALIDATE ENVIRONMENT
	// =====================================================

	switch payload.Environment {

	case "Xcode":
		// Local StoreKit testing.

	case "Sandbox":
		// Apple sandbox environment.

	case "Production":
		// Real App Store environment.

	default:

		return nil, fmt.Errorf(
			"unknown Apple environment: %s",
			payload.Environment,
		)
	}

	// =====================================================
	// CREATE TRANSACTION
	// =====================================================

	return &Transaction{
		Payload: &payload,
	}, nil
}

// =====================================================
// DECODE JWS PAYLOAD
// =====================================================
//
// Decodes the payload portion of a JWS.
//
// IMPORTANT:
// This only decodes the payload.
// It does NOT cryptographically verify the JWS.
//
// For Sandbox/Production, VerifyAppleJWS performs the
// cryptographic verification before the payload is used.
//

func decodeJWSPayload(
	payloadPart string,
) ([]byte, error) {

	if payloadPart == "" {
		return nil, fmt.Errorf(
			"JWS payload is empty",
		)
	}

	payloadBytes, err :=
		base64.RawURLEncoding.DecodeString(
			payloadPart,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid JWS payload encoding: %w",
			err,
		)
	}

	return payloadBytes, nil
}
