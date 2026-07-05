package service

import (
	"context"

	"starledger/ent"
	"starledger/ent/tenant"
)

type TenantService struct {
	client *ent.Client
}

func NewTenantService(client *ent.Client) *TenantService {
	return &TenantService{client: client}
}

func (s *TenantService) List(ctx context.Context, page, pageSize int) ([]*ent.Tenant, int, error) {
	offset := (page - 1) * pageSize
	total, err := s.client.Tenant.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := s.client.Tenant.Query().
		Offset(offset).
		Limit(pageSize).
		Order(ent.Desc(tenant.FieldCreatedAt)).
		All(ctx)
	return items, total, err
}

func (s *TenantService) Create(ctx context.Context, name, contact, phone, email string) (*ent.Tenant, error) {
	return s.client.Tenant.Create().
		SetName(name).
		SetContact(contact).
		SetPhone(phone).
		SetEmail(email).
		Save(ctx)
}

func (s *TenantService) Update(ctx context.Context, id int, name, contact, phone, email, status string) (*ent.Tenant, error) {
	updater := s.client.Tenant.UpdateOneID(id)
	if name != "" {
		updater = updater.SetName(name)
	}
	if contact != "" {
		updater = updater.SetContact(contact)
	}
	if phone != "" {
		updater = updater.SetPhone(phone)
	}
	if email != "" {
		updater = updater.SetEmail(email)
	}
	if status != "" {
		updater = updater.SetStatus(tenant.Status(status))
	}
	return updater.Save(ctx)
}

func (s *TenantService) Delete(ctx context.Context, id int) error {
	_, err := s.client.Tenant.UpdateOneID(id).SetStatus(tenant.StatusSuspended).Save(ctx)
	return err
}
