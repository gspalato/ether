FROM golang:1.22 AS base

    WORKDIR /app

    EXPOSE 80

    EXPOSE 443

FROM golang:1.22 AS build

    WORKDIR /app

    ADD ./services/ether .

    RUN ls -R

    RUN go mod download

    RUN CGO_ENABLED=0 GOOS=linux go build unreal.sh/ether/cmd/ether

FROM base AS final

    WORKDIR /src

    COPY --from=build /app/ether .

    ENTRYPOINT ["/ether"]