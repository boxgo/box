/*
Package app manages application build information and runtime information. Name and version can be preset when build or run.

Build demo:
	go build -ldflags "-X 'github.com/boxgo/box/pkg/app.Name=box' -X 'github.com/boxgo/box/pkg/app.Version=1.0.0'" *.go

Run demo:
	./main
*/
package app
