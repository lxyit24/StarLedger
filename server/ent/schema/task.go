package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Mixin of the Task.
func (Task) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("所属租户ID"),
		field.String("title").NotEmpty().Comment("任务标题"),
		field.String("description").Default("").Comment("任务描述"),
		field.Int("assignee_id").Optional().Comment("负责人ID"),
		field.Int("creator_id").Comment("创建人ID"),
		field.Enum("status").
			Values("pending", "in_progress", "completed", "cancelled").
			Default("pending").
			Comment("状态: pending/in_progress/completed/cancelled"),
		field.Enum("priority").
			Values("low", "medium", "high", "urgent").
			Default("medium").
			Comment("优先级: low/medium/high/urgent"),
		field.Time("due_date").Optional().Comment("截止日期"),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("tasks").Unique().Required().Field("tenant_id"),
		edge.From("assignee", User.Type).Ref("assigned_tasks").Unique().Field("assignee_id"),
		edge.From("creator", User.Type).Ref("created_tasks").Unique().Required().Field("creator_id"),
	}
}

// Indexes of the Task.
func (Task) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("assignee_id"),
		index.Fields("creator_id"),
		index.Fields("status"),
		index.Fields("priority"),
	}
}
