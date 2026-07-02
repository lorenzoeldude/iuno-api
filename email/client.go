package email

import (
	"os"

	"github.com/resend/resend-go/v2"
)

var Client *resend.Client

func Init() {
	Client = resend.NewClient(os.Getenv("RESEND_API_KEY"))
}