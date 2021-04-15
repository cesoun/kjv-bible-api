# Parse
FROM node:12 as node
WORKDIR /usr/kjvapi

COPY ./utility ./utility

RUN node ./utility/book-parser.js

# Build
FROM golang:1.16 as go
WORKDIR /usr/kjvapi

ENV GOOS linux 
ENV GOARCH amd64

COPY ./backend ./backend

RUN cd ./backend && go build -o bin/kjv-bible-api .

# Run
FROM ubuntu:20.04 as ubuntu
WORKDIR /usr/kjvapi

COPY --from=node /usr/kjvapi/data ./data
COPY --from=go /usr/kjvapi/backend/bin .

CMD [ "./kjv-bible-api" ]
# CMD [ "./kjv-bible-api", "-p=3000" ]

# This is the app's default port, you could use the line above and change the port here to for exposing the port.
# There is probably a way to make this a shared varaible/argument but it's whatever.
EXPOSE 8080