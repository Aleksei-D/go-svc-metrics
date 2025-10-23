# internal

В данной директории и её поддиректориях будет содержаться имплементация вашего сервиса

#  Указание версии при компиляции или запуске приложения
Для компиляции и запуска приложения с указанием buildVersion, buildDate, buildCommit необходимо запустить программу с использованием флага -X


## для сборки приложения
`go build -ldflags "-X main.buildVersion=v0.0.0 -X main.buildDate=24/09/2025 -X main.buildCommit=commit" main.go`

## для запуска приложения
`go run -ldflags "-X main.buildVersion=v0.0.0 -X main.buildDate=24/09/2025 -X main.buildCommit=commit" main.go`