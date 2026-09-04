//
//  apple_server_notifications.go
//  iuno-api
//

package storekit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// =====================================================
// APP STORE SERVER NOTIFICATION V2
// =====================================================

type AppleServerNotification struct {
	NotificationType string                       `json:"notificationType"`
	Subtype          string                       `json:"subtype"`
	NotificationUUID string                       `json:"notificationUUID"`
	Version          string                       `json:"version"`
	SignedDate       int64                        `json:"signedDate"`
	Data             *AppleServerNotificationData `json:"data"`
}

type AppleServerNotificationData struct {
	AppAppleID            int64  `json:"appAppleId"`
	BundleID              string `json:"bundleId"`
	BundleVersion         string `json:"bundleVersion"`
	Environment           string `json:"environment"`
	SignedRenewalInfo     string `json:"signedRenewalInfo"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
	Status                int    `json:"status"`
}

// =====================================================
// VERIFY / DECODE SERVER NOTIFICATION
// =====================================================
//
// IMPORTANT:
//
// This currently decodes the Apple JWS payload.
//
// It does NOT yet cryptographically verify Apple's
// certificate chain and signature.
//
// Cryptographic verification will be added in
// apple_jws.go.
//
// =====================================================

func VerifyServerNotification(
	signedPayload string,
) (*AppleServerNotification, error) {

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
	// DECODE JSON
	// =====================================================

	var notification AppleServerNotification

	if err :=
		json.Unmarshal(
			payloadBytes,
			&notification,
		); err != nil {

		return nil, fmt.Errorf(
			"failed to decode notification payload: %w",
			err,
		)
	}

	// =====================================================
	// BASIC VALIDATION
	// =====================================================

	if notification.Version != "" &&
		notification.Version != "2.0" {

		return nil, fmt.Errorf(
			"unsupported notification version: %s",
			notification.Version,
		)
	}

	if notification.NotificationType == "" {

		return nil, fmt.Errorf(
			"notification type is missing",
		)
	}

	return &notification, nil
}

// =====================================================
// DATA ENVIRONMENT
// =====================================================

func (n *AppleServerNotification) DataEnvironment() string {

	if n.Data == nil {
		return ""
	}

	return n.Data.Environment
}
