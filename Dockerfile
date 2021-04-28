# Frontend
FROM node:15.14.0 AS frontend

WORKDIR /usr/src/app
COPY package.json yarn.lock ./
COPY frontend frontend

RUN yarn
RUN yarn build

# Backend
FROM golang:1.16.3-alpine

WORKDIR /go/src/github.com/broothie/jottr
COPY go.mod go.sum main.go ./
COPY --from=frontend /usr/src/app/public public

RUN go build
CMD ./jottr
