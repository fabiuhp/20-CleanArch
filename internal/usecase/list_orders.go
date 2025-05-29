package usecase

import (
	"context"
	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type ListOrdersInput struct {
	Page    int
	PerPage int
}

type ListOrdersOutput struct {
	Orders []entity.Order
	Total  int
}

type ListOrdersUseCase interface {
	Execute(ctx context.Context, input ListOrdersInput) (*ListOrdersOutput, error)
}

type listOrdersUseCase struct {
	orderRepository entity.OrderRepository
}

func NewListOrdersUseCase(orderRepository entity.OrderRepository) ListOrdersUseCase {
	return &listOrdersUseCase{
		orderRepository: orderRepository,
	}
}

func (u *listOrdersUseCase) Execute(ctx context.Context, input ListOrdersInput) (*ListOrdersOutput, error) {
	orders, total, err := u.orderRepository.List(ctx, input.Page, input.PerPage)
	if err != nil {
		return nil, err
	}
	return &ListOrdersOutput{
		Orders: orders,
		Total:  total,
	}, nil
}
