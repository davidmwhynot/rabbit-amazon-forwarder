package sns

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/davidmwhynot/rabbit-amazon-forwarder/config"
	"github.com/davidmwhynot/rabbit-amazon-forwarder/forwarder"
)

var badRequest = "Bad request"

func TestCreateForwarder(t *testing.T) {
	entry := config.AmazonEntry{Type: "SNS",
		Name:   "sns-test",
		Target: "arn",
	}
	forwarder := CreateForwarder(entry)
	if forwarder.Name() != entry.Name {
		t.Errorf("wrong forwarder name, expected:%s, found: %s", entry.Name, forwarder.Name())
	}
}

func TestPush(t *testing.T) {
	topicName := "topic1"
	entry := config.AmazonEntry{Type: "SNS",
		Name:   "sns-test",
		Target: topicName,
	}
	scenarios := []struct {
		name    string
		mock    snsiface.SNSAPI
		message string
		topic   string
		err     error
	}{
		{
			name:    "empty message",
			mock:    mockAmazonSNS{resp: sns.PublishOutput{MessageId: aws.String("messageId")}, topic: topicName, message: ""},
			message: "",
			topic:   topicName,
			err:     errors.New(forwarder.EmptyMessageError),
		},
		{
			name:    "bad request",
			mock:    mockAmazonSNS{resp: sns.PublishOutput{MessageId: aws.String("messageId")}, topic: topicName, message: badRequest},
			message: badRequest,
			topic:   topicName,
			err:     errors.New(badRequest),
		},
		{
			name:    "success",
			mock:    mockAmazonSNS{resp: sns.PublishOutput{MessageId: aws.String("messageId")}, topic: topicName, message: "abc"},
			message: "abc",
			topic:   topicName,
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

type mockAmazonSNS struct {
	snsiface.SNSAPI
	resp    sns.PublishOutput
	topic   string
	message string
}

func (m mockAmazonSNS) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	if *input.TargetArn != m.topic {
		return nil, errors.New("Wrong topic name")
	}
	if *input.Message != m.message {
		return nil, errors.New("Wrong message body")
	}
	if *input.Message == badRequest {
		return nil, errors.New(badRequest)
	}
	return &m.resp, nil
}
