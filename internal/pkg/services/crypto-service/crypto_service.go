package cryptoservice

import (
	"context"
	"fmt"
	"github.com/fidesy-pay/facade/internal/pkg/model"
	crypto_service "github.com/fidesy-pay/facade/pkg/crypto-service"
	external_api "github.com/fidesy-pay/facade/pkg/external-api"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	cryptoServiceClient crypto_service.CryptoServiceClient
	externalAPI         external_api.ExternalAPIClient
}

func New(
	cryptoServiceClient crypto_service.CryptoServiceClient,
	externalAPI external_api.ExternalAPIClient,
) *Service {
	return &Service{
		cryptoServiceClient: cryptoServiceClient,
		externalAPI:         externalAPI,
	}
}

type GetBalanceParams struct {
	Address string
	Chain   string
	Token   string
}

func (s *Service) GetBalance(ctx context.Context, params GetBalanceParams) (*model.Balance, error) {
	errGroup := errgroup.Group{}

	var (
		balance    float64
		tokenPrice float64
	)

	errGroup.Go(func() error {
		balanceResp, err := s.cryptoServiceClient.GetBalance(ctx, &crypto_service.GetBalanceRequest{
			Address: params.Address,
			Chain:   params.Chain,
			Token:   params.Token,
		})
		if err != nil {
			return fmt.Errorf("cryptoServiceClient.GetBalance: %w", err)
		}

		balance = balanceResp.GetBalance()

		return nil
	})

	errGroup.Go(func() error {
		tokenPriceResp, err := s.externalAPI.GetPrice(ctx, &external_api.GetPriceRequest{
			Symbol: params.Token,
		})
		if err != nil {
			return fmt.Errorf("coinGeckoAPIClient.GetPrice: %w", err)
		}

		tokenPrice = tokenPriceResp.GetPriceUsd()

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	return &model.Balance{
		Token:      params.Token,
		Balance:    balance,
		UsdBalance: balance * tokenPrice,
	}, nil
}
