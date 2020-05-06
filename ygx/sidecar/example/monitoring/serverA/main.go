package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/ygx/schemas/example"

	//_ "github.com/go-chassis/go-chassis-plugins/tracing/zipkin"
	//tracers
	_ "github.com/go-chassis/go-chassis/plugins/tracing/zipkin"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/serverA/

func main() {
	chassis.RegisterSchema("rest", &example.TracingHello{}, server.WithSchemaID("TracingHello"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
