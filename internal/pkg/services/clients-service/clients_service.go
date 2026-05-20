package clientsservice

import (
	"context"
	"fmt"
	desc "github.com/fidesy-pay/facade/pkg/clients-service"
)

type Service struct {
	clientsServiceClient desc.ClientsServiceClient
}

func New(
	clientsServiceClient desc.ClientsServiceClient,
) *Service {
	return &Service{
		clientsServiceClient: clientsServiceClient,
	}
}

func (s *Service) ListClients(ctx context.Context, ids []string) ([]*desc.Client, error) {
	clientsResp, err := s.clientsServiceClient.ListClients(ctx, &desc.ListClientsRequest{
		Filter: &desc.ListClientsRequest_Filter{
			IdIn: ids,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("clientsServiceClient.ListClients: %w", err)
	}

	return clientsResp.GetClients(), nil
}
