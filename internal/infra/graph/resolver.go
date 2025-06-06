package graph

import (
	"github.com/devfullcycle/20-CleanArch/internal/infra/graph/model"
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewResolver(createOrderUseCase usecase.CreateOrderUseCase, listOrdersUseCase usecase.ListOrdersUseCase) *Resolver {
	return &Resolver{
		CreateOrderUseCase: &createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

type OrderRepository interface {
	CreateOrder(price float64, tax float64) (*model.Order, error)
	GetOrders() ([]*model.Order, error)
	GetOrder(id string) (*model.Order, error)
}
