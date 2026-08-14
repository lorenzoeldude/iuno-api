package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// =====================================================
// APPLE TRANSACTION PAYLOAD
// =====================================================
//
// Represents the decoded payload contained inside Apple's
// signed transaction JWS.
//
// IMPORTANT:
//
// This file currently DECODEs the JWS payload but does NOT
// cryptographically verify Apple's signature.
//
// That is acceptable for Xcode StoreKit development.
//
// Before accepting real production purchases, the JWS must
// be cryptographically verified against Apple's signing
// certificate chain.
// =====================================================

type AppleTransactionPayload struct {
	TransactionID         string `json:"transactionId"`
	OriginalTransactionID string `json:"originalTransactionId"`

	ProductID string `json:"productId"`

	PurchaseDate int64 `json:"purchaseDate"`
	ExpiresDate  int64 `json:"expiresDate"`

	RevocationDate *int64 `json:"revocationDate"`

	Type string `json:"type"`

	InAppOwnershipType string `json:"inAppOwnershipType"`

	Environment string `json:"environment"`

	Price    *int64  `json:"price"`
	Currency *string `json:"currency"`

	Storefront   string `json:"storefront"`
	StorefrontID string `json:"storefrontId"`

	TransactionReason string `json:"transactionReason"`

	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier"`
}

// =====================================================
// DECODE APPLE TRANSACTION
// =====================================================
//
// Apple's signed transaction is a JWS:
//
//     header.payload.signature
//
// We decode the payload portion and convert it into
// AppleTransactionPayload.
//
// The signature is NOT verified here yet.
//
// =====================================================

func DecodeAppleTransaction(
	jws string,
) (*AppleTransactionPayload, error) {

	// =====================================================
	// VALIDATE JWS STRUCTURE
	// =====================================================

	parts := strings.Split(
		jws,
		".",
	)

	if len(parts) != 3 {

		return nil, fmt.Errorf(
			"invalid JWS: expected 3 parts, got %d",
			len(parts),
		)
	}

	// =====================================================
	// VALIDATE PARTS ARE NOT EMPTY
	// =====================================================

	if parts[0] == "" {

		return nil, fmt.Errorf(
			"invalid JWS: header is empty",
		)
	}

	if parts[1] == "" {

		return nil, fmt.Errorf(
			"invalid JWS: payload is empty",
		)
	}

	if parts[2] == "" {

		return nil, fmt.Errorf(
			"invalid JWS: signature is empty",
		)
	}

	// =====================================================
	// DECODE PAYLOAD
	// =====================================================

	payloadBytes, err :=
		base64.RawURLEncoding.DecodeString(
			parts[1],
		)

	if err != nil {

		return nil, fmt.Errorf(
			"failed to decode JWS payload: %w",
			err,
		)
	}

	// =====================================================
	// PARSE JSON
	// =====================================================

	var payload AppleTransactionPayload

	if err := json.Unmarshal(
		payloadBytes,
		&payload,
	); err != nil {

		return nil, fmt.Errorf(
			"failed to parse Apple transaction payload: %w",
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
	// RETURN PAYLOAD
	// =====================================================

	return &payload, nil
}
