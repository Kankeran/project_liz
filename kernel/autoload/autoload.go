package autoload

import (
	"Liz/kernel/event"
	"Liz/kernel/services"
)

func init() {
	services.Build()
	event.RunListener()
}
