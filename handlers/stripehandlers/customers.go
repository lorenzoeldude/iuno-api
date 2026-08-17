package stripehandlers

import (
	"context"
	"log"
	"errors"
	"fmt"

	"iuno-api/db"
	"iuno-api/stripe"

	"github.com/jackc/pgx/v5"
	stripeSDK "github.com/stripe/stripe-go/v86"
)

// =====================================================
// GET STRIPE CUSTOMER ID
//
// Given an IUNO user:
//
// user_id → Stripe customer ID
//
// Returns an empty string if no Stripe customer exists.
// =====================================================

func GetStripeCustomerID(
	ctx context.Context,
	userID int,
) (string, error) {

	if userID <= 0 {
		return "", fmt.Errorf("invalid user ID: %d", userID)
	}

	var customerID string

	err := db.Pool.QueryRow(
		ctx,
		`
		SELECT provider_customer_id
		FROM billing_customers
		WHERE user_id = $1
		  AND provider = 'stripe'
		LIMIT 1
		`,
		userID,
	).Scan(&customerID)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf(
			"failed to get Stripe customer for user %d: %w",
			userID,
			err,
		)
	}

	return customerID, nil
}

// =====================================================
// SAVE STRIPE CUSTOMER
//
// Creates/updates:
//
// IUNO user → Stripe customer
// =====================================================

func SaveStripeCustomer(
	ctx context.Context,
	userID int,
	customerID string,
) error {

	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if customerID == "" {
		return errors.New("Stripe customer ID is empty")
	}

	_, err := db.Pool.Exec(
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
		return fmt.Errorf(
			"failed to save Stripe customer: %w",
			err,
		)
	}

	return nil
}

// =====================================================
// GET OR CREATE STRIPE CUSTOMER
//
// This is the main function used by checkout.
//
// 1. Look for an existing customer.
// 2. If found, return it.
// 3. Otherwise create one in Stripe.
// 4. Save the relationship in our database.
// 5. Return the Stripe customer ID.
// =====================================================
func GetOrCreateStripeCustomer(
		ctx context.Context,
		userID int,
	) (string, error) {

		log.Printf(
			"GetOrCreateStripeCustomer: user_id=%d",
			userID,
		)

		// =================================================
		// CHECK DATABASE FIRST
		// =================================================

		customerID, err := GetStripeCustomerID(
			ctx,
			userID,
		)

		if err != nil {
			return "", err
		}

		if customerID != "" {
			return customerID, nil
		}

		// =================================================
		// CREATE CUSTOMER IN STRIPE
		// =================================================

		params := &stripeSDK.CustomerCreateParams{
			Metadata: map[string]string{
				"user_id": fmt.Sprintf("%d", userID),
			},
		}

		stripeCustomer, err := stripe.Client.V1Customers.Create(
			ctx,
			params,
		)

		if err != nil {
			return "", fmt.Errorf(
				"failed to create Stripe customer: %w",
				err,
			)
		}

		if stripeCustomer == nil || stripeCustomer.ID == "" {
			return "", errors.New(
				"Stripe returned an empty customer",
			)
		}

		// =================================================
		// SAVE CUSTOMER MAPPING
		// =================================================

		if err := SaveStripeCustomer(
			ctx,
			userID,
			stripeCustomer.ID,
		); err != nil {
			return "", err
		}

		return stripeCustomer.ID, nil
	}

// =====================================================
// GET USER ID FROM STRIPE CUSTOMER
//
// Stripe customer ID → IUNO user ID
// =====================================================

func GetUserIDByStripeCustomerID(
	ctx context.Context,
	customerID string,
) (int, error) {

	if customerID == "" {
		return 0, errors.New(
			"Stripe customer ID is empty",
		)
	}

	var userID int

	err := db.Pool.QueryRow(
		ctx,
		`
		SELECT user_id
		FROM billing_customers
		WHERE provider = 'stripe'
		  AND provider_customer_id = $1
		LIMIT 1
		`,
		customerID,
	).Scan(&userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf(
			"no IUNO user found for Stripe customer %s",
			customerID,
		)
	}

	if err != nil {
		return 0, fmt.Errorf(
			"failed to find user for Stripe customer %s: %w",
			customerID,
			err,
		)
	}

	return userID, nil
}
