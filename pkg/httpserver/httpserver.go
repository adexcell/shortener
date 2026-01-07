// Package httpserver пакет для настройки и конфигурирования сервера.
package httpserver

type Config struct {
	Port string `mapstructure:"port"`
}
