FROM golang:1.22 AS base

    WORKDIR /app

    EXPOSE 80

    EXPOSE 443

FROM golang:1.22 AS build

    WORKDIR /src

    ADD ./services/ether .

    RUN ls

    RUN go mod download

    RUN CGO_ENABLED=0 GOOS=linux go build -o ./ unreal.sh/ether/cmd/ether

FROM base AS final

    WORKDIR /app

    COPY --from=build /src/ether .
    
    ENTRYPOINT ["/ether"]