package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("所属租户ID"),
		field.String("username").NotEmpty().Unique().Comment("用户名"),
		field.String("password_hash").NotEmpty().Comment("密码哈希"),
		field.String("real_name").Default("").Comment("真实姓名"),
		field.String("email").Default("").Comment("邮箱"),
		field.String("phone").Default("").Comment("手机号"),
		field.Enum("status").Values("active", "disabled").Default("active").Comment("状态: active/disabled"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("users").Unique().Required().Field("tenant_id"),
		edge.From("roles", Role.Type).Ref("users"),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("username"),
		index.Fields("status"),
	}
}
