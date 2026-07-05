package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Bill holds the schema definition for the Bill entity.
type Bill struct {
	ent.Schema
}

// Mixin of the Bill.
func (Bill) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Bill.
func (Bill) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("所属租户ID"),
		field.String("bill_no").NotEmpty().Unique().Comment("账单编号"),
		field.String("bill_type").NotEmpty().Comment("账单类型: 服务器租赁/域名租赁/IDC交易等"),
		field.Int("related_resource_id").Optional().Nillable().Comment("关联资源ID"),
		field.Float("amount").Default(0).Comment("金额"),
		field.Time("bill_date").Default(time.Now).Comment("账单日期"),
		field.Time("due_date").Comment("到期日"),
		field.Time("paid_date").Optional().Nillable().Comment("支付日期"),
		field.Enum("payment_status").Values("pending", "paid", "overdue", "cancelled").Default("pending").Comment("支付状态: 待支付/已支付/逾期/已取消"),
		field.String("payment_method").Default("").Comment("支付方式"),
		field.String("invoice_no").Default("").Comment("发票号"),
		field.String("remark").Default("").Comment("备注"),
	}
}

// Edges of the Bill.
func (Bill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("bills").Unique().Required().Field("tenant_id"),
		edge.From("server_lease", ServerLease.Type).Ref("bills").Unique().Field("related_resource_id"),
		edge.To("invoices", Invoice.Type),
	}
}

// Indexes of the Bill.
func (Bill) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("bill_no"),
		index.Fields("bill_type"),
		index.Fields("payment_status"),
		index.Fields("due_date"),
	}
}
