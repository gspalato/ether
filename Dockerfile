FROM golang:1.19 AS base

    WORKDIR /app

    EXPOSE 80

    EXPOSE 443

FROM golang:1.19 AS build

    WORKDIR /app

    ADD ./services/ether .

    RUN go mod download

    RUN CGO_ENABLED=0 GOOS=linux go build -o unreal.sh/ether/cmd/ether

FROM base AS final

    WORKDIR /src

    COPY --from=build /app/ether .
    
    ENTRYPOINT ["/ether"]