.PHONY: generate build build-linux build-windows build-darwin clean run dist frontend backend

# 生成Ent代码
generate:
	cd internal/ent && go generate ./...

# 构建前端
frontend:
	npm install
	npm run build

# 构建后端（Linux）
backend-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/home-vue-go main.go
	chmod +x dist/home-vue-go

# 构建后端（Windows）
backend-windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/home-vue-go.exe main.go

# 构建后端（macOS）
backend-darwin:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/home-vue-go main.go
	chmod +x dist/home-vue-go

# 构建后端（当前平台）
backend:
	CGO_ENABLED=1 go build -ldflags="-s -w" -o dist/home-vue-go main.go
	chmod +x dist/home-vue-go 2>/dev/null || true

# 完整构建到dist目录（Linux）- 单一可执行文件
build-linux: clean generate
	@echo "========================================"
	@echo "构建 Linux 版本"
	@echo "========================================"
ifeq ($(OS),Windows_NT)
	@if not exist dist mkdir dist
	@call npm install && call npm run build
	@if not exist dist\index.html (echo 错误: 前端构建失败 && exit 1)
	@echo [提示] 在Windows上交叉编译Linux版本需要gcc工具链
	@echo [提示] 如果遇到错误，推荐使用WSL或在Linux系统上直接构建
	@echo [开始编译]...
	@set CGO_ENABLED=1
	@set GOOS=linux
	@set GOARCH=amd64
	@go build -ldflags="-s -w" -o dist\home-vue-go main.go
	@if errorlevel 1 (
		@echo.
		@echo [错误] 交叉编译失败
		@echo [原因] 在Windows上交叉编译Linux版本需要Linux的gcc工具链
		@echo.
		@echo [解决方案1] 使用WSL（推荐）:
		@echo    wsl
		@echo    cd /mnt/d/Desktop/Home-Vue-go
		@echo    make build-linux
		@echo.
		@echo [解决方案2] 在Linux服务器上直接构建:
		@echo    git clone ^<your-repo^>
		@echo    cd Home-Vue-go
		@echo    make build-linux
		@echo.
		@echo [解决方案3] 安装gcc工具链（复杂）:
		@echo    - 使用MSYS2: pacman -S mingw-w64-x86_64-gcc
		@echo    - 或使用TDM-GCC
		@exit /b 1
	)
	@if exist dist\static rmdir /s /q dist\static 2>nul
	@if exist dist\index.html del /f /q dist\index.html 2>nul
	@if exist dist\favicon.ico del /f /q dist\favicon.ico 2>nul
	@echo.
	@echo 构建完成！单一可执行文件: dist\home-vue-go
	@echo 后端API: http://localhost:1551
	@echo 前端界面: http://localhost:1552
	@echo 1Panel配置: 只需配置1551和1552端口，运行命令: ./home-vue-go
else
	@mkdir -p dist
	@npm install && npm run build
	@if [ ! -f "dist/index.html" ]; then echo "错误: 前端构建失败"; exit 1; fi
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/home-vue-go main.go
	@chmod +x dist/home-vue-go
	@rm -rf dist/static dist/index.html dist/favicon.ico 2>/dev/null || true
	@echo ""
	@echo "构建完成！单一可执行文件: dist/home-vue-go"
	@echo "后端API: http://localhost:1551"
	@echo "前端界面: http://localhost:1552"
	@echo "1Panel配置: 只需配置1551和1552端口，运行命令: ./home-vue-go"
endif

# 完整构建到dist目录（Windows）- 单一可执行文件
build-windows: clean generate
	@echo ========================================
	@echo 构建 Windows 版本
	@echo ========================================
	@if not exist dist mkdir dist
	@npm install && npm run build
	@if not exist dist\index.html (echo 错误: 前端构建失败 && exit 1)
	@set CGO_ENABLED=1 && go build -ldflags="-s -w" -o dist\home-vue-go.exe main.go
	@if exist dist\static rmdir /s /q dist\static
	@if exist dist\index.html del /f /q dist\index.html
	@if exist dist\favicon.ico del /f /q dist\favicon.ico
	@echo.
	@echo 构建完成！单一可执行文件: dist\home-vue-go.exe
	@echo 后端API: http://localhost:1551
	@echo 前端界面: http://localhost:1552
	@echo 1Panel配置: 只需配置1551和1552端口，运行命令: ./home-vue-go.exe

