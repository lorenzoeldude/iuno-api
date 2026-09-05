//
// apple_api.go
// iuno-api
//

package storekit

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"iuno-api/utils"
)

const (
	appleProductionURL = "https://api.storekit.apple.com"
	appleSandboxURL    = "https://api.storekit-sandbox.apple.com"

	appleBundleID = "com.iunoni.IUNONI"
)

// =====================================================
// APPLE API CLIENT
// =====================================================

type AppleAPIClient struct {
	KeyID      string
	IssuerID   string
	PrivateKey *ecdsa.PrivateKey
	BundleID   string
}

// =====================================================
// CREATE CLIENT
// =====================================================

func NewAppleAPIClient() (*AppleAPIClient, error) {

	keyID := os.Getenv("APPLE_IAP_KEY_ID")
	if keyID == "" {
		return nil, fmt.Errorf(
			"APPLE_IAP_KEY_ID is not configured",
		)
	}

	issuerID := os.Getenv("APPLE_IAP_ISSUER_ID")
	if issuerID == "" {
		return nil, fmt.Errorf(
			"APPLE_IAP_ISSUER_ID is not configured",
		)
	}

	privateKeyString := os.Getenv("APPLE_IAP_PRIVATE_KEY")
	if privateKeyString == "" {
		return nil, fmt.Errorf(
			"APPLE_IAP_PRIVATE_KEY is not configured",
		)
	}

	privateKey, err := parseApplePrivateKey(
		privateKeyString,
	)
	if err != nil {
		return nil, err
	}

	return &AppleAPIClient{
		KeyID:      keyID,
		IssuerID:   issuerID,
		PrivateKey: privateKey,
		BundleID:   appleBundleID,
	}, nil
}

// =====================================================
// PARSE PRIVATE KEY
// =====================================================

func parseApplePrivateKey(
	key string,
) (*ecdsa.PrivateKey, error) {

	key = strings.TrimSpace(key)

	block, _ := pem.Decode(
		[]byte(key),
	)

	if block == nil {
		return nil, fmt.Errorf(
			"failed to decode Apple private key PEM",
		)
	}

	privateKey, err :=
		x509.ParsePKCS8PrivateKey(
			block.Bytes,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse Apple private key: %w",
			err,
		)
	}

	ecdsaKey, ok :=
		privateKey.(*ecdsa.PrivateKey)

	if !ok {
		return nil, fmt.Errorf(
			"Apple private key is not an ECDSA private key",
		)
	}

	return ecdsaKey, nil
}

// =====================================================
// CREATE JWT
// =====================================================
//
// Apple App Store Server API authentication uses an
// ES256 JWT containing:
//
//   iss = Issuer ID
//   iat = issued-at timestamp
//   exp = expiration timestamp
//   aud = appstoreconnect-v1
//   bid = Bundle ID
//
// The key ID is placed in the JWT header as "kid".
//

func (c *AppleAPIClient) createJWT() (
	string,
	error,
) {

	now := time.Now()

	claims := jwt.MapClaims{
		"iss": c.IssuerID,
		"iat": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
		"aud": "appstoreconnect-v1",
		"bid": c.BundleID,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		claims,
	)

	token.Header["kid"] = c.KeyID
	token.Header["typ"] = "JWT"

	signedToken, err :=
		token.SignedString(
			c.PrivateKey,
		)

	if err != nil {
		return "", fmt.Errorf(
			"failed to sign Apple API JWT: %w",
			err,
		)
	}

	return signedToken, nil
}

// =====================================================
// GET TRANSACTION
// =====================================================
//
// Fetches a transaction from Apple's App Store Server API.
//
// Production:
//     https://api.storekit.apple.com
//
// Sandbox:
//     https://api.storekit-sandbox.apple.com
//
// Returns Apple's signedTransactionInfo JWS.
//

