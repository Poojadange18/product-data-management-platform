FROM golang:1.26 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o producthub .
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /app/producthub ./producthub
COPY --from=build /app/frontend ./frontend
COPY --from=build /app/docs ./docs
EXPOSE 8080
ENTRYPOINT ["./producthub"]