# 完整构建到dist目录（macOS）- 单一可执行文件
build-darwin: clean generate
	@echo "========================================"
	@echo "构建 macOS 版本"
	@echo "========================================"
ifeq ($(OS),Windows_NT)
	@if not exist dist mkdir dist
	@call npm install && call npm run build
	@if not exist dist\index.html (echo 错误: 前端构建失败 && exit 1)
	@set CGO_ENABLED=1 && set GOOS=darwin && set GOARCH=amd64 && go build -ldflags="-s -w" -o dist\home-vue-go main.go
	@if exist dist\static rmdir /s /q dist\static 2>nul
	@if exist dist\index.html del /f /q dist\index.html 2>nul
	@if exist dist\favicon.ico del /f /q dist\favicon.ico 2>nul
	@echo.
	@echo 构建完成！单一可执行文件: dist\home-vue-go
	@echo 后端API: http://localhost:1551
	@echo 前端界面: http://localhost:1552
else
	@mkdir -p dist
	@npm install && npm run build
	@if [ ! -f "dist/index.html" ]; then echo "错误: 前端构建失败"; exit 1; fi
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/home-vue-go main.go
	@chmod +x dist/home-vue-go
	@rm -rf dist/static dist/index.html dist/favicon.ico 2>/dev/null || true
	@echo ""
	@echo "构建完成！单一可执行文件: dist/home-vue-go"
	@echo "后端API: http://localhost:1551"
	@echo "前端界面: http://localhost:1552"
endif

# 完整构建到dist目录（当前平台）- 单一可执行文件
build: clean generate
	@echo "========================================"
	@echo "构建当前平台版本"
	@echo "========================================"
ifeq ($(OS),Windows_NT)
	@if not exist dist mkdir dist
	@call npm install && call npm run build
	@if not exist dist\index.html (echo 错误: 前端构建失败 && exit 1)
	@set CGO_ENABLED=1 && go build -ldflags="-s -w" -o dist\home-vue-go.exe main.go
	@if exist dist\static rmdir /s /q dist\static 2>nul
	@if exist dist\index.html del /f /q dist\index.html 2>nul
	@if exist dist\favicon.ico del /f /q dist\favicon.ico 2>nul
	@echo.
	@echo 构建完成！单一可执行文件: dist\home-vue-go.exe
	@echo 后端API: http://localhost:1551
	@echo 前端界面: http://localhost:1552
	@echo 1Panel配置: 只需配置1551和1552端口，运行命令: ./home-vue-go.exe
else
	@mkdir -p dist
	@npm install && npm run build
	@if [ ! -f "dist/index.html" ]; then echo "错误: 前端构建失败"; exit 1; fi
	@CGO_ENABLED=1 go build -ldflags="-s -w" -o dist/home-vue-go main.go
	@chmod +x dist/home-vue-go 2>/dev/null || true
	@rm -rf dist/static dist/index.html dist/favicon.ico 2>/dev/null || true
	@echo ""
	@echo "构建完成！单一可执行文件: dist/home-vue-go"
	@echo "后端API: http://localhost:1551"
	@echo "前端界面: http://localhost:1552"
	@echo "1Panel配置: 只需配置1551和1552端口，运行命令: ./home-vue-go"
endif

# 打包到dist目录（推荐使用）
dist: build

# 运行开发服务器
run:
	go run main.go

# 清理
clean:
	@echo "清理构建文件..."
ifeq ($(OS),Windows_NT)
	@if exist home-vue-go.exe del /f /q home-vue-go.exe 2>nul
	@if exist home-vue-go del /f /q home-vue-go 2>nul
	@if exist dist rmdir /s /q dist 2>nul
	@if exist node_modules rmdir /s /q node_modules 2>nul
else
	@rm -f home-vue-go home-vue-go.exe 2>/dev/null || true
	@rm -rf dist 2>/dev/null || true
	@rm -rf node_modules 2>/dev/null || true
endif
	@echo "清理完成"
