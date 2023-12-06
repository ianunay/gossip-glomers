package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type slot struct {
	offset int
	msg    int
}

type instance struct {
	n                *maelstrom.Node
	log              map[string][]slot
	committedOffsets map[string]int
}

func main() {

	n := maelstrom.NewNode()
	i := &instance{n: n}

	n.Handle("send", i.sendHandler)
	n.Handle("poll", i.pollHandler)
	n.Handle("commit_offsets", i.commitOffsetsHandler)
	n.Handle("list_committed_offsets", i.listCommitOffsetsHandler)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (i *instance) sendHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	key := body["key"].(string)
	message := int(body["message"].(float64))

	// get last offset for the key in the log
	lastOffset := i.log[key][len(i.log[key])-1].offset

	slot := slot{
		offset: lastOffset + 1,
		msg:    message,
	}

	i.log[key] = append(i.log[key], slot)

	return i.n.Reply(msg, map[string]any{
		"type":   "add_ok",
		"offset": slot.offset,
	})
}

func (i *instance) pollHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	messages := map[string][][2]int{}

	offsets := body["offsets"].(map[string]int)
	for key, offset := range offsets {
		for _, slot := range i.log[key] {
			if slot.offset >= offset {
				messages[key] = append(messages[key], [2]int{slot.offset, slot.msg})
			}
		}
	}

	return i.n.Reply(msg, map[string]any{
		"type": "poll_ok",
		"msgs": messages,
	})
}

func (i *instance) commitOffsetsHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	offsets := body["offsets"].(map[string]int)

	for key, offset := range offsets {
		i.committedOffsets[key] = offset
	}

	return i.n.Reply(msg, map[string]any{
		"type": "commit_offsets_ok",
	})
}

func (i *instance) listCommitOffsetsHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	keys := body["keys"].([]string)
	offsets := map[string]int{}

	for _, v := range keys {
		offsets[v] = i.committedOffsets[v]
	}

	return i.n.Reply(msg, map[string]any{
		"type":    "list_committed_offsets_ok",
		"offsets": offsets,
	})
}
