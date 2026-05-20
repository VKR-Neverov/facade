package graph

import (
	cryptoservice "github.com/fidesy-pay/facade/internal/pkg/services/crypto-service"
	invoicesservice "github.com/fidesy-pay/facade/internal/pkg/services/invoices-service"
	auth_service "github.com/fidesy-pay/facade/pkg/auth-service"
	clients_service "github.com/fidesy-pay/facade/pkg/clients-service"
	crypto_service "github.com/fidesy-pay/facade/pkg/crypto-service"
	email_service "github.com/fidesy-pay/facade/pkg/email-service"
	external_api "github.com/fidesy-pay/facade/pkg/external-api"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	clientsClient   clients_service.ClientsServiceClient
	cryptoClient    crypto_service.CryptoServiceClient
	externalAPI     external_api.ExternalAPIClient
	authClient      auth_service.AuthServiceClient
	invoicesService *invoicesservice.Service
	cryptoService   *cryptoservice.Service
	emailClient     email_service.EmailServiceClient
}

func NewResolver(
	clientsClient clients_service.ClientsServiceClient,
	cryptoServiceClient crypto_service.CryptoServiceClient,
	externalAPI external_api.ExternalAPIClient,
	authClient auth_service.AuthServiceClient,
	invoicesService *invoicesservice.Service,
	cryptoService *cryptoservice.Service,
	emailClient email_service.EmailServiceClient,
) *Resolver {
	return &Resolver{
		clientsClient:   clientsClient,
		cryptoClient:    cryptoServiceClient,
		externalAPI:     externalAPI,
		authClient:      authClient,
		invoicesService: invoicesService,
		cryptoService:   cryptoService,
		emailClient:     emailClient,
	}
}
