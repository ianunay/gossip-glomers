package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type transaction struct {
	Type  string  `json:"type"`
	MsgID int     `json:"msg_id"`
	Txn   [][]any `json:"txn"`
}

type instance struct {
	n     *maelstrom.Node
	store map[int]int
}

func main() {

	n := maelstrom.NewNode()

	i := &instance{n: n, store: make(map[int]int)}

	n.Handle("txn", i.handleTxn)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (i *instance) handleTxn(msg maelstrom.Message) error {
	var body transaction
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	response := make([][]any, 0, len(body.Txn))

	for _, txn := range body.Txn {
		op := make([]any, len(txn))
		copy(op, txn)

		if txn[0] == "w" {
			i.store[int(txn[1].(float64))] = int(txn[2].(float64))
		} else if txn[0] == "r" {
			op[2] = i.store[int(txn[1].(float64))]
		}

		response = append(response, op)
	}

	return i.n.Reply(msg, map[string]any{
		"type": "txn_ok",
		"txn":  response,
	})
}
