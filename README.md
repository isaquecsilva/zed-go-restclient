# zed-go-restclient

## Whats this?

Simple go application to make **HTTP requests** created for integration with **Zed Tasks**.

## How it works?

Application gens an http request file schema, where you put your custom requests. It reads the file, and searchs for the name of request you want to execute.

<br/>
<p>Example:</p>

First, generate your schema file:
```
zed-go-restclient --gen_schema
```

A file, named requests_schema.json, with the following contents shall be created in your current directory:
```json
[
  {
    "name": "",
    "url": "",
    "method": "",
    "headers": {},
    "body": null
  }
]
```

_Fill the schema file according with your needs._

To execute a request, simply type:
```
zed-go-restclient --requests_file requests_schema.json --name "Name of your request present in the schema."
```

## Integrating with Zed Tasks:

To integrate with Zed tasks feature, first, open the tasks settings in Zed.

Shortcut

### _windows:_
```
ctrl + shift + P
```

Search for 'open tasks'.

Append the following contents into array of tasks.json file:

```json
{
  "label": "RestClient - Make Request",
  "command": "zed-go-restclient --requests_file $ZED_FILE --name \"$ZED_SELECTED_TEXT\"",
  "show_command": false,
  "show_summary": true,
  "shell": {
    "program": "your_shell_here"
  }
},
{
  "label": "RestClient - Generate Schema",
  "command": "zed-go-restclient --gen_schema",
  "show_command": false,
  "show_summary": false,
  "shell": {
    "program": "your_shell_here"
  }
}
```

After this setup, to make a request through Zed tasks, go to _requests_schema.json_, with your mouse cursor, select the name of the task to execute and type `shift + alt + T` and run task `RestClient - Make Request`

## Installation:

```
go install github.com/isaquecsilva/zed-go-restclient
```

Or to build from source:

```
go build -trimpath -o ./zed-go-restclient -ldflags='-s -w' ./rest.go
```

In case of building, remember to put the binary somewhere on your filesystem exposed by **PATH** env:

## _TODOS:_

- Application shall handle multipart/form-data uploads;
