define build
	mkdir -p bin && cd ./bin; \
	go build ../cmd/$1/$1.go;
endef

serve:
	./third_party/maelstrom/maelstrom serve

# 1
test-echo:
	$(call build,echo)
	./third_party/maelstrom/maelstrom test -w echo --bin ./bin/echo --node-count 1 --time-limit 10

# 2
test-unique-ids:
	$(call build,unique-ids)
	./third_party/maelstrom/maelstrom test -w unique-ids --bin ./bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

# 3-a
test-broadcast-single-node:
	$(call build,broadcast)
	./third_party/maelstrom/maelstrom test -w broadcast --bin ./bin/broadcast --node-count 1 --time-limit 20 --rate 10

# 3-b
test-broadcast-multi-node:
	$(call build,broadcast)
	./third_party/maelstrom/maelstrom test -w broadcast --bin ./bin/broadcast --node-count 5 --time-limit 20 --rate 10

# 3-c
test-broadcast-multi-node-fault-tolerant:
	$(call build,broadcast)
	./third_party/maelstrom/maelstrom test -w broadcast --bin ./bin/broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition

# 3-d
test-broadcast-efficiency:
	$(call build,broadcast)
	./third_party/maelstrom/maelstrom test -w broadcast --bin ./bin/broadcast --node-count 25 --time-limit 20 --rate 100 --latency 100

# 4
test-grow-only-counter:
	$(call build,grow-only-counter)
	./third_party/maelstrom/maelstrom test -w g-counter --bin ./bin/grow-only-counter --node-count 3 --rate 100 --time-limit 20

# 5-a
test-kafka-log-single-node:
	$(call build,kafka-log)
	./third_party/maelstrom/maelstrom test -w kafka --bin ./bin/kafka-log --node-count 1 --concurrency 2n --time-limit 20 --rate 1000

# 5-b

# 6-a
test-txn-single-node:
	$(call build,txn-rw-register)
	./third_party/maelstrom/maelstrom test -w txn-rw-register --bin ./bin/txn-rw-register --node-count 1 --time-limit 20 --rate 1000 --concurrency 2n --consistency-models read-uncommitted --availability total

# 6-b