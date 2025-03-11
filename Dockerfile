FROM golang:1.24-alpine AS build 
WORKDIR /src

COPY . .
RUN go mod download && \
	go build -o /reddit-autopost ./main.go

FROM alpine:latest 
WORKDIR /
COPY --from=build /reddit-autopost /reddit-autopost

ENTRYPOINT ["/reddit-autopost"]
