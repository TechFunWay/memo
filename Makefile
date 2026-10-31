.PHONY: help clean-dist build-frontend deploy-frontend run deploy dev

.DEFAULT_GOAL := help

PORT ?= 8905

# 显示帮助信息
help:
	@echo "可用命令:"
	@echo ""
	@echo "  make clean-dist            删除前端 dist 目录"
	@echo "  make build-frontend        打包前端 (npm run build)"
	@echo "  make deploy-frontend       清空后端 static 并复制前端 dist 到后端"
	@echo "  make run                   运行后端服务 (go run main.go)"
	@echo "  make deploy                完整流程：删除dist→打包→复制到后端→运行"
	@echo "  make dev                   构建正式结构到 dev/ 并重启服务（默认端口 8905）"
	@echo "  make dev PORT=8080         指定端口构建并重启 dev 服务"
	@echo ""

# 1. 删除前端 dist 目录
clean-dist:
	@echo "🗑️  删除前端 dist 目录..."
	rm -rf frontend/dist
	@echo "✅ 前端 dist 已删除"

# 2. 打包前端 (输出到 frontend/dist)
build-frontend:
	@echo " 打包前端..."
	cd frontend && npm run build -- --outDir dist
	@echo "✅ 前端打包完成，产物在 frontend/dist"

# 3. 清空后端 static 并复制前端 dist 到后端 static
deploy-frontend:
	@echo "🗑️  清空后端 static 目录..."
	rm -rf backend/static/*
	@echo "📂 创建后端 static 目录..."
	mkdir -p backend/static
	@echo "📋 复制前端 dist 到后端 static..."
	cp -r frontend/dist/* backend/static/
	@echo "✅ 前端资源已部署到 backend/static"

# 4. 运行 go run main.go
run:
	@if lsof -Pi :8905 -sTCP:LISTEN -t >/dev/null 2>&1; then \
		echo "⚠️  端口 8905 已被占用，正在停止旧服务..."; \
		lsof -ti:8905 | xargs kill -9 2>/dev/null; \
		sleep 1; \
	fi
	@echo "🚀 启动后端服务..."
	cd backend && go run main.go -port 8905

# 本地正式结构测试：编译前端和后端，组装到 dev/，停止旧端口服务后运行
dev:
	@PORT="$(PORT)" ./scripts/dev.sh

# 完整流程：删除dist→打包→复制到后端→运行
deploy: clean-dist build-frontend deploy-frontend run
