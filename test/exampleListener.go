package test

import (
	"fmt"

	"Liz/kernel/event"
)

type ExampleListener struct {

}

func (e *ExampleListener) ShowInfo(data *event.Data) {
	fmt.Println("[event name: " + data.Name + "] my example info :)")
}

type MyListener struct {

}

func (m *MyListener) Show(data *event.Data) {
	fmt.Println("[event name: " + data.Name + "] my second info :)")
}
