package storekit

import (
	"context"

	"github.com/jackc/pgx/v5"

	"iuno-api/db"
)

// =====================================================
// GET APPLE SUBSCRIPTION OWNER
// =====================================================

func GetSubscriptionOwner(
	ctx context.Context,
	tx pgx.Tx,
	providerSubscriptionID string,
) (int, error) {

	var userID int

	err := tx.QueryRow(
		ctx,
		`
		SELECT user_id
		FROM subscriptions
		WHERE provider = 'apple'
		  AND provider_subscription_id = $1
		`,
		providerSubscriptionID,
	).Scan(&userID)

	return userID, err
}

// =====================================================
// CHECK WHETHER USER HAS ACTIVE APPLE SUBSCRIPTION
// =====================================================

func HasActiveAppleSubscription(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
) (bool, error) {

	var exists bool

	err := tx.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM subscriptions
			WHERE user_id = $1
			  AND provider = 'apple'
			  AND status = 'active'
			  AND (
				  current_period_end IS NULL
				  OR current_period_end > now()
			  )
		)
		`,
		userID,
	).Scan(&exists)

	return exists, err
}

// =====================================================
// SAVE / UPDATE APPLE SUBSCRIPTION
// =====================================================

func SaveSubscription(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	payload *Transaction,
	isPremium bool,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO subscriptions (
			user_id,
			provider,
			provider_subscription_id,
			product_id,
			status,
			current_period_start,
			current_period_end,
			cancel_at_period_end
		)
		VALUES (
			$1,
			'apple',
			$2,
			$3,
			$4,
			$5,
			$6,
			false
		)
		ON CONFLICT (
			provider,
			provider_subscription_id
		)
		DO UPDATE SET

			product_id =
				EXCLUDED.product_id,

			status =
				EXCLUDED.status,

			current_period_start =
				EXCLUDED.current_period_start,

			current_period_end =
				EXCLUDED.current_period_end,

			cancel_at_period_end =
				EXCLUDED.cancel_at_period_end,

			updated_at =
				now()
		`,
		userID,
		payload.Payload.OriginalTransactionID,
		payload.Payload.ProductID,
		SubscriptionStatus(isPremium),
		PurchaseDate(
			payload.Payload.PurchaseDate,
		),
		ExpirationDate(
			payload.Payload.ExpiresDate,
		),
	)

	return err
}

// =====================================================
// SAVE / UPDATE APPLE PAYMENT TRANSACTION
// =====================================================

func SavePaymentTransaction(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	payload *Transaction,
	isPremium bool,
) error {

	// Xcode local StoreKit transactions use
	// transaction ID "0".
	//
	// These are not real payment transactions.

	if IsXcodeTransaction(
		payload.Payload,
	) {
		return nil
	}

	amount :=
		PaymentAmount(
			payload.Payload,
		)

	currency :=
		PaymentCurrency(
			payload.Payload,
		)

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO payment_transactions (
			user_id,
			provider,
			provider_transaction_id,
			product_id,
			amount,
			currency,
			status
		)
		VALUES (
			$1,
			'apple',
			$2,
			$3,
			$4,
			$5,
			$6
		)
		ON CONFLICT (
			provider,
			provider_transaction_id
		)
		DO UPDATE SET

			product_id =
				EXCLUDED.product_id,

			amount =
				EXCLUDED.amount,

			currency =
				EXCLUDED.currency,

			status =
				EXCLUDED.status
		`,
		userID,
		payload.Payload.TransactionID,
		payload.Payload.ProductID,
		amount,
		currency,
		SubscriptionStatus(isPremium),
	)

	return err
}

// =====================================================
// UPDATE USER PREMIUM
// =====================================================

func UpdateUserPremium(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	isPremium bool,
) error {

	_, err := tx.Exec(
		ctx,
		`
		UPDATE users
		SET is_premium = $1
		WHERE id = $2
		`,
		isPremium,
		userID,
	)

	return err
}

// =====================================================
// BEGIN TRANSACTION
// =====================================================

func BeginTransaction(
	ctx context.Context,
) (pgx.Tx, error) {

	return db.Pool.Begin(ctx)
}
