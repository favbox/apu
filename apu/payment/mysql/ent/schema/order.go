package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Order struct {
	ent.Schema
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.String("merchant_id"),
		field.String("channel"),
		field.String("pay_way"),

		field.String("out_order_no"),
		field.Uint64("total_amount"),
		field.String("body"),
		field.Int8("order_status"),

		field.String("auth_code"),

		field.String("wx_appid"),
		field.String("sub_openid"),

		field.String("jump_url"),
		field.String("notify_url"),

		field.String("client_ip"),
		field.String("attach"),
		field.String("order_expiration"),
		field.String("extend_params"),
	}
}

func (Order) Edges() []ent.Edge {
	return nil
}
