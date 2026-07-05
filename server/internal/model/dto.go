package model

import (
	"encoding/json"
	"time"
)

// Date is a custom time type that accepts both "2006-01-02" and RFC3339 formats.
type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	// Try "2006-01-02" format first
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		// Fall back to RFC3339
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(d.Time.Format("2006-01-02"))
}

// LoginReq is the request for login.
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterReq is the request for registration.
type RegisterReq struct {
	TenantName string `json:"tenant_name" binding:"required"`
	TenantType string `json:"tenant_type" binding:"required,oneof=personal enterprise team"`
	Contact    string `json:"contact"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required,min=6"`
	RealName   string `json:"real_name"`
}

// LoginResp is the response for login.
type LoginResp struct {
	Token      string `json:"token"`
	UserID     int    `json:"user_id"`
	TenantID   int    `json:"tenant_id"`
	TenantType string `json:"tenant_type"`
	Username   string `json:"username"`
}

// ChangePasswordReq is the request for changing password.
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// CreateUserReq is the request for creating a user.
type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleIDs  []int  `json:"role_ids"`
}

// UpdateUserReq is the request for updating a user.
type UpdateUserReq struct {
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
	RoleIDs  []int  `json:"role_ids"`
}

// CreateTenantReq is the request for creating a tenant.
type CreateTenantReq struct {
	Name    string `json:"name" binding:"required"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

// UpdateTenantReq is the request for updating a tenant.
type UpdateTenantReq struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Status  string `json:"status"`
}

// CreateRoleReq is the request for creating a role.
type CreateRoleReq struct {
	Name        string   `json:"name" binding:"required"`
	Permissions []string `json:"permissions"`
	Description string   `json:"description"`
}

// UpdateRoleReq is the request for updating a role.
type UpdateRoleReq struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Description string   `json:"description"`
}

// PageReq is the common pagination request.
type PageReq struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// PageResp is the common pagination response.
type PageResp struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

// CreateContractReq is the request for creating a contract.
type CreateContractReq struct {
	Title    string `json:"title" binding:"required"`
	PartyA   string `json:"party_a"`
	PartyB   string `json:"party_b"`
	Amount   float64 `json:"amount"`
	StartDate *Date  `json:"start_date"`
	EndDate   *Date  `json:"end_date"`
	FileURL  string `json:"file_url"`
	Remark   string `json:"remark"`
}

// UpdateContractReq is the request for updating a contract.
type UpdateContractReq struct {
	Title    string `json:"title"`
	PartyA   string `json:"party_a"`
	PartyB   string `json:"party_b"`
	Amount   float64 `json:"amount"`
	StartDate *Date  `json:"start_date"`
	EndDate   *Date  `json:"end_date"`
	Status   string `json:"status"`
	FileURL  string `json:"file_url"`
	Remark   string `json:"remark"`
}

// CreateTaskReq is the request for creating a task.
type CreateTaskReq struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	AssigneeID  int    `json:"assignee_id"`
	Priority    string `json:"priority"`
	DueDate     *Date  `json:"due_date"`
}

// UpdateTaskReq is the request for updating a task.
type UpdateTaskReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeID  int    `json:"assignee_id"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     *Date  `json:"due_date"`
}
