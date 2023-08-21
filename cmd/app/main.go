package main

import (
	"SocialNetCounters/internal/counters"
	"SocialNetCounters/internal/handler"
	"SocialNetCounters/internal/helper"
	"SocialNetCounters/internal/router"
	"SocialNetCounters/internal/service"
	"SocialNetCounters/internal/session"
	"SocialNetCounters/internal/store"
	"SocialNetCounters/internal/store/tarantool"
	"SocialNetCounters/internal/tracing"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
)

var (
	tarantoolMaster    *tarantool.Tarantool
	tarantoolReadNodes store.ReadNodes[store.CountersStore]
)

func main() {
	port := helper.GetEnvValue("PORT", "8080")

	traceServer := helper.GetEnvValue("TRACE_SERVER", "trace")
	tracePort := helper.GetEnvValue("TRACE_PORT", "14268")
	tracer, err := tracing.TracerProvider(fmt.Sprintf("http://%s:%s/api/traces", traceServer, tracePort))
	if err != nil {
		log.Fatal(err)
	}
	defer tracer.Shutdown(context.Background())

	initTarantoolDb()

	tokenService := service.NewTokenServiceClient(tarantoolMaster, tracer)
	app := handler.NewInstance(
		tarantoolMaster,
		tarantoolMaster,
		&tarantoolReadNodes,
		tokenService,
		tracer,
	)
	sessionConsumer := session.NewSessionConsumer(tarantoolMaster)
	go sessionConsumer.ReadSessionInfo(context.Background())

	counterConsumer := counters.NewCountersConsumer(tarantoolMaster)
	go counterConsumer.ReadCountersInfo(context.Background())

	r := router.NewRouter(app)

	log.Fatalln(http.ListenAndServe(":"+port, r))
}

func initTarantoolDb() {
	thost := helper.GetEnvValue("TARANTOOL_HOST", "localhost")
	tport := helper.GetEnvValue("TARANTOOL_PORT", "3301")
	tuser := helper.GetEnvValue("TARANTOOL_USER_NAME", "user")
	tpassword := helper.GetEnvValue("TARANTOOL_USER_PASSWORD", "password")
	tarantoolMaster, _ = tarantool.NewTarantoolMaster(thost, tport, tuser, tpassword)
	//TODO инициализировать настоящий слэйв
	tarantoolSlave, _ := tarantool.NewTarantoolMaster(thost, tport, tuser, tpassword)
	nodes := []store.Backend[store.CountersStore]{
		{
			Id:     uuid.Must(uuid.NewV4()).String(),
			IsDead: false,
			Store:  tarantoolMaster,
		},
		{
			Id:     uuid.Must(uuid.NewV4()).String(),
			IsDead: false,
			Store:  tarantoolSlave,
		},
	}
	tarantoolReadNodes = store.NewReadNode[store.CountersStore]()
	tarantoolReadNodes.Current = 0
	tarantoolReadNodes.Nodes = nodes
}
