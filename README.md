# go-hasocket

Wrapper around the Home Assistant WebSocket API.
Giving you the ability to easily call services and listen to events.

## Getting started

### Installation

```bash
go get github.com/moonen-home-automation/go-hasocket
```

### Instantiate an app instance

```go
package main

import (
	gohasocket "github.com/moonen-home-automation/go-hasocket"
	"log"
)

func main() {
	app, err := gohasocket.NewApp("ws://myhassinstance:8123/api/websocket", "supersecrettoken")
	if err != nil {
		log.Fatal(err)
	}
	
	// Code here
}
```

### Calling services

#### Function without return value

```go
package main

import (
	gohasocket "github.com/moonen-home-automation/go-hasocket"
	"github.com/moonen-home-automation/go-hasocket/pkg/services"
	"log"
)

func main() {
	app, err := gohasocket.NewApp("ws://myhassinstance:8123/api/websocket", "supersecrettoken")
	if err != nil {
		log.Fatal(err)
	}

	sr := services.NewServiceRequest()
	sr.Domain = "light"
	sr.Service = "toggle"
	sr.Target = services.ServiceTarget{
		EntityId: "light.desklamp",
	}
	_, err = app.CallService(sr)
	if err != nil {
		log.Fatal(err)
	}
}
```

#### Function with return value

```go
package main

import (
	"fmt"
	gohasocket "github.com/moonen-home-automation/go-hasocket"
	"github.com/moonen-home-automation/go-hasocket/pkg/services"
	"log"
	"os"
)

func main() {
	app, err := gohasocket.NewApp("ws://myhassinstance:8123/api/websocket", "supersecrettoken")
	if err != nil {
		log.Fatal(err)
	}

	sr := services.NewServiceRequest()
	sr.Domain = "todo"
	sr.Service = "get_items"
	sr.Target = services.ServiceTarget{
		EntityId: "todo.daily",
	}
	sr.ReturnResponse = true
	res, err := app.CallService(sr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Result.Response)
}
```

### Listening for events

```go
package main

import (
	"fmt"
	gohasocket "github.com/moonen-home-automation/go-hasocket"
	"github.com/moonen-home-automation/go-hasocket/pkg/services"
	"log"
	"os"
)

func main() {
	app, err := gohasocket.NewApp("ws://myhassinstance:8123/api/websocket", "supersecrettoken")
	if err != nil {
		log.Fatal(err)
	}

	el, err := app.RegisterListener("test")
	if err != nil {
		log.Fatal(err)
	}

	err = el.Register()
	if err != nil {
		log.Fatal(err)
	}

	elChan := make(chan events.EventData, 10)
	go el.Listen(elChan)
	defer func() {
		err := el.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		msg, ok := <-elChan
		if !ok {
			continue
		}

		fmt.Println(string(msg.RawEventJSON))
	}
}
```