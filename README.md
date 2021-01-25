# Passwordless Auth with Amazon Cognito
for POC purpose

# Testing
> aws cognito-idp initiate-auth \
  --auth-flow CUSTOM_AUTH \
  --auth-parameters "USERNAME=<user_in_cognito>" \
  --client-id <CLIENT_ID> \
  --profile <AWS_PROFILE>
