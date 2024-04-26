package config

// 配置文件模板
type ConfigTemplate struct {
	Port         string
	MailSrvIp    string
	MailSrvPort  string
	MailAccount  string
	MailPassword string
	Filter       [][2]string
}

// 配置文件变量
var Config = ConfigTemplate{
	Port:         "5678",
	MailSrvIp:    "127.0.0.1",
	MailSrvPort:  "993",
	MailAccount:  "",
	MailPassword: "",
	Filter: [][2]string{
		{"您的", "临时访问代码"},
		{"如何更新", "同户设备"},
	},
}
