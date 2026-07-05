package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AuditLog holds the schema definition for the AuditLog entity.
type AuditLog struct {
	ent.Schema
}

// Mixin of the AuditLog.
func (AuditLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the AuditLog.
func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("租户ID"),
		field.Int("user_id").Optional().Default(0).Comment("操作用户ID"),
		field.String("username").Default("").Comment("操作用户名"),
		field.String("action").NotEmpty().Comment("操作类型: create/update/delete/login/logout"),
		field.String("resource_type").Default("").Comment("资源类型: bill/contract/task/server/user"),
		field.Int("resource_id").Optional().Default(0).Comment("资源ID"),
		field.String("detail").Default("").Comment("操作详情"),
		field.String("ip_address").Default("").Comment("操作IP地址"),
		field.String("user_agent").Default("").Comment("用户代理"),
		field.String("status").Default("success").Comment("操作状态: success/failure"),
	}
}

// Edges of the AuditLog.
func (AuditLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("audit_logs").Unique().Required().Field("tenant_id"),
	}
}

// Indexes of the AuditLog.
func (AuditLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "created_at"),
		index.Fields("tenant_id", "action"),
		index.Fields("tenant_id", "resource_type"),
		index.Fields("user_id"),
	}
}
