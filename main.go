package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	promql := flag.String("promql", "", "The prometheus query language to be ensure the service_id label")
	// component_id and service_id are the same. service_id is the old name, and component_id is the new name
	serviceID := flag.String("component_id", "", "The value of label service_id")

	flag.Parse()

	if *promql == "" {
		log.Fatal("The arg 'promql' is required")
	}
	if *serviceID == "" {
		log.Fatal("The arg 'component_id' is required")
	}

	newPromQL, err := ensureServiceID(*promql, *serviceID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(newPromQL)
}
