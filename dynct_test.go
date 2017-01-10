package main

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestReadsAtPaceFromBuffer(t *testing.T) {
	incomingChannel := make(chan *WriteEntry)
	writingChannel := make(chan *WriteEntry, 20)
	for i := 0; i < 20; i++ {
		writingChannel <- &WriteEntry{
			Json:"",
			WriteUnitesConsumed:1,
		}
	}
	startTime := time.Now()
	go start(2, incomingChannel, writingChannel)
	for i := 0; i < 20; i++ {
		<-incomingChannel
	}
	elapsed := time.Since(startTime)
	assert.True(t, elapsed > 10 * time.Second, "The elapsed time should be enough for writes per seconds", elapsed)
}