package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Path struct {
	VideoPath  string `mapstructure:"video_path"`
	AvatarPath string `mapstructure:"avatar_path"`
	EnvPath    string `mapstructure:"env_path"`
}
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
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}
type JwtConfig struct {
	AccessSecret  string `mapstructure:"access_secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
}
type ApiConfig struct {
	ApiKey  string `mapstructure:"api_key"`
	BaseUrl string `mapstructure:"base_url"`
}
type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	Jwt   JwtConfig   `mapstructure:"jwt"`
	Api   ApiConfig   `mapstructure:"api"`
	Path  Path        `mapstructure:"path"`
}

var Cfg *Config

func Load(confPath []string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("load env error:", err)
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	for _, p := range confPath {
		v.AddConfigPath(p)
	}

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

	v.AutomaticEnv()
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
	Cfg = &cfg
	log.Println("cfg", Cfg)
	return &cfg, nil
}
