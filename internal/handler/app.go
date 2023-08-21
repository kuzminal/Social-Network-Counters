package handler

import (
	"SocialNetCounters/internal/service"
	"SocialNetCounters/internal/store"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Instance struct {
	sessionStore       store.SessionStore
	countersWriteStore store.CountersStore
	countersReadStore  *store.ReadNodes[store.CountersStore]
	tokenService       *service.Client
	tracer             *trace.TracerProvider
}

func NewInstance(
	sessionStore store.SessionStore,
	countersWriteStore store.CountersStore,
	countersReadStore *store.ReadNodes[store.CountersStore],
	tokenService *service.Client,
	tracer *trace.TracerProvider,
) *Instance {
	return &Instance{
		sessionStore:       sessionStore,
		countersWriteStore: countersWriteStore,
		countersReadStore:  countersReadStore,
		tokenService:       tokenService,
		tracer:             tracer,
	}
}
