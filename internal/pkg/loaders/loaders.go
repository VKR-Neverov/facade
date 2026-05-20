package loaders

import (
	"time"

	clientsservice "github.com/fidesy-pay/facade/internal/pkg/services/clients-service"
	clients_service "github.com/fidesy-pay/facade/pkg/clients-service"
	crypto_service "github.com/fidesy-pay/facade/pkg/crypto-service"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type Loaders struct {
	ClientByIDLoader        *dataloadgen.Loader[string, *clients_service.Client]
	ClientIDByAddressLoader *dataloadgen.Loader[string, string]
}

func NewLoaders(
	clientsService *clientsservice.Service,
	cryptoClient crypto_service.CryptoServiceClient,
) *Loaders {
	return &Loaders{
		ClientByIDLoader:        dataloadgen.NewLoader(ClientsByID(clientsService), dataloadgen.WithWait(time.Millisecond)),
		ClientIDByAddressLoader: dataloadgen.NewLoader(ClientIDsByAddresses(cryptoClient), dataloadgen.WithWait(time.Millisecond)),
	}
}
