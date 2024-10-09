package main

import (
	"context"
)

/*
main is the entrypoint of our app.
*/
func main() {

	// Create and start the service.
	err := NewAndStart(context.Background())
	if err != nil {
		panic(err)
	}

	// Try to close the service when done.
	err = app.Close(context.Background())
	if err != nil {
		panic(err)
	}
}
