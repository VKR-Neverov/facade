package invoicesservice

import (
	"context"
	"fmt"

	"github.com/fidesy-pay/facade/internal/pkg/model"
	desc "github.com/fidesy-pay/facade/pkg/invoices-service"
	"github.com/samber/lo"
)

type Service struct {
	invoicesClient desc.InvoicesServiceClient
}

func New(
	invoicesClient desc.InvoicesServiceClient,
) *Service {
	return &Service{
		invoicesClient: invoicesClient,
	}
}

func (s *Service) CreateInvoice(ctx context.Context, clientID string, usdAmount float64) (string, error) {
	invoiceResp, err := s.invoicesClient.CreateInvoice(ctx, &desc.CreateInvoiceRequest{
		ClientId:  clientID,
		UsdAmount: usdAmount,
	})
	if err != nil {
		return "", fmt.Errorf("invoicesClient.CreateInvoice: %w", err)
	}

	return invoiceResp.GetId(), nil
}

type UpdateInvoiceParams struct {
	ID            string
	Chain         string
	Token         string
	PayerClientID *string
}

func (s *Service) UpdateInvoice(ctx context.Context, params UpdateInvoiceParams) (*desc.Invoice, error) {
	invoiceResp, err := s.invoicesClient.UpdateInvoice(ctx, &desc.UpdateInvoiceRequest{
		Id:            params.ID,
		Chain:         params.Chain,
		Token:         params.Token,
		PayerClientId: params.PayerClientID,
	})
	if err != nil {
		return nil, fmt.Errorf("invoicesClient.UpdateInvoice: %w", err)
	}

	return invoiceResp.GetInvoice(), nil
}

func (s *Service) CheckInvoice(ctx context.Context, input model.CheckInvoiceInput) (*desc.Invoice, error) {
	invoiceResp, err := s.invoicesClient.CheckInvoice(ctx, &desc.CheckInvoiceRequest{
		Id: input.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("invoicesClient.CheckInvoice: %w", err)
	}

	return invoiceResp.GetInvoice(), nil
}

type ListInvoicesFilter struct {
	IDIn       []string
	ClientIDIn []string
	StatusIn   []model.InvoiceStatus
}

func (s *Service) ListInvoices(ctx context.Context, filter ListInvoicesFilter, page, perPage uint64) ([]*desc.Invoice, error) {
	reqFilter := &desc.ListInvoicesRequest_Filter{}
	if len(filter.IDIn) > 0 {
		reqFilter.IdIn = filter.IDIn
	}

	if len(filter.ClientIDIn) > 0 {
		reqFilter.ClientIdIn = filter.ClientIDIn
	}

	if len(filter.StatusIn) > 0 {
		reqFilter.InvoiceStatusIn = lo.Map(filter.StatusIn, func(status model.InvoiceStatus, index int) desc.InvoiceStatus {
			return desc.InvoiceStatus(desc.InvoiceStatus_value[status.String()])
		})
	}

	invoicesResp, err := s.invoicesClient.ListInvoices(ctx, &desc.ListInvoicesRequest{
		Filter:  reqFilter,
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, fmt.Errorf("invoicesClient.ListInvoices: %w", err)
	}

	return invoicesResp.GetInvoices(), nil
}
