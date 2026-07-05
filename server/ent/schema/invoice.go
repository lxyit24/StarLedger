package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Invoice holds the schema definition for the Invoice entity.
type Invoice struct {
	ent.Schema
}

// Mixin of the Invoice.
func (Invoice) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Invoice.
func (Invoice) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("所属租户ID"),
		field.Int("bill_id").Optional().Comment("关联账单ID"),
		field.String("invoice_no").NotEmpty().Unique().Comment("发票号码"),
		field.String("invoice_code").Default("").Comment("发票代码"),
		field.Enum("invoice_type").
			Values("vat_normal", "vat_special").
			Default("vat_normal").
			Comment("发票类型: vat_normal(增值税普通发票)/vat_special(增值税专用发票)"),
		field.Float("amount").Default(0).Comment("不含税金额"),
		field.Float("tax_amount").Default(0).Comment("税额"),
		field.Float("total_amount").Default(0).Comment("价税合计"),
		field.Float("tax_rate").Default(0).Comment("税率(%)"),
		field.String("buyer_name").Default("").Comment("购买方名称"),
		field.String("buyer_tax_no").Default("").Comment("购买方纳税人识别号"),
		field.String("buyer_address").Default("").Comment("购买方地址电话"),
		field.String("buyer_bank").Default("").Comment("购买方开户行及账号"),
		field.String("seller_name").Default("").Comment("销售方名称"),
		field.String("seller_tax_no").Default("").Comment("销售方纳税人识别号"),
		field.String("seller_address").Default("").Comment("销售方地址电话"),
		field.String("seller_bank").Default("").Comment("销售方开户行及账号"),
		field.Time("invoice_date").Optional().Nillable().Comment("开票日期"),
		field.Enum("status").
			Values("draft", "issued", "cancelled", "red").
			Default("draft").
			Comment("发票状态: draft(草稿)/issued(已开具)/cancelled(已作废)/red(红冲)"),
		field.String("items_detail").Default("").Comment("发票明细(JSON格式)"),
		field.String("remark").Default("").Comment("备注"),
	}
}

// Edges of the Invoice.
func (Invoice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("invoices").Unique().Required().Field("tenant_id"),
		edge.From("bill", Bill.Type).Ref("invoices").Unique().Field("bill_id"),
	}
}

// Indexes of the Invoice.
func (Invoice) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("invoice_no"),
		index.Fields("status"),
		index.Fields("bill_id"),
		index.Fields("tenant_id", "status"),
		index.Fields("tenant_id", "invoice_date"),
	}
}
