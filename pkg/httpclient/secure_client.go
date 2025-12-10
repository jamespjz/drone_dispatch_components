package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SecureHTTPClient struct {
	client *http.Client
}

func NewSecureHTTPClient() *SecureHTTPClient {
	return &SecureHTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
					// InsecureSkipVerify: true, // 根据需要启用，但不推荐在生产环境中使用
				},
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// DoRequest 执行HTTP请求并返回响应
func (shc *SecureHTTPClient) DoRequest(ctx context.Context, method, requestURL string, body io.Reader, headers map[string]string) (*http.Response, error) {
	// 验证URL
	if err := shc.validateURL(requestURL); err != nil {
		return nil, fmt.Errorf("URL验证失败: %w", err)
	}
	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	// 设置请求头
	for key, value := range headers {
		// 验证请求头
		if err := shc.validateHeader(key, value); err != nil {
			return nil, fmt.Errorf("请求头验证失败: %w", err)
		}
		req.Header.Set(key, value)
	}
	// 设置默认请求头
	req.Header.Set("User-Agent", "DroneDispatch/1.0")
	req.Header.Set("Accept", "application/json")
	// 执行请求
	resp, err := shc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("执行HTTP请求失败: %w", err)
	}

	return resp, nil
}

// validateURL 验证URL的安全性
func (shc *SecureHTTPClient) validateURL(requestURL string) error {
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return fmt.Errorf("URL解析失败: %w", err)
	}

	// 只允许HTTPS协议 （开发环境可以运行http）
	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return fmt.Errorf("不支持的URL协议: %s", parsedURL.Scheme)
	}

	// 验证主机名
	if parsedURL.Host == "" {
		return fmt.Errorf("URL主机名不能为空")
	}

	// 防止内网访问
	if shc.isPrivateIP(parsedURL.Hostname()) {
		//return fmt.Errorf("禁止访问内网地址: %s", parsedURL.Hostname())
	}

	return nil
}

// validateHeader 验证请求头的安全性
func (shc *SecureHTTPClient) validateHeader(key, value string) error {
	// 防止请求头注入
	if strings.ContainsAny(key, "\r\n") || strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("请求头包含非法字符")
	}
	// 限制请求头长度
	if len(key) > 100 || len(value) > 1000 {
		return fmt.Errorf("请求头长度超限")
	}

	return nil
}

// isPrivateIP 检查主机名是否为内网IP
func (shc *SecureHTTPClient) isPrivateIP(host string) bool {
	// 简单的内网IP检查，实际应用中可能需要更完善的实现
	privateRanges := []string{
		"127.",
		"10.",
		"172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.",
		"172.24.", "172.25.", "172.26.", "172.27.",
		"172.28.", "172.29.", "172.30.", "172.31.",
		"192.168.",
		"localhost",
	}

	for _, prefix := range privateRanges {
		if strings.HasPrefix(host, prefix) {
			return true
		}
	}

	return false
}
