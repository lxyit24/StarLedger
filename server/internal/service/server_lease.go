package service

import (
	"context"
	"fmt"
	"time"

	"starledger/ent"
	"starledger/ent/bill"
	"starledger/ent/serverlease"
)

type ServerLeaseService struct {
	client *ent.Client
}

func NewServerLeaseService(client *ent.Client) *ServerLeaseService {
	return &ServerLeaseService{client: client}
}

func (s *ServerLeaseService) List(ctx context.Context, tenantID, page, pageSize int, provider, status, keyword string) ([]*ent.ServerLease, int, error) {
	offset := (page - 1) * pageSize
	q := s.client.ServerLease.Query().Where(serverlease.TenantID(tenantID))

	if provider != "" {
		q = q.Where(serverlease.Provider(provider))
	}
	if status != "" {
		q = q.Where(serverlease.StatusEQ(serverlease.Status(status)))
	}
	if keyword != "" {
		q = q.Where(serverlease.ServerNameContains(keyword))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := q.
		Offset(offset).
		Limit(pageSize).
		Order(ent.Desc(serverlease.FieldCreatedAt)).
		All(ctx)
	return items, total, err
}

func (s *ServerLeaseService) Get(ctx context.Context, id, tenantID int) (*ent.ServerLease, error) {
	return s.client.ServerLease.Query().
		Where(serverlease.ID(id), serverlease.TenantID(tenantID)).
		WithBills().
		Only(ctx)
}

func (s *ServerLeaseService) Create(ctx context.Context, tenantID int, serverName, provider, config, ipAddress, location string, monthlyCost float64, startDate, endDate time.Time, renewCycle, contractNo, remark string) (*ent.ServerLease, error) {
	return s.client.ServerLease.Create().
		SetTenantID(tenantID).
		SetServerName(serverName).
		SetProvider(provider).
		SetConfig(config).
		SetIPAddress(ipAddress).
		SetLocation(location).
		SetMonthlyCost(monthlyCost).
		SetStartDate(startDate).
		SetEndDate(endDate).
		SetRenewCycle(serverlease.RenewCycle(renewCycle)).
		SetContractNo(contractNo).
		SetRemark(remark).
		SetStatus(serverlease.StatusActive).
		Save(ctx)
}

func (s *ServerLeaseService) Update(ctx context.Context, id, tenantID int, serverName, provider, config, ipAddress, location string, monthlyCost float64, endDate time.Time, contractNo, remark string) (*ent.ServerLease, error) {
	updater := s.client.ServerLease.UpdateOneID(id).
		SetServerName(serverName).
		SetProvider(provider).
		SetConfig(config).
		SetIPAddress(ipAddress).
		SetLocation(location).
		SetMonthlyCost(monthlyCost).
		SetEndDate(endDate).
		SetContractNo(contractNo).
		SetRemark(remark)

	return updater.Save(ctx)
}

func (s *ServerLeaseService) Delete(ctx context.Context, id, tenantID int) error {
	_, err := s.client.ServerLease.UpdateOneID(id).
		SetStatus(serverlease.StatusTerminated).
		Save(ctx)
	return err
}

// Renew extends the lease and creates a billing record.
func (s *ServerLeaseService) Renew(ctx context.Context, id, tenantID int, months int) (*ent.ServerLease, error) {
	lease, err := s.client.ServerLease.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if months <= 0 {
		months = 1
	}

	newEndDate := lease.EndDate.AddDate(0, months, 0)
	amount := lease.MonthlyCost * float64(months)

	// Update lease
	updated, err := s.client.ServerLease.UpdateOneID(id).
		SetEndDate(newEndDate).
		SetStatus(serverlease.StatusActive).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Create bill
	billNo := fmt.Sprintf("BL%s%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
	_, err = s.client.Bill.Create().
		SetTenantID(tenantID).
		SetBillNo(billNo).
		SetBillType("服务器租赁").
		SetRelatedResourceID(id).
		SetAmount(amount).
		SetBillDate(time.Now()).
		SetDueDate(lease.EndDate).
		SetPaymentStatus(bill.PaymentStatusPending).
		SetRemark(fmt.Sprintf("续租 %d 个月", months)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *ServerLeaseService) Expiring(ctx context.Context, tenantID int, days int) ([]*ent.ServerLease, error) {
	deadline := time.Now().AddDate(0, 0, days)
	return s.client.ServerLease.Query().
		Where(
			serverlease.TenantID(tenantID),
			serverlease.StatusIn(serverlease.StatusActive, serverlease.StatusExpiring),
			serverlease.EndDateLTE(deadline),
		).
		Order(ent.Asc(serverlease.FieldEndDate)).
		All(ctx)
}

// CheckExpiring updates server lease status for expiring/expired leases.
func (s *ServerLeaseService) CheckExpiring(ctx context.Context) error {
	now := time.Now()
	warningDate := now.AddDate(0, 0, 30)

	// Mark as expiring (within 30 days)
	_, err := s.client.ServerLease.Update().
		Where(
			serverlease.StatusEQ(serverlease.StatusActive),
			serverlease.EndDateLTE(warningDate),
			serverlease.EndDateGT(now),
		).
		SetStatus(serverlease.StatusExpiring).
		Save(ctx)
	if err != nil {
		return err
	}

	// Mark as expired
	_, err = s.client.ServerLease.Update().
		Where(
			serverlease.StatusIn(serverlease.StatusActive, serverlease.StatusExpiring),
			serverlease.EndDateLTE(now),
		).
		SetStatus(serverlease.StatusExpired).
		Save(ctx)
	return err
}
