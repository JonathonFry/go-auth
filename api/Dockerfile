FROM golang:1.10 AS builder

ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

COPY . /go/src/github.com/jonathonfry/login

WORKDIR /go/src/github.com/jonathonfry/login

RUN dep ensure --vendor-only

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /api .

CMD if [ ${ENV} = production ]; \
	then \
	/api; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi