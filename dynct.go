package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
)

func main() {
	writeThroughput := envToInt("WRITE_THROUGHPUT")
	defaultBuffers := writeThroughput * 20
	configuration := &Configuration{
		BufferSize:envToIntWithDefault("BUFFER_SIZE", defaultBuffers),
		QueueUrl:os.Getenv("QUEUE_URL"),
		TableName:os.Getenv("TABLE_NAME"),
	}
	fmt.Println("Starting dynct with configuration", configuration)
	puller, err := NewPuller(configuration)
	writer := NewDynamo(configuration)
	if (err != nil) {
		fmt.Println("failed to create puller,", err)
		panic(err)
	}
	writingChannel := make(chan *WriteEntry, defaultBuffers)
	deletingChannel := make(chan *WriteEntry, defaultBuffers)
	go writer.pipeThrough(writingChannel, deletingChannel)
	start(writeThroughput, writingChannel, puller.messagesBuffer)
}

func start(writeThroughput int, writingChannel chan *WriteEntry, incomingChannel chan *WriteEntry) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			transferWrites(writeThroughput, writingChannel, incomingChannel)
		}
	}
}

func transferWrites(throughput int, writingChannel chan *WriteEntry, messages chan *WriteEntry) {
	for i := 0; i < throughput; {
		message := <-messages
		writingChannel <- message
		i += message.WriteUnitesConsumed
	}
}
func envToInt(envName string) int {
	result, err := strconv.Atoi(os.Getenv(envName))
	if (err != nil) {
		panic("failed to parse value for config " + envName)
	}
	return result
}

func envToIntWithDefault(envName string, defaultValue int) int {
	propertyValue := os.Getenv(envName)
	if (propertyValue == "") {
		return defaultValue
	}
	result, err := strconv.Atoi(propertyValue)
	if (err != nil) {
		fmt.Println("failed to parse value for config " + envName + ", using default value")
	}
	return result
}

type Configuration struct {
	TableName  string
	BufferSize int //number of messages to keep in buffer, does not guarantee is write units size. Should change in the future to set write units size
	QueueUrl   string
}