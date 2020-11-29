package main

import (
	"flag"
	"github.com/spf13/viper"
	"github.com/tb/image-catalog/common-packages/conf"

	"github.com/tb/image-catalog/common-packages/system"
	"github.com/tb/image-catalog/imageapp/routes"
	"github.com/zenazn/goji"
)

func main() {
	//Load configuration file
	conf.LoadConfigFile()
	var application = &system.Application{}


	//Apply authentication filter
	goji.Use(application.ApplyAuth)
	//Prepare routes
	routes.PrepareRoutes(application)
	//Setting server address
	flag.Set("bind", viper.GetString("apps.imageapp.address"))
	goji.Serve()
}
