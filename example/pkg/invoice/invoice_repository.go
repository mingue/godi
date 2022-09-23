package invoice

import (
	"context"
	"fmt"
)

type InvoiceRepository interface {
	GetAll() []Invoice
}

var _ InvoiceRepository = &InvoiceRepositoryImpl{}

type InvoiceRepositoryImpl struct {
	ctx context.Context
}

func NewInvoiceRepositoryImpl(ctx context.Context) *InvoiceRepositoryImpl {
	return &InvoiceRepositoryImpl{
		ctx: ctx,
	}
}

func (r *InvoiceRepositoryImpl) GetAll() []Invoice {
	val := r.ctx.Value("Some")
	fmt.Printf("Printing value from repo: %v\n", val)
	return []Invoice{
		{},
		{},
	}
}
