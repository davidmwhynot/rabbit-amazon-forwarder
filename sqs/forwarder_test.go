package sqs

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/davidmwhynot/rabbit-amazon-forwarder/config"
	"github.com/davidmwhynot/rabbit-amazon-forwarder/forwarder"
)

var badRequest = "Bad request"

func TestCreateForwarder(t *testing.T) {
	entry := config.AmazonEntry{Type: "SQS",
		Name:   "sqs-test",
		Target: "arn",
	}
	forwarder := CreateForwarder(entry)
	if forwarder.Name() != entry.Name {
		t.Errorf("wrong forwarder name, expected:%s, found: %s", entry.Name, forwarder.Name())
	}
}

func TestPush(t *testing.T) {
	queueName := "queue1"
	entry := config.AmazonEntry{Type: "SQS",
		Name:   "sqs-test",
		Target: queueName,
	}
	scenarios := []struct {
		name    string
		mock    sqsiface.SQSAPI
		message string
		queue   string
		err     error
	}{
		{
			name:    "empty message",
			mock:    mockAmazonSQS{resp: sqs.SendMessageOutput{MessageId: aws.String("messageId")}, queue: queueName, message: ""},
			message: "",
			queue:   queueName,
			err:     errors.New(forwarder.EmptyMessageError),
		},
		{
			name:    "bad request",
			mock:    mockAmazonSQS{resp: sqs.SendMessageOutput{MessageId: aws.String("messageId")}, queue: queueName, message: badRequest},
			message: badRequest,
			queue:   queueName,
			err:     errors.New(badRequest),
		},
		{
			name:    "success",
			mock:    mockAmazonSQS{resp: sqs.SendMessageOutput{MessageId: aws.String("messageId")}, queue: queueName, message: "abc"},
			message: "abc",
			queue:   queueName,
			err:     nil,
		},
	}
	for _, scenario := range scenarios {
		t.Log("Scenario name: ", scenario.name)
		forwarder := CreateForwarder(entry, scenario.mock)
		err := forwarder.Push(scenario.message)
		if scenario.err == nil && err != nil {
			t.Errorf("Error should not occur")
			return
		}
		if scenario.err == err {
			return
		}
		if err.Error() != scenario.err.Error() {
			t.Errorf("Wrong error, expecting:%v, got:%v", scenario.err, err)
		}
	}
}

type mockAmazonSQS struct {
	sqsiface.SQSAPI
	resp    sqs.SendMessageOutput
	queue   string
	message string
}

func (m mockAmazonSQS) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if *input.QueueUrl != m.queue {
		return nil, errors.New("Wrong queue name")
	}
	if *input.MessageBody != m.message {
		return nil, errors.New("Wrong message body")
	}
	if *input.MessageBody == badRequest {
		return nil, errors.New(badRequest)
	}
	return &m.resp, nil
}
