package main_test

import (
	"github.com/aws/aws-lambda-go/events"
	main "github.com/davidmhewitt/NamecheapDDNSRegisterLambda"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: ""},
			expect:  "",
			err:     main.ErrIPNotProvided,
		},
		{
			request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"ip":       "",
					"domain":   "",
					"password": "",
				},
			},
			expect: "911",
			err:    nil,
		},
	}

	for _, test := range tests {
		response, err := main.Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
