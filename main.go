package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	promql := os.Getenv("PROMQL")
	// component_id and service_id are the same. service_id is the old name, and component_id is the new name
	pod := flag.String("pod", "", "The value of label pod")

	flag.Parse()

	if promql == "" {
		log.Fatal("The arg 'promql' is required")
	}
	if *pod == "" {
		log.Fatal("The arg 'pod' is required")
	}

	newPromQL, err := ensureServiceID(promql, *pod)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(newPromQL)
}
