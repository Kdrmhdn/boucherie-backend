package handler

import (
	"context"
	"database/sql"
	"net/http"
)

// DashboardHandler handles the dashboard stats endpoint.
type DashboardHandler struct {
	db *sql.DB
}

// NewDashboardHandler creates a new dashboard handler.
func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

type dashboardStats struct {
	TodayRevenue  float64      `json:"todayRevenue"`
	TodayCash     float64      `json:"todayCash"`
	TodayCredit   float64      `json:"todayCredit"`
	TodaySales    int          `json:"todaySales"`
	TotalClients  int          `json:"totalClients"`
	PendingCredit float64      `json:"pendingCredit"`
	OverdueCount  int          `json:"overdueCount"`
	TopDebtors    []debtorInfo `json:"topDebtors"`
}

type debtorInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	TotalCredit float64 `json:"totalCredit"`
}

// ServeHTTP handles GET /api/v1/dashboard.
func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats := dashboardStats{}

	// Today's revenue
	h.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(total),0), COALESCE(SUM(paid_amount),0), COALESCE(SUM(credit_amount),0), COUNT(*) FROM sales WHERE date(date) = date('now')`,
	).Scan(&stats.TodayRevenue, &stats.TodayCash, &stats.TodayCredit, &stats.TodaySales)

	// Total clients
	h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM clients`).Scan(&stats.TotalClients)

	// Pending credits
	h.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(remaining_amount),0) FROM credits WHERE status != 'paye'`,
	).Scan(&stats.PendingCredit)

	// Overdue count
	h.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM credits WHERE status = 'en_retard'`,
	).Scan(&stats.OverdueCount)

	// Top debtors
	stats.TopDebtors = h.loadTopDebtors(ctx)

	JSON(w, http.StatusOK, stats)
}

func (h *DashboardHandler) loadTopDebtors(ctx context.Context) []debtorInfo {
	rows, err := h.db.QueryContext(ctx,
		`SELECT id, name, total_credit FROM clients WHERE total_credit > 0 ORDER BY total_credit DESC LIMIT 5`)
	if err != nil {
		return []debtorInfo{}
	}
	defer rows.Close()

	var debtors []debtorInfo
	for rows.Next() {
		var d debtorInfo
		if err := rows.Scan(&d.ID, &d.Name, &d.TotalCredit); err != nil {
			continue
		}
		debtors = append(debtors, d)
	}
	if debtors == nil {
		return []debtorInfo{}
	}
	return debtors
}
