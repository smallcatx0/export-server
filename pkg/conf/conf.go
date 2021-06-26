package conf

import "github.com/spf13/viper"

var AppConf *viper.Viper

func InitAppConf(filePath *string) error {
	AppConf = viper.New()
	AppConf.SetConfigFile(*filePath)
	// 设置默认
	AppConf.SetDefault("base.env", "dev")
	AppConf.SetDefault("base.debug", true)
	AppConf.SetDefault("base.http_port", "80")
	AppConf.Set("flag_param.c", *filePath)

	err := AppConf.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func Env() string {
	return AppConf.GetString("base.env")
}

func IsDebug() bool {
	return AppConf.GetBool("base.debug")
}

func HttpPort() string {
	return AppConf.GetString("base.http_port")
}
