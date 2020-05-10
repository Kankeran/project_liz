package builder

import (
	"fmt"
	"os/exec"

	"Liz/domain"
	"github.com/go-git/go-git/v5"
)

const configYamlData = "services:\nlisteners:"
const serviceData = "package services\n // Build building container container\n func Build() {\n\nevent.PrepareDispatcher(map[string]func(d *event.Data){})\n\n}"
const containerData = `package container

type containerStruct struct {
	services         map[string]interface{}
	servicesCreators map[string]func() interface{}
}

var containerInstance = &containerStruct{
	services:         make(map[string]interface{}),
	servicesCreators: make(map[string]func() interface{}),
}

// Get getting searched service instance
func Get(serviceName string) interface{} {
	service, ok := containerInstance.services[serviceName]
	if !ok {
		containerInstance.services[serviceName] = containerInstance.servicesCreators[serviceName]()
		service = containerInstance.services[serviceName]
	}

	return service
}

// Has check service exists
func Has(serviceName string) bool {
	_, ok := containerInstance.servicesCreators[serviceName]

	return ok
}

// Set sets function to invoking service with specified name
func Set(serviceName string, serviceCreator func() interface{}) {
	containerInstance.servicesCreators[serviceName] = serviceCreator
}
`

const dispatcherData = `package event

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
`

const appData = `package main

func main() {
	services.Build()
}
`

// ProjectStarter used to generate common files
type ProjectStarter struct {
	configYamlFileWriter, serviceFileWriter, containerFileWriter, dispatcherFileWriter *domain.FileWriter
	codeFormatter                                                                      *domain.CodeFormatter
}

// NewProjectStarter initialize new ProjectStarter
func NewProjectStarter(
	configYamlFileWriter *domain.FileWriter,
	serviceFileWriter *domain.FileWriter,
	containerFileWriter *domain.FileWriter,
	dispatcherFileWriter *domain.FileWriter,
	codeFormatter *domain.CodeFormatter,
) *ProjectStarter {
	return &ProjectStarter{
		configYamlFileWriter: configYamlFileWriter,
		serviceFileWriter:    serviceFileWriter,
		containerFileWriter:  containerFileWriter,
		dispatcherFileWriter: dispatcherFileWriter,
		codeFormatter:        codeFormatter,
	}
}

// Build generates project common files
func (ps *ProjectStarter) Build(projectName string) {
	var err error

	_, err = git.PlainInit("./", false)

	if err != nil && err != git.ErrRepositoryAlreadyExists {
		panic(err)
	}

	cmd := exec.Command("go", "mod", "init", projectName)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = ps.configYamlFileWriter.Write([]byte(configYamlData))
	if err != nil {
		panic(err)
	}

	var output []byte

	output, err = ps.codeFormatter.Format(containerData)
	if err != nil {
		panic(err)
	}
	err = ps.containerFileWriter.Write(output)
	if err != nil {
		panic(err)
	}

	output, err = ps.codeFormatter.Format(dispatcherData)
	if err != nil {
		panic(err)
	}
	err = ps.dispatcherFileWriter.Write(output)
	if err != nil {
		panic(err)
	}

	output, err = ps.codeFormatter.Format(serviceData)
	if err != nil {
		panic(err)
	}
	err = ps.serviceFileWriter.Write(output)
	if err != nil {
		panic(err)
	}

	appWriter := domain.NewFileWriter(".", projectName+".go")

	output, err = ps.codeFormatter.Format(appData)
	if err != nil {
		panic(err)
	}

	err = appWriter.Write(output)
	if err != nil {
		panic(err)
	}
}
