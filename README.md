# upgraded-telegram

# installation 

docker pull peterjbishop/go-http

# contents

API's and Implementations

Go Http/Net server 

- JWT Authenticated Routes
- rate limited endpoints

- DynamoDB
    + users
        POST http://0.0.0.0:8080/users/new
        POST http://0.0.0.0:8080/users/login
        GET http://0.0.0.0:8080/users/all
        GET http://0.0.0.0:8080/users/id/:id
        PUT http://0.0.0.0:8080/users/update
        DELETE http://0.0.0.0:8080/users/delete/:id
    + messaging
        POST http://0.0.0.0:8080/chats/new
        POST http://0.0.0.0:8080/chats/chat/:id/messages/new
        GET http://0.0.0.0:8080/chats/all
        GET http://0.0.0.0:8080/chats/chat/:id
        GET http://0.0.0.0:8080/chats/chat/:id/messages
        PUT http://0.0.0.0:8080/chats/chat/update
        DELETE http://0.0.0.0:8080/chats/chat/:id/delete
        DELETE http://0.0.0.0:8080/chats/chat/:id/messages/message/:id/delete
    + events
        POST http://0.0.0.0:8080/new
        GET http://0.0.0.0:8080/events/event/:id
        GET http://0.0.0.0:8080/events/all
        PUT http://0.0.0.0:8080/events/event/update
        DELETE http://0.0.0.0:8080/events/event/:id/delete

- S3
    + upload
        POST http://0.0.0.0:8080/upload
    + download
        POST http://0.0.0.0:8080/download

- Google Maps
    + Routing (directions)
        GET http://0.0.0.0:8080/maps/from/:origin/to/:destination
    + Geocode
        GET http://0.0.0.0:8080/maps/reversegeocode/lat/:lat/long/:long
    + Reverse Geocode
        GET http://0.0.0.0:8080/maps/geocode/:address

# notes

