package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"go/build"
	"log"
	"os"
)

func LoadConfigFile() error {

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	fmt.Println("my GOPATH -- > ", gopath)

	err := os.Setenv("IMAGE_CONF_FILE", gopath+"/src/github.com/tb/task-logger/backend/golang/common-packages/conf")
	if err != nil {
		log.Println(err)
		log.Println("Could not find IMAGE_CONF_FILE enviroment variable, which should point to conf directory path")
		return err
	}

	//Read the configuration file from environment variable and provide application wide access
	filePath := os.Getenv("IMAGE_CONF_FILE")
	log.Println(filePath)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filePath)

	viper.SetConfigName("config")

	err = viper.ReadInConfig()

	if err != nil {
		log.Println(err)
		log.Println("Could not find IMAGE_CONF_FILE enviroment variable, which should point to conf directory path")
		return err
	}
	return nil

}
