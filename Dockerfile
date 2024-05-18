# Build tha application from source
FROM golang:1.22-alpine3.19 as build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o shortify main.go

# Deploy the application binary into a distroless image
FROM gcr.io/distroless/static-debian12 as build-release-stage

WORKDIR /app

COPY --from=build-stage /app/shortify .
COPY --from=build-stage /app/config.yml .

EXPOSE 8513

ENTRYPOINT ["./shortify"]

CMD ["serve"]