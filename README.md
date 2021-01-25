# Passwordless Auth with Amazon Cognito
for POC purpose

# Testing
calling initial-auth with CUSTOM_AUTH flow. this command will return a SESSION. 
also, you will recieve SECRET_CODE via SMS.
> aws --profile <AWS_PROFILE> cognito-idp initiate-auth \
  --auth-flow CUSTOM_AUTH \
  --auth-parameters "USERNAME=<user_in_cognito>" \
  --client-id <CLIENT_ID> \
  
onec you got both session and secret code, input them into command below
> aws --profile <AWS_PROFILE> cognito-idp respond-to-auth-challenge \
  --client-id 2pgajrug6lsuv8mvgauf2e7cg5 \
  --challenge-name CUSTOM_CHALLENGE \
  --challenge-responses USERNAME=sila,ANSWER=<SECRET_CODE> \
  --session "SESSION"
