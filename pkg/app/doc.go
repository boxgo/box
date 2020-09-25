/*
Package app manages application build information and runtime information. Name and version can be preset when build or run.

Name & Version fetch priority:
	go build > BOX_APP_{{XXXX}} > box_{{RandomString(6)}}

Build demo:
	go build -ldflags "-X 'github.com/boxgo/box/pkg/app.Name=box' -X 'github.com/boxgo/box/pkg/app.Version=1.0.0'" *.go

Run demo:
	BOX_APP_NAME=box BOX_APP_VERSION=1.0.0 ./main
*/
package app
