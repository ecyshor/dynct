package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedPutItem struct {
	dynamodbiface.DynamoDBAPI
	Resp  dynamodb.PutItemOutput
	Error error
}

func (m mockedPutItem) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	// Only need to return mocked response output
	return &m.Resp, m.Error
}

func TestDrainsBuffer(t *testing.T) {
	dyn := &DynamoDb{
		Configuration: &Configuration{},
		ApiClient:mockedPutItem{
			Error:errors.New("failed"),
		},
	}
	writingEntries := make(chan *WriteEntry, 20)
	output := make(chan *WriteEntry, 20)
	dyn.pipeThrough(writingEntries, output)
	for i := 0; i < 50; i++ {
		writingEntries <- &WriteEntry{
			Json:"{}",
		}
	}
	time.Sleep(20 * time.Millisecond)
	assert.Empty(t, writingEntries)
}