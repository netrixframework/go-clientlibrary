# go-replicaclient

[![API Reference](https://img.shields.io/badge/api-reference-blue.svg)](https://netrixframework.github.io/docs/framework/api.html)

Netrix client to be run on each replica. Uses the Netrix API to send/receive messages. Refer [here](https://netrixframework.github.io/) for Netrix documentation.

## Installation

Fetch and install library using `go get`,

```sh
go get github.com/netrixframework/go-replicaclient
```

## Usage

Define a directive handler to handle `Directive` messages from Netrix. For example,

```go
type Replica struct {
    ...
}

func (n *Replica) Start() error {}

func (n *Replica) Stop() error {}

func (n *Replica) Restart() error {}
```

To initialize the client, pass the ID, directive handler and optionally a logger that the client can use to log debug messages.

```go
import (
    "github.com/netrixframework/go-replicaclient"
    "github.com/netrixframework/netrix/types"
)

func main() {

    replica := NewReplica()
    ...
    err := replicaclient.Init(&replicaclient.Config{
        ReplicaID: types.ReplicaID(Replica.ID),
        NetrixAddr: "<netrix_addr>",
        ClientServerAddr: "localhost:7074",
    }, replica, logger)
    if err != nil {
        ...
    }
}
```

The library maintains a single instance of `ReplicaClient` which can be used to send/receive messages or publish events

```go
client, _ := replicaclient.GetClient()
client.SendMessage(type, to, data, true)
...
message, ok := client.ReceiveMessage()
...

client.PublishEvent("Commit", map[string]string{
    "param1": "val1",
})

```

## Configuration

`Init` accepts a `Config` instance as the first argument which is defined as,

```go
type Config struct {
    ReplicaID        types.ReplicaID
    NetrixAddr       string
    ClientServerAddr string
    ClientAdvAddr    string
    Info             map[string]interface{}
}
```
