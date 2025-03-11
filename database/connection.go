package connection


type DBconfig struct {
	User string		`yaml:"User"`
	DBname string	`yaml:"DBname"`
	Port int		`yaml:"Port"`
	Password string	`yaml:"Password"`
	Host string		`yaml:"Host"`
}

