package lambda

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/davidmwhynot/rabbit-amazon-forwarder/config"
	"github.com/davidmwhynot/rabbit-amazon-forwarder/forwarder"
)

const (
	badRequest     = "Bad request"
	handlerError   = "Handled"
	unhandledError = "Unhandled"
)

func TestCreateForwarder(t *testing.T) {
	entry := config.AmazonEntry{Type: "Lambda",
		Name:   "lambda-test",
		Target: "function1-test",
	}
	forwarder := CreateForwarder(entry)
	if forwarder.Name() != entry.Name {
		t.Errorf("wrong forwarder name, expected:%s, found: %s", entry.Name, forwarder.Name())
	}
}

func TestPush(t *testing.T) {
	functionName := "function1-test"
	entry := config.AmazonEntry{Type: "Lambda",
		Name:   "lambda-test",
		Target: functionName,
	}
	scenarios := []struct {
		name     string
		mock     lambdaiface.LambdaAPI
		message  string
		function string
		err      error
	}{
		{
			name:     "empty message",
			mock:     mockAmazonLambda{resp: lambda.InvokeOutput{StatusCode: aws.Int64(202)}, function: functionName, message: ""},
			message:  "",
			function: functionName,
			err:      errors.New(forwarder.EmptyMessageError),
		},
		{
			name:     "bad request",
			mock:     mockAmazonLambda{resp: lambda.InvokeOutput{StatusCode: aws.Int64(202)}, function: functionName, message: badRequest},
			message:  badRequest,
			function: functionName,
			err:      errors.New(badRequest),
		},
		{
			name:     "handled error",
			mock:     mockAmazonLambda{resp: lambda.InvokeOutput{StatusCode: aws.Int64(202), FunctionError: aws.String(handlerError)}, function: functionName, message: handlerError},
			message:  handlerError,
			function: functionName,
			err:      errors.New(handlerError),
		},
		{
			name:     "unhandled error",
			mock:     mockAmazonLambda{resp: lambda.InvokeOutput{StatusCode: aws.Int64(202), FunctionError: aws.String(unhandledError)}, function: functionName, message: unhandledError},
			message:  unhandledError,
			function: functionName,
			err:      errors.New(unhandledError),
		},
		{
			name:     "success",
			mock:     mockAmazonLambda{resp: lambda.InvokeOutput{StatusCode: aws.Int64(202)}, function: functionName, message: "abc"},
			message:  "abc",
			function: functionName,
			err:      nil,
		},
	}
	for _, scenario := range scenarios {
		t.Log("Scenario name: ", scenario.name)
		forwarder := CreateForwarder(entry, scenario.mock)
		err := forwarder.Push(scenario.message)
		if scenario.err == nil && err != nil {
			t.Errorf("Error should not occur. Error: %s", err.Error())
			return
		}
		if scenario.err == err {
			return
		}
		if err == nil {
			t.Errorf("Error should occur. Expected: %s", scenario.err.Error())
			return
		}
		if err.Error() != scenario.err.Error() {
			t.Errorf("Wrong error, expecting:%v, got:%v", scenario.err, err)
		}
	}
}

type mockAmazonLambda struct {
	lambdaiface.LambdaAPI
	resp     lambda.InvokeOutput
	function string
	message  string
}

func (m mockAmazonLambda) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	if *input.FunctionName != m.function {
		return nil, errors.New("Wrong function name")
	}
	if string(input.Payload) != m.message {
		return nil, errors.New("Wrong message body")
	}
	if string(input.Payload) == badRequest {
		return nil, errors.New(badRequest)
	}
	return &m.resp, nil
}
