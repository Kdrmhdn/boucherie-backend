package repository

import (
	"boucherie-api/internal/domain"
	"context"
	"database/sql"
)

// SQLiteClientRepo implements port.ClientRepository using SQLite.
type SQLiteClientRepo struct {
	db *sql.DB
}

// NewClientRepo creates a new SQLite-backed client repository.
func NewClientRepo(db *sql.DB) *SQLiteClientRepo {
	return &SQLiteClientRepo{db: db}
}

// FindAll returns every client ordered by creation date (newest first).
func (r *SQLiteClientRepo) FindAll(ctx context.Context) ([]domain.Client, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, phone, email, avatar, total_credit, created_at FROM clients ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []domain.Client
	for rows.Next() {
		var c domain.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.Avatar, &c.TotalCredit, &c.CreatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

// FindByID returns a single client by ID.
func (r *SQLiteClientRepo) FindByID(ctx context.Context, id string) (*domain.Client, error) {
	var c domain.Client
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, phone, email, avatar, total_credit, created_at FROM clients WHERE id = ?`, id,
	).Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.Avatar, &c.TotalCredit, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Create inserts a new client.
func (r *SQLiteClientRepo) Create(ctx context.Context, client *domain.Client) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO clients (id, name, phone, email, avatar, total_credit, created_at) VALUES (?,?,?,?,?,?,?)`,
		client.ID, client.Name, client.Phone, client.Email, client.Avatar, client.TotalCredit, client.CreatedAt,
	)
	return err
}

// Update modifies an existing client.
func (r *SQLiteClientRepo) Update(ctx context.Context, client *domain.Client) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE clients SET name=?, phone=?, email=?, avatar=? WHERE id=?`,
		client.Name, client.Phone, client.Email, client.Avatar, client.ID,
	)
	return err
}

// UpdateTotalCredit atomically adds delta to a client's total credit.
func (r *SQLiteClientRepo) UpdateTotalCredit(ctx context.Context, clientID string, delta float64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE clients SET total_credit = total_credit + ? WHERE id = ?`,
		delta, clientID,
	)
	return err
}
