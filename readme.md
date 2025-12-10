## 🚀 特性 - 插件说明
- [lecheng-drone](https://gitee.com/jamespi/lecheng-drone) - 这是一个无人机控制的核心库，提供了无人机的基本操作接口。
- 插件调用：go get gitee.com/jamespi/drone_dispatch@v0.1.1
- 整体插件采用六边形架构思想构建

## 🏗️ 架构设计
- 架构图![架构图](2.jpg)
- dji_dock2插件业务架构图![架构图](1.png)

## 🚀 主要特性

### 1. 安全配置管理

- **加密存储**: 敏感信息（如token、密钥）现在使用AES-GCM加密存储
- **环境变量**: 支持从环境变量读取加密密钥
- **配置验证**: 增强的配置验证机制

```go
// 使用安全配置
secureConfig := config.GetSecureConfig()
err := secureConfig.SetSecureToken("fh2_token", "your-sensitive-token")
token, err := secureConfig.GetSecureToken("fh2_token")
```

### 2. 安全HTTP客户端

- **TLS配置**: 强制使用TLS 1.2+
- **连接池**: 优化的连接池管理
- **超时控制**: 30秒请求超时
- **URL验证**: 防止SSRF攻击

```go
// 使用安全HTTP客户端
client := httpclient.NewSecureHTTPClient()
resp, err := client.DoRequest(ctx, "GET", url, nil, headers)
```

### 3. 多租户支持

- **租户上下文**: 基于Context的租户信息传递
- **状态隔离**: 每个租户独立的状态管理
- **权限控制**: 租户级别的权限验证

```go
// 创建租户上下文
tenantInfo := tenant.NewTenantInfo("tenant_001", "token", "org_001")
ctx := tenant.WithTenant(context.Background(), tenantInfo)

// 使用多租户API
projectList, err := fh2.GetProjectListWithContext(ctx)
```

### 4. 输入验证

- **UUID验证**: 严格的UUID格式验证
- **参数清理**: 自动清理和转义输入参数
- **长度限制**: 防止过长输入导致的攻击
- **注入防护**: 防止XSS和SQL注入

```go
// 输入验证
validator := validator.GetValidator()
err := validator.ValidateUUID(projectUuid)
err := validator.ValidateDeviceSN(deviceSn)
```



## 🚀 快速开始 - 插件调用示例

```azure
func main() {
	// 初始化配置 - 这是必须的第一步
if err := config.InitDefaultConfig(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	// 司空2调用
	// 启用指定插件
	plugin.LoadEnableList([]string{"fh2", "dji_dock2"})
	// 多租户使用
	demonstrateMultiTenant()

	// 单租户使用（向后兼容  已废除）
	//demonstrateBackwardCompatible()
}

// demonstrateMultiTenant 多租户使用
func demonstrateMultiTenant() {
	log.Println("多租户使用示例开始...")
	// 创建多租户
	tenantInfo := tenant.NewTenantInfo(1, "eyJhbGciOiJIUzUxMiIsImNyaXQiOlsidHlwIiwiYWxnIiwia2lkIl0sImtpZCI6IjhiZmRiZmRkLWM4OGYtNGE5Yi04NzI3LWQ0ZGYzYWE5OTJlOSIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50IjoiMTU1MDMwMzc3NjAiLCJleHAiOjIwNjU4NzM3NDEsIm5iZiI6MTc1MDM0MDk0MSwib3JnYW5pemF0aW9uX3V1aWQiOiJhODg3ZjRkMy0wMTg2LTQ1OGMtOTBhMC1jMWQ1MGU4ZjM4ZjciLCJwcm9qZWN0X3V1aWQiOiIiLCJzdWIiOiJmaDIiLCJ1c2VyX2lkIjoiMTU3NDcwMzk4NDY5MTQxMjk5MiJ9.HwRfFQUXT3vGdElPCLFg06d-BzRcRREyvWJfFtzvdYrmVHB-zy9bZEN08BSYKxCpGPKp8F2_vO39U9-zY9E1uA", "c33595a4-3996-481d-9d81-459d435ade84")
	tenantInfo.Permissions = []string{"fh2:read", "fh2:write"}
	tenantInfo.ExpiresAt = time.Now().Add(12 * time.Hour) // 设置过期时间为2小时后

	// 创建带租户信息得上下文
	ctx := context.Background()
	ctx = tenant.WithTenant(ctx, tenantInfo)
	ctx = tenant.WithRequestID(ctx, uuid.New().String())

	// 获取插件实例
if fh2, ok := plugin.Get[service.FH2DroneAdapter](plugin.FH2Plugin); ok {
		log.Println("开始使用多租户插件...")
		// 在这里可以传递ctx给插件的方法，以便插件内部使用租户信息
		// 获取组织下的项目列表
if projectList, err := fh2.GetprojectList(); err != nil {
			log.Println("获取项目列表失败:", err)
		} else {
			log.Println("获取项目列表成功, 租户 %s 的项目列表: %s\n", tenantInfo.TenantId, projectList)
		}
	} else {
		fmt.Println("插件未启用或不存在")
	}
}

```


## 📖 依赖插件
- go get gopkg.in/yaml.v3 （废弃）
- go get github.com/spf13/viper


## 🐳 安全建议

### 生产环境部署

1. **设置强加密密钥**:
   ```bash
   export DRONE_ENCRYPTION_KEY="your-very-secure-32-character-key"
   ```

2. **启用HTTPS验证**:
   ```go
   // 在生产环境中设置
   InsecureSkipVerify: false
   ```

3. **限制网络访问**:
  - 配置防火墙规则
  - 使用VPN或专网

4. **定期更新密钥**:
  - 定期轮换API密钥
  - 监控异常访问

### 监控和告警

1. **日志监控**: 监控错误日志和异常访问
2. **性能监控**: 监控API响应时间
3. **安全告警**: 配置安全事件告警

## 迁移步骤

### 1. 更新依赖

```bash
go mod tidy
```

### 2. 配置环境变量

```bash
export DRONE_ENCRYPTION_KEY="your-32-character-encryption-key"
```

### 3. 更新配置文件

将敏感信息从配置文件中移除，使用环境变量或加密存储。

### 4. 测试验证

```bash
go test ./...
go run main.go
```

