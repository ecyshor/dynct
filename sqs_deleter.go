package main

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (deleter *SqsClient) HandleDeletes(deleteChannel chan WriteEntry) {
	for {
		entryToDelete := <-deleteChannel
		log.Debug("Deleting sqs entry")
		deleter.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl: &deleter.configuration.QueueUrl,
			ReceiptHandle: &entryToDelete.ReceiptHandle,
		})
	}
}