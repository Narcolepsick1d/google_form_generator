FROM golang:1.19-alpine AS builder

#RUN mkdir cmd/main.go
WORKDIR /app
RUN apk add --no-cache gcc musl-dev

COPY ./ ./
RUN go build -o main cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=0 /app/configs/app.env ./configs/
COPY --from=builder /app/main .

CMD ["/app/main"]
#RUN chmod +x /app/google-gen
#RUN CGO_ENABLED=0 go build -o  google-gen ./cmd
#RUN chmod +x cmd/main.go
#FROM alpine:latest
#
#RUN mkdir /home/garik/GolandProjects/google-gen-back/cmd/main.go
#
#COPY --from=builder /app/google-gen /app/
#COPY --from=0 /app/configs/app.env ./configs/
#
#CMD ["./app/google-gen"]