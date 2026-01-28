package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// LoginHistory 登录历史记录实体
type LoginHistory struct {
	ent.Schema
}

func (LoginHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("username").Comment("用户名"),
		field.String("ip").Comment("登录IP地址"),
		field.String("location").Optional().Comment("IP地理位置"),
		field.String("user_agent").Optional().Comment("用户代理"),
		field.Time("login_time").Default(time.Now).Comment("登录时间"),
		field.Bool("success").Default(true).Comment("登录是否成功"),
	}
}

func (LoginHistory) Edges() []ent.Edge {
	return nil
}
