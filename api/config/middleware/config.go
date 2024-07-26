package middleware

import (
	"fmt"
	"net/url"
	"strings"

	"apu/config/middleware/cache"
	"apu/config/middleware/storage"
	"apu/config/middleware/vdb"
)

type StorageConfig struct {
	StorageType      string
	StorageLocalPath string
}

// NewStorageConfig 创建一个带有默认值的 StorageConfig 实例的函数
func NewStorageConfig() *StorageConfig {
	return &StorageConfig{
		StorageType:      "local",
		StorageLocalPath: "storage",
	}
}

type VectorStoreConfig struct {
	VectorStore string
}

// NewVectorStoreConfig 创建带有默认值的 VectorStoreConfig 实例的函数
func NewVectorStoreConfig() *VectorStoreConfig {
	return &VectorStoreConfig{}
}

type KeywordStoreConfig struct {
	KeywordStore string
}

// NewKeywordStoreConfig 创建带有默认值的 KeywordStoreConfig 实例的函数
func NewKeywordStoreConfig() *KeywordStoreConfig {
	return &KeywordStoreConfig{
		KeywordStore: "jieba",
	}
}

type DatabaseConfig struct {
	DBHost                      string
	DBPort                      int
	DBUsername                  string
	DBPassword                  string
	DBDatabase                  string
	DBCharset                   string
	DBExtras                    string
	SQLALCHEMYDatabaseURIScheme string
	SQLALCHEMYPoolSize          int
	SQLALCHEMYMaxOverflow       int
	SQLALCHEMYPoolRecycle       int
	SQLALCHEMYPoolPrePing       bool
	SQLALCHEMYEcho              bool
}

// NewDatabaseConfig 创建一个带有默认值的 DatabaseConfig 实例的函数
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DBHost:                      "localhost",
		DBPort:                      5432,
		DBUsername:                  "postgres",
		DBPassword:                  "",
		DBDatabase:                  "dify",
		DBCharset:                   "",
		DBExtras:                    "",
		SQLALCHEMYDatabaseURIScheme: "postgresql",
		SQLALCHEMYPoolSize:          30,
		SQLALCHEMYMaxOverflow:       10,
		SQLALCHEMYPoolRecycle:       3600,
		SQLALCHEMYPoolPrePing:       false,
		SQLALCHEMYEcho:              false,
	}
}

// SQLALCHEMYDatabaseURI 计算并获取 SQLALCHEMYDatabaseURI 的方法
func (c *DatabaseConfig) SQLALCHEMYDatabaseURI() string {
	var dbExtras string
	if c.DBCharset != "" {
		dbExtras = fmt.Sprintf("%s&client_encoding=%s", c.DBExtras, c.DBCharset)
	} else {
		dbExtras = c.DBExtras
	}
	dbExtras = strings.Trim(dbExtras, "&")
	if dbExtras != "" {
		dbExtras = "?" + dbExtras
	}
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s%s", c.SQLALCHEMYDatabaseURIScheme, url.QueryEscape(c.DBUsername), url.QueryEscape(c.DBPassword), c.DBHost, c.DBPort, c.DBDatabase, dbExtras)
}

// SQLALCHEMYEngineOptions 计算并获取 SQLALCHEMYEngineOptions 的方法
func (c *DatabaseConfig) SQLALCHEMYEngineOptions() map[string]interface{} {
	return map[string]interface{}{
		"pool_size":     c.SQLALCHEMYPoolSize,
		"max_overflow":  c.SQLALCHEMYMaxOverflow,
		"pool_recycle":  c.SQLALCHEMYPoolRecycle,
		"pool_pre_ping": c.SQLALCHEMYPoolPrePing,
		"connect_args":  map[string]interface{}{"options": "-c timezone=UTC"},
	}
}

type AsynqConfig struct {
	RedisAddr string
}

func NewAsynqConfig() *AsynqConfig {
	return &AsynqConfig{
		RedisAddr: "127.0.0.1:6379",
	}
}

type Config struct {
	// 数据配置
	AsyncqConfig       *AsynqConfig
	DatabaseConfig     *DatabaseConfig
	KeywordStoreConfig *KeywordStoreConfig
	RedisConfig        *cache.RedisConfig

	// 对象存储及供应商配置
	StorageConfig *StorageConfig
	AliyunConfig  *storage.AliyunConfig

	// 向量数据库及供应商配置
	VectorStoreConfig *VectorStoreConfig
	QdrantConfig      *vdb.QdrantConfig
}

func NewConfig() *Config {
	cfg := &Config{
		DatabaseConfig:     NewDatabaseConfig(),
		KeywordStoreConfig: NewKeywordStoreConfig(),
		RedisConfig:        cache.NewRedisConfig(),
		StorageConfig:      NewStorageConfig(),
		AliyunConfig:       storage.NewAliyunConfig(),
		VectorStoreConfig:  NewVectorStoreConfig(),
		QdrantConfig:       vdb.NewQdrantConfig(),
	}
	return cfg
}
