FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN GOOS=linux CGO_ENABLED=0 go build -o /dist/server ./main.go

FROM alpine:edge

COPY --from=build /dist/server /sbin/server

EXPOSE 80

CMD /sbin/server
