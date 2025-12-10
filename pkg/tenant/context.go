package tenant

import (
	"context"
	"fmt"
	"time"
)

type contextKey string

const (
	// TenantContextKey 租户上下文键
	TenantContextKey contextKey = "tenant"
	// RequestIDContextKey 请求ID上下文键
	RequestIDContextKey contextKey = "request_id"
)

// TenantInfo 租户信息
type TenantInfo struct {
	TenantId    int64             `json:"tenant_id"`
	UserToken   string            `json:"user_token"`
	OrgID       string            `json:"org_id"`
	ProjectUUID string            `json:"project_uuid"`
	Permissions []string          `json:"permissions"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	ExpiresAt   time.Time         `json:"expires_at"`
}

// NewTenantInfo 创建新的租户信息
func NewTenantInfo(tenantID int64, userToken, projectUUID string) *TenantInfo {
	return &TenantInfo{
		TenantId:    tenantID,
		UserToken:   userToken,
		ProjectUUID: projectUUID,
		Permissions: []string{},
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour), // 默认过期时间为24小时
	}
}

// WithTenant 添加租户信息到上下文
func WithTenant(ctx context.Context, tenantInfo *TenantInfo) context.Context {
	return context.WithValue(ctx, TenantContextKey, tenantInfo)

}

// WithRequestID 添加请求ID到上下文
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDContextKey, requestID)

}

// GetTenantFromContext 从上下文获取租户信息
func GetTenantFromContext(ctx context.Context) (*TenantInfo, error) {
	tenant, ok := ctx.Value(TenantContextKey).(*TenantInfo)
	if !ok {
		return nil, fmt.Errorf("租户信息未找到")
	}

	if !tenant.IsValid() {
		return nil, fmt.Errorf("租户信息无效或已过期")
	}

	return tenant, nil
}

// GetRequestIDFromContext 从上下文获取请求ID
func GetRequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDContextKey).(string)
	if !ok {
		return ""
	}
	return requestID
}

// IsValid 检查租户信息是否有效
func (ti *TenantInfo) IsValid() bool {
	if ti.TenantId == 0 || ti.UserToken == "" {
		return false
	}
	// 检查是否过期
	if !ti.ExpiresAt.IsZero() && time.Now().After(ti.ExpiresAt) {
		return false
	}
	return true
}
