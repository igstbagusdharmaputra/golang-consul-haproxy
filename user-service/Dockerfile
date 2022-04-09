FROM golang:1.16-alpine AS build-env
RUN apk --no-cache add git gcc bash
WORKDIR /go/src/app
COPY . ./
RUN go get github.com/hashicorp/consul/api
RUN CGO_ENABLED=0 go build -o main .

FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/app/main /app/
EXPOSE 8200
CMD ["./main"]