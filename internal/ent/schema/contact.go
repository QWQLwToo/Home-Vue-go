package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Contact 联系方式实体
type Contact struct {
	ent.Schema
}

func (Contact) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("type").Comment("联系方式类型（Email, Github, 支付宝, 微信等）"),
		field.String("icon").Comment("图标类名"),
		field.String("url").Optional().Comment("链接URL（mailto:或https://）"),
		field.String("qr_code").Optional().Comment("二维码图片URL或路径"),
		field.String("hover_color").Optional().Comment("悬停颜色"),
		field.Int("sort_order").Default(0).Comment("排序顺序"),
	}
}

func (Contact) Edges() []ent.Edge {
	return nil
}
