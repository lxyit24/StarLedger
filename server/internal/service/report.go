package service

import (
	"context"
	"time"

	"starledger/ent"
	"starledger/ent/bill"
	"starledger/ent/contract"
	"starledger/ent/serverlease"
	"starledger/ent/task"
)

type ReportService struct {
	client *ent.Client
}

func NewReportService(client *ent.Client) *ReportService {
	return &ReportService{client: client}
}

// MonthlyTrend returns monthly bill totals for the last 12 months.
func (s *ReportService) MonthlyTrend(ctx context.Context, tenantID int) ([]map[string]interface{}, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month()-11, 1, 0, 0, 0, 0, now.Location())

	rows, err := s.client.Bill.Query().
		Where(
			bill.TenantID(tenantID),
			bill.BillDateGTE(start),
		).
		Order(ent.Asc(bill.FieldBillDate)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Group by month
	monthMap := make(map[string]map[string]interface{})
	for i := 0; i < 12; i++ {
		t := start.AddDate(0, i, 0)
		key := t.Format("2006-01")
		monthMap[key] = map[string]interface{}{
			"month":  key,
			"total":  0.0,
			"count":  0,
			"paid":   0.0,
			"unpaid": 0.0,
		}
	}

	for _, b := range rows {
		key := b.BillDate.Format("2006-01")
		if m, ok := monthMap[key]; ok {
			m["total"] = m["total"].(float64) + b.Amount
			m["count"] = m["count"].(int) + 1
			if b.PaymentStatus == bill.PaymentStatusPaid {
				m["paid"] = m["paid"].(float64) + b.Amount
			} else {
				m["unpaid"] = m["unpaid"].(float64) + b.Amount
			}
		}
	}

	result := make([]map[string]interface{}, 0, 12)
	for i := 0; i < 12; i++ {
		t := start.AddDate(0, i, 0)
		key := t.Format("2006-01")
		result = append(result, monthMap[key])
	}
	return result, nil
}

