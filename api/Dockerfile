FROM golang:1.11.0-stretch AS builder

WORKDIR /login

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o api

CMD if [ ${ENV} = production ]; \
	then \
	/api; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi