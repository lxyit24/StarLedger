package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ServerLease holds the schema definition for the ServerLease entity.
type ServerLease struct {
	ent.Schema
}

// Mixin of the ServerLease.
func (ServerLease) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the ServerLease.
func (ServerLease) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Comment("所属租户ID"),
		field.String("server_name").NotEmpty().Comment("服务器名称"),
		field.String("provider").Default("").Comment("服务商: 阿里云/腾讯云/华为云等"),
		field.String("config").Default("").Comment("配置: CPU/内存/硬盘"),
		field.String("ip_address").Default("").Comment("IP地址"),
		field.String("location").Default("").Comment("机房位置"),
		field.Float("monthly_cost").Default(0).Comment("月费用"),
		field.Time("start_date").Default(time.Now).Comment("租赁开始日期"),
		field.Time("end_date").Comment("租赁到期日期"),
		field.Enum("renew_cycle").Values("monthly", "quarterly", "yearly").Default("monthly").Comment("续租周期: 月/季/年"),
		field.Enum("status").Values("active", "expiring", "expired", "terminated").Default("active").Comment("状态: 正常/即将到期/已过期/已退租"),
		field.String("contract_no").Default("").Comment("合同编号"),
		field.String("remark").Default("").Comment("备注"),
	}
}

// Edges of the ServerLease.
func (ServerLease) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("server_leases").Unique().Required().Field("tenant_id"),
		edge.To("bills", Bill.Type),
	}
}

// Indexes of the ServerLease.
func (ServerLease) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("status"),
		index.Fields("end_date"),
		index.Fields("provider"),
	}
}
