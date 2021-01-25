package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ChallengeStatus -
type ChallengeStatus struct {
	AttemptCount int
	Passed       bool
}

const (
	attempLimit = 3
)

// HandleLambdaEvent -
func HandleLambdaEvent(ctx context.Context, event *events.CognitoEventUserPoolsDefineAuthChallenge) (*events.CognitoEventUserPoolsDefineAuthChallenge, error) {

	event.Response = calculateEventResponse(event.Request.Session)

	b, _ := json.Marshal(&event)
	log.Println(string(b))

	return event, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

func calculateEventResponse(session []*events.CognitoEventUserPoolsChallengeResult) events.CognitoEventUserPoolsDefineAuthChallengeResponse {
	var eventResponse = events.CognitoEventUserPoolsDefineAuthChallengeResponse{}
	var sessionLen = len(session)

	if sessionLen > 0 {
		log.Println("We only accept custom challenges; fail auth")
		if session[sessionLen-1].ChallengeName != "CUSTOM_CHALLENGE" {
			eventResponse = events.CognitoEventUserPoolsDefineAuthChallengeResponse{
				IssueTokens:        false,
				FailAuthentication: true,
			}
		} else if sessionLen >= attempLimit && session[sessionLen-1].ChallengeResult == false {
			log.Printf("The user provided a wrong answer %d times; fail auth", sessionLen)
			eventResponse = events.CognitoEventUserPoolsDefineAuthChallengeResponse{
				IssueTokens:        false,
				FailAuthentication: true,
			}
		} else if session[sessionLen-1].ChallengeName != "CUSTOM_CHALLENGE" && session[sessionLen-1].ChallengeResult == true {
			log.Println("The user provided the right answer; succeed auth")
			eventResponse = events.CognitoEventUserPoolsDefineAuthChallengeResponse{
				IssueTokens:        true,
				FailAuthentication: false,
			}
		}
	} else {
		log.Println("The user did not provide a correct answer yet; present challenge")
		eventResponse = events.CognitoEventUserPoolsDefineAuthChallengeResponse{
			IssueTokens:        false,
			FailAuthentication: false,
			ChallengeName:      "CUSTOM_CHALLENGE",
		}
	}

	return eventResponse
}
