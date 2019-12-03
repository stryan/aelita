package main

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	PROTOV = "0.1"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/aelita/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.SetConfigType("yaml")
	port := viper.GetString("port")
	pluginsDir := viper.GetString("plugins")

	ael := NewController()
	RegisterInternal(ael)
	RegisterExternal(ael,pluginsDir)
	ael.LogCommands()
	StartServer(":"+port, ael)
}

func RegisterExternal(ael *Controller,dir string) {
	ael.AddCommand(parseYAMLCommand(dir+"/fortune.yaml"))
}
