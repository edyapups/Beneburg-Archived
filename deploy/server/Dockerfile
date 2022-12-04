## Build
FROM node:19-alpine AS frontend-build

WORKDIR /app
ARG ENVIRONMENT
ENV ENVIRONMENT=${ENVIRONMENT}

COPY ./frontend/package.json ./
COPY ./frontend/src ./src/
COPY ./frontend/tsconfig.json ./
COPY ./frontend/.env.${ENVIRONMENT} ./
COPY ./frontend/tsconfig.node.json ./
COPY ./frontend/vite.config.ts ./
COPY ./frontend/index.html ./
COPY ./frontend/svelte.config.js ./

RUN npm install
RUN npm run build -- --mode ${ENVIRONMENT}

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

COPY --from=frontend-build /app/dist /static
COPY --from=build /main /main

EXPOSE 443
EXPOSE 80

ENTRYPOINT ["/main"]