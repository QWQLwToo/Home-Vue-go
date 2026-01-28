package database

import (
	"context"
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"home-vue-go/internal/ent"
	"home-vue-go/internal/ent/migrate"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	Client *ent.Client
}

func Init(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_fk=1")
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))

	// 运行数据库迁移
	ctx := context.Background()
	if err := client.Schema.Create(ctx, migrate.WithForeignKeys(false)); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化默认数据
	if err := initDefaultData(ctx, client); err != nil {
		log.Printf("初始化默认数据失败: %v", err)
	}

	return &Database{Client: client}, nil
}

func (d *Database) Close() error {
	return d.Client.Close()
}

func initDefaultData(ctx context.Context, client *ent.Client) error {
	// 检查是否已有站点配置
	_, err := client.SiteConfig.Get(ctx, 1)
	if err == nil {
		// 已存在配置，不初始化
		return nil
	}

	// 创建默认站点配置
	_, err = client.SiteConfig.Create().
		SetSiteName("个人主页").
		SetSiteURL("https://example.com").
		SetSiteIcon("/favicon.ico").
		SetSiteDescription("一个基于Vue3的个人主页").
		SetSiteKeywords("个人主页,Vue3").
		SetUserName("用户").
		SetProfileImageURL("").
		SetIcpNumber("").
		SetPoliceNumber("").
		SetPageTitle("个人主页").
		SetFavicon("/favicon.ico").
		SetUmamiScript("").
		SetUmamiWebsiteID("").
		SetIconLibrary("//lib.baomitu.com/font-awesome/6.5.0/css/all.min.css").
		SetFontLibrary("").
		Save(ctx)
	if err != nil {
		return err
	}

	// 检查是否已有用户，如果没有才创建默认管理员用户
	userCount, err := client.User.Query().Count(ctx)
	if err != nil {
		return err
	}
	
	// 只有在没有任何用户时才创建默认管理员
	if userCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = client.User.Create().
			SetUsername("admin").
			SetPassword(string(hashedPassword)).
			Save(ctx)
		return err
	}

	return nil
}
