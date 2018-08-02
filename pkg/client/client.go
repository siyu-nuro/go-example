package client

import (
	"context"

)

// CommunicationClient is a client that handles all customer facing communications
type CommunicationClient interface {
	ValidatePhoneNumber(ctx context.Context, phoneNumber string) bool
}
