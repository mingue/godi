package invoice

import "log"

type InvoiceRepository interface {
	GetAll() []Invoice
}

var _ InvoiceRepository = &InvoiceRepositoryImpl{}

type InvoiceRepositoryImpl struct {
	requestContext RequestContext
}

func NewInvoiceRepositoryImpl(requestContext RequestContext) *InvoiceRepositoryImpl {
	return &InvoiceRepositoryImpl{
		requestContext: requestContext,
	}
}

func (r *InvoiceRepositoryImpl) GetAll() []Invoice {
	log.Printf("Printing value: %v, User Agent: %v, Request count: %v", r.requestContext.SomeValue, r.requestContext.UserAgent, r.requestContext.Counter)
	return []Invoice{
		{},
		{},
	}
}
