package loaders

import (
	"context"
	"fmt"
	clientsservice "github.com/fidesy-pay/facade/internal/pkg/services/clients-service"
	desc "github.com/fidesy-pay/facade/pkg/clients-service"
)

func ClientsByID(clientsService *clientsservice.Service) func(ctx context.Context, ids []string) ([]*desc.Client, []error) {
	return func(ctx context.Context, ids []string) ([]*desc.Client, []error) {
		clients, err := clientsService.ListClients(ctx, ids)
		if err != nil {
			return nil, []error{fmt.Errorf("clientsService.ListClients: %w", err)}
		}

		clientByID := make(map[string]*desc.Client, len(clients))
		for _, client := range clients {
			clientByID[client.Id] = client
		}

		output := make([]*desc.Client, len(ids))
		errors := make([]error, 0)
		for index, id := range ids {
			client, ok := clientByID[id]
			if !ok {
				errors = append(errors, fmt.Errorf("analytic not found %s", id))
				output[index] = &desc.Client{}
				continue
			}

			output[index] = client
		}

		return output, nil
	}
}
