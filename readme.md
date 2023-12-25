# google登录

使用 Google 帐号登录，提供了一种对用户进行身份验证的方式

[文档](https://developers.google.com/identity)

## [Sign In with Google for Web](https://developers.google.com/identity/gsi/web/guides/overview)

### [设置](https://developers.google.com/identity/gsi/web/guides/get-google-api-clientid)

已获授权的 JavaScript 来源 - 添加登录页地址

对于本地测试或开发，请同时添加 `http://localhost` 和 `http://localhost:<port_number>`

### [流程](https://developers.google.com/identity/gsi/web/guides/integrate)

用户操作后，前端请求后端，用token获取用户信息

## 运行

1. 克隆项目 `git clone -b google https://github.com/violetpupil/demos demos-google`
2. 下载依赖 `go mod tidy`
3. 复制 `.env.exp` 创建 `.env`，填写环境变量
4. 运行项目 `go run .`
