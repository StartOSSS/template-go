FROM golang:1.25 as build
WORKDIR /src
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot
ENV PORT=8080
COPY --from=build /out/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
