build:
	GOOS=linux GOARCH=amd64 go build -o echo main.go
	docker build -t kureikain/echo .

push:
	docker push kureikain/echo:latest

run:
	docker run -p 5678:5678 kureikain/echo

all: build push
