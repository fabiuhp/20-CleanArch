package database

import (
	"context"
	"database/sql"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	return r.Create(context.Background(), order)
}

func (r *OrderRepository) Create(ctx context.Context, order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) List(ctx context.Context, page int, perPage int) ([]entity.Order, int, error) {
	var orders []entity.Order

	totalRows, err := r.Db.Query("Select count(*) from orders")
	if err != nil {
		return nil, 0, err
	}
	defer totalRows.Close()
	var total int
	totalRows.Next()
	err = totalRows.Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.Db.Query("Select id, price, tax, final_price from orders LIMIT ? OFFSET ?", perPage, (page-1)*perPage)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		err = rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}
	return orders, total, nil
}
