# Microservices

## Broker Service
### Add another microservice for authentications
    1. User tries to authenticate thru the broker, 
    2. the broker calls the authentication service, 
    3. determine if that user is able to be authenticated
    4. Send back the appropriate response
### this means authentications has to have some kind of persistence, 
> We need to be able to store some information
#### So we can add a Postgres(in Docker)

## Authentication Service
1. Will listen for a `POST` request with a json body of a username nad pw
2. Then it will use our `data.Models`, to check to see if the password and username combo exist
    - a. Will send an appropriate response back
3. Broker Service is a good example to use

### Update the Broker for a stndard JSON format, and connect to `AUTH`
Modify the broker application
- listen for a request from the frontend
- then fire a request off to the authentication microservice
- recieve the response from the microservice and send some kind of response back to the end user