bold=\033[1m
normal=\033[0m

help:
	@echo "Варианты выполнения команды make:"
	@echo "\t${bold}make test${normal} - запуск тестов"
	@echo "\t${bold}make cover${normal} - вывод покрытия тестами кода в браузер"

test:
	@go test ./internal/service
	@go test ./internal/dto/validators
	@go test ./internal/dto

cover:
	@go test -coverprofile cover.out ./... -covermode atomic
	@go tool cover -html=cover.out
	@rm ./cover.out