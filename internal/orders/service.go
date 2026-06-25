package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	repo "github.com/kodoktempur666/go-ecom-api/internal/adapters/postgresql/sqlc"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock = errors.New("product not enough stock")
)

type svc struct {
	repo	*repo.Queries
	db		*pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db: db,
	}
}

func (s *svc) CreateOrder(ctx context.Context,tempOrder createOrderParams) (repo.Order, error) {
	// validasi payload
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer id is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one item is required")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	//create order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	// jika produk exists
	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrProductNoStock
		}

		// create order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID: order.ID,
			ProductID: item.ProductID,
			Quantity: item.Quantity,
			PriceCent: product.PriceInCent,
		})
		if err != nil {
			return repo.Order{}, err
		}
	}

	tx.Commit(ctx)

	return order, nil
}
