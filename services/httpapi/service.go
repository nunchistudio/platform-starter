package main

import (
	"context"
	"net/http"

	"go.nunchi.studio/helix/event"
	"go.nunchi.studio/helix/integration/nats"
	"go.nunchi.studio/helix/integration/rest"
	"go.nunchi.studio/helix/service"
)

/*
App holds the different components needed to run our Go service. In this
case, it holds a REST router and NATS JetStream context.
*/
type App struct {
	REST      rest.REST
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

	// First, create a new REST router. We keep empty config but feel free to
	// dive more later for configuring OpenAPI behavior.
	router, err := rest.New(rest.Config{})
	if err != nil {
		return err
	}

	// Then, create a new NATS JetStream context. We keep empty config but feel
	// free to dive more later for advanced configuration.
	js, err := nats.Connect(nats.Config{})
	if err != nil {
		return err
	}

	// Build app with the router created.
	app = &App{
		REST:      router,
		JetStream: js,
	}

	// Add a route, returning a 202 HTTP response. When a request is made, publish
	// a message. We use static data for demo purpose. In the demo below, we create
	// an Event object using the event package. We then create a new context.Context
	// by calling event.ContextWithEvent. This returns a new context including the
	// event created. helix integrations automatically read/write an Event from/into
	// a context when possible. The integration then passes the Event in the
	// appropriate headers. In this case, the NATS JetStream integration achieves
	// this by passing and reading an Event from the messages' headers.
	router.POST("/anything", func(rw http.ResponseWriter, req *http.Request) {
		var e = event.Event{
			Name:   "post.anything",
			UserID: "7469e788-617a-4b6a-8a26-a61f6acd01d3",
			Subscriptions: []event.Subscription{
				{
					CustomerID:  "2658da04-7c8f-4a7e-9ab0-d5d555b8173e",
					PlanID:      "7781028b-eb48-410d-8cae-c36cffed663d",
					Usage:       "api.requests",
					IncrementBy: 1.0,
				},
			},
		}

		ctx := event.ContextWithEvent(req.Context(), e)
		msg := &nats.Msg{
			Subject: "demo",
			Sub: &nats.Subscription{
				Queue: "demo-queue",
			},
			Data: []byte(`{ "hello": "world" }`),
		}

		js.Publish(ctx, msg)

		rest.WriteAccepted(rw, req)
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
