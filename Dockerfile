FROM golang:1.14.0-alpine3.11 as chromium

RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ca-certificates \
    ttf-freefont

FROM golang:1.14 AS gobuilder

WORKDIR /app
COPY . .

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64
RUN go build -o main wikipdf

FROM chromium

COPY --from=gobuilder /app/main .
ENTRYPOINT ["./main"]