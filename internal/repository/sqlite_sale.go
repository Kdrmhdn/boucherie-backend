package repository

import (
	"boucherie-api/internal/domain"
	"context"
	"database/sql"
)

// SQLiteSaleRepo implements port.SaleRepository.
type SQLiteSaleRepo struct {
	db *sql.DB
}

// NewSaleRepo creates a new SQLite-backed sale repository.
func NewSaleRepo(db *sql.DB) *SQLiteSaleRepo {
	return &SQLiteSaleRepo{db: db}
}

// FindAll returns sales with optional filtering by client and/or date.
func (r *SQLiteSaleRepo) FindAll(ctx context.Context, clientID *string, date *string) ([]domain.Sale, error) {
	query := `SELECT id, client_id, client_name, total, paid_amount, credit_amount, date FROM sales WHERE 1=1`
	var args []interface{}

	if clientID != nil {
		query += ` AND client_id = ?`
		args = append(args, *clientID)
	}
	if date != nil {
		query += ` AND date(date) = ?`
		args = append(args, *date)
	}
	query += ` ORDER BY date DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []domain.Sale
	for rows.Next() {
		var s domain.Sale
		if err := rows.Scan(&s.ID, &s.ClientID, &s.ClientName, &s.Total, &s.PaidAmount, &s.CreditAmount, &s.Date); err != nil {
			return nil, err
		}
		// Load items for this sale
		items, err := r.findItemsBySaleID(ctx, s.ID)
		if err != nil {
			return nil, err
		}
		s.Items = items
		sales = append(sales, s)
	}
	return sales, rows.Err()
}

// FindByID returns a single sale with its items.
func (r *SQLiteSaleRepo) FindByID(ctx context.Context, id string) (*domain.Sale, error) {
	var s domain.Sale
	err := r.db.QueryRowContext(ctx,
		`SELECT id, client_id, client_name, total, paid_amount, credit_amount, date FROM sales WHERE id = ?`, id,
	).Scan(&s.ID, &s.ClientID, &s.ClientName, &s.Total, &s.PaidAmount, &s.CreditAmount, &s.Date)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	items, err := r.findItemsBySaleID(ctx, s.ID)
	if err != nil {
		return nil, err
	}
	s.Items = items
	return &s, nil
}

// Create inserts a sale and its items in a single transaction.
func (r *SQLiteSaleRepo) Create(ctx context.Context, sale *domain.Sale) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO sales (id, client_id, client_name, total, paid_amount, credit_amount, date) VALUES (?,?,?,?,?,?,?)`,
		sale.ID, sale.ClientID, sale.ClientName, sale.Total, sale.PaidAmount, sale.CreditAmount, sale.Date,
	)
	if err != nil {
		return err
	}

	for _, item := range sale.Items {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO sale_items (id, sale_id, product_id, product_name, quantity, subtotal) VALUES (?,?,?,?,?,?)`,
			item.ID, sale.ID, item.ProductID, item.ProductName, item.Quantity, item.Subtotal,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *SQLiteSaleRepo) findItemsBySaleID(ctx context.Context, saleID string) ([]domain.SaleItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, sale_id, product_id, product_name, quantity, subtotal FROM sale_items WHERE sale_id = ?`, saleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.SaleItem
	for rows.Next() {
		var item domain.SaleItem
		if err := rows.Scan(&item.ID, &item.SaleID, &item.ProductID, &item.ProductName, &item.Quantity, &item.Subtotal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
