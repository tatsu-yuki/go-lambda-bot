build:
	GOOS=linux GOARCH=amd64 go build -o summary_bot summary_bot/main.go
zip:
	GOOS=linux GOARCH=amd64 go build -o main summary_bot/main.go
	zip handler.zip main



