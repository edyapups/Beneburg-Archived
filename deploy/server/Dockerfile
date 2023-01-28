FROM golang:1.19.2-buster AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY pkg ./pkg
COPY assets ./assets
COPY templates ./templates
COPY main.go ./

RUN go build -o /main


## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /main /main

EXPOSE 443
EXPOSE 80

ENTRYPOINT ["/main"]