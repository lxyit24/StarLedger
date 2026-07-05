package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Mixin of the Role.
func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("所属租户ID"),
		field.String("name").NotEmpty().Comment("角色名称: admin/finance/operator"),
		field.JSON("permissions", []string{}).Optional().Comment("权限列表(JSON)"),
		field.String("description").Default("").Comment("角色描述"),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("roles").Unique().Required().Field("tenant_id"),
		edge.To("users", User.Type),
	}
}

// Indexes of the Role.
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
	}
}
