package main

import (
	"context"
	"errors"

	"go.nunchi.studio/helix/integration/nats"
	"go.nunchi.studio/helix/service"
	"go.nunchi.studio/helix/telemetry/trace"

	"github.com/nats-io/nats.go/jetstream"
)

/*
App holds the different components needed to run our Go service. In this
case, it only holds a NATS JetStream context.
*/
type App struct {
	JetStream nats.JetStream
}

/*
app is the instance of App currently running.
*/
var app *App

/*
NewAndStart creates a new helix service and starts it.
*/
func NewAndStart(ctx context.Context) error {

	// First, create a new NATS JetStream context. We keep empty config but feel
	// free to dive more later for advanced configuration.
	js, err := nats.Connect(nats.Config{})
	if err != nil {
		return err
	}

	// Build app with the NATS JetStream context created.
	app = &App{
		JetStream: js,
	}

	// Create a new stream in NATS JetStream called "demo-stream", for subject "demo".
	stream, _ := js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:     "demo-stream",
		Subjects: []string{"demo"},
	})

	// Create a new NATS JetStream consumer called "demo-queue".
	consumer, _ := stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Name: "demo-queue",
	})

	// Create a new, empty context.
	ctxDetached := context.Background()

	// Start consuming messages from the queue "demo-queue" on subject "demo". We
	// pass the empty context previously created. The context in the callback
	// function is a copy of one the passed, but now contains the Event object at
	// the origin of the trace (if any). You can also create your own span, which
	// will be a child span of the trace found in the context (if any). In our case,
	// the context includes Event created during the HTTP request, as well as the
	// trace handled by the REST router. At any point in time, you can record an
	// error in the span, which will be reported back to the root span.
	consumer.Consume(ctxDetached, func(ctx context.Context, msg jetstream.Msg) {
		_, span := trace.Start(ctx, trace.SpanKindConsumer, "Custom Span")
		defer span.End()

		if 2+2 == 4 {
			span.RecordError("this is a demo error based on a dummy condition", errors.New("any error"))
		}

		msg.Ack()
	})

	// Start the service using the helix's service package. Only one helix service
	// must be running per process. This is a blocking operation.
	err = service.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

/*
Close tries to gracefully close the helix service. This will automatically close
all connections of each integration when applicable. You can add other logic as
well here.
*/
func (app *App) Close(ctx context.Context) error {
	err := service.Close(ctx)
	if err != nil {
		return err
	}

	return nil
}
