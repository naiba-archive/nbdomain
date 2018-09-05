package main

import (
	"git.cm/nb/domain-panel"
	"git.cm/nb/domain-panel/controller"
)

func init() {
	panel.DB.AutoMigrate(panel.User{}, panel.Panel{}, panel.Cat{}, panel.Domain{})
}

func main() {
	controller.Web()
	select {}
}
