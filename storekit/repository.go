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
//
// The expiration date is the authoritative value for whether
// an Apple subscription is currently active.
//
// This intentionally checks current_period_end > now()
// rather than relying solely on the stored status field.
//
// That means an expired subscription cannot remain active
// simply because its status was not updated yet.
//

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
//
// One Apple subscription is represented by its
// OriginalTransactionID.
//
// Apple sends a new transaction ID for every renewal,
// while the OriginalTransactionID remains the same.
//
// Therefore we update the existing subscription rather
// than creating a new subscription for every renewal.
//
// IMPORTANT:
// An older transaction must never overwrite a newer
// current_period_end.
//

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
				CASE

					-- If the incoming transaction has no
					-- expiration date, preserve the existing
					-- expiration when one exists.
					--
					-- Otherwise use the incoming value.

					WHEN EXCLUDED.current_period_end IS NULL
						AND subscriptions.current_period_end IS NULL
					THEN EXCLUDED.status

					-- Incoming transaction has no expiration,
					-- but the existing subscription does.
					--
					-- Keep the existing subscription state.

					WHEN EXCLUDED.current_period_end IS NULL
					THEN
						CASE
							WHEN subscriptions.current_period_end > now()
							THEN 'active'
							ELSE 'expired'
						END

					-- Incoming transaction has an expiration.
					-- Use the incoming expiration only if it is
					-- newer than the currently stored expiration.

					WHEN subscriptions.current_period_end IS NULL
						OR EXCLUDED.current_period_end >
						   subscriptions.current_period_end
					THEN
						CASE
							WHEN EXCLUDED.current_period_end > now()
							THEN 'active'
							ELSE 'expired'
						END

					-- Incoming transaction is older.
					-- Keep the state corresponding to the
					-- existing newer expiration.

					ELSE
						CASE
							WHEN subscriptions.current_period_end > now()
							THEN 'active'
							ELSE 'expired'
						END

				END,

			current_period_start =
				CASE

					-- No existing expiration means there is
					-- nothing newer to protect.

					WHEN subscriptions.current_period_end IS NULL
						OR (
							EXCLUDED.current_period_end IS NOT NULL
							AND EXCLUDED.current_period_end >
								subscriptions.current_period_end
						)
					THEN EXCLUDED.current_period_start

					ELSE subscriptions.current_period_start

				END,

			current_period_end =
				CASE

					-- Incoming transaction has no expiration.
					-- Preserve an existing expiration.

					WHEN EXCLUDED.current_period_end IS NULL
						AND subscriptions.current_period_end IS NOT NULL
					THEN subscriptions.current_period_end

					-- Incoming transaction has a newer expiration.

					WHEN subscriptions.current_period_end IS NULL
						OR (
							EXCLUDED.current_period_end IS NOT NULL
							AND EXCLUDED.current_period_end >
								subscriptions.current_period_end
						)
					THEN EXCLUDED.current_period_end

					-- Incoming transaction is older.
					-- Keep the newer expiration.

					ELSE subscriptions.current_period_end

				END,

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
//
// Every Apple renewal has its own transaction ID.
//
// This table therefore stores each individual payment,
// while subscriptions stores the current subscription
// state.
//

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
	// These are not real App Store payment transactions.

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
//
// users.is_premium is a cached value.
//
// The actual Apple subscription state is stored in
// subscriptions.current_period_end.
//

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
