package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Log      LogConfig      `mapstructure:"log"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Redis    RedisConfig    `mapstructure:"redis"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
}

type AuthConfig struct {
	JWTSecret      string        `mapstructure:"jwt_secret"`
	AccessTokenTTL time.Duration `mapstructure:"access_token_ttl"`
}

// ServerConfig 服务器配置结构体
type ServerConfig struct {
	Addr    string `mapstructure:"addr"`
	GinMode string `mapstructure:"gin_mode"`
}

// LogConfig 日志配置结构体
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// PostgresConfig PostgreSQL配置结构体
type PostgresConfig struct {
	DSN string `mapstructure:"dsn"`
}

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// MinIOConfig MinIO配置结构体
type MinIOConfig struct {
	// MinIO端点
	Endpoint string `mapstructure:"endpoint"`
	// MinIO访问密钥
	AccessKey string `mapstructure:"access_key"`
	// MinIO密密钥
	SecretKey string `mapstructure:"secret_key"`
	// MinIO是否使用SSL
	UseSSL bool `mapstructure:"use_ssl"`
	// MinIO桶名
	Bucket string `mapstructure:"bucket"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	// 创建Viper实例
	v := viper.New()
	// 设置默认值
	setDefaults(v)

	v.SetEnvPrefix("WB")
	// 设置环境变量前缀为WB
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 设置环境变量键值对分隔符为下划线
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 自动加载环境变量
	v.AutomaticEnv()
	// 加载配置文件
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// 设置默认配置文件名和类型
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		//在没有指定配置文件路径时，默认从当前目录、上一级目录、应用根目录加载配置文件，从上到下找，找到了就停止
		v.AddConfigPath("./configs")
		v.AddConfigPath("../configs")
		v.AddConfigPath("/app/configs")
	}
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return nil, err
		}
	}
	// 解析配置文件
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if cfg.Auth.JWTSecret == "" {
		cfg.Auth.JWTSecret = "dev-secret-change-me"
	}
	if cfg.Auth.AccessTokenTTL <= 0 {
		cfg.Auth.AccessTokenTTL = 24 * time.Hour
	}
	// 返回解析后的配置
	return &cfg, nil
}

// setDefaults 设置默认值
// 用于设置应用配置的默认值
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.addr", ":8080")
	v.SetDefault("server.gin_mode", "debug")

	v.SetDefault("log.level", "debug")
	v.SetDefault("log.format", "console")

	v.SetDefault("auth.jwt_secret", "dev-secret-change-me")
	v.SetDefault("auth.access_token_ttl", "24h")

	v.SetDefault("postgres.dsn", "postgres://whiteboard:whiteboard@localhost:5432/whiteboard?sslmode=disable")

	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("minio.endpoint", "localhost:9000")
	v.SetDefault("minio.access_key", "minioadmin")
	v.SetDefault("minio.secret_key", "minioadmin")
	v.SetDefault("minio.use_ssl", false)
	v.SetDefault("minio.bucket", "whiteboard")
}
