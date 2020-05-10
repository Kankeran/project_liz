package event

import "sync"

// Data holds dispatched event data
type Data struct {
	Name  string
	Value interface{}
}

type dispatchData struct {
	event     *Data
	waitGroup *sync.WaitGroup
}

type dispatcher struct {
	eventChannel chan *dispatchData
	listeners    map[string]func(*Data)
}

var dispatcherInstance *dispatcher

// Dispatch asynchronously call event by name
func Dispatch(name string, data interface{}) {
	dispatcherInstance.dispatch(name, data)
}

// DispatchSync synchronously call event by name
func DispatchSync(name string, data interface{}) {
	dispatcherInstance.dispatchSync(name, data)
}

// PrepareDispatcher prepares and operates the dispatch system
func PrepareDispatcher(listeners map[string]func(*Data)) {
	dispatcherInstance = &dispatcher{
		make(chan *dispatchData),
		listeners,
	}
	dispatcherInstance.run()
}

func (d *dispatcher) run() {
	go func() {
		for {
			go func(dat *dispatchData) {
				d.listeners[dat.event.Name](dat.event)
				if dat.waitGroup != nil {
					dat.waitGroup.Done()
				}
			}(<-d.eventChannel)
		}
	}()
}

func (d *dispatcher) dispatch(name string, data interface{}) {
	d.eventChannel <- &dispatchData{&Data{name, data}, nil}
}

func (d *dispatcher) dispatchSync(name string, data interface{}) {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	d.eventChannel <- &dispatchData{&Data{name, data}, waitGroup}
	waitGroup.Wait()
}
