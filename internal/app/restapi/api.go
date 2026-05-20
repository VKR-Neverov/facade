package restapi

import (
	"encoding/json"
	"net/http"

	invoicesservice "github.com/fidesy-pay/facade/internal/pkg/services/invoices-service"
	clients_service "github.com/fidesy-pay/facade/pkg/clients-service"
	"github.com/go-chi/chi"
)

const APIKey = "api-key"

type Service struct {
	clientsClient   clients_service.ClientsServiceClient
	invoicesService *invoicesservice.Service
}

func New(
	router *chi.Mux,
	clientsClient clients_service.ClientsServiceClient,
	invoicesService *invoicesservice.Service,
) *Service {
	s := &Service{
		clientsClient:   clientsClient,
		invoicesService: invoicesService,
	}

	router.Post("/api/invoice", s.createInvoice)
	router.Get("/api/invoice/{invoice_id}", s.getInvoice)

	return s
}

type createInvoiceBody struct {
	USDAmount float64 `json:"usd_amount"`
}

func (s *Service) createInvoice(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get(APIKey)
	if apiKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var invoiceBody *createInvoiceBody
	err := json.NewDecoder(r.Body).Decode(&invoiceBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	clients, err := s.clientsClient.ListClients(r.Context(), &clients_service.ListClientsRequest{
		Filter: &clients_service.ListClientsRequest_Filter{
			ApiKeyIn: []string{apiKey},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	if len(clients.Clients) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	client := clients.Clients[0]

	invoiceID, err := s.invoicesService.CreateInvoice(r.Context(), client.Id, invoiceBody.USDAmount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	body, _ := json.Marshal(map[string]string{
		"invoice_id": invoiceID,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Service) getInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "invoice_id")
	if invoiceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invoices, err := s.invoicesService.ListInvoices(r.Context(), invoicesservice.ListInvoicesFilter{
		IDIn: []string{invoiceID},
	}, 1, 1)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	if len(invoices) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, _ := json.Marshal(InvoiceModelFromResp(invoices[0]))
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
