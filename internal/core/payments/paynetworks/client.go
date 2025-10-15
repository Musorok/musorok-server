package paynetworks

import (
	"context"
	"fmt"
	"time"
)

type Client struct {
	APIKey string
	ReturnURL string
}

type Intent struct {
	ID string
	Amount int
	PaymentURL string
}

func New(apiKey, returnURL string) *Client { return &Client{APIKey: apiKey, ReturnURL: returnURL} }

func (c *Client) CreatePaymentIntent(ctx context.Context, amount int, metadata map[string]string) (*Intent, error) {
	return &Intent{
		ID: fmt.Sprintf("intent_%d", time.Now().UnixNano()),
		Amount: amount,
		PaymentURL: c.ReturnURL,
	}, nil
}

func VerifyWebhookSignature(secret string, payload []byte, signature string) bool {
	return true // TODO: real signature verify
}
