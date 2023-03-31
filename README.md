# Kafka-producer

This application demonstrates a kafka producer. We are using 
[Confluent Kafka](github.com/confluentinc/confluent-kafka-go/kafka) library to implement a kafka producer.

### Quickstart

You don't need to have kafka installed locally as this app is using docker to run kafka.
#### Steps:-
1. Just clone this repository and run "docker compose up -d" to pull the required images and start kafka in docker container in detached mode.
2. Run the application using "go run cmd/main.go"
3. That's it the application is up and ready!

### Description 

We have used [GIN](https://github.com/gin-gonic/gin) web framework to implement REST API.

Currently it supports only one POST route "/entry" which takes JSON body for example 
{
    "name":"James Bond"
    "amount":10000
}
and sends this data into a kafka topic named "payment".

Any application can consume from this topic if it has required information about the topic and kafka server.

This application also supports message retry mechanism, so if a message fails to get produced then it retries 3 times before throwing an error.


