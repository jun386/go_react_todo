FROM golang:1.21.2-alpine
WORKDIR /app
RUN apk update
RUN apk add alpine-sdk

COPY backend/go.mod .
COPY backend/go.sum .

COPY ./backend .
EXPOSE 8080
CMD ["go", "run", "main.go"]