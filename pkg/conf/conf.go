package conf

import "github.com/spf13/viper"

var AppConf *viper.Viper

func InitAppConf(filePath *string) error {
	AppConf = viper.New()
	AppConf.SetConfigFile(*filePath)
	AppConf.SetConfigType("yaml")
	// 设置默认
	AppConf.SetDefault("env", "dev")
	AppConf.SetDefault("debug", true)
	AppConf.SetDefault("http_port", "80")
	AppConf.Set("flag_param.c", *filePath)

	err := AppConf.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func Env() string {
	return AppConf.GetString("env")
}

func IsDebug() bool {
	return AppConf.GetBool("debug")
}

func HttpPort() string {
	return AppConf.GetString("http_port")
}
