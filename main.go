package main

import (
	"path/filepath"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	PROTOV = "0.2"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/aelita/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %v \n", err)
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
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
	        files = append(files, path)
		return nil
	})
	if err != nil {
		log.Print("Failed to add external commands")
		return
	}
	for _, file := range files {
		if strings.Contains(file, "yaml") {
			ael.AddCommand(parseYAMLCommand(file))
		}
	}
}
