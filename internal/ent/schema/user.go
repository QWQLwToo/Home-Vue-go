package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User 用户实体（用于JWT认证）
type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("username").Unique().Comment("用户名"),
		field.String("password").Comment("密码（bcrypt哈希）"),
	}
}

func (User) Edges() []ent.Edge {
	return nil
}
