package model

//Config 全局配置
type Config struct {
	Debug bool
	Web   struct {
		Addr   string
		Domain string
	}
	Mail struct {
		SMTP string `mapstructure:"smtp"`
		Port int
		User string
		Pass string
		SSL  bool `mapstructure:"ssl"`
	}
	Database struct {
		User   string
		Pass   string
		Server string
		Name   string
		Loc    string
	}
	ReCaptcha string `mapstructure:"recaptcha"`
}
