package invoice

type InvoiceService interface {
	GetAll() []Invoice
}

var _ InvoiceService = &InvoiceServiceImpl{}

type InvoiceServiceImpl struct {
	repo InvoiceRepository
}

func NewInvoiceServiceImpl(repo InvoiceRepository) *InvoiceServiceImpl{
	return &InvoiceServiceImpl{
		repo: repo,
	}
}

func (s *InvoiceServiceImpl) GetAll() []Invoice {
	return s.repo.GetAll()
}
