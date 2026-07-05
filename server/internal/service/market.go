package service

import (
	"context"
	"fmt"

	"starledger/ent"
	"starledger/ent/tenantmodule"
)

type MarketService struct {
	client *ent.Client
}

func NewMarketService(client *ent.Client) *MarketService {
	return &MarketService{client: client}
}

// GetTenantModules returns all module records for a tenant.
func (s *MarketService) GetTenantModules(ctx context.Context, tenantID int) (map[string]bool, error) {
	mods, err := s.client.TenantModule.Query().
		Where(tenantmodule.TenantID(tenantID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool)
	for _, m := range mods {
		result[m.ModuleName] = m.Enabled
	}
	return result, nil
}

// EnableModule enables a module for a tenant (upsert).
func (s *MarketService) EnableModule(ctx context.Context, tenantID int, moduleName string) error {
	existing, err := s.client.TenantModule.Query().
		Where(
			tenantmodule.TenantID(tenantID),
			tenantmodule.ModuleName(moduleName),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			_, err = s.client.TenantModule.Create().
				SetTenantID(tenantID).
				SetModuleName(moduleName).
				SetEnabled(true).
				Save(ctx)
			return err
		}
		return err
	}
	_, err = s.client.TenantModule.UpdateOne(existing).
		SetEnabled(true).
		Save(ctx)
	return err
}

// DisableModule disables a module for a tenant.
func (s *MarketService) DisableModule(ctx context.Context, tenantID int, moduleName string) error {
	existing, err := s.client.TenantModule.Query().
		Where(
			tenantmodule.TenantID(tenantID),
			tenantmodule.ModuleName(moduleName),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil
		}
		return err
	}
	_, err = s.client.TenantModule.UpdateOne(existing).
		SetEnabled(false).
		Save(ctx)
	return err
}

// IsModuleEnabled checks if a specific module is enabled for a tenant.
func (s *MarketService) IsModuleEnabled(ctx context.Context, tenantID int, moduleName string) bool {
	exists, err := s.client.TenantModule.Query().
		Where(
			tenantmodule.TenantID(tenantID),
			tenantmodule.ModuleName(moduleName),
			tenantmodule.Enabled(true),
		).
		Exist(ctx)
	if err != nil {
		return false
	}
	return exists
}

// InitTenantModules enables all given modules for a new tenant.
func (s *MarketService) InitTenantModules(ctx context.Context, tenantID int, moduleNames []string) error {
	for _, name := range moduleNames {
		err := s.EnableModule(ctx, tenantID, name)
		if err != nil {
			return fmt.Errorf("failed to enable module %s: %w", name, err)
		}
	}
	return nil
}
