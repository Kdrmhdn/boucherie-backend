package repository

import (
	"boucherie-api/internal/domain"
	"context"
	"database/sql"
)

// SQLiteCreditRepo implements port.CreditRepository.
type SQLiteCreditRepo struct {
	db *sql.DB
}

// NewCreditRepo creates a new SQLite-backed credit repository.
func NewCreditRepo(db *sql.DB) *SQLiteCreditRepo {
	return &SQLiteCreditRepo{db: db}
}

// FindAll returns all credits, optionally filtered by status.
func (r *SQLiteCreditRepo) FindAll(ctx context.Context, status *domain.CreditStatus) ([]domain.Credit, error) {
	query := `SELECT id, client_id, client_name, sale_id, amount, remaining_amount, status, created_at, due_date FROM credits WHERE 1=1`
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

	var credits []domain.Credit
	for rows.Next() {
		c, err := r.scanCredit(rows)
		if err != nil {
			return nil, err
		}
		// Load payments
		payments, err := r.findPaymentsByCreditID(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		c.Payments = payments
		credits = append(credits, *c)
	}
	return credits, rows.Err()
}

// FindByID returns a single credit with its payments.
func (r *SQLiteCreditRepo) FindByID(ctx context.Context, id string) (*domain.Credit, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, client_id, client_name, sale_id, amount, remaining_amount, status, created_at, due_date FROM credits WHERE id = ?`, id)

	var c domain.Credit
	var dueDate sql.NullTime
	err := row.Scan(&c.ID, &c.ClientID, &c.ClientName, &c.SaleID, &c.Amount, &c.RemainingAmount, &c.Status, &c.CreatedAt, &dueDate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if dueDate.Valid {
		c.DueDate = &dueDate.Time
	}

	payments, err := r.findPaymentsByCreditID(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	c.Payments = payments
	return &c, nil
}

// FindByClientID returns credits for a specific client.
func (r *SQLiteCreditRepo) FindByClientID(ctx context.Context, clientID string) ([]domain.Credit, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, client_id, client_name, sale_id, amount, remaining_amount, status, created_at, due_date FROM credits WHERE client_id = ? ORDER BY created_at DESC`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []domain.Credit
	for rows.Next() {
		c, err := r.scanCredit(rows)
		if err != nil {
			return nil, err
		}
		payments, err := r.findPaymentsByCreditID(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		c.Payments = payments
		credits = append(credits, *c)
	}
	return credits, rows.Err()
}

// Create inserts a new credit record.
func (r *SQLiteCreditRepo) Create(ctx context.Context, credit *domain.Credit) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO credits (id, client_id, client_name, sale_id, amount, remaining_amount, status, created_at, due_date) VALUES (?,?,?,?,?,?,?,?,?)`,
		credit.ID, credit.ClientID, credit.ClientName, credit.SaleID, credit.Amount, credit.RemainingAmount, credit.Status, credit.CreatedAt, credit.DueDate,
	)
	return err
}

// Update modifies an existing credit (remaining_amount, status).
func (r *SQLiteCreditRepo) Update(ctx context.Context, credit *domain.Credit) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE credits SET remaining_amount=?, status=? WHERE id=?`,
		credit.RemainingAmount, credit.Status, credit.ID,
	)
	return err
}

// AddPayment inserts a payment record for a credit.
func (r *SQLiteCreditRepo) AddPayment(ctx context.Context, payment *domain.Payment) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO payments (id, credit_id, amount, date, method) VALUES (?,?,?,?,?)`,
		payment.ID, payment.CreditID, payment.Amount, payment.Date, payment.Method,
	)
	return err
}

func (r *SQLiteCreditRepo) scanCredit(rows *sql.Rows) (*domain.Credit, error) {
	var c domain.Credit
	var dueDate sql.NullTime
	if err := rows.Scan(&c.ID, &c.ClientID, &c.ClientName, &c.SaleID, &c.Amount, &c.RemainingAmount, &c.Status, &c.CreatedAt, &dueDate); err != nil {
		return nil, err
	}
	if dueDate.Valid {
		c.DueDate = &dueDate.Time
	}
	return &c, nil
}

func (r *SQLiteCreditRepo) findPaymentsByCreditID(ctx context.Context, creditID string) ([]domain.Payment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, credit_id, amount, date, method FROM payments WHERE credit_id = ? ORDER BY date DESC`, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.CreditID, &p.Amount, &p.Date, &p.Method); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, rows.Err()
}
