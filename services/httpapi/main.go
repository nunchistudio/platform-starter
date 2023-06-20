package main

/*
main is the entrypoint of our app.
*/
func main() {

	// Create and start the service.
	err := NewAndStart()
	if err != nil {
		panic(err)
	}

	// Try to close the service when done.
	err = app.Close()
	if err != nil {
		panic(err)
	}
}
