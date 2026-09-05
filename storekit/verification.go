package storekit

import (
	"encoding/json"
	"fmt"

	"iuno-api/utils"
)

// =====================================================
// VERIFY APPLE TRANSACTION
// =====================================================
//
// This is the StoreKit entry point for validating a
// signed transaction.
//
// The JWS is cryptographically verified before its
// payload is trusted.
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
	// CRYPTOGRAPHICALLY VERIFY APPLE JWS
	// =====================================================

	payloadBytes, err :=
		VerifyAppleJWS(
			signedTransaction,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"Apple JWS verification failed: %w",
			err,
		)
	}

	// =====================================================
	// DECODE VERIFIED PAYLOAD
	// =====================================================

	var payload utils.AppleTransactionPayload

	if err :=
		json.Unmarshal(
			payloadBytes,
			&payload,
		); err != nil {

		return nil, fmt.Errorf(
			"failed to parse verified Apple transaction payload: %w",
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
