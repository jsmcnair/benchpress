FROM golang:1.24-alpine3.21
COPY main.go /tmp/main.go
RUN go build -o /bin/bp /tmp/main.go && rm /tmp/main.go
ENTRYPOINT ["/bin/bp"]