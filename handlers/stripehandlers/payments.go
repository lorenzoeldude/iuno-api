package stripehandlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"iuno-api/stripe"

	"github.com/jackc/pgx/v5"
	stripeSDK "github.com/stripe/stripe-go/v86"
)

// =====================================================
// FIND PAYMENT OWNER
// =====================================================

func FindPaymentOwner(
		ctx context.Context,
		tx pgx.Tx,
		customerID string,
		metadata map[string]string,
	) (int, error) {

		log.Printf(
			"FindPaymentOwner: customerID=%s metadata=%v",
			customerID,
			metadata,
		)

		var userID int

		// =================================================
		// 1. LOOK UP EXISTING CUSTOMER MAPPING
		// =================================================

		if customerID != "" {
			err := tx.QueryRow(
				ctx,
				`
				SELECT user_id
				FROM billing_customers
				WHERE provider = 'stripe'
				AND provider_customer_id = $1
				`,
				customerID,
			).Scan(&userID)

			if err == nil {
				return userID, nil
			}

			if !errors.Is(err, pgx.ErrNoRows) {
				return 0, err
			}
		}

		// =================================================
		// 2. TRY EVENT METADATA
		// =================================================

		userIDString := ""

		if metadata != nil {
			userIDString = metadata["user_id"]
		}

		if userIDString != "" {

			userID, err := strconv.Atoi(userIDString)

			if err != nil {
				return 0, errors.New(
					"invalid user_id in payment metadata",
				)
			}

			// Verify that the user actually exists.
			var exists bool

			err = tx.QueryRow(
				ctx,
				`
				SELECT EXISTS (
					SELECT 1
					FROM users
					WHERE id = $1
				)
				`,
				userID,
			).Scan(&exists)

			if err != nil {
				return 0, err
			}

			if !exists {
				return 0, errors.New(
					"user_id from payment metadata does not exist",
				)
			}

			// Repair the Stripe customer mapping.
			if customerID != "" {
				_, err = tx.Exec(
					ctx,
					`
					INSERT INTO billing_customers (
						user_id,
						provider,
						provider_customer_id
					)
					VALUES (
						$1,
						'stripe',
						$2
					)
					ON CONFLICT (user_id, provider)
					DO UPDATE SET
						provider_customer_id = EXCLUDED.provider_customer_id
					`,
					userID,
					customerID,
				)

				if err != nil {
					return 0, err
				}
			}

			return userID, nil
		}

		// =================================================
		// 3. RETRIEVE STRIPE CUSTOMER
		// =================================================

		if customerID == "" {
			return 0, errors.New(
				"payment owner not found",
			)
		}

		customer, err := stripe.Client.V1Customers.Retrieve(
			ctx,
			customerID,
			&stripeSDK.CustomerRetrieveParams{},
		)

		if err != nil {
			return 0, err
		}

		if customer == nil {
			return 0, errors.New(
				"Stripe customer not found",
			)
		}

		log.Printf(
			"Stripe customer retrieved: id=%s metadata=%v",
			customer.ID,
			customer.Metadata,
		)

		// =================================================
		// 4. GET USER ID FROM STRIPE CUSTOMER METADATA
		// =================================================

		userIDString = customer.Metadata["user_id"]

		if userIDString == "" {
			return 0, errors.New(
				"Stripe customer has no user_id metadata",
			)
		}

		userID, err = strconv.Atoi(userIDString)

		if err != nil {
			return 0, errors.New(
				"invalid user_id in Stripe customer metadata",
			)
		}

		// Verify that the user actually exists.
		var exists bool

		err = tx.QueryRow(
			ctx,
			`
			SELECT EXISTS (
				SELECT 1
				FROM users
				WHERE id = $1
			)
			`,
			userID,
		).Scan(&exists)

		if err != nil {
			return 0, err
		}

		if !exists {
			return 0, errors.New(
				"user_id from Stripe customer metadata does not exist",
			)
		}

		// =================================================
		// 5. REPAIR CUSTOMER MAPPING
		// =================================================

		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO billing_customers (
				user_id,
				provider,
				provider_customer_id
			)
			VALUES (
				$1,
				'stripe',
				$2
			)
			ON CONFLICT (user_id, provider)
			DO UPDATE SET
				provider_customer_id = EXCLUDED.provider_customer_id
			`,
			userID,
			customerID,
		)

		if err != nil {
			return 0, err
		}

		log.Printf(
			"FindPaymentOwner: resolved user_id=%d for stripe_customer_id=%s",
			userID,
			customerID,
		)

		return userID, nil
	}

// =====================================================
// GET PAYMENT PRODUCT
// =====================================================

func GetPaymentProduct(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	metadata map[string]string,
) (string, error) {

	if metadata != nil {
		if priceID := metadata["price_id"]; priceID != "" {
			return priceID, nil
		}
	}

	var productID string

	err := tx.QueryRow(
		ctx,
		`
		SELECT product_id
		FROM subscriptions
		WHERE user_id = $1
		  AND provider = 'stripe'
		ORDER BY created_at DESC
		LIMIT 1
		`,
		userID,
	).Scan(&productID)

	if err == nil && productID != "" {
		return productID, nil
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	return "stripe_payment", nil
}

// =====================================================
// RECORD SUCCESSFUL PAYMENT
// =====================================================

func RecordSuccessfulPayment(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	paymentIntentID string,
	productID string,
	amount int64,
	currency string,
) error {

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
			'stripe',
			$2,
			$3,
			$4,
			$5,
			'paid'
		)
		ON CONFLICT (
			provider,
			provider_transaction_id
		)
		DO UPDATE SET
			status = 'paid',
			product_id = EXCLUDED.product_id,
			amount = EXCLUDED.amount,
			currency = EXCLUDED.currency
		`,
		userID,
		paymentIntentID,
		productID,
		amount,
		currency,
	)

	return err
}

