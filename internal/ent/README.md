# Ent代码生成说明

## 生成Ent代码

在项目根目录运行以下命令来生成Ent代码：

```bash
cd internal/ent
go generate ./...
```

或者使用项目根目录的Makefile：

```bash
make generate
```

这将根据 `schema/` 目录中的schema定义自动生成所有必要的Ent代码。

## Schema文件

- `site.go` - 站点实体
- `contact.go` - 联系方式实体
- `site_config.go` - 站点配置实体
- `user.go` - 用户实体
