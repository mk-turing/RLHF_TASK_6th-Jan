package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:12345")
	if err != nil {
		log.Fatalf("dialing: %s", err)
	}
	defer client.Close()

	// Initialize a client-side version
	var clientVersion int64

	// Perform some operations
	var value string
	if err := client.Call("KV.Get", "key1", &value); err != nil {
		log.Fatalf("get: %s", err)
	}
	fmt.Println("Initial Value:", value)

	// Check if the version is stale before writing
	serverVersion := getServerVersion(client)
	if serverVersion > clientVersion {
		log.Println("Data is stale. Reloading...")
		// Reload data from the server
		// ...
		clientVersion = serverVersion
	}

	if err := client.Call("KV.Set", "key1", "new_value"); err != nil {
		log.Fatalf("set: %s", err)
	}

	// Check again for stale data before reading
	serverVersion = getServerVersion(client)
	if serverVersion > clientVersion {
		log.Println("Data is stale. Reloading...")
		// Reload data from the server
		// ...
		clientVersion = serverVersion
	}

	if err := client.Call("KV.Get", "key1", &value); err != nil {
		log.Fatalf("get: %s", err)
	}
	fmt.Println("Updated Value:", value)
}

func getServerVersion(client *rpc.Client) int64 {
	var version int64
	if err := client.Call("KV.GetVersion", new(int64), &version); err != nil {
		log.Fatalf("get version: %s", err)
	}
	return version
}
