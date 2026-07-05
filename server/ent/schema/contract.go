package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Contract holds the schema definition for the Contract entity.
type Contract struct {
	ent.Schema
}

// Mixin of the Contract.
func (Contract) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Contract.
func (Contract) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("所属租户ID"),
		field.String("title").NotEmpty().Comment("合同标题"),
		field.String("party_a").Default("").Comment("甲方"),
		field.String("party_b").Default("").Comment("乙方"),
		field.Float("amount").Default(0).Comment("合同金额"),
		field.Time("start_date").Default(time.Now).Comment("开始日期"),
		field.Time("end_date").Comment("结束日期"),
		field.Enum("status").
			Values("draft", "active", "expired", "terminated").
			Default("draft").
			Comment("状态: draft/active/expired/terminated"),
		field.String("file_url").Default("").Comment("合同文件URL"),
		field.String("remark").Default("").Comment("备注"),
	}
}

// Edges of the Contract.
func (Contract) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("contracts").Unique().Required().Field("tenant_id"),
	}
}

// Indexes of the Contract.
func (Contract) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("status"),
		index.Fields("end_date"),
	}
}
