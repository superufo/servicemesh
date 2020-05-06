package main

import (
	//tracers
	//_ "github.com/go-chassis/go-chassis-plugins/tracing/zipkin"/
	//_ "github.com/go-chassis/go-chassis-plugins/tracing/zipkin"
	"github.com/go-chassis/ygx/schemas/example"

	"github.com/go-mesh/openlogging"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/server"

	_ "github.com/go-chassis/go-chassis/plugins/tracing/zipkin"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/serverB/

func main() {
	chassis.RegisterSchema("rest", &example.RestFulHello{}, server.WithSchemaID("RestHelloService"))
	if err := chassis.Init(); err != nil {
		openlogging.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
