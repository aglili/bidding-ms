# A Backend service for an auction/bidding platform built in go

## Prerequisites

- Docker

## Running This Project

- clone repository and change directory into the project directory
- copy the .env.example into a .env file and edit the env variables
- run `docker compose up --build -d`

## How it works

1. Users Sign Up and Create auctions
2. Users Cannot bid on their own auctions
3. Other users bid on open auctions
4. When users are out bid we send notifications via websockets to alert them
5. When the auction closed we also send notifications via websockets to the winner of the auction

## This Project uses

1. Redis Pub/Sub - For real time message communication
2. Session Based Authentication - provided by gin sessions [Docs]("github.com/gin-contrib/sessions")
3. Postgresql for storage
