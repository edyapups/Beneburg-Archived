## Build

FROM golang:1.19.2-buster AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY pkg ./pkg

COPY main.go ./

RUN go build -o /main


## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY migrations /migrations
COPY templates /templates
COPY static /static
COPY --from=build /main /main

EXPOSE 443
EXPOSE 80

ENTRYPOINT ["/main"]