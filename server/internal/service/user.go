package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"starledger/ent"
	"starledger/ent/role"
	"starledger/ent/user"
)

type UserService struct {
	client *ent.Client
}

func NewUserService(client *ent.Client) *UserService {
	return &UserService{client: client}
}

func (s *UserService) List(ctx context.Context, tenantID, page, pageSize int) ([]*ent.User, int, error) {
	offset := (page - 1) * pageSize
	q := s.client.User.Query().Where(user.TenantID(tenantID))
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := q.
		Offset(offset).
		Limit(pageSize).
		WithRoles().
		Order(ent.Desc(user.FieldCreatedAt)).
		All(ctx)
	return items, total, err
}

func (s *UserService) Create(ctx context.Context, tenantID int, username, password, realName, email, phone string, roleIDs []int) (*ent.User, error) {
	exists, _ := s.client.User.Query().Where(user.Username(username)).Exist(ctx)
	if exists {
		return nil, errors.New("用户名已存在")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	creator := s.client.User.Create().
		SetTenantID(tenantID).
		SetUsername(username).
		SetPasswordHash(string(hash)).
		SetRealName(realName).
		SetEmail(email).
		SetPhone(phone).
		SetStatus(user.StatusActive)

	if len(roleIDs) > 0 {
		creator = creator.AddRoleIDs(roleIDs...)
	}

	return creator.Save(ctx)
}

func (s *UserService) Update(ctx context.Context, id, tenantID int, realName, email, phone, status string, roleIDs []int) (*ent.User, error) {
	updater := s.client.User.UpdateOneID(id).
		SetRealName(realName).
		SetEmail(email).
		SetPhone(phone)

	if status != "" {
		updater = updater.SetStatus(user.Status(status))
	}

	if len(roleIDs) > 0 {
		// Clear existing roles and set new ones
		updater = updater.ClearRoles().AddRoleIDs(roleIDs...)
	}

	return updater.Save(ctx)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	_, err := s.client.User.UpdateOneID(id).SetStatus(user.StatusDisabled).Save(ctx)
	return err
}

func (s *UserService) GetProfile(ctx context.Context, userID int) (*ent.User, error) {
	return s.client.User.Get(ctx, userID)
}

// LoadUserRoles loads user roles and permissions into context-friendly format.
func (s *UserService) GetUserRolesAndPerms(ctx context.Context, userID int) ([]string, []string, error) {
	u, err := s.client.User.Query().
		Where(user.ID(userID)).
		WithRoles().
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	var roleNames []string
	permSet := make(map[string]bool)
	for _, r := range u.Edges.Roles {
		roleNames = append(roleNames, r.Name)
		for _, p := range r.Permissions {
			permSet[p] = true
		}
	}

	var perms []string
	for p := range permSet {
		perms = append(perms, p)
	}
	return roleNames, perms, nil
}

// RoleService handles role operations.
type RoleService struct {
	client *ent.Client
}

func NewRoleService(client *ent.Client) *RoleService {
	return &RoleService{client: client}
}

func (s *RoleService) List(ctx context.Context, tenantID int) ([]*ent.Role, error) {
	return s.client.Role.Query().
		Where(role.TenantID(tenantID)).
		Order(ent.Desc(role.FieldCreatedAt)).
		All(ctx)
}

func (s *RoleService) Create(ctx context.Context, tenantID int, name, description string, permissions []string) (*ent.Role, error) {
	return s.client.Role.Create().
		SetTenantID(tenantID).
		SetName(name).
		SetDescription(description).
		SetPermissions(permissions).
		Save(ctx)
}

func (s *RoleService) Update(ctx context.Context, id int, name, description string, permissions []string) (*ent.Role, error) {
	updater := s.client.Role.UpdateOneID(id)
	if name != "" {
		updater = updater.SetName(name)
	}
	if description != "" {
		updater = updater.SetDescription(description)
	}
	if permissions != nil {
		updater = updater.SetPermissions(permissions)
	}
	return updater.Save(ctx)
}
