package main

import (
	"log"

	"git.cm/nb/domain-panel"
	"git.cm/nb/domain-panel/controller"
)

func init() {
	log.Println("load config", panel.CF)
	panel.DB.AutoMigrate(panel.User{}, panel.Panel{})
}

func main() {
	controller.Web()
	select {}
}
