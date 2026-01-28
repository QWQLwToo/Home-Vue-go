# 构建说明

本项目支持打包为**单一可执行文件**，包含前后端，启动后同时提供：
- **1551端口**：后端API服务
- **1552端口**：前端界面服务（自动代理API请求到1551）

所有前端文件已嵌入到二进制文件中，无需额外文件。

## 构建方式

### 使用 Makefile（推荐）

本项目使用 `Makefile` 进行跨平台构建，支持 Windows、Linux 和 macOS。

**前置要求：**
- 安装 Node.js 和 npm
- 安装 Go 1.23+
- 安装 make 工具
  - **Windows**: 可通过 Git for Windows、Chocolatey 或 WSL 安装
  - **Linux/macOS**: 通常已预装

### Windows 上使用 Make

如果你安装了 Git for Windows，make 工具通常已经包含在内，但可能不在 PATH 中。

**方式1：使用提供的 make.bat 包装脚本（推荐）**
```bash
make.bat build-linux
make.bat build-windows
make.bat build
```

这个脚本会自动查找 Git 自带的 make 工具。

**方式2：将 Git 的 make 添加到 PATH**
1. 找到 Git 安装目录（通常是 `C:\Program Files\Git`）
2. 将 `C:\Program Files\Git\usr\bin` 添加到系统 PATH 环境变量
3. 重启终端后即可直接使用 `make` 命令

**方式3：使用 Chocolatey 安装 make**
```powershell
choco install make
```

**方式4：使用 WSL**
在 WSL 中运行 make 命令。

### 构建命令

#### 构建当前平台版本
```bash
make build
```

#### 构建 Linux 版本（用于服务器部署）
```bash
make build-linux
```

#### 构建 Windows 版本
```bash
make build-windows
```

#### 构建 macOS 版本
```bash
make build-darwin
```

### 其他 Makefile 命令

**只构建前端：**
```bash
make frontend
```

**只构建后端（当前平台）：**
```bash
make backend
```

**只构建后端（指定平台）：**
```bash
make backend-linux    # Linux
make backend-windows  # Windows
make backend-darwin   # macOS
```

**生成Ent代码（首次构建前需要）：**
```bash
make generate
```

**清理构建文件：**
```bash
make clean
```

**运行开发服务器：**
```bash
make run
```

## 构建输出

构建完成后，`dist` 目录将只包含**单一可执行文件**：

```
dist/
└── home-vue-go.exe    # Windows单一可执行文件（包含前后端）
或
└── home-vue-go        # Linux/macOS单一可执行文件（包含前后端）
```

**注意**：所有前端文件（HTML、CSS、JavaScript等）都已嵌入到二进制文件中，构建脚本会自动清理dist目录中的前端源文件。

## 运行服务器

### Windows
```bash
cd dist
home-vue-go.exe
```

### Linux/macOS
```bash
cd dist
./home-vue-go
```

## 访问地址

启动后，服务器会同时提供两个服务：

- **后端API**: http://localhost:1551
  - API接口：http://localhost:1551/api
  - 管理接口：http://localhost:1551/api/admin

- **前端界面**: http://localhost:1552
  - 主页：http://localhost:1552
  - 管理界面：http://localhost:1552/admin
  - 登录页面：http://localhost:1552/login

**注意**：前端会自动将 `/api` 请求代理到 `http://localhost:1551`，无需额外配置。

## 1Panel 配置

在1Panel中配置非常简单：

1. **上传文件**：将 `home-vue-go`（Linux版本）上传到服务器

2. **配置端口**：
   - 后端API端口：`1551`
   - 前端服务端口：`1552`

3. **运行命令**：
   ```bash
   ./home-vue-go
   ```

4. **完成**：无需其他配置，单一文件即可运行完整服务

## 配置说明

- **后端API端口**：默认 `1551`，可通过环境变量 `API_PORT` 修改
- **前端服务端口**：默认 `1552`，可通过环境变量 `FRONTEND_PORT` 修改
- **数据目录**：运行时会自动在二进制文件同目录下创建 `data` 目录
- **默认管理员账号**：`admin` / `admin123`（首次启动时显示）

### 修改端口示例

**Windows:**
```bash
set API_PORT=8080
set FRONTEND_PORT=8081
home-vue-go.exe
```

**Linux/macOS:**
```bash
export API_PORT=8080
export FRONTEND_PORT=8081
./home-vue-go
```

## 注意事项

1. **CGO依赖**：本项目使用SQLite数据库，需要启用CGO（`CGO_ENABLED=1`）
2. **前端构建**：构建时必须先运行 `npm run build` 生成dist目录，Go编译时会嵌入这些文件
3. **Go版本**：需要Go 1.23或更高版本（支持embed功能）
4. **端口占用**：确保1551和1552端口未被占用
5. **单一文件**：构建完成后，只需一个可执行文件即可运行，无需其他依赖
6. **Make工具**：Windows用户需要安装make工具（Git for Windows自带，或使用Chocolatey安装）
7. **跨平台编译**：在Windows上交叉编译Linux版本需要gcc工具链，推荐使用WSL或在Linux系统上直接构建

## Windows 上使用 Make

### 方式1：使用提供的 make.bat 包装脚本（最简单）

项目根目录提供了 `make.bat` 脚本，它会自动查找 Git 自带的 make 工具：

```bash
# 构建 Linux 版本
make.bat build-linux

# 构建 Windows 版本
make.bat build-windows

# 构建当前平台版本
make.bat build
```

### 方式2：将 Git 的 make 添加到 PATH（推荐）

