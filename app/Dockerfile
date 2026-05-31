FROM golang:1.22-alpine3.21 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /feline-registry .

FROM scratch
COPY --from=build /feline-registry /feline-registry
EXPOSE 8080
ENTRYPOINT ["/feline-registry"]
