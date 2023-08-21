package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
)

func (i *Instance) GetTotalMessages(w http.ResponseWriter, r *http.Request) {
	// Extract TraceID from header
	md, _ := metadata.FromOutgoingContext(r.Context())
	traceIdString := md["x-trace-id"][0]
	// Convert string to byte array
	traceId, err := trace.TraceIDFromHex(traceIdString)
	if err != nil {
		return
	}
	// Creating a span context with a predefined trace-id
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
	})
	// Embedding span config into the context
	ctx := trace.ContextWithSpanContext(r.Context(), spanContext)

	ctx, span := i.tracer.Tracer("counters-service").Start(ctx, "GetTotalMessages")
	defer span.End()

	id := chi.URLParam(r, "user_id")
	userId := r.Context().Value("userId").(string)
	if len(userId) == 0 {
		log.Println("Could not count messages for empty user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	readNode := i.countersReadStore.GetReadNode()
	msg, _ := readNode.GetTotalMessages(id)
	msgDTO, _ := json.Marshal(msg)
	w.Header().Add("Content-Type", "application/json")
	w.Write(msgDTO)
}