Git for Windows 自带 make 工具，通常位于：
- `C:\Program Files\Git\usr\bin\make.exe`
- `C:\Program Files (x86)\Git\usr\bin\make.exe`

**添加步骤：**
1. 右键"此电脑" → "属性" → "高级系统设置" → "环境变量"
2. 在"系统变量"中找到 `Path`，点击"编辑"
3. 添加 Git 的 bin 目录：`C:\Program Files\Git\usr\bin`
4. 点击"确定"保存
5. 重启终端后即可直接使用 `make` 命令

### 方式3：使用 Chocolatey 安装 make
```powershell
choco install make
```

### 方式4：使用 WSL
在 WSL 中运行 make 命令。

### 方式5：手动下载 make for Windows
从 https://sourceforge.net/projects/gnuwin32/files/make/ 下载并安装

## 架构说明

### 单一可执行文件架构

```
┌─────────────────────────────────┐
│   home-vue-go (单一可执行文件)   │
│   包含: 后端代码 + 前端文件(嵌入) │
├─────────────────────────────────┤
│                                  │
│  ┌──────────────┐  ┌──────────┐ │
│  │  后端API服务  │  │ 前端服务 │ │
│  │  端口: 1551  │  │ 端口:1552│ │
│  └──────┬───────┘  └────┬─────┘ │
│         │               │        │
│         │               │        │
│         └───────┬───────┘        │
│                 │                │
│            API代理               │
│    (前端/api/* → 后端1551)        │
└─────────────────────────────────┘
```

- **后端服务（1551）**：提供所有API接口
- **前端服务（1552）**：
  - 从嵌入的文件系统提供前端静态文件（HTML、CSS、JS）
  - 自动代理 `/api/*` 请求到后端1551端口
  - 支持SPA路由

## 完整构建流程示例

### 在 Windows 上构建 Linux 版本（用于1Panel）

**方式1：使用WSL（推荐）**

**前提条件：** 确保WSL中已安装Go和Node.js

**安装Go（如果未安装）：**
```bash
# 在WSL中运行
# 下载Go（替换版本号为最新版本，当前需要Go 1.23+）
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz

# 删除旧版本（如果存在）
sudo rm -rf /usr/local/go

# 解压到/usr/local
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# 添加到PATH（添加到~/.bashrc或~/.zshrc）
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

**安装Node.js（如果未安装）：**
```bash
# 使用nvm安装（推荐）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 18
nvm use 18

# 或使用apt安装
sudo apt update
sudo apt install nodejs npm
```

**构建步骤：**
```bash
# 在WSL中运行（不需要再运行wsl命令，如果已经在WSL中）
cd /mnt/d/Desktop/Home-Vue-go
make build-linux
# 构建完成后，dist/home-vue-go 就是Linux可执行文件
```

**方式2：直接在Windows上构建（需要gcc工具链）**
```bash
# 如果遇到交叉编译错误，需要安装gcc工具链
# 使用MSYS2安装：
# pacman -S mingw-w64-x86_64-gcc

# 然后运行
make build-linux
```

**方式3：在Linux服务器上直接构建（最简单，推荐）**

这是最推荐的方式，无需交叉编译工具链，构建速度快且稳定。

**前提条件：**
- Linux服务器已安装 Go 1.23+ 和 Node.js 16.16.0+
- 已安装 make 工具（通常已预装）

**构建步骤：**
```bash
# 1. 克隆项目到服务器
git clone <your-repo>
cd Home-Vue-go

# 2. 确保环境已安装（如果未安装）
# 安装Go（示例，根据你的发行版调整）
# wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
# sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
# echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
# source ~/.bashrc

# 安装Node.js（示例，使用nvm）
# curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
# source ~/.bashrc
# nvm install 18

# 3. 构建Linux版本
make build-linux

# 4. 构建完成后，可执行文件位于 dist/home-vue-go
# 可以直接运行，或配置到1Panel
./dist/home-vue-go
```

**1Panel配置：**
- 运行命令：`./home-vue-go`（或完整路径）
- 端口映射：1551（后端API）、1552（前端界面）
- 工作目录：可执行文件所在目录

**注意**：由于项目使用SQLite（需要CGO），在Windows上交叉编译Linux版本需要额外的工具链。推荐使用WSL或在Linux系统上直接构建。

### 在 Linux 上构建 Windows 版本

```bash
make build-windows
# 输出: dist/home-vue-go.exe
```

### 在 macOS 上构建 Linux 版本

```bash
make build-linux
# 输出: dist/home-vue-go
```

## 清理构建文件

```bash
make clean
```

这会删除：
- `dist/` 目录
- `node_modules/` 目录
- 可执行文件

## 开发模式

开发时，可以分别运行前后端：

**终端1 - 后端：**
```bash
make run
# 或
go run main.go
# 后端运行在 http://localhost:1551
```

**终端2 - 前端：**
```bash
npm run dev
# 前端运行在 http://localhost:1552，自动代理API到1551
```

## 部署优势

✅ **单一可执行文件**：前后端一体化，所有文件嵌入在二进制中  
✅ **无需依赖**：不需要Node.js、npm或其他运行时  
✅ **端口分离**：API和前端服务分离，便于管理和扩展  
✅ **自动代理**：前端自动代理API请求，无需额外配置  
✅ **1Panel友好**：只需配置两个端口，运行一个命令即可  
✅ **部署简单**：上传一个文件，配置端口，即可运行  
✅ **跨平台构建**：使用make统一构建流程，支持多平台
