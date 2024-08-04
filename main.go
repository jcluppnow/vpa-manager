package main

import (
	"vpa-manager/controller"
)

func main() {
	env := controller.LoadEnv()
	controller.CreateInformers(env)
	select {}
}
