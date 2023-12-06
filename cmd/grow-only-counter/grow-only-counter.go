package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)

	n.Handle("add", func(msg maelstrom.Message) error {

		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		delta := int(body["delta"].(float64))

		readCtx, readCancel := context.WithTimeout(context.Background(), time.Second)
		defer readCancel()

		sum, err := kv.ReadInt(readCtx, n.ID())
		if err != nil {
			sum = 0
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err = kv.Write(ctx, n.ID(), sum+delta)
		if err != nil {
			return err
		}

		return n.Reply(msg, map[string]any{
			"type": "add_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		sum := 0

		for _, nodeID := range n.NodeIDs() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			value, err := kv.ReadInt(ctx, nodeID)
			if err != nil {
				sum += value
			}
		}

		return n.Reply(msg, map[string]any{
			"type":  "read_ok",
			"value": sum,
		})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
