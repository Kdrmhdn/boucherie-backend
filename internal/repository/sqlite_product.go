package repository

import (
	"boucherie-api/internal/domain"
	"context"
	"database/sql"
)

// SQLiteProductRepo implements port.ProductRepository.
type SQLiteProductRepo struct {
	db *sql.DB
}

// NewProductRepo creates a new SQLite-backed product repository.
func NewProductRepo(db *sql.DB) *SQLiteProductRepo {
	return &SQLiteProductRepo{db: db}
}

// FindAll returns products, optionally filtered by category.
func (r *SQLiteProductRepo) FindAll(ctx context.Context, category *domain.MeatCategory) ([]domain.Product, error) {
	query := `SELECT id, name, category, price_per_kg, image, in_stock FROM products`
	var args []interface{}
	if category != nil {
		query += ` WHERE category = ?`
		args = append(args, string(*category))
	}
	query += ` ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		var inStock int
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.PricePerKg, &p.Image, &inStock); err != nil {
			return nil, err
		}
		p.InStock = inStock == 1
		products = append(products, p)
	}
	return products, rows.Err()
}

// FindByID returns a single product.
func (r *SQLiteProductRepo) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	var p domain.Product
	var inStock int
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, category, price_per_kg, image, in_stock FROM products WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Category, &p.PricePerKg, &p.Image, &inStock)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.InStock = inStock == 1
	return &p, nil
}

// Create inserts a new product.
func (r *SQLiteProductRepo) Create(ctx context.Context, product *domain.Product) error {
	inStock := 0
	if product.InStock {
		inStock = 1
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO products (id, name, category, price_per_kg, image, in_stock) VALUES (?,?,?,?,?,?)`,
		product.ID, product.Name, product.Category, product.PricePerKg, product.Image, inStock,
	)
	return err
}

// Update modifies an existing product.
func (r *SQLiteProductRepo) Update(ctx context.Context, product *domain.Product) error {
	inStock := 0
	if product.InStock {
		inStock = 1
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE products SET name=?, category=?, price_per_kg=?, image=?, in_stock=? WHERE id=?`,
		product.Name, product.Category, product.PricePerKg, product.Image, inStock, product.ID,
	)
	return err
}

// Delete removes a product by ID.
func (r *SQLiteProductRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id)
	return err
}
