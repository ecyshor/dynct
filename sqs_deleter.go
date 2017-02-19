package main

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (deleter *SqsClient) HandleDeletes(deleteChannel chan *WriteEntry) {
	for {
		entryToDelete := <-deleteChannel
		deleter.deleteItem(entryToDelete)
	}
}
func (client *SqsClient) deleteItem(entry *WriteEntry) {
	log.Debug("Deleting sqs entry")
	_, err := client.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl: &client.configuration.QueueUrl,
		ReceiptHandle: &entry.ReceiptHandle,
	})
	if (err != nil) {
		log.Errorf("Could not delete message from sqs", err)
	}
}