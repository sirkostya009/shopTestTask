FROM golang:1.22
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -tags prod -o /app/app /app/cmd/main.go
EXPOSE 8080
CMD [ "/app/app" ]
