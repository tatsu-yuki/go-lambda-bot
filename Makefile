build:
	GOOS=linux GOARCH=amd64 go build -o summary_bot_build summary_bot/main.go
zip-mac:
  zip summary_bot_build.zip summary_bot_build
zip-win:
	GOOS=linux GOARCH=amd64 build-lambda-zip -o summary_bot_build.zip summary_bot_build



