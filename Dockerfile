FROM golang:latest

RUN mkdir -p /usr/local/iridium
COPY . /usr/local/iridium
WORKDIR /usr/local/iridium
RUN mv scripts/entrypoint.sh .
RUN chmod +x entrypoint.sh
RUN go build -o iridium cmd/iridium/main.go

EXPOSE 8080
ENTRYPOINT [ "./entrypoint.sh" ]
