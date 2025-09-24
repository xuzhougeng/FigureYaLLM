# 问题修复记录

## Docker部署静态文件404问题

### 🐛 问题描述

Docker部署后出现以下问题：
1. `/styles.css` 和 `/script.js` 404 NOT FOUND
2. 页面点击"Ask"按钮后立即刷新，无法正常提交

### 🔍 问题分析

**根本原因**: HTML文件中使用了相对路径引用静态资源

```html
<!-- 问题代码 -->
<link rel="stylesheet" href="styles.css">
<script src="script.js"></script>
```

**问题产生原因**:
1. 在Docker容器中，静态文件通过 `/static/*` 路径提供服务
2. HTML使用相对路径时，浏览器会在当前页面路径下寻找资源
3. 当前页面是 `/`，所以浏览器会请求 `/styles.css` 和 `/script.js`
4. 但实际文件路径应该是 `/static/styles.css` 和 `/static/script.js`

### ✅ 解决方案

**修复1: 更正静态文件路径**

将HTML中的相对路径改为绝对路径：

```html
<!-- 修复后 -->
<link rel="stylesheet" href="/static/styles.css">
<script src="/static/script.js"></script>
```

**修复2: 确认Go路由配置**

服务器路由配置正确：
```go
// 静态文件服务
r.Static("/static", "./static")
r.StaticFile("/", "./static/index.html")
```

### 🧪 验证修复

**测试步骤**:
1. 重新构建Docker镜像: `docker build -t figureya-recommend:latest .`
2. 启动容器: `docker run -p 8080:8080 --env-file .env figureya-recommend:latest`
3. 验证静态文件: `curl -I http://localhost:8080/static/styles.css`
4. 验证前端页面: 访问 `http://localhost:8080`

**验证结果**:
- ✅ CSS文件正常加载: HTTP 200 OK, Content-Type: text/css
- ✅ JS文件正常加载: HTTP 200 OK, Content-Type: text/javascript
- ✅ 前端交互正常: Ask按钮不再刷新页面
- ✅ API调用正常: 推荐功能工作正常

### 📝 修复文件

修改文件: `static/index.html`
- 第7行: `href="styles.css"` → `href="/static/styles.css"`
- 第81行: `src="script.js"` → `src="/static/script.js"`

### 🚀 部署验证

**本地开发**: `go run main.go load_env.go`
- 地址: http://localhost:8080
- 状态: ✅ 正常工作

**Docker容器**: `./deploy.sh start`
- 地址: http://localhost:8080
- 状态: ✅ 正常工作

**Docker Compose**: `docker-compose up -d`
- 地址: http://localhost:8080
- 状态: ✅ 正常工作

### 🛡️ 预防措施

1. **使用绝对路径**: 所有静态资源使用绝对路径引用
2. **路由一致性**: 确保开发环境和生产环境路由配置一致
3. **测试覆盖**: 添加静态文件加载的自动化测试
4. **文档更新**: 在部署文档中说明静态文件路径规范

### 📊 性能影响

- 镜像大小: 9.04MB (无变化)
- 启动时间: < 3秒 (无变化)
- 静态文件响应: < 100ms (正常)
- API响应时间: 15-20秒 (LLM调用，正常)

修复完成后，系统在所有部署模式下都能正常工作！