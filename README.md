# :incoming_envelope::bust_in_silhouette: guerrilla

A command-line tool (and Go module) for [https://www.guerrillamail.com/](https://www.guerrillamail.com/).

Create a temporary email address in the terminal to quickly sign up for services and receive verification emails.

![Screenshot of Guerrilla command-line receiving emails](demo.gif)

Built based on the [official API documention](https://docs.google.com/document/d/1Qw5KQP1j57BPTDmms5nspe-QAjNEsNg8cQHpAAycYNM/edit?hl=en).

## Usage: CLI

Install with Go:

```
go install github.com/liamg/guerrilla/cmd/guerrilla@latest
```

No configuration or authentication required, just run:

```
$ guerrilla
```

...and start receiving emails!

## Usage: Go Module

```go
package main

import (
    "fmt"
    "github.com/liamg/guerrilla/pkg/guerrilla"
)

func main() {
    
    client, _ := guerrilla.Init()
    poller := guerrilla.NewPoller(client)

    for email := range poller.Poll() {
        fmt.Printf("Email received: Subject=%s\nBody=%s\n\n", email.Subject, email.Body)
    }
}
```

## TODO

- Link replacement
- Social image
