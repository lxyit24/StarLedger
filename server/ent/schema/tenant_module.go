package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TenantModule holds the schema definition for the TenantModule entity.
type TenantModule struct {
	ent.Schema
}

// Mixin of the TenantModule.
func (TenantModule) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the TenantModule.
func (TenantModule) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("租户ID"),
		field.String("module_name").NotEmpty().Comment("模块标识"),
		field.Bool("enabled").Default(true).Comment("是否启用"),
	}
}

// Edges of the TenantModule.
func (TenantModule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("modules").Unique().Required().Field("tenant_id"),
	}
}

// Indexes of the TenantModule.
func (TenantModule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "module_name").Unique(),
	}
}
