package main

import (
	"encoding/json"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	n := maelstrom.NewNode()
	ids := []int{}

	n.Handle("topology", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": ids,
		})
	})

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		id := int(body["message"].(float64))

		// If id already exists in this node then don't gossip.
		if contains(ids, id) {
			return nil
		}

		ids = append(ids, id)

		log.Println("broadcasting", msg.Src, n.ID(), id, n.NodeIDs())
		for _, nodeId := range n.NodeIDs() {
			if nodeId == msg.Src || nodeId == n.ID() {
				continue
			}
			nodeId := nodeId
			go func() {
				retry(100, 100*time.Millisecond, func() error {
					return n.Send(nodeId, body)
				})
			}()
		}

		return n.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func contains(arr []int, value int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error:", err)
			time.Sleep(sleep)
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}
