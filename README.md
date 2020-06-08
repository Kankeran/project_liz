[![Go Report Card](https://goreportcard.com/badge/github.com/Kankeran/project_liz)](https://goreportcard.com/report/github.com/Kankeran/project_liz)

# Introduction
Liz is an assistant to write the project in Go.
She prepares the framework and builds its components.

# How to use

### 1. Creating new project.
To start a new project with Liz, start the application with the parameter ```new```

the parameter accepts 2 additional flags:

```-name``` - name of the new project

```-path``` - path to new project

### 2. Building a project.
To build a project, start the application with the parameter ```build```

the parameter accepts 1 additional flag:

```-path``` - path to the project to be built

### 3. Using in code.
Use in the code is limited to filling the file services.yaml
There are 2 branches in the file ```services``` and ```listeners```

services are declared as follows:
```
service_name:
    struct: "returned_service_type"
    constructor: "funtion_returns_pointer_to_new_service" (optional)
    arguments: (array of arguments for constructor) (optional)
    calls: (array of specified struct named call)
    return: (argument to return when service is called typed in parametr "struct") (default: $this) (only when construct defined)
```

params of call structure:
```
    method: "method_name_to_call"
    arguments: (array of arguments passed to method)
```

listeners are declared as follows:
```
event_to_listening: (array of listener struct)
    service: "@service_name_where_method_defined"
    method: "method_to_call_after_event"
```

additional functionality in services:

```$this``` - pointer to services returned by a construct or a struct

```@serwis_name``` - returns specified service by name

```~``` - conjunction of strings

```$env(env_variable)``` - insert specified value by name from env file

```{@(serwis), $this}.function(...)``` - calling function od specified service
