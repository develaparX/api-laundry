package service

import (
	"clean-laundry/model"
	"clean-laundry/model/dto"
	"clean-laundry/repository"
)

// BillService interface defines the methods that our service must implement.
type BillService interface {
	CreateNewBill(payload dto.BillRequest) (model.Bill, error)
}

// billService struct contains the repositories and other services needed by BillService.
type billService struct {
	repo            repository.BillRepository
	userService     UserService
	productService  ProductService
	customerService CustomerService
}

// CreateNewBill implements the BillService interface.
func (b *billService) CreateNewBill(payload dto.BillRequest) (model.Bill, error) {
	// Check if customer exists
	customer, err := b.customerService.FindById(payload.CustomerId)
	if err != nil {
		return model.Bill{}, err
	}

	// Check if user exists
	user, err := b.userService.FindById(payload.UserId)
	if err != nil {
		return model.Bill{}, err
	}

	// Check if products exist and create bill details
	var billDetails []model.BillDetail
	for _, bd := range payload.BillDetails {
		product, err := b.productService.FindById(bd.Product.Id)
		if err != nil {
			return model.Bill{}, err
		}
		billDetails = append(billDetails, model.BillDetail{Product: product, Qty: bd.Qty, Price: product.Price})
	}
	newPayload := model.Bill{
		Customer:    customer,
		User:        user,
		BillDetails: billDetails,
	}

	// Call repository to create a new bill
	bill, err := b.repo.Create(newPayload)
	if err != nil {
		return model.Bill{}, err
	}
	return bill, nil
}

// NewBillService creates a new instance of BillService.
func NewBillService(repo repository.BillRepository, uS UserService, pS ProductService, cS CustomerService) BillService {
	return &billService{
		repo:            repo,
		userService:     uS,
		productService:  pS,
		customerService: cS,
	}
}
