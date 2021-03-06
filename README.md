# Passwordless Auth with Amazon Cognito
for POC purpose

![diagram1](/images/lambda-challenges.png)

## Testing
1. Calling initial-auth with CUSTOM_AUTH flow. this command will return a SESSION. Also, you will recieve SECRET_CODE via SMS.
> aws --profile <AWS_PROFILE> cognito-idp initiate-auth \
      --auth-flow CUSTOM_AUTH \
      --auth-parameters "USERNAME=<user_in_cognito>" \
      --client-id <CLIENT_ID> 
  
2. Onec you got both session and secret codes, input them into command below.
> aws --profile <AWS_PROFILE> cognito-idp respond-to-auth-challenge 
      --client-id <CLIENT_ID> \
      --challenge-name CUSTOM_CHALLENGE \
      --challenge-responses USERNAME=sila,ANSWER=<SECRET_CODE> \
      --session <SESSION>
  
## References
- https://dev.to/duarten/passwordless-authentication-with-cognito-13c
- https://aws.amazon.com/blogs/mobile/implementing-passwordless-email-authentication-with-amazon-cognito
- https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-challenge.html
- https://docs.aws.amazon.com/cli/latest/reference/cognito-idp/index.html#cli-aws-cognito-idp
  
