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
)

const (
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
	// sendSMS()

	event.Response = events.CognitoEventUserPoolsCreateAuthChallengeResponse{
		PrivateChallengeParameters: map[string]string{
			"code": code,
		},
		PublicChallengeParameters: map[string]string{
			"phone": event.Request.UserAttributes["phone_number"],
		},
		ChallengeMetadata: code,
	}

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

// func sendSMS()
