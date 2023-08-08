FROM golang:1.19-alpine AS builder

RUN mkdir /home/garik/GolandProjects/google-gen-back/cmd/main.go

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o  google-gen ./cmd

RUN chmod +x /app/google-gen

FROM alpine:latest

RUN mkdir /home/garik/GolandProjects/google-gen-back/cmd/main.go

COPY --from=builder /app/google-gen /app/
COPY --from=0 /app/configs/app.env ./configs/

CMD ["./app/google-gen"]