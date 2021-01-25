package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HandleLambdaEvent -
func HandleLambdaEvent(ctx context.Context, event *events.CognitoEventUserPoolsVerifyAuthChallenge) (*events.CognitoEventUserPoolsVerifyAuthChallenge, error) {
	challengeAnswer := event.Request.ChallengeAnswer
	log.Println("challengeAnswer: ", challengeAnswer)

	expectedAnswer := event.Request.PrivateChallengeParameters["code"]
	log.Println("expectedAnswer: ", expectedAnswer)

	if challengeAnswer == expectedAnswer {
		event.Response.AnswerCorrect = true
	} else {
		event.Response.AnswerCorrect = false
	}

	b, _ := json.Marshal(event)
	log.Println(string(b))

	return event, nil

}

func main() {
	lambda.Start(HandleLambdaEvent)
}
