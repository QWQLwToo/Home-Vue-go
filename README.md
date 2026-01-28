## 个人主页 - 动态版本

基于Vue3 + Go的个人主页项目，支持动态更新配置，无需修改代码即可更新站点信息。

### 技术栈

**前端：**
- Vue3 + Vite
- CSS3 + HTML5 + JavaScript
- Vue Router
- Axios

**后端：**
- Go + Gin
- Ent ORM
- JWT认证
- SQLite数据库

### 功能特性

- ✅ 动态配置站点信息（名称、URL、图标等）
- ✅ 可视化管理界面
- ✅ 支持多种图片格式上传（jpg、png、jpeg、webp、avif等）
- ✅ 联系方式管理（Email、GitHub、支付宝、微信等）
- ✅ JWT认证保护管理接口
- ✅ 数据存储在SQLite，轻量级数据库
- ✅ 支持Linux服务器二进制打包部署

### 项目结构

```
Home-Vue-go/
├── docs/          # 文档目录
├── scripts/       # 构建脚本
├── src/           # 前端源代码
├── internal/      # Go后端代码
└── public/        # 静态资源
```

详细的项目结构说明请查看 [PROJECT_STRUCTURE.md](./docs/PROJECT_STRUCTURE.md)

### 快速开始

#### 1. 环境要求

- Go >= 1.21
- Node.js >= 16.16.0
- npm >= 8.15.0

#### 2. 安装依赖

```bash
# 安装前端依赖
npm install

# 安装Go依赖
go mod download
```

#### 3. 生成Ent代码（必须）

**重要：** 在运行项目之前，必须先生成Ent代码，否则Go代码无法编译。

```bash
# 进入ent目录
cd internal/ent

# 生成Ent代码
go generate ./...

# 返回项目根目录
cd ../..
```

**如果遇到版本兼容性错误**（如 `invalid array length` 或 `Deprecated undefined`），使用：
```bash
cd internal/ent
# 使用与go.mod一致的版本（当前是v0.13.1）
go run -mod=mod entgo.io/ent/cmd/ent@v0.13.1 generate ./schema
cd ../..
```

**重要**：版本号必须与 `go.mod` 中的 `entgo.io/ent` 版本一致。

**说明：**
- `go generate` 是Go的内置命令，用于运行代码生成工具
- `./...` 表示当前目录及所有子目录
- 这会根据 `schema/` 目录中的定义生成Ent ORM代码
- 如果Go版本较新（如go1.25+），可能需要使用 `go run` 方式

#### 4. 运行项目

**开发模式：**

需要打开两个终端窗口：

**终端1 - 启动Go后端：**
```bash
# 在项目根目录运行
go run main.go
```

**终端2 - 启动前端开发服务器：**
```bash
# 在项目根目录运行
npm run dev
```

**说明：**
- `go run main.go` 会编译并运行Go程序
- 后端默认运行在 `http://localhost:1551`
- 前端默认运行在 `http://localhost:1552`
- 首次运行会自动创建 `data/` 目录和数据库
- 可以通过环境变量 `PORT` 修改后端端口（默认1551）

**Windows用户注意：** 如果遇到中文乱码，在PowerShell中运行：
```powershell
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
chcp 65001
```

访问：
- 前端：http://localhost:1552
- 管理界面：http://localhost:1552/admin
- 登录页面：http://localhost:1552/login

**默认管理员账号：**
- 用户名：`admin`
- 密码：`admin123`

⚠️ **重要：** 首次运行后请立即修改默认密码！

#### 5. 构建部署

**构建Go后端（Linux）：**
```bash
# 生成Ent代码（必须）
cd internal/ent
go generate ./...
cd ../..

# 构建Linux二进制文件
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o home-vue-go main.go
```

**构建Go后端（Windows）：**
```bash
# 生成Ent代码（必须）
cd internal\ent
go generate ./...
cd ..\..

# 构建Windows二进制文件
go build -o home-vue-go.exe main.go
```

**构建前端：**
```bash
npm install
npm run build
```

构建完成后：
- Go二进制文件：`./home-vue-go` (Linux) 或 `./home-vue-go.exe` (Windows)
- 前端构建文件：`./dist`

**说明：**
- `go build` 编译Go程序为二进制文件
- `CGO_ENABLED=1` 启用CGO（SQLite需要）
- `GOOS=linux GOARCH=amd64` 指定目标平台和架构
- `-o` 指定输出文件名

### 部署说明

1. **上传文件到服务器：**
   - 上传 `home-vue-go` 二进制文件
   - 上传 `dist` 目录（前端构建文件）

2. **运行二进制文件：**
   ```bash
   ./home-vue-go
   ```

3. **数据存储：**
   - 数据库文件：`./data/home.db`
   - 上传的图片：`./data/uploads/`
   - 所有数据存储在二进制文件同级目录的 `data` 文件夹中

4. **环境变量（可选）：**
   ```bash
   export PORT=1551          # 服务端口，默认1551
   export JWT_SECRET=your-secret-key  # JWT密钥，建议修改
   ```

### 管理界面使用

1. 访问 `/admin` 进入管理界面
2. 使用默认账号登录
3. 在管理界面中可以：
   - **站点配置**：修改站点名称、URL、图标、描述等
   - **站点管理**：添加、编辑、删除站点链接
   - **联系方式管理**：管理Email、GitHub、支付宝、微信等联系方式

### 数据迁移

如果你之前使用的是静态版本（使用JSON配置文件），可以：

1. 启动新版本服务
2. 登录管理界面
3. 手动导入原有数据，或使用API批量导入

### API文档

#### 公开API（无需认证）

- `GET /api/sites` - 获取站点列表
- `GET /api/contacts` - 获取联系方式列表
- `GET /api/config` - 获取站点配置

#### 管理API（需要JWT认证）

- `POST /api/auth/login` - 登录获取token
- `GET /api/admin/sites` - 获取站点列表（管理）
- `POST /api/admin/sites` - 创建站点
- `PUT /api/admin/sites/:id` - 更新站点
- `DELETE /api/admin/sites/:id` - 删除站点
- `GET /api/admin/contacts` - 获取联系方式列表（管理）
- `POST /api/admin/contacts` - 创建联系方式
- `PUT /api/admin/contacts/:id` - 更新联系方式
- `DELETE /api/admin/contacts/:id` - 删除联系方式
- `GET /api/admin/config` - 获取站点配置（管理）
- `PUT /api/admin/config` - 更新站点配置
- `POST /api/admin/upload` - 上传图片文件

### 注意事项

1. **Email格式验证**：Email类型的联系方式URL必须是 `mailto:` 格式（如：`mailto:i@bsgun.cn`）
2. **图片上传**：支持 jpg、jpeg、png、gif、webp、avif、svg、bmp 格式
3. **数据库安全**：SQLite数据库文件存储在服务器本地，不对外暴露，防止数据库注入
4. **JWT密钥**：生产环境请务必修改JWT_SECRET环境变量
5. **默认密码**：首次部署后请立即修改管理员密码

### 开发命令

**Go相关命令：**
```bash
# 更新依赖
go mod tidy

# 下载依赖
go mod download

# 生成Ent代码
cd internal/ent && go generate ./... && cd ../..

# 运行后端（开发模式）
go run main.go

# 构建后端（当前平台）
go build -o home-vue-go main.go

# 构建后端（Linux）
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o home-vue-go main.go

# 查看Go版本
go version

# 查看模块信息
go mod graph
```

**前端相关命令：**
```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build

# 预览构建结果
npm run preview
```

### 许可证

MIT License
