package repository

import (
	"boucherie-api/internal/domain"
	"context"
	"database/sql"
)

// SQLiteOrderRepo implements port.OrderRepository.
type SQLiteOrderRepo struct {
	db *sql.DB
}

// NewOrderRepo creates a new SQLite-backed order repository.
func NewOrderRepo(db *sql.DB) *SQLiteOrderRepo {
	return &SQLiteOrderRepo{db: db}
}

// FindAll returns orders, optionally filtered by status.
func (r *SQLiteOrderRepo) FindAll(ctx context.Context, status *domain.OrderStatus) ([]domain.Order, error) {
	query := `SELECT id, client_id, client_name, client_phone, pickup_date, notes, status, created_at FROM orders WHERE 1=1`
	var args []interface{}
	if status != nil {
		query += ` AND status = ?`
		args = append(args, string(*status))
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.ClientID, &o.ClientName, &o.ClientPhone, &o.PickupDate, &o.Notes, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		items, err := r.findItemsByOrderID(ctx, o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

// FindByID returns a single order with its items.
func (r *SQLiteOrderRepo) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	var o domain.Order
	err := r.db.QueryRowContext(ctx,
		`SELECT id, client_id, client_name, client_phone, pickup_date, notes, status, created_at FROM orders WHERE id = ?`, id,
	).Scan(&o.ID, &o.ClientID, &o.ClientName, &o.ClientPhone, &o.PickupDate, &o.Notes, &o.Status, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	items, err := r.findItemsByOrderID(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return &o, nil
}

// Create inserts an order and its items in a transaction.
func (r *SQLiteOrderRepo) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (id, client_id, client_name, client_phone, pickup_date, notes, status, created_at) VALUES (?,?,?,?,?,?,?,?)`,
		order.ID, order.ClientID, order.ClientName, order.ClientPhone, order.PickupDate, order.Notes, order.Status, order.CreatedAt,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO order_items (id, order_id, product_id, product_name, quantity) VALUES (?,?,?,?,?)`,
			item.ID, order.ID, item.ProductID, item.ProductName, item.Quantity,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Update modifies an existing order (status, pickup_date, notes).
func (r *SQLiteOrderRepo) Update(ctx context.Context, order *domain.Order) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status=?, pickup_date=?, notes=? WHERE id=?`,
		order.Status, order.PickupDate, order.Notes, order.ID,
	)
	return err
}

// Delete removes an order and its items (CASCADE).
func (r *SQLiteOrderRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM orders WHERE id = ?`, id)
	return err
}

func (r *SQLiteOrderRepo) findItemsByOrderID(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, order_id, product_id, product_name, quantity FROM order_items WHERE order_id = ?`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
