# dynct
Ensure constant write throughput to dynamodb using sqs so you can insure the provisioned throughput is not exceeded 
## How?
Every item which need to be written to the table is instead being pushed to a sqs queue, from which dynct will read in an internal buffer and consume the buffer at the desired constant rate.
###In depth
##Requirements
##Configuration
###Dynct
####Environment variables:
|Variable Name|Description|Required|
|-------------|-----------------------------|----|
|WRITE_THROUGHPUT|The value of throughput which the system must try to respect in writing.|true|
|BUFFER_SIZE|The number of messages to keep in the sqs pullers buffer. This is updated async.Defaults to WRITE_THROUGHPUT * 20|false|
|QUEUE_URL|The url of the queue where the write messages are pushed|false|
|TABLE_NAME|The dynamodb table where the messages are being written|false|
###AWS
The default aws configuration system is used. The required configuration are:
- AWS_REGION

You can also set it to load the config file from `~/.aws/config` by setting the environment variable:
- AWS_SDK_LOAD_CONFIG to `true`