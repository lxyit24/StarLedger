package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tenant holds the schema definition for the Tenant entity.
type Tenant struct {
	ent.Schema
}

// Mixin of the Tenant.
func (Tenant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Tenant.
func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("租户名称"),
		field.Enum("type").
			Values("personal", "enterprise", "team").
			Default("personal").
			Comment("租户类型: personal/enterprise/team"),
		field.String("contact").Default("").Comment("联系人"),
		field.String("phone").Default("").Comment("联系电话"),
		field.String("email").Default("").Comment("联系邮箱"),
		field.Enum("status").Values("active", "suspended").Default("active").Comment("状态: active/suspended"),
	}
}

// Edges of the Tenant.
func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
		edge.To("roles", Role.Type),
		edge.To("server_leases", ServerLease.Type),
		edge.To("bills", Bill.Type),
		edge.To("modules", TenantModule.Type),
		edge.To("contracts", Contract.Type),
		edge.To("tasks", Task.Type),
		edge.To("audit_logs", AuditLog.Type),
		edge.To("invoices", Invoice.Type),
	}
}

// Indexes of the Tenant.
func (Tenant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
	}
}
