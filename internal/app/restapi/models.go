package restapi

import (
	"time"

	invoices_service "github.com/fidesy-pay/facade/pkg/invoices-service"
)

type Invoice struct {
	ID          string    `json:"id,omitempty"`
	ClientID    string    `json:"client_id,omitempty"`
	UsdAmount   float64   `json:"usd_amount,omitempty"`
	TokenAmount float64   `json:"token_amount,omitempty"`
	Chain       string    `json:"chain,omitempty"`
	Token       string    `json:"token,omitempty"`
	Status      string    `json:"status,omitempty"`
	Address     string    `json:"address,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func InvoiceModelFromResp(invoice *invoices_service.Invoice) *Invoice {
	if invoice == nil {
		return nil
	}

	return &Invoice{
		ID:          invoice.Id,
		ClientID:    invoice.ClientId,
		UsdAmount:   invoice.UsdAmount,
		TokenAmount: invoice.TokenAmount,
		Chain:       invoice.Chain,
		Token:       invoice.Token,
		Status:      invoice.Status.String(),
		Address:     invoice.Address,
		CreatedAt:   invoice.CreatedAt.AsTime(),
	}
}
