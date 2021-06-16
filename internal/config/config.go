package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	WebAddress         string        `mapstructure:"web_address"`
	WebReadTimeout     time.Duration `mapstructure:"web_read_timeout"`
	WebWriteTimeout    time.Duration `mapstructure:"web_write_timeout"`
	WebShutdownTimeout time.Duration `mapstructure:"web_shutdown_timeout"`
	PgHost             string        `mapstructure:"pg_host"`
	PgPort             string        `mapstructure:"pg_port"`
	PgUser             string        `mapstructure:"pg_user"`
	PgPassword         string        `mapstructure:"pg_password"`
	PgName             string        `mapstructure:"pg_name"`
	RdrHost            string        `mapstructure:"rdr_host"`
	RdrPort            string        `mapstructure:"rdr_port"`
	RdrDb              int           `mapstructure:"rdr_db"`
	RdrPool            int           `mapstructure:"rdr_pool"`
	MailHost           string        `mapstructure:"mail_host"`
	MailPort           int           `mapstructure:"mail_port"`
	MailUser           string        `mapstructure:"mail_user"`
	MailPassword       string        `mapstructure:"mail_password"`
	IsDev              bool          `mapstructure:"is_dev"`
}

func load() Config {
	var config Config
	v := viper.New()

	v.SetConfigFile(".env")
	v.AddConfigPath(".")

	v.AutomaticEnv()
	v.ReadInConfig()
	v.Unmarshal(&config)

	return config
}

var cfg = load()

func Cfg() *Config {
	return &cfg
}

func (c *Config) String() string {
	return fmt.Sprintf("Host: %s User: %s Password: %s DbName: %s", c.PgHost, c.PgUser, c.PgPassword, c.PgName)
}
