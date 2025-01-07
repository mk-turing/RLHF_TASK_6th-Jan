package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
)

type KVClient struct {
	client  *rpc.Client
	version int64
	store   map[string]string
}

func NewKVClient(network, address string) (*KVClient, error) {
	client, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &KVClient{client: client, version: 0, store: make(map[string]string)}, nil
}

func (c *KVClient) Synchronize() error {
	var mapVersion int64
	var keyValueMap map[string]string
	err := c.client.Call("KVService.GetMap", c.version, &mapVersion, &keyValueMap)
	if err != nil {
		return err
	}
	if mapVersion > c.version {
		c.version = mapVersion
		c.store = keyValueMap
		log.Println("Synchronized to version", mapVersion)
	}
	return nil
}

func (c *KVClient) Set(key, value string) error {
	newVersion := 0
	for {
		err := c.client.Call("KVService.Set", key, value, c.version, &newVersion)
		if err != nil {
			return err
		}
		if newVersion > c.version {
			c.version = newVersion
			log.Println("Set", key, "to", value, "at version", newVersion)
			break
		}
		// Synchronize and retry if conflict
		err = c.Synchronize()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *KVClient) Get(key string) (string, error) {
	var value string
	var version int64
	for {
		err := c.client.Call("KVService.Get", key, &value, &version)
		if err != nil {
			return "", err
		}
		if version > c.version {
			c.version = version
			log.Println("Synchronized to version", version)
		}
		if value != "" {
			log.Println("Got", key, "=", value, "at version", version)
			return value, nil
		}
		// Synchronize if key not found
		err = c.Synchronize()
		if err != nil {
			return "", err
		}
	}
}

func main() {
	client, err := NewKVClient("tcp", ":12345")
	if err != nil {
		log.Fatalf("failed to connect to RPC server: %v", err)
	}

	go func() {
		for {
			client.Synchronize()
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case command := <-chan string(os.Stdin):
			parts := strings.Split(command, " ")
			switch parts[0] {
			case "set":
				if len(parts) != 3 {
					fmt.Println("Usage: set <key> <value>")
					continue
				}
				err := client.Set(parts[1], parts[2])
				if err != nil {
					fmt.Println("Error setting key:", err)
				}
			case "get":
				if len(parts) != 2 {
					fmt.Println("Usage: get <key>")
					continue
				}
				value, err := client.Get(parts[1])
				if err != nil {
					fmt.Println("Error getting key:", err)
				} else {
					fmt.Println(value)
				}
			default:
				fmt.Println("Unknown command. Use 'set <key> <value>' or 'get <key>'.")
			}
		}
	}
}
