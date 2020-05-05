build:
	go build producer.go
	go build consumer.go
clean:
	rm -rf producer
	rm -rf consumer
runp:
	./producer
runc:
	./consumer
