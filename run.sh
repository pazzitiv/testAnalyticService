#!/bin/sh -e

# Подтянуть зависимости
deps(){
  go get ./...
}

build(){
  deps
  GOARCH=amd64 GOOS=linux go build -o ./bin/analytics-service ./cmd
}

run(){
  echo "Сборка приложения"
  build
  echo "Запуск контейнеров"
  docker-compose -f ./docker-compose.yml build
}

using(){
  echo "Укажите команду при запуске: ./run.sh [command]"
  echo "Список команд:"
  echo "  build - сборка сервиса"
  echo "  run - сборка и запуск сервиса"
}

command="$1"
if [ -z "$command" ]
then
 using
 exit 0;
else
 $command $@
fi
