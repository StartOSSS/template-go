FROM golang:1.25 as build
WORKDIR /src
COPY go.mod go.sum ./
COPY integration/main.go ./integration/main.go
RUN go build -o /out/integration ./integration

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=build /out/integration /integration
ENTRYPOINT ["/integration"]
