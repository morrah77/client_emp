package main

import (
	"flag"
	"fmt"
	ct "morrah77.com/client_emp/client"
	it "morrah77.com/client_emp/interactor"
	"morrah77.com/client_emp/validators"
	"os"
)

var listUrl, itemUrl string
var maxAttempt int

func init()  {
	flag.StringVar(&listUrl, "list_url", "http://dummy.restapiexample.com/api/v1/employees", "API URL to fetch list, " +
		"like http://dummy.restapiexample.com/api/v1/employees")
	flag.StringVar(&itemUrl, "item_url", "http://dummy.restapiexample.com/api/v1/employee/%v",
		"API URL pattern to fetch item, like http://dummy.restapiexample.com/api/v1/employee/%v")
	flag.IntVar(&maxAttempt, "max_attempt", 0, "Maximum attempts count for fetching data")
}

func main() {
	flag.Parse()
	var interactor *it.Interactor
	client, err := ct.NewClient(listUrl, itemUrl, maxAttempt, validators.IsValidNumberId)
	if err != nil {
		fmt.Printf("Could not create client: %s", err.Error())
		os.Exit(1)
	}
	interactor, err = it.NewInteractor(client, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Printf("Could not create interactor: %s", err.Error())
		os.Exit(1)
	}

	os.Exit(interactor.Run())
}
