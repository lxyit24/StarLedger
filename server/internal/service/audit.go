package service

import (
	"context"

	"starledger/ent"
	"starledger/ent/auditlog"
)

type AuditLogService struct {
	client *ent.Client
}

func NewAuditLogService(client *ent.Client) *AuditLogService {
	return &AuditLogService{client: client}
}

// Create logs an audit entry.
func (s *AuditLogService) Create(ctx context.Context, tenantID, userID int, username, action, resourceType string, resourceID int, detail, ip, ua string) error {
	_, err := s.client.AuditLog.Create().
		SetTenantID(tenantID).
		SetUserID(userID).
		SetUsername(username).
		SetAction(action).
		SetResourceType(resourceType).
		SetResourceID(resourceID).
		SetDetail(detail).
		SetIPAddress(ip).
		SetUserAgent(ua).
		SetStatus("success").
		Save(ctx)
	return err
}

// List returns audit logs with pagination.
func (s *AuditLogService) List(ctx context.Context, tenantID, page, pageSize int, action, resourceType string) ([]*ent.AuditLog, int, error) {
	q := s.client.AuditLog.Query().Where(auditlog.TenantID(tenantID))

	if action != "" {
		q = q.Where(auditlog.Action(action))
	}
	if resourceType != "" {
		q = q.Where(auditlog.ResourceType(resourceType))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	items, err := q.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order(ent.Desc(auditlog.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
