FROM golang:1.23.0-alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /app

# Install dependencies
RUN apk --no-cache add ca-certificates tzdata

COPY . .

RUN go env -w GOPROXY=https://goproxy.io
RUN go mod download

RUN go build -o build cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata 

COPY --from=builder /app/build /release
COPY --from=builder /app/templates /templates

ARG APP_ENV
ARG CONFIG_PATH
ENV APP_ENV=${APP_ENV}
ENV CONFIG_PATH=${CONFIG_PATH}

ENTRYPOINT [ "/release" ]