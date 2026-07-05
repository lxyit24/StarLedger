package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"starledger/ent"
	"starledger/ent/tenant"
	"starledger/ent/user"
)

type AuthService struct {
	client *ent.Client
}

func NewAuthService(client *ent.Client) *AuthService {
	return &AuthService{client: client}
}

// Register creates a new tenant and admin user.
func (s *AuthService) Register(ctx context.Context, tenantName, contact, phone, email, username, password, realName string) (*ent.User, error) {
	// Check if username already exists
	exists, _ := s.client.User.Query().Where(user.Username(username)).Exist(ctx)
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// Create tenant + admin user in a transaction
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("开启事务失败: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Create tenant
	t, err := tx.Tenant.Create().
		SetName(tenantName).
		SetContact(contact).
		SetPhone(phone).
		SetEmail(email).
		SetStatus(tenant.StatusActive).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建租户失败: %w", err)
	}

	// Create admin role
	adminRole, err := tx.Role.Create().
		SetTenantID(t.ID).
		SetName("admin").
		SetDescription("超级管理员").
		SetPermissions([]string{"*"}).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	// Create admin user
	u, err := tx.User.Create().
		SetTenantID(t.ID).
		SetUsername(username).
		SetPasswordHash(string(hash)).
		SetRealName(realName).
		SetEmail(email).
		SetPhone(phone).
		SetStatus(user.StatusActive).
		AddRoles(adminRole).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// Initialize tenant modules: enable all modules by default
	defaultModules := []string{"server_lease", "billing"}
	for _, modName := range defaultModules {
		_, err = tx.TenantModule.Create().
			SetTenantID(t.ID).
			SetModuleName(modName).
			SetEnabled(true).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("初始化模块失败: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	return u, nil
}

// Login authenticates a user and returns the user entity.
func (s *AuthService) Login(ctx context.Context, username, password string) (*ent.User, error) {
	u, err := s.client.User.Query().
		Where(user.Username(username)).
		WithTenant().
		WithRoles().
		Only(ctx)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if u.Status == user.StatusDisabled {
		return nil, errors.New("账号已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if u.Edges.Tenant != nil && u.Edges.Tenant.Status == tenant.StatusSuspended {
		return nil, errors.New("租户已被暂停")
	}

	return u, nil
}

// ChangePassword changes a user's password.
func (s *AuthService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	u, err := s.client.User.Get(ctx, userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	_, err = s.client.User.UpdateOneID(userID).SetPasswordHash(string(hash)).Save(ctx)
	return err
}
