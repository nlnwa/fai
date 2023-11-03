FROM golang:1.21 AS build

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" .


FROM gcr.io/distroless/static-debian12:latest

COPY --from=build /go/src/app/fai /fai

EXPOSE 8081

ENTRYPOINT ["/fai"]
