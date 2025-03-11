package config

import (
	"os"
	"fmt"
	"gopkg.in/yaml.v3"
	"urlshozim/database"
)

type Config struct{
	connection.DBconfig	`yaml:",inline"`
}

func Configure() Config {
	f, err := os.Open("local.yaml")
	if err!=nil{
		fmt.Errorf("%v", err)
	}
	defer f.Close()
	var cfg Config 
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)   //I would definately revisit this place
	if err!=nil {
		fmt.Println("config files deconfigation naaaaaaa")
		//potom chtoto pridumaem
	}
	return cfg
}