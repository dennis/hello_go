# Compile and run it

```
git clone https://github.com/dennis/hello_go
cd hello_go
go build main.go
./main
```

# Using it via Docker

To build the image and start it:
```
docker build -t hello_go .           
docker run -p "8080:8080" hello_go
```

# Using the service

The service does not persist any data, so any changes are only effective while
it is running. Credentials are hardcoded in source, so ther is no security
either.

The services is pre-populated with two messages. One from each user.

## Authentication

The service requires basic authentication on all requests. Use
`authtokendennis` or `authtokenmarianne` for this. This is a token, so there
is no password

## Via postman

For you convience I've created a collection for
[Postman](https://www.postman.com/), you can
[import](https://learning.postman.com/docs/postman/collections/importing-and-exporting-data/#importing-data-into-postman)
the collection and configure a `api_token` variable (I suggest you do it via
the environment). It needs to be either `authtokendennis` or
`authtokenmarianne`. 

## Overview

| Verb   | URL                                  | Description                                        |
|--------|--------------------------------------|----------------------------------------------------|
| GET    | http://localhost:8080/api/messages   | Get all messages                                   |
| GET    | http://localhost:8080/api/messages/1 | Get a single mesages                               |
| POST   | http://localhost:8080/api/messages   | Creates a new message                              |
| DELETE | http://localhost:8080/api/messages/1 | Deletes a mesages (only if user wrote the message) |
| PUT    | http://localhost:8080/api/messages/1 | Updates a mesages (only if user wrote the message) |

## Examples

```
$ curl -u authtokendennis: http://localhost:8080/api/messages
[{"id":"1","topic":"Hello World","body":"Lorem lipsum","author":"Dennis"},{"id":"2","topic":"re: Hello World","body":"Really?","author":"Marianne"}]

$ curl -u authtokendennis: http://localhost:8080/api/messages/1
{"id":"1","topic":"Hello World","body":"Lorem lipsum","author":"Dennis"}

$ curl -u authtokendennis: -X PUT http://localhost:8080/api/messages/1 --data '{"id":"1","topic":"Changed via CURL","body":"Lorem lipsum","author":"Dennis"}'
{"id":"1","topic":"Changed via CURL","body":"Lorem lipsum","author":"Dennis"}

$ curl -u authtokendennis: http://localhost:8080/api/messages/1
{"id":"1","topic":"Changed via CURL","body":"Lorem lipsum","author":"Dennis"}

$ curl -u authtokendennis: -X POST http://localhost:8080/api/messages --data '{"topic":"Added via CURL","body":"Lorem lipsum"}'
{"id":"3","topic":"Added via CURL","body":"Lorem lipsum","author":"Dennis"}

$ curl -u authtokendennis: -X DELETE http://localhost:8080/api/messages/1
$ curl -u authtokendennis: -X DELETE http://localhost:8080/api/messages/2

# id 2 wasn't delete, as dennis (which got authtokendennis) isn't the owner of that message
$ curl -u authtokendennis: http://localhost:8080/api/messages
[{"id":"2","topic":"re: Hello World","body":"Really?","author":"Marianne"},{"id":"3","topic":"Added via CURL","body":"Lorem lipsum","author":"Dennis"}]

$ curl -u authtokenmarianne: -X DELETE http://localhost:8080/api/messages/2
$ curl -u authtokendennis: http://localhost:8080/api/messages
[{"id":"3","topic":"Added via CURL","body":"Lorem lipsum","author":"Dennis"}]
```
