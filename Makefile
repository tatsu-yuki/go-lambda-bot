build:
	GOOS=linux GOARCH=amd64 go build -o summary_bot_build summary_bot/main.go
	GOOS=linux GOARCH=amd64 go build -o question_bot question_bot/main.go

zip-mac:
	GOOS=linux GOARCH=amd64 go build -o summary_bot_build summary_bot/main.go
	zip summary_bot_build.zip summary_bot_build
	GOOS=linux GOARCH=amd64 go build -o question_bot_build question_bot/main.go
	zip question_bot_build.zip question_bot_build

zip-win:
	GOOS=linux GOARCH=amd64 go build -o summary_bot_build summary_bot/main.go
	GOOS=linux GOARCH=amd64 build-lambda-zip -o summary_bot_build.zip summary_bot_build
