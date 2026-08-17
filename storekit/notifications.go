package storekit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"iuno-api/utils"
)

// =====================================================
// APPLE SERVER NOTIFICATION
// =====================================================
//
// Apple sends App Store Server Notifications as a JSON
// object containing a signedPayload.
//
// The signedPayload is itself a JWS.
//
// This file currently:
//
//   1. Parses the notification JSON
//   2. Decodes the signed notification payload
//   3. Decodes the embedded signed transaction
//   4. Returns the useful StoreKit data
//
// IMPORTANT:
//
// Cryptographic verification of Apple's notification JWS
// is NOT implemented here yet.
//
// Before Production notifications are trusted, the JWS
// signature and Apple's certificate chain must be verified.
// =====================================================


// =====================================================
// APPLE NOTIFICATION REQUEST
// =====================================================

type Notification struct {
	SignedPayload string `json:"signedPayload"`
}


// =====================================================
// DECODED NOTIFICATION
// =====================================================

type NotificationPayload struct {
	NotificationType string `json:"notificationType"`

	Subtype string `json:"subtype"`

	NotificationUUID string `json:"notificationUUID"`

	Version string `json:"version"`

	Environment string `json:"environment"`

	Data *NotificationData `json:"data"`
}


// =====================================================
// NOTIFICATION DATA
// =====================================================

type NotificationData struct {
	AppAppleID string `json:"appAppleId"`

	BundleID string `json:"bundleId"`

	BundleVersion string `json:"bundleVersion"`

	Environment string `json:"environment"`

	SignedTransactionInfo string `json:"signedTransactionInfo"`

	SignedRenewalInfo string `json:"signedRenewalInfo"`
}


// =====================================================
// DECODE APPLE NOTIFICATION
// =====================================================
//
// Decodes the signedPayload contained in Apple's
// App Store Server Notification.
//
// Cryptographic verification is intentionally left for
// the verification layer.
// =====================================================

func DecodeNotification(
	signedPayload string,
) (*NotificationPayload, error) {

	if signedPayload == "" {
		return nil, fmt.Errorf(
			"signed notification payload is empty",
		)
	}

	// =====================================================
	// JWS STRUCTURE
	// =====================================================

	parts := strings.Split(
		signedPayload,
		".",
	)

	if len(parts) != 3 {
		return nil, fmt.Errorf(
			"invalid notification JWS: expected 3 parts, got %d",
			len(parts),
		)
	}

	if parts[0] == "" {
		return nil, fmt.Errorf(
			"invalid notification JWS: header is empty",
		)
	}

	if parts[1] == "" {
		return nil, fmt.Errorf(
			"invalid notification JWS: payload is empty",
		)
	}

	if parts[2] == "" {
		return nil, fmt.Errorf(
			"invalid notification JWS: signature is empty",
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
			"failed to decode notification payload: %w",
			err,
		)
	}

	// =====================================================
	// PARSE JSON
	// =====================================================

	var payload NotificationPayload

	if err := json.Unmarshal(
		payloadBytes,
		&payload,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to parse notification payload: %w",
			err,
		)
	}

	// =====================================================
	// VALIDATE REQUIRED FIELDS
	// =====================================================

	if payload.NotificationType == "" {
		return nil, fmt.Errorf(
			"notificationType missing",
		)
	}

	if payload.NotificationUUID == "" {
		return nil, fmt.Errorf(
			"notificationUUID missing",
		)
	}

	if payload.Environment == "" {
		return nil, fmt.Errorf(
			"notification environment missing",
		)
	}

	return &payload, nil
}


// =====================================================
// DECODE TRANSACTION FROM NOTIFICATION
// =====================================================
//
// Apple includes signedTransactionInfo inside the
// notification's data object.
//
// We reuse the same transaction decoder used by the
// normal StoreKit transaction endpoint.
// =====================================================

func DecodeNotificationTransaction(
	payload *NotificationPayload,
) (*Transaction, error) {

	if payload == nil {
		return nil, fmt.Errorf(
			"notification payload is nil",
		)
	}

	if payload.Data == nil {
		return nil, fmt.Errorf(
			"notification data is missing",
		)
	}

	if payload.Data.SignedTransactionInfo == "" {
		return nil, fmt.Errorf(
			"signedTransactionInfo is missing",
		)
	}

	transactionPayload, err :=
		decodeSignedTransaction(
			payload.Data.SignedTransactionInfo,
		)

	if err != nil {
		return nil, err
	}

	return &Transaction{
		Payload: transactionPayload,
	}, nil
}


// =====================================================
// DECODE SIGNED TRANSACTION
// =====================================================
//
// Internal helper used by notification processing.
//
// This only decodes the JWS payload.
//
// Cryptographic verification is handled separately.
// =====================================================

func decodeSignedTransaction(
	signedTransaction string,
) (*utils.AppleTransactionPayload, error) {

	if signedTransaction == "" {
		return nil, fmt.Errorf(
			"signed transaction is empty",
		)
	}

	payload, err :=
		utils.DecodeAppleTransaction(
			signedTransaction,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode signed transaction: %w",
			err,
		)
	}

	return payload, nil
}


// =====================================================
// NOTIFICATION HELPERS
// =====================================================

func IsSubscriptionNotification(
	notificationType string,
) bool {

	switch notificationType {

	case "SUBSCRIBED",
		"DID_RENEW",
		"DID_FAIL_TO_RENEW",
		"EXPIRED",
		"GRACE_PERIOD_EXPIRED",
		"REVOKE",
		"OFFER_REDEEMED":

		return true

	default:

		return false
	}
}


// =====================================================
// PREMIUM NOTIFICATION
// =====================================================
//
// Returns whether this notification represents a state
// where the subscription should normally remain premium.
//
// The actual expiration/revocation check should still be
// performed using the signed transaction information.
// =====================================================

func IsPremiumNotification(
	notificationType string,
) bool {

	switch notificationType {

	case "SUBSCRIBED",
		"DID_RENEW",
		"OFFER_REDEEMED":

		return true

	default:

		return false
	}
}