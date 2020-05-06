package main

import (
	"github.com/go-mesh/openlogging"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/ygx/schemas/example/metrics"

	_ "github.com/go-chassis/go-chassis/plugins/tracing/zipkin"
)

func main() {
	chassis.RegisterSchema("rest", &metrics.User{}, server.WithSchemaID("MetricsServer")) //server
	err := chassis.Init()
	if err != nil {
		openlogging.GetLogger().Errorf("chassis init failed ,err :%v", err.Error())
	}

	chassis.Run()
}
