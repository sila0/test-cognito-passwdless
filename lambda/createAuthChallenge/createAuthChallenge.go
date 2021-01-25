package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"crypto/rand"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
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
func HandleLambdaEvent(ctx context.Context, event *events.CognitoEventUserPoolsCreateAuthChallenge) (*events.CognitoEventUserPoolsCreateAuthChallenge, error) {
	code := calculateChallengeCode(event.Request.Session)
	
	event.Response = events.CognitoEventUserPoolsCreateAuthChallengeResponse{
		PrivateChallengeParameters: map[string]string{
			"code": code,
		},
		PublicChallengeParameters: map[string]string{
			"phone": event.Request.UserAttributes["phone_number"],
		},
		ChallengeMetadata: code,
	}

	// e.Response.PublicChallengeParameters = map[string]string{
	// 	"email": e.Request.UserAttributes["email"],
	// }

	// e.Response.PrivateChallengeParameters = map[string]string{
	// 	"secret": secretLoginCode,
	// }

	// e.Response.ChallengeMetadata = fmt.Sprintf("CODE-%s", secretLoginCode)

	b, _ := json.Marshal(event)
	log.Println(string(b))

	return event, nil

}

func main() {
	lambda.Start(HandleLambdaEvent)
}

func genSecret() string {
	b := make([]byte, SecretLen)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
	}

	return base64.URLEncoding.EncodeToString(b)[:SecretLen]
}

func calculateChallengeCode(session []*events.CognitoEventUserPoolsChallengeResult) string {
	var sessionLen = len(session)
	var secretLoginCode string

	if sessionLen == 0 {
		log.Println("this is a new auth session, generate a new secret login code")
		secretLoginCode = genSecret()
		log.Println("secretLoginCode:", secretLoginCode)
	} else {
		log.Println("there is existing session. Dont generate new secret, re-use code from current session")
		secretLoginCode = session[sessionLen-1].ChallengeMetadata
	}

	return secretLoginCode
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
