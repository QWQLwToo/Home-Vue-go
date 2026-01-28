package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Site 站点实体
type Site struct {
	ent.Schema
}

func (Site) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("name").Comment("站点名称"),
		field.String("url").Comment("站点URL"),
		field.String("icon").Comment("站点图标类名"),
		field.Int("sort_order").Default(0).Comment("排序顺序"),
	}
}

func (Site) Edges() []ent.Edge {
	return nil
}
