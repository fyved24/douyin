package configs

// SigningKey  string `mapstructure:"signing-key" json:"signing-key" yaml:"signing-key"`    // jwt签名
import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var (
	Settings ServerConfig
	//RedisConfig RedisConfig
)

type ServerConfig struct {
	Name         string      `mapstructure:"name"`
	Port         int         `mapstructure:"port"`
	MysqlConfigs MysqlConfig `mapstructure:"mysql"`
	RedisConfigs RedisConfig `mapstructure:"redis"`
	LogsAddress  string      `mapstructure:"logsAddress"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"DB"`
}

func InitConfig() {
	// 实例化viper
	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile("./configs-dev.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := ServerConfig{}
	//给serverConfig初始值
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	// 传递给全局变量
	Settings = serverConfig
	color.Blue("11111111", Settings.LogsAddress)
	fmt.Sprintf("%+v", Settings)
}