// =====================================================
// RECORD FAILED PAYMENT
// =====================================================

func RecordFailedPayment(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	paymentIntentID string,
	productID string,
	amount int64,
	currency string,
) error {

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
			'stripe',
			$2,
			$3,
			$4,
			$5,
			'failed'
		)
		ON CONFLICT (
			provider,
			provider_transaction_id
		)
		DO UPDATE SET
			status = 'failed',
			product_id = EXCLUDED.product_id,
			amount = EXCLUDED.amount,
			currency = EXCLUDED.currency
		`,
		userID,
		paymentIntentID,
		productID,
		amount,
		currency,
	)

	return err
}

// =====================================================
// ACTIVATE PREMIUM
// =====================================================

func ActivatePremium(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
) error {

	_, err := tx.Exec(
		ctx,
		`
		UPDATE users
		SET is_premium = true
		WHERE id = $1
		`,
		userID,
	)

	return err
}

// =====================================================
// PAYMENT INTENT SUCCEEDED
// =====================================================

func handlePaymentIntentSucceeded(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {

	var paymentIntentData struct {
		ID       string            `json:"id"`
		Amount   int64             `json:"amount"`
		Currency string            `json:"currency"`
		Customer string            `json:"customer"`
		Metadata map[string]string `json:"metadata"`
	}

	if err := json.Unmarshal(
		event.Data.Raw,
		&paymentIntentData,
	); err != nil {
		return err
	}

	if paymentIntentData.ID == "" {
		return errors.New(
			"payment_intent.succeeded has no payment intent ID",
		)
	}

	userID, err := FindPaymentOwner(
		ctx,
		tx,
		paymentIntentData.Customer,
		paymentIntentData.Metadata,
	)
	if err != nil {
		return err
	}

	productID, err := GetPaymentProduct(
		ctx,
		tx,
		userID,
		paymentIntentData.Metadata,
	)
	if err != nil {
		return err
	}

	if err := RecordSuccessfulPayment(
		ctx,
		tx,
		userID,
		paymentIntentData.ID,
		productID,
		paymentIntentData.Amount,
		paymentIntentData.Currency,
	); err != nil {
		return err
	}

	return ActivatePremium(
		ctx,
		tx,
		userID,
	)
}

// =====================================================
// PAYMENT INTENT FAILED
// =====================================================

func handlePaymentIntentFailed(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {

	var paymentIntentData struct {
		ID       string            `json:"id"`
		Amount   int64             `json:"amount"`
		Currency string            `json:"currency"`
		Customer string            `json:"customer"`
		Metadata map[string]string `json:"metadata"`
	}

	if err := json.Unmarshal(
		event.Data.Raw,
		&paymentIntentData,
	); err != nil {
		return err
	}

	if paymentIntentData.ID == "" {
		return errors.New(
			"payment_intent.payment_failed has no payment intent ID",
		)
	}

	userID, err := FindPaymentOwner(
		ctx,
		tx,
		paymentIntentData.Customer,
		paymentIntentData.Metadata,
	)
	if err != nil {
		return err
	}

	productID, err := GetPaymentProduct(
		ctx,
		tx,
		userID,
		paymentIntentData.Metadata,
	)
	if err != nil {
		return err
	}

	return RecordFailedPayment(
		ctx,
		tx,
		userID,
		paymentIntentData.ID,
		productID,
		paymentIntentData.Amount,
		paymentIntentData.Currency,
	)
}

// =====================================================
// INVOICE PAYMENT PAID
// =====================================================

func handleInvoicePaymentPaid(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {

	var invoicePayment struct {
		ID         string `json:"id"`
		Invoice    string `json:"invoice"`
		AmountPaid int64  `json:"amount_paid"`
		Currency   string `json:"currency"`
		Status     string `json:"status"`

		Payment struct {
			Type          string `json:"type"`
			PaymentIntent string `json:"payment_intent"`
			Charge        string `json:"charge"`
			PaymentRecord string `json:"payment_record"`
		} `json:"payment"`
	}

	if err := json.Unmarshal(
		event.Data.Raw,
		&invoicePayment,
	); err != nil {
		return err
	}

	if invoicePayment.Payment.Type != "payment_intent" ||
		invoicePayment.Payment.PaymentIntent == "" {
		return nil
	}

	paymentIntentID := invoicePayment.Payment.PaymentIntent

	var existingUserID int

	err := tx.QueryRow(
		ctx,
		`
		SELECT user_id
		FROM payment_transactions
		WHERE provider = 'stripe'
		  AND provider_transaction_id = $1
		`,
		paymentIntentID,
	).Scan(&existingUserID)

	if err == nil {
		return nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	pi, err := getStripePaymentIntent(ctx, paymentIntentID)
	if err != nil {
		return err
	}

	if pi == nil {
		return errors.New("Stripe returned nil PaymentIntent")
	}

	customerID := ""

	if pi.Customer != nil {
		customerID = pi.Customer.ID
	}

	userID, err := FindPaymentOwner(
		ctx,
		tx,
		customerID,
		pi.Metadata,
	)
	if err != nil {
		return err
	}

	productID, err := GetPaymentProduct(
		ctx,
		tx,
		userID,
		pi.Metadata,
	)
	if err != nil {
		return err
	}

	return RecordSuccessfulPayment(
		ctx,
		tx,
		userID,
		paymentIntentID,
		productID,
		pi.Amount,
		string(pi.Currency),
	)
}

// =====================================================
// INVOICE PAID
// =====================================================

func handleInvoicePaid(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {

	var invoice struct {
		ID           string `json:"id"`
		Subscription string `json:"subscription"`
		Customer     string `json:"customer"`
		AmountPaid   int64  `json:"amount_paid"`
		Currency     string `json:"currency"`
		Status       string `json:"status"`
	}

	if err := json.Unmarshal(
		event.Data.Raw,
		&invoice,
	); err != nil {
		return err
	}

	var userID int

	if invoice.Subscription != "" {
		err := tx.QueryRow(
			ctx,
			`
			SELECT user_id
			FROM subscriptions
			WHERE provider = 'stripe'
			  AND provider_subscription_id = $1
			`,
			invoice.Subscription,
		).Scan(&userID)

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	if userID == 0 && invoice.Customer != "" {
		err := tx.QueryRow(
			ctx,
			`
			SELECT user_id
			FROM billing_customers
			WHERE provider = 'stripe'
			  AND provider_customer_id = $1
			`,
			invoice.Customer,
		).Scan(&userID)

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	if userID == 0 {
		return errors.New("cannot identify invoice owner")
	}

	productID := "stripe_subscription"

	if invoice.Subscription != "" {
		var subscriptionProduct string

		err := tx.QueryRow(
			ctx,
			`
			SELECT product_id
			FROM subscriptions
			WHERE provider = 'stripe'
			  AND provider_subscription_id = $1
			`,
			invoice.Subscription,
		).Scan(&subscriptionProduct)

		if err == nil && subscriptionProduct != "" {
			productID = subscriptionProduct
		} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

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
			'stripe',
			$2,
			$3,
			$4,
			$5,
			'paid'
		)
		ON CONFLICT (
			provider,
			provider_transaction_id
		)
		DO UPDATE SET
			status = 'paid',
			product_id = EXCLUDED.product_id,
			amount = EXCLUDED.amount,
			currency = EXCLUDED.currency
		`,
		userID,
		invoice.ID,
		productID,
		invoice.AmountPaid,
		invoice.Currency,
	)

	if err != nil {
		return err
	}

	return ActivatePremium(
		ctx,
		tx,
		userID,
	)
}

// =====================================================
// INVOICE PAYMENT FAILED
// =====================================================

func handleInvoicePaymentFailed(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {

	var invoice struct {
		ID           string `json:"id"`
		Subscription string `json:"subscription"`
		AmountDue    int64  `json:"amount_due"`
		Currency     string `json:"currency"`
	}

	if err := json.Unmarshal(
		event.Data.Raw,
		&invoice,
	); err != nil {
		return err
	}

	if invoice.Subscription == "" {
		return nil
	}

	var userID int

	err := tx.QueryRow(
		ctx,
		`
		SELECT user_id
		FROM subscriptions
		WHERE provider = 'stripe'
		  AND provider_subscription_id = $1
		`,
		invoice.Subscription,
	).Scan(&userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

// =====================================================
// GET STRIPE PAYMENT INTENT
// =====================================================

func getStripePaymentIntent(
	ctx context.Context,
	paymentIntentID string,
) (*stripeSDK.PaymentIntent, error) {
	return stripe.Client.V1PaymentIntents.Retrieve(
		ctx,
		paymentIntentID,
		nil,
	)
}
