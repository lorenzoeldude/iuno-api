//
//  apple_server_notifications.go
//  iuno-api
//

package storekit

import (
	"encoding/json"
	"fmt"
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

	payloadBytes, err :=
		VerifyAppleJWS(signedPayload)

	if err != nil {
		return nil, fmt.Errorf(
			"Apple server notification JWS verification failed: %w",
			err,
		)
	}

	var notification AppleServerNotification

	if err :=
		json.Unmarshal(
			payloadBytes,
			&notification,
		); err != nil {

		return nil, fmt.Errorf(
			"failed to decode Apple server notification payload: %w",
			err,
		)
	}

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

	if notification.NotificationUUID == "" {
		return nil, fmt.Errorf(
			"notification UUID is missing",
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
