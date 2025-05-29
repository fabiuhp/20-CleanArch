package graph

import (
	"context"
	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
)

func NewResolver(createOrderUseCase usecase.CreateOrderUseCase, listOrdersUseCase usecase.ListOrdersUseCase) *Resolver {
	return &Resolver{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (r *Resolver) Orders(ctx context.Context, page int, perPage int) ([]*entity.Order, error) {
	input := usecase.ListOrdersInput{
		Page:    page,
		PerPage: perPage,
	}

	output, err := r.ListOrdersUseCase.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	orders := make([]*entity.Order, len(output.Orders))
	for i, order := range output.Orders {
		orders[i] = &order
	}

	return orders, nil
}
