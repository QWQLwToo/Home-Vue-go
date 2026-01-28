package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// SiteConfig 站点配置实体（单例）
type SiteConfig struct {
	ent.Schema
}

func (SiteConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Default(1),
		field.String("site_name").Comment("站点名称"),
		field.String("site_url").Comment("站点URL"),
		field.String("site_icon").Comment("站点图标"),
		field.String("site_description").Comment("站点描述"),
		field.String("site_keywords").Comment("站点关键词"),
		field.String("user_name").Comment("用户名"),
		field.String("profile_image_url").Optional().Comment("头像URL"),
		field.String("icp_number").Optional().Comment("ICP备案号"),
		field.String("police_number").Optional().Comment("公安备案号"),
		// 新增字段：对应.env文件中的配置
		field.String("page_title").Optional().Comment("网页标题"),
		field.String("favicon").Optional().Comment("网页图标路径"),
		field.String("umami_script").Optional().Comment("Umami统计脚本地址"),
		field.String("umami_website_id").Optional().Comment("Umami统计网站ID"),
		field.String("icon_library").Optional().Comment("图标库CDN地址"),
		field.String("font_library").Optional().Comment("字体库CDN地址"),
	}
}

func (SiteConfig) Edges() []ent.Edge {
	return nil
}
