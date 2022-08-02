package main

import (
	"fmt"

	"github.com/liamg/guerrilla/pkg/guerrilla"
)

func main() {
	client, err := guerrilla.Init()
	if err != nil {
		panic(err)
	}

	summaries, err := client.GetAllEmails()
	if err != nil {
		panic(err)
	}

	for _, summary := range summaries {
		email, err := client.GetEmail(summary.ID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Email received: %#v\n", email)
	}

}
