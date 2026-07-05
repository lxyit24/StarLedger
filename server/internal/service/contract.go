package service

import (
	"context"
	"errors"
	"time"

	"starledger/ent"
	"starledger/ent/contract"
)

type ContractService struct {
	client *ent.Client
}

func NewContractService(client *ent.Client) *ContractService {
	return &ContractService{client: client}
}

// List returns contracts for a tenant with pagination.
func (s *ContractService) List(ctx context.Context, tenantID, page, pageSize int) ([]*ent.Contract, int, error) {
	query := s.client.Contract.Query().Where(contract.TenantID(tenantID))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	items, err := query.
		Order(ent.Desc(contract.FieldCreatedAt)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Get returns a contract by ID.
func (s *ContractService) Get(ctx context.Context, id, tenantID int) (*ent.Contract, error) {
	return s.client.Contract.Query().
		Where(contract.ID(id), contract.TenantID(tenantID)).
		Only(ctx)
}

// Create creates a new contract.
func (s *ContractService) Create(ctx context.Context, tenantID int, title, partyA, partyB string, amount float64, startDate, endDate *time.Time, fileURL, remark string) (*ent.Contract, error) {
	builder := s.client.Contract.Create().
		SetTenantID(tenantID).
		SetTitle(title).
		SetPartyA(partyA).
		SetPartyB(partyB).
		SetAmount(amount).
		SetFileURL(fileURL).
		SetRemark(remark).
		SetStatus(contract.StatusDraft)

	if startDate != nil {
		builder.SetStartDate(*startDate)
	}
	if endDate != nil {
		builder.SetEndDate(*endDate)
	}

	return builder.Save(ctx)
}

// Update updates a contract.
func (s *ContractService) Update(ctx context.Context, id, tenantID int, title, partyA, partyB string, amount float64, startDate, endDate *time.Time, status, fileURL, remark string) (*ent.Contract, error) {
	c, err := s.client.Contract.Query().
		Where(contract.ID(id), contract.TenantID(tenantID)).
		Only(ctx)
	if err != nil {
		return nil, errors.New("合同不存在")
	}

	updater := c.Update()
	if title != "" {
		updater.SetTitle(title)
	}
	if partyA != "" {
		updater.SetPartyA(partyA)
	}
	if partyB != "" {
		updater.SetPartyB(partyB)
	}
	if amount > 0 {
		updater.SetAmount(amount)
	}
	if startDate != nil {
		updater.SetStartDate(*startDate)
	}
	if endDate != nil {
		updater.SetEndDate(*endDate)
	}
	if status != "" {
		updater.SetStatus(contract.Status(status))
	}
	if fileURL != "" {
		updater.SetFileURL(fileURL)
	}
	if remark != "" {
		updater.SetRemark(remark)
	}

	return updater.Save(ctx)
}

// Delete deletes a contract.
func (s *ContractService) Delete(ctx context.Context, id, tenantID int) error {
	n, err := s.client.Contract.Delete().
		Where(contract.ID(id), contract.TenantID(tenantID)).
		Exec(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("合同不存在")
	}
	return nil
}

// CheckExpiring checks contracts that are about to expire.
func (s *ContractService) CheckExpiring(ctx context.Context) error {
	now := time.Now()
	threshold := now.AddDate(0, 0, 30) // 30 days from now

	_, err := s.client.Contract.Update().
		Where(
			contract.StatusEQ(contract.StatusActive),
			contract.EndDateLTE(threshold),
			contract.EndDateGT(now),
		).
		SetStatus(contract.StatusExpired).
		Save(ctx)
	return err
}
