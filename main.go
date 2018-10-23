package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ErrIPNotProvided       = errors.New("no IP address was provided as a query parameter")
	ErrHostnameNotProvided = errors.New("no hostname was provided as a query parameter")
	ErrPasswordNotProvided = errors.New("no password was provided as a query parameter")
)

const (
	namecheapUpdateURI = "https://dynamicdns.park-your-domain.com/update?domain=%s&password=%s&ip=%s"
)

// <?xml version="1.0" encoding="UTF-8"?>
// <interface-response>
//    <Command>SETDNSHOST</Command>
//    <Language>eng</Language>
//    <ErrCount>1</ErrCount>
//    <errors>
//       <Err1>Domain name not found</Err1>
//    </errors>
//    <ResponseCount>1</ResponseCount>
//    <responses>
//       <response>
//          <ResponseNumber>316153</ResponseNumber>
//          <ResponseString>Validation error; not found; domain name(s)</ResponseString>
//       </response>
//    </responses>
//    <Done>true</Done>
//    <debug />
// </interface-response>

type Email struct {
	Where string `xml:"where,attr"`

	Addr string
}

type Address struct {
	City, State string
}

type Result struct {
	XMLName xml.Name `xml:"Person"`

	Name string `xml:"FullName"`

	Phone string

	Email []Email

	Groups []string `xml:"Group>Value"`

	Address
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	ip, ok := request.QueryStringParameters["ip"]
	if !ok {
		return events.APIGatewayProxyResponse{}, ErrIPNotProvided
	}

	domain, ok := request.QueryStringParameters["domain"]
	if !ok {
		return events.APIGatewayProxyResponse{}, ErrHostnameNotProvided
	}

	password, ok := request.QueryStringParameters["password"]
	if !ok {
		return events.APIGatewayProxyResponse{}, ErrPasswordNotProvided
	}

	resp, err := http.Get(fmt.Sprintf(namecheapUpdateURI, domain, password, ip))
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "911",
			StatusCode: 200,
		}, nil
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("%s", content)

	return events.APIGatewayProxyResponse{
		Body:       "Hello " + request.Body,
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
