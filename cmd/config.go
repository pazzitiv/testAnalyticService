package main

type Config struct {
	// LogLevel уровень логирования
	LogLevel string `long:"log-level" description:"Log level: panic, fatal, warn, info,debug" env:"LOG_LEVEL" default:"warn"`
	// Host Хост для Web-сервера
	Host string `long:"host" description:"Listen host" env:"HOST" default:"0.0.0.0"`
	// Port Порт для Web-сервера
	Port int `long:"port" description:"Listen port" env:"PORT" default:"8888"`

	// DBHost Хост для DB
	DBHost string `long:"dbhost" description:"DB host" env:"DBHOST" require:"true" default:"127.0.0.1"`
	// DBPort Порт для DB
	DBPort int `long:"dbport" description:"DB port" env:"DBPORT" require:"true" default:"5432"`
	// DBName Имя базы данных для DB
	DBName string `long:"dbname" description:"DB name" env:"DBNAME" require:"true"`
	// DBUsername Пользователь для DB
	DBUsername string `long:"dbusername" description:"DB username" env:"DBUSERNAME" require:"true"`
	// DBPassword Пароль для DB
	DBPassword string `long:"dbpassword" description:"DB password" env:"DBPASSWORD" require:"true"`
}
