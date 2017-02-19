package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type mockedDeleteMsgs struct {
	sqsiface.SQSAPI
}

func (m mockedDeleteMsgs) DeleteMessage(in *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	// Only need to return mocked response output
	return nil, nil
}

func TestSqsClient_HandleDeletes(t *testing.T) {
	client := &SqsClient{
		sqsClient:mockedDeleteMsgs{},
		configuration:&Configuration{
			BufferSize:20,
		},
	}
	deletingChannel := make(chan *WriteEntry, 20)
	go client.HandleDeletes(deletingChannel)
	deletingChannel <- &WriteEntry{}
	deletingChannel <- &WriteEntry{}
	time.Sleep(200 * time.Millisecond)
	assert.Empty(t, deletingChannel)
}
