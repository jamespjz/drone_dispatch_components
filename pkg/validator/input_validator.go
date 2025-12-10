package validator

import (
	"fmt"
	"net/url"
	"regexp"
)

// InputValidator 输入验证器
type InputValidator struct {
}

// NewInputValidator 创建输入验证器
func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

// ValidateUUID 验证UUID格式
func (iv *InputValidator) ValidateUUID(uuid string) error {
	if uuid == "" {
		return fmt.Errorf("UUID不能为空")
	}
	// UUID格式验证
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	if !uuidRegex.MatchString(uuid) {
		return fmt.Errorf("UUID格式无效: %s", uuid)
	}
	return nil
}

// ValidateDeviceSN 验证设备序列号
func (iv *InputValidator) ValidateDeviceSN(sn string) error {
	if sn == "" {
		return fmt.Errorf("设备序列号不能为空")
	}
	// 设备序列号长度限制
	if len(sn) > 35 {
		return fmt.Errorf("设备序列号长度超限: %d", len(sn))
	}
	// 只允许字母、数字和特定符号
	snRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !snRegex.MatchString(sn) {
		return fmt.Errorf("设备序列号包含非法字符: %s", sn)
	}
	return nil
}

// ValidateProjectName 验证项目名称
func (iv *InputValidator) ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("名称不能为空")
	}
	// 长度限制
	if len(name) > 35 {
		return fmt.Errorf("名称长度超限: %d", len(name))
	}
	return nil
}

// ValidateQueryParam 验证查询参数
func (iv *InputValidator) ValidateQueryParam(param string) error {
	if param == "" {
		return fmt.Errorf("查询参数不能为空")
	}
	// 长度限制
	if len(param) > 100 {
		return fmt.Errorf("查询参数长度超限: %d", len(param))
	}
	// url编码验证
	if _, err := url.QueryUnescape(param); err != nil {
		return fmt.Errorf("查询参数编码无效: %s", param)
	}

	return nil
}

// 全局验证器实例
var globalInputValidator *InputValidator

// 初始化全局验证器
func GetValidator() *InputValidator {
	if globalInputValidator == nil {
		globalInputValidator = NewInputValidator()
	}
	return globalInputValidator
}
