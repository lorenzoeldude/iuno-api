package storekit

import (
	"time"

	"iuno-api/utils"
)

const (
	MonthlyProductID = "com.iunoni.premium.monthly"
	YearlyProductID  = "com.iunoni.premium.yearly"
)

// =====================================================
// APPLE TRANSACTION
// =====================================================

type Transaction struct {
	Payload *utils.AppleTransactionPayload
}

// =====================================================
// VALIDATE PRODUCT
// =====================================================

func IsValidProduct(productID string) bool {

	return productID == MonthlyProductID ||
		productID == YearlyProductID
}

// =====================================================
// DETERMINE PREMIUM STATUS
// =====================================================

func IsPremium(
	payload *utils.AppleTransactionPayload,
) bool {

	if payload.RevocationDate != nil {
		return false
	}

	if payload.ExpiresDate <= 0 {
		return true
	}

	expirationDate := time.UnixMilli(
		payload.ExpiresDate,
	)

	return expirationDate.After(time.Now())
}

// =====================================================
// SUBSCRIPTION STATUS
// =====================================================

func SubscriptionStatus(
	isPremium bool,
) string {

	if isPremium {
		return "active"
	}

	return "expired"
}

// =====================================================
// PURCHASE DATE
// =====================================================

func PurchaseDate(
	milliseconds int64,
) *time.Time {

	if milliseconds <= 0 {
		return nil
	}

	t := time.UnixMilli(
		milliseconds,
	)

	return &t
}

// =====================================================
// EXPIRATION DATE
// =====================================================

func ExpirationDate(
	milliseconds int64,
) *time.Time {

	if milliseconds <= 0 {
		return nil
	}

	t := time.UnixMilli(
		milliseconds,
	)

	return &t
}

// =====================================================
// PAYMENT AMOUNT
// =====================================================
//
// Apple provides price in milliunits.
//
// Example:
//
// 4990 milliunits = 499 cents = $4.99
//
// Our payment_transactions.amount is stored in cents.
//

func PaymentAmount(
	payload *utils.AppleTransactionPayload,
) int {

	if payload.Price == nil {
		return 0
	}

	return int(
		*payload.Price / 10,
	)
}

// =====================================================
// PAYMENT CURRENCY
// =====================================================

func PaymentCurrency(
	payload *utils.AppleTransactionPayload,
) string {

	if payload.Currency == nil ||
		*payload.Currency == "" {

		return "USD"
	}

	return *payload.Currency
}

// =====================================================
// XCODE TRANSACTION
// =====================================================
//
// Xcode local StoreKit transactions use transaction ID "0".
//
// These should not be stored as real payment transactions.
//

func IsXcodeTransaction(
	payload *utils.AppleTransactionPayload,
) bool {

	return payload.TransactionID == "0"
}