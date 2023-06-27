package main

import (
	"capten/pkg/cert"
	"fmt"
)

func main() {
	// err := cert.GenerateCerts()
	err := cert.GenerateCerts("certs", "config/capten.yaml")
	if err != nil {
		fmt.Println("Err - " + err.Error())
	}
}
