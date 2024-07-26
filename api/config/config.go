package config

import (
	"fmt"

	"apu/config/middleware"
	"apu/config/packaging"
)

var ApuConfig = newConfig()

type config struct {
	// 打包信息
	PackagingInfo *packaging.Info

	// 部署配置

	// 功能配置

	// 中间件配置
	MiddlewareConfig *middleware.Config

	// 扩展服务配置

	// 企业功能配置

	Debug                        bool
	CodeMaxNumber                int
	CodeMinNumber                int
	CodeMaxStringLength          int
	CodeMaxStringArrayLength     int
	CodeMaxObjectArrayLength     int
	CodeMaxNumberArrayLength     int
	HttpRequestMaxConnectTimeout int
	HttpRequestMaxReadTimeout    int
	HttpRequestMaxWriteTimeout   int
	HttpRequestNodeMaxBinarySize int
	HttpRequestNodeMaxTextSize   int
	SsrfProxyHttpUrl             string
	SsrfProxyHttpsUrl            string
}

// newConfig 创建带有默认值的 Config 实例的函数
func newConfig() *config {
	return &config{
		PackagingInfo: packaging.NewPackagingInfo(),
		//DeploymentConfig:                  DeploymentConfig{},
		//FeatureConfig:                     FeatureConfig{},
		MiddlewareConfig: middleware.NewConfig(),
		//ExtraServiceConfig:                ExtraServiceConfig{},
		//EnterpriseFeatureConfig:           EnterpriseFeatureConfig{},
		Debug:                        false,
		CodeMaxNumber:                9223372036854775807,
		CodeMinNumber:                -9223372036854775808,
		CodeMaxStringLength:          80000,
		CodeMaxStringArrayLength:     30,
		CodeMaxObjectArrayLength:     30,
		CodeMaxNumberArrayLength:     1000,
		HttpRequestMaxConnectTimeout: 300,
		HttpRequestMaxReadTimeout:    600,
		HttpRequestMaxWriteTimeout:   600,
		HttpRequestNodeMaxBinarySize: 1024 * 1024 * 10,
		HttpRequestNodeMaxTextSize:   1024 * 1024,
		SsrfProxyHttpUrl:             "",
		SsrfProxyHttpsUrl:            "",
	}
}

// HTTPRequestNodeReadableMaxBinarySize 计算并获取 HttpRequestNodeMaxBinarySize 的方法
func (c *config) HTTPRequestNodeReadableMaxBinarySize() string {

	return fmt.Sprintf("%.2fMB", float64(c.HttpRequestNodeMaxBinarySize)/1024/1024)
}

// HTTPRequestNodeReadableMaxTextSize 计算并获取 HTTP_REQUEST_NODE_READABLE_MAX_TEXT_SIZE 的方法
func (c *config) HTTPRequestNodeReadableMaxTextSize() string {
	return fmt.Sprintf("%.2fMB", float64(c.HttpRequestNodeMaxTextSize)/1024/1024)
}
