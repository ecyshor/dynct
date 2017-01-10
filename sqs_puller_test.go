package main

import (
	"testing"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"time"
)

type mockedReceiveMsgs struct {
	sqsiface.SQSAPI
	Resp sqs.ReceiveMessageOutput
}

func (m mockedReceiveMsgs) ReceiveMessage(in *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	// Only need to return mocked response output
	return &m.Resp, nil
}

func TestFillsBufferToExpectedCapacityOnly(t *testing.T) {
	puller := &SqsClient{
		sqsClient:mockedReceiveMsgs{Resp: sqs.ReceiveMessageOutput{
			Messages: []*sqs.Message{
				{Body: aws.String(`{"from":"user_1","to":"room_1","msg":"Hello!"}`), ReceiptHandle:aws.String(``)},
				{Body: aws.String(`{"from":"user_2","to":"room_1","msg":"Hi user_1 :)"}`), ReceiptHandle:aws.String(``)},
			},
		}, },
		configuration:&Configuration{
			BufferSize:20,
		},
	}
	puller.start()
	buffer := puller.messagesBuffer
	time.Sleep(2 * time.Second)
	assert.Equal(t, 20, len(buffer))
}

func TestAssignsCorrectWriteCapacityForMessages(t *testing.T) {
	cases := []struct {
		Body               string
		ExpectedWriteUnits int
	}{
		{
			Body:string(make([]byte, 128)),
			ExpectedWriteUnits:1,
		},
		{
			Body:string(make([]byte, 1024)),
			ExpectedWriteUnits:1,
		},
		{
			Body:string(make([]byte, 1600)),
			ExpectedWriteUnits:2,
		},
		{
			Body:string(make([]byte, 5100)),
			ExpectedWriteUnits:5,
		},
	}
	for _, testCase := range cases {
		puller := &SqsClient{
			sqsClient:mockedReceiveMsgs{Resp: sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					{Body: aws.String(testCase.Body), ReceiptHandle:aws.String(``)},
				},
			},
			},
			configuration:&Configuration{
				BufferSize:1,
			},
		}
		puller.start()
		writeEntry := <-puller.messagesBuffer
		assert.Equal(t, testCase.ExpectedWriteUnits, writeEntry.WriteUnitesConsumed)
	}
}