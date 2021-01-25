package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/google/uuid"
)

// Email -
type Email struct {
	Name      string
	Sender    string
	Recipient string
	Subject   string
	TextBody  string
	HTMLBody  string
	PDFFile   string
	PDFBase64 string
	MIME      string
}

const (
	// Sender - Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "sila@colonandcurve.co.th"

	// SesRegion - Replace us-west-2 with the AWS Region you're using for Amazon SES.
	SesRegion = "ap-south-1"

	// CharSet - The character encoding for the email.
	CharSet = "UTF-8"

	// Datetime format
	layout = "Mon, 2 Jan 2006 15:04:05 MST"

	// SecretLen - secret len
	SecretLen = 6
)

// HandleLambdaEvent -
func HandleLambdaEvent(ctx context.Context, e events.CognitoEventUserPoolsCreateAuthChallenge) {
	var secretLoginCode string
	var previousChallenge string

	log.Println("secretLoginCode:", secretLoginCode)

	if len(e.Request.Session) > 0 {
		secretLoginCode = SecretGen()
		// sendEmail(e.Request.UserAttributes["email"])
	} else {
		e.Request.Session
	}

	e.Response.PublicChallengeParameters = map[string]string{
		"email": e.Request.UserAttributes["email"],
	}

	e.Response.PrivateChallengeParameters = map[string]string{
		"secret": secretLoginCode,
	}

	e.Response.ChallengeMetadata = fmt.Sprintf("CODE-%s", secretLoginCode)

}

func main() {
	lambda.Start(HandleLambdaEvent)
}

// SecretGen - secret generator
func SecretGen() string {
	u, _ := uuid.NewRandom()
	uuidIndex := len(u.String()) - SecretLen
	return u.String()[uuidIndex:]
}

// Send -
func (e *Email) Send() {
	// Create a new session and specify an AWS Region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(SesRegion),
	})

	// debuging purpose
	log.Println("check session region, ", *sess.Config.Region)

	// create an SES client in the session.
	log.Println("create ses client")
	svc := ses.New(sess)

	// set SES input
	input := &ses.SendRawEmailInput{
		Source:       aws.String(e.Sender),
		Destinations: []*string{aws.String(e.Recipient)},
		RawMessage:   &ses.RawMessage{Data: []byte(e.MIME)},
	}

	log.Printf("send to %s\n", e.Recipient)
	result, err := svc.SendRawEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
