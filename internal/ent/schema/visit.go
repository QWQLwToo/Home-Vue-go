package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Visit 访问记录实体
type Visit struct {
	ent.Schema
}

func (Visit) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("path").Comment("访问路径"),
		field.String("ip").Comment("访问IP"),
		field.String("user_agent").Optional().Comment("用户代理"),
		field.String("referer").Optional().Comment("来源页面"),
		field.Time("visit_time").Default(time.Now).Comment("访问时间"),
	}
}

func (Visit) Edges() []ent.Edge {
	return nil
}
