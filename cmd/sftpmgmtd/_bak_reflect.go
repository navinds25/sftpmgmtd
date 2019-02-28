package main

import (
	"log"
	"reflect"

	"github.com/navinds25/sftpmgmt/pkg/sftpconfig"
)

// RunTask runs a task
func RunTask() error {
	log.Println("hello")
	return nil
}

func main() {
	config, err := sftpconfig.GetConfig("_config/pull_config.yml")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config)
	for _, struc := range config["pull_config"] {
		struc.Source.Remote.Auth.Username = "testusr"
		struc.Source.Remote.Auth.Password = "mypass"
		fields := reflect.TypeOf(struc.Source.Remote.Auth)
		values := reflect.ValueOf(struc.Source.Remote.Auth)
		fields.NumField()
		log.Println(fields)
		log.Println(values)
		log.Println(fields.Field(0).Name)
		value := values.Field(0)
		if value.String() == "testusr" {
			log.Println("pass")
		}
	}
}