func (c *AppleAPIClient) GetTransaction(
	transactionID string,
	environment string,
) (string, error) {

	if transactionID == "" {
		return "", fmt.Errorf(
			"transaction ID is empty",
		)
	}

	var baseURL string

	switch environment {

	case "Production":
		baseURL = appleProductionURL

	case "Sandbox":
		baseURL = appleSandboxURL

	default:
		return "", fmt.Errorf(
			"unsupported Apple environment: %s",
			environment,
		)
	}

	// =====================================================
	// CREATE JWT
	// =====================================================

	token, err := c.createJWT()
	if err != nil {
		return "", err
	}

	// =====================================================
	// CREATE REQUEST
	// =====================================================

	url := baseURL +
		"/inApps/v1/transactions/" +
		transactionID

	request, err :=
		http.NewRequest(
			http.MethodGet,
			url,
			nil,
		)

	if err != nil {
		return "", fmt.Errorf(
			"failed to create Apple API request: %w",
			err,
		)
	}

	request.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	request.Header.Set(
		"Accept",
		"application/json",
	)

	// =====================================================
	// SEND REQUEST
	// =====================================================

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	response, err :=
		client.Do(request)

	if err != nil {
		return "", fmt.Errorf(
			"Apple API request failed: %w",
			err,
		)
	}

	defer response.Body.Close()

	// =====================================================
	// READ RESPONSE
	// =====================================================

	body, err :=
		io.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf(
			"failed to read Apple API response: %w",
			err,
		)
	}

	// =====================================================
	// CHECK STATUS
	// =====================================================

	if response.StatusCode < 200 ||
		response.StatusCode >= 300 {

		return "", fmt.Errorf(
			"Apple API returned HTTP %d: %s",
			response.StatusCode,
			string(body),
		)
	}

	// =====================================================
	// PARSE RESPONSE
	// =====================================================

	var result struct {
		SignedTransactionInfo string `json:"signedTransactionInfo"`
	}

	if err :=
		json.Unmarshal(
			body,
			&result,
		); err != nil {

		return "", fmt.Errorf(
			"failed to parse Apple API response: %w",
			err,
		)
	}

	if result.SignedTransactionInfo == "" {
		return "", fmt.Errorf(
			"Apple API response contains no signedTransactionInfo",
		)
	}

	return result.SignedTransactionInfo, nil
}

// =====================================================
// VERIFY TRANSACTION WITH APPLE
// =====================================================
//
// This performs an authoritative server-side lookup:
//
// 1. Ask Apple for the transaction.
// 2. Verify Apple's returned JWS.
// 3. Decode the verified payload.
// 4. Verify the transaction belongs to IUNONI.
// 5. Verify the returned transaction ID matches.
//

func (c *AppleAPIClient) VerifyTransactionWithApple(
	transactionID string,
	environment string,
) (*Transaction, error) {

	// =====================================================
	// GET TRANSACTION FROM APPLE
	// =====================================================

	signedTransactionInfo, err :=
		c.GetTransaction(
			transactionID,
			environment,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get transaction from Apple: %w",
			err,
		)
	}

	// =====================================================
	// VERIFY APPLE JWS
	// =====================================================

	payloadBytes, err :=
		VerifyAppleJWS(
			signedTransactionInfo,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"Apple API transaction JWS verification failed: %w",
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
			"failed to decode Apple API transaction: %w",
			err,
		)
	}

	// =====================================================
	// VALIDATE BUNDLE ID
	// =====================================================

	if payload.BundleID != appleBundleID {
		return nil, fmt.Errorf(
			"Apple transaction belongs to bundle %q, expected %q",
			payload.BundleID,
			appleBundleID,
		)
	}

	// =====================================================
	// VALIDATE TRANSACTION ID
	// =====================================================

	if payload.TransactionID != transactionID {
		return nil, fmt.Errorf(
			"Apple returned transaction %q, expected %q",
			payload.TransactionID,
			transactionID,
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

	if payload.Environment != environment {
		return nil, fmt.Errorf(
			"Apple transaction environment %q does not match expected %q",
			payload.Environment,
			environment,
		)
	}

	// =====================================================
	// CREATE TRUSTED TRANSACTION
	// =====================================================

	return &Transaction{
		Payload: &payload,
	}, nil
}
