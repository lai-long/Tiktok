package config

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Path 常用路径
type Path struct {
	VideoPath  string `mapstructure:"video_path"`
	AvatarPath string `mapstructure:"avatar_path"`
	EnvPath    string `mapstructure:"env_path"`
}

// MySQLConfig Mysql配置
type MySQLConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parse_time"`
	Loc       string `mapstructure:"loc"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

// JwtConfig JWT密钥
type JwtConfig struct {
	AccessSecret  string `mapstructure:"access_secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
}

// APIConfig a
type APIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	MapAPI  string `mapstructure:"map_api"`
}

// Config 总配置
type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	Jwt   JwtConfig   `mapstructure:"jwt"`
	API   APIConfig   `mapstructure:"api"`
	Path  Path        `mapstructure:"filepath"`
}

// Cfg 调用配置
var Cfg *Config
var lock sync.RWMutex

// Load 加载配置
func Load(confPath []string) (*Config, error) {
	if err := godotenv.Load("/home/lai-long/Tiktok/.env"); err != nil {
		log.Println("load env error:", err)
	}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	for _, p := range confPath {
		v.AddConfigPath(p)
	}
	v.AutomaticEnv()
	err := v.BindEnv("mysql.password", "MYSQL_PASSWORD")
	if err != nil {
		return nil, errors.Wrap(err, "mysql password bind env error")
	}
	err = v.BindEnv("redis.password", "REDIS_PASSWORD")
	if err != nil {
		return nil, errors.Wrap(err, "redis password bind env error")
	}
	err = v.BindEnv("api.api_key", "OPENAI_API_KEY")
	if err != nil {
		return nil, errors.Wrap(err, "openai_api_key bind env error")
	}
	err = v.BindEnv("jwt.access_secret", "JWT_ACCESS_SECRET")
	if err != nil {
		return nil, errors.Wrap(err, "jwt_access_secret bind env error")
	}
	err = v.BindEnv("jwt.refresh_secret", "JWT_REFRESH_SECRET")
	if err != nil {
		return nil, errors.Wrap(err, "jwt_refresh_secret bind env error")
	}

	v.SetDefault("mysql.host", "localhost")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.user", "test")
	v.SetDefault("mysql.password", "123456")
	v.SetDefault("re.host", "localhost")
	v.SetDefault("re.port", 6379)
	v.SetDefault("re.password", "123456")

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}
	lock.Lock()
	Cfg = &cfg
	lock.Unlock()

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		var newCfg Config
		if err := v.Unmarshal(&newCfg); err != nil {
			log.Println("failed to unmarshal config")
			return
		}
		lock.Lock()
		Cfg = &newCfg
		lock.Unlock()
		log.Println("config changed successfully")
	})
	log.Println("config init successfully")
	return &cfg, nil
}
