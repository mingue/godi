package invoice

import (
	"encoding/json"
	"net/http"
)

var _ http.Handler = &GetInvoicesHandler{}

type GetInvoicesHandler struct{
	service InvoiceService
}

func NewInvoiceHandler(service InvoiceService) *GetInvoicesHandler {
	return &GetInvoicesHandler{
		service: service,
	}
}

func (h *GetInvoicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	invoices := h.service.GetAll()
	body, _ := json.Marshal(invoices)
	w.Write(body)
}

func (h *GetInvoicesHandler) SomethingSilly() {

}