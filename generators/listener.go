package generators

import (
	"Liz/elements"
)

type Listener struct{}

func (l *Listener) Generate(listeners []*elements.Listener) string {
	var code = "func(d *event.Data) {\n"
	for _, listener := range listeners {
		code += listener.ServiceGetter + "." + listener.Method + "(d)\n"
	}
	code += "},\n"
	return code
}
