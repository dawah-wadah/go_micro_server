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