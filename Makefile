bold=\033[1m
normal=\033[0m

help:
	@echo "Варианты выполнения команды make:"
	@echo "\t${bold}make test${normal}\t\t - запуск тестов"
	@echo "\t${bold}make cover${normal}\t\t - вывод покрытия кода тестами в браузер"

test:
	@go test ./internal/service
	@go test ./internal/dto/validators
	@go test ./internal/dto
	@go test ./internal/adapters/rest/handlers

cover:
	@go test -coverprofile cover.out ./... -covermode atomic
	@go tool cover -html=cover.out
	@rm ./cover.out