package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"

	"github.com/go-chassis/ygx/schemas/admin"
)

// if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/server/
func main() {
	chassis.RegisterSchema("rest", &admin.RestFulAdminLogin{}, server.WithSchemaID("RestAdminLoginService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
