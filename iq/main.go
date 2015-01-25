package main

import "github.com/synapse-garden/iq/web"

func main() {
	iqr := web.CreateRunner(25000, nil)
	iqr.StartRun()
	println("Runner running on port 25000")
	err := <-iqr.Errors()
	if err != nil {
		panic(err)
	}
}
