package autoload

import (
	"Liz/kernel/event"
	"Liz/kernel/services"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	services.Build()
	event.RunListener()
}
