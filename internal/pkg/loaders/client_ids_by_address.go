package loaders

import (
	"context"
	"fmt"

	crypto_service "github.com/fidesy-pay/facade/pkg/crypto-service"
)

func ClientIDsByAddresses(cryptoClient crypto_service.CryptoServiceClient) func(ctx context.Context, addresses []string) ([]string, []error) {
	return func(ctx context.Context, addresses []string) ([]string, []error) {
		walletsResp, err := cryptoClient.ListWallets(ctx, &crypto_service.ListWalletsRequest{
			Filter: &crypto_service.ListWalletsRequest_Filter{
				AddressIn: addresses,
			},
		})
		if err != nil {
			return nil, []error{fmt.Errorf("cryptoClient.ListWallets: %w", err)}
		}

		wallets := walletsResp.Wallets

		clientIDByAddress := make(map[string]string, len(wallets))
		for _, wallet := range wallets {
			clientIDByAddress[wallet.Address] = wallet.ClientId
		}

		output := make([]string, len(addresses))
		errors := make([]error, 0)
		for index, address := range addresses {
			clientID, ok := clientIDByAddress[address]
			if !ok {
				errors = append(errors, fmt.Errorf("clientID not found %s", address))
				output[index] = ""
				continue
			}

			output[index] = clientID
		}

		return output, nil
	}
}
