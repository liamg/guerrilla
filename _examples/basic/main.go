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

	emails, err := client.GetAllEmails()
	if err != nil {
		panic(err)
	}

	for _, email := range emails {
		fmt.Println(email.MailSubject)
		mail, err := client.GetEmail(email.MailID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", mail)
	}

}
