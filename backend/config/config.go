package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Admin    AdminConfig    `mapstructure:"admin"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	AutoMigrate  bool   `mapstructure:"auto_migrate"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expire time.Duration `mapstructure:"expire"`
}

type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Email    string `mapstructure:"email"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

var globalConfig *Config

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	// 设置环境变量前缀
	viper.SetEnvPrefix("BLOG")

	// 替换环境变量中的分隔符
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 处理环境变量替换
	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)
		if strings.HasPrefix(val, "${") && strings.HasSuffix(val, "}") {
			envVar := strings.TrimSuffix(strings.TrimPrefix(val, "${"), "}")
			parts := strings.SplitN(envVar, ":", 2)
			envKey := parts[0]
			defaultValue := ""
			if len(parts) > 1 {
				defaultValue = parts[1]
			}
			envValue := os.Getenv(envKey)
			if envValue == "" {
				envValue = defaultValue
			}
			viper.Set(key, envValue)
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	globalConfig = config
	return config, nil
}

func Get() *Config {
	return globalConfig
}

// InitTestConfig 为测试初始化默认配置
func InitTestConfig() {
	if globalConfig == nil {
		globalConfig = &Config{
			JWT: JWTConfig{
				Secret: "test-secret-key",
				Expire: 24 * time.Hour,
			},
		}
	}
}
