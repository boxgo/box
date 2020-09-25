package app

/*
`Name` application name. You can preset name when go build. And also, you can set name by env `BOX_APP_NAME`.
`Version` application version. You can preset name when go build. And also, you can set name by env `BOX_APP_VERSION`.
`Hostname` hostname
`IP` application runtime host's ip.
`StartAt` application start time.

Example:
```sh
go build -ldflags "-X 'github.com/boxgo/box/pkg/app.Name=box' -X 'github.com/boxgo/box/pkg/app.Version=1.0.0'" *.go
```
or
```sh
BOX_APP_NAME=box BOX_APP_VERSION=1.0.0 ./main
```
*/
