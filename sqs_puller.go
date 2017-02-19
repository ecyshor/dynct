package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"time"
	"math"
)

type SqsClient struct {
	sqsClient      sqsiface.SQSAPI
	configuration  *Configuration
	messagesBuffer chan *WriteEntry
}
type WriteEntry struct {
	Json                string
	ReceiptHandle       string
	WriteUnitesConsumed int
}

func NewSqsClient(configuration *Configuration) (*SqsClient, error) {
	log.Infof("Starting puller with configuration", configuration)
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("failed to create session,", err)
		return nil, err
	}
	puller := &SqsClient{
		sqsClient: sqs.New(sess),
		configuration:configuration,
	}
	puller.start()
	return puller, nil
}

func (puller *SqsClient) start() {
	ticker := time.NewTicker(100 * time.Millisecond)
	messagesBuffer := make(chan *WriteEntry, puller.configuration.BufferSize + 10)
	log.Infof("Starting sqs pulling")
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Debug("Running update on ticker config")
				if len(messagesBuffer) < puller.configuration.BufferSize {
					//TODO ensure puller does just one long pulling call, but when buffer is not full and queue has entries then do on multiple routines
					log.Debug("Running retrieve from sqs using long pulling.")
					puller.retrieve(messagesBuffer)
				}
			}
		}
	}()
	puller.messagesBuffer = messagesBuffer
}

func (puller *SqsClient) retrieve(output chan *WriteEntry) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(puller.configuration.QueueUrl), // Required
		AttributeNames: []*string{
			aws.String("QueueAttributeName"), // Required
		},
		MaxNumberOfMessages: aws.Int64(10),
		MessageAttributeNames: []*string{
			aws.String("MessageAttributeName"), // Required
		},
		WaitTimeSeconds:         aws.Int64(20),
	}
	resp, err := puller.sqsClient.ReceiveMessage(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Fatalf("Error while calling receive message from sqs", err)
		return
	}
	log.Debugf("Send %s write entries from sqs", len(resp.Messages))
	for _, message := range resp.Messages {
		output <- &WriteEntry{
			Json:*message.Body,
			WriteUnitesConsumed:int(math.Ceil(float64(len(*message.Body)) / 1024)),
			ReceiptHandle:*message.ReceiptHandle,
		}
	}
}