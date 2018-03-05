FROM golang:1.10 AS builder

ARG app_env
ENV APP_ENV $app_env

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

COPY . /go/src/github.com/jonathonfry/login

WORKDIR /go/src/github.com/jonathonfry/login

RUN dep ensure --vendor-only

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

CMD if [ ${APP_ENV} = production ]; \
	then \
	app; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi