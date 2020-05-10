package test

import (
	"fmt"

	"Liz/kernel/event"
)

// ExampleListener example listener
type ExampleListener struct {
}

// ShowInfo shows example info
func (e *ExampleListener) ShowInfo(data *event.Data) {
	fmt.Println("[event name: " + data.Name + "] my example info :)")
}

// MyListener example listener
type MyListener struct {
}

// Show shows example info
func (m *MyListener) Show(data *event.Data) {
	fmt.Println("[event name: " + data.Name + "] my second info :)")
}
