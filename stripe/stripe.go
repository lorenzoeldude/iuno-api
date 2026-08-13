package stripe

import (
	"log"
	"os"

	stripeSDK "github.com/stripe/stripe-go/v86"
)

var Client *stripeSDK.Client

func Init() {

	key := os.Getenv("STRIPE_SECRET_KEY")

	if key == "" {
		log.Fatal("STRIPE_SECRET_KEY is missing")
	}

	Client = stripeSDK.NewClient(key)

	log.Println("Stripe initialized")
}