# Comdel

automatic deleting comments that contains information about promoting gambling sites on YouTube Video's

## comdel folder

contains the frontend application build using svelte, the aim is to create clear and user friendly interface

## comdel-backend folder

an application Logic or backend applicatin that serve and store information to databases, build using Go programming language and help of Fiber Go, pgx for PostgreSQL Communication, and also google API library for OAuth and accessing YouTube video and comments

## Running

in order to run this application, we need Docker installed in your machine and also you need an endpoint of my AI (which is secret and i wont told you that :v)

you need to clone this repo and run this simple command:
```
docker compose up -d --build
```

and see if it running or not:
```
docker ps
```

if you see the container running, congrats!

if you want to stop the container just run this command
```
docker compose down
```

**You may need root access if you using Linux (just allow user in the same group as Docker)**
