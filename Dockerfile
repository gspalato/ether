FROM golang:1.19 AS base

    WORKDIR /app

    COPY go.mod go.sum ./

    RUN go mod download

    ADD . .

    RUN CGO_ENABLED=0 GOOS=linux go build -o unreal.sh/ether/cmd/ether

FROM base AS final

    WORKDIR /app

    ENTRYPOINT ["/ether"]