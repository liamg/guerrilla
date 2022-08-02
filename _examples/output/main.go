package main

import (
	"os"
	"time"

	"github.com/liamg/guerrilla/pkg/guerrilla"

	"github.com/liamg/guerrilla/internal/app/output"
)

func main() {
	printer := output.New(os.Stdout)

	printer.PrintSummary("abcdefghijk@guerrillamailblock.com")

	time.Sleep(time.Second * 3)

	printer.PrintEmail(guerrilla.Email{
		EmailSummary: guerrilla.EmailSummary{
			From:      "welcome@example.service",
			ID:        "1",
			Subject:   "Verify your email!",
			Timestamp: time.Now(),
		},
		Body: `
Hi loser,

Please verify your email address by clicking the link below:

<a href="https://example.service/verify?token=abc" target="_blank">Verify Email</a>

Thanks,
An Example Service

`,
	})

	<-(make(chan bool))
}
