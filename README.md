# barnard

barnard is a terminal-based client for the [Mumble](https://mumble.info) voice
chat software.

![Screenshot](https://i.imgur.com/B8ldT5k.png)

## Installation

Requirements:

1. [Go](https://golang.org/)
2. [Git](https://git-scm.com/)
3. [Opus](https://opus-codec.org/) development headers
4. [OpenAL](http://kcat.strangesoft.net/openal.html) development headers

To fetch and build:

````
apt install libalut-dev libopus-dev
go install layeh.com/barnard@latest
````

After running the command above, `barnard` will be compiled as `$(go env GOPATH)/bin/barnard`.

## Development Environment Setup

````
apt install libalut-dev libopus-dev
go mod init
go get
go build
````

## Manual

### Key bindings

- <kbd>F1</kbd>: show help
- <kbd>F2</kbd>: toggle voice transmission
- <kbd>Ctrl+L</kbd>: clear chat log
- <kbd>Tab</kbd>: toggle focus between chat and user tree
- <kbd>Page Up</kbd>: scroll chat up
- <kbd>Page Down</kbd>: scroll chat down
- <kbd>Home</kbd>: scroll chat to the top
- <kbd>End</kbd>: scroll chat to the bottom
- <kbd>F10</kbd>: quit

## License

GPLv2

## Author

Tim Cooper (<tim.cooper@layeh.com>)
