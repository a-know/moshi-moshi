FROM golang:1.10-alpine AS build

RUN apk update && apk upgrade \
    && apk add curl git

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/a-know/moshi-moshi
COPY . .

RUN dep ensure
RUN go build -o moshi-moshi

# Final Output ... (2)
FROM golang:1.10-alpine
COPY --from=build /go/src/github.com/a-know/moshi-moshi/moshi-moshi /bin/moshi-moshi

EXPOSE 8080
RUN chmod +x /bin/moshi-moshi
CMD /bin/moshi-moshi