// BillTypeDistribution returns bill amount grouped by bill_type.
func (s *ReportService) BillTypeDistribution(ctx context.Context, tenantID int) ([]map[string]interface{}, error) {
	bills, err := s.client.Bill.Query().
		Where(bill.TenantID(tenantID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	typeMap := make(map[string]float64)
	for _, b := range bills {
		typeMap[b.BillType] += b.Amount
	}

	result := make([]map[string]interface{}, 0, len(typeMap))
	for k, v := range typeMap {
		result = append(result, map[string]interface{}{
			"type":   k,
			"amount": v,
		})
	}
	return result, nil
}

// BillStatusSummary returns bill count and amount grouped by payment_status.
func (s *ReportService) BillStatusSummary(ctx context.Context, tenantID int) ([]map[string]interface{}, error) {
	bills, err := s.client.Bill.Query().
		Where(bill.TenantID(tenantID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]map[string]interface{})
	for _, status := range []string{"pending", "paid", "overdue", "cancelled"} {
		statusMap[status] = map[string]interface{}{
			"status": status,
			"count":  0,
			"amount": 0.0,
		}
	}

	for _, b := range bills {
		s := string(b.PaymentStatus)
		if m, ok := statusMap[s]; ok {
			m["count"] = m["count"].(int) + 1
			m["amount"] = m["amount"].(float64) + b.Amount
		}
	}

	result := make([]map[string]interface{}, 0, len(statusMap))
	for _, v := range statusMap {
		result = append(result, v)
	}
	return result, nil
}

// Overview returns dashboard overview statistics.
func (s *ReportService) Overview(ctx context.Context, tenantID int) (map[string]interface{}, error) {
	// Total revenue (paid bills)
	totalRevenueInt, _ := s.client.Bill.Query().
		Where(bill.TenantID(tenantID), bill.PaymentStatusEQ(bill.PaymentStatusPaid)).
		Aggregate(ent.Sum(bill.FieldAmount)).
		Int(ctx)
	totalRevenue := float64(totalRevenueInt)

	// Pending amount
	pendingAmountInt, _ := s.client.Bill.Query().
		Where(bill.TenantID(tenantID), bill.PaymentStatusEQ(bill.PaymentStatusPending)).
		Aggregate(ent.Sum(bill.FieldAmount)).
		Int(ctx)
	pendingAmount := float64(pendingAmountInt)

	// Overdue count and amount
	overdueCount, _ := s.client.Bill.Query().
		Where(bill.TenantID(tenantID), bill.PaymentStatusEQ(bill.PaymentStatusOverdue)).
		Count(ctx)
	overdueAmountInt, _ := s.client.Bill.Query().
		Where(bill.TenantID(tenantID), bill.PaymentStatusEQ(bill.PaymentStatusOverdue)).
		Aggregate(ent.Sum(bill.FieldAmount)).
		Int(ctx)
	overdueAmount := float64(overdueAmountInt)

	// Server count
	serverCount, _ := s.client.ServerLease.Query().
		Where(serverlease.TenantID(tenantID), serverlease.StatusEQ(serverlease.StatusActive)).
		Count(ctx)

	// Monthly server cost
	serverCostsInt, _ := s.client.ServerLease.Query().
		Where(serverlease.TenantID(tenantID), serverlease.StatusEQ(serverlease.StatusActive)).
		Aggregate(ent.Sum(serverlease.FieldMonthlyCost)).
		Int(ctx)
	serverCosts := float64(serverCostsInt)

	// Active contracts
	activeContracts, _ := s.client.Contract.Query().
		Where(contract.TenantID(tenantID), contract.StatusEQ(contract.StatusActive)).
		Count(ctx)

	// Contract total amount
	contractAmountInt, _ := s.client.Contract.Query().
		Where(contract.TenantID(tenantID), contract.StatusEQ(contract.StatusActive)).
		Aggregate(ent.Sum(contract.FieldAmount)).
		Int(ctx)
	contractAmount := float64(contractAmountInt)

	// Pending tasks
	pendingTasks, _ := s.client.Task.Query().
		Where(task.TenantID(tenantID), task.StatusIn(task.StatusPending, task.StatusInProgress)).
		Count(ctx)

	// This month revenue
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthRevenueInt, _ := s.client.Bill.Query().
		Where(
			bill.TenantID(tenantID),
			bill.PaymentStatusEQ(bill.PaymentStatusPaid),
			bill.PaidDateGTE(startOfMonth),
		).
		Aggregate(ent.Sum(bill.FieldAmount)).
		Int(ctx)
	monthRevenue := float64(monthRevenueInt)

	return map[string]interface{}{
		"total_revenue":     totalRevenue,
		"pending_amount":    pendingAmount,
		"overdue_count":     overdueCount,
		"overdue_amount":    overdueAmount,
		"server_count":      serverCount,
		"server_month_cost": serverCosts,
		"active_contracts":  activeContracts,
		"contract_amount":   contractAmount,
		"pending_tasks":     pendingTasks,
		"month_revenue":     monthRevenue,
	}, nil
}

// ServerCostAnalysis returns server cost breakdown by provider.
func (s *ReportService) ServerCostAnalysis(ctx context.Context, tenantID int) ([]map[string]interface{}, error) {
	servers, err := s.client.ServerLease.Query().
		Where(serverlease.TenantID(tenantID), serverlease.StatusIn(serverlease.StatusActive, serverlease.StatusExpiring)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	providerMap := make(map[string]map[string]interface{})
	for _, sv := range servers {
		p, ok := providerMap[sv.Provider]
		if !ok {
			p = map[string]interface{}{
				"provider":     sv.Provider,
				"count":        0,
				"monthly_cost": 0.0,
			}
			providerMap[sv.Provider] = p
		}
		p["count"] = p["count"].(int) + 1
		p["monthly_cost"] = p["monthly_cost"].(float64) + sv.MonthlyCost
	}

	result := make([]map[string]interface{}, 0, len(providerMap))
	for _, v := range providerMap {
		result = append(result, v)
	}
	return result, nil
}

// TaskStatistics returns task completion statistics.
func (s *ReportService) TaskStatistics(ctx context.Context, tenantID int) (map[string]interface{}, error) {
	total, _ := s.client.Task.Query().Where(task.TenantID(tenantID)).Count(ctx)
	pending, _ := s.client.Task.Query().Where(task.TenantID(tenantID), task.StatusEQ(task.StatusPending)).Count(ctx)
	inProgress, _ := s.client.Task.Query().Where(task.TenantID(tenantID), task.StatusEQ(task.StatusInProgress)).Count(ctx)
	completed, _ := s.client.Task.Query().Where(task.TenantID(tenantID), task.StatusEQ(task.StatusCompleted)).Count(ctx)
	cancelled, _ := s.client.Task.Query().Where(task.TenantID(tenantID), task.StatusEQ(task.StatusCancelled)).Count(ctx)

	return map[string]interface{}{
		"total":       total,
		"pending":     pending,
		"in_progress": inProgress,
		"completed":   completed,
		"cancelled":   cancelled,
	}, nil
}
