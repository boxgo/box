package rest

func init() {
	registerMain()
	registerConfig()
	registerServer()
	registerAPIS()
	registerMiddleware()
	registerSchedule()
	registerService()
	registerProject()
}

func registerMain() {
	register(sourceInfo{
		path: "src/main.go",
		content: `
package main

import (
	"github.com/boxgo/box"
	"{{.Repo}}/{{.ProjectName}}/src/middlewares"
	"{{.Repo}}/{{.ProjectName}}/src/schedules"
	"github.com/boxgo/kit/logger"
	"github.com/boxgo/kit/request"
)

func main() {
	server := newServer()
	config := newConfig()

	box := box.
		NewBox(
			box.WithConfig(config),
			box.WithBoxes(server),
			box.WithConfigHook(newConfigHook()),
		).
		Mount(
			logger.Default,
			request.GlobalOptions,
		).
		Mount(middlewares.MiniBoxes()...).
		Mount(schedules.MiniBoxes()...)

	box.Serve()
}
		`,
	})
}

func registerServer() {
	register(sourceInfo{
		path: "src/server.go",
		content: `
package main

import (
	"context"

	"{{.Repo}}/{{.ProjectName}}/src/apis"
	"github.com/boxgo/kit/middlewares"
	"github.com/boxgo/kit/middlewares/session"
	"github.com/boxgo/rest"
)

type (
	Server struct {
		*rest.Server
	}
)

func (server *Server) ServerWillReady(ctx context.Context) {
	server.Use(middlewares.DefaultLogger.Logger())
	server.Use(session.Default.Session())
	server.Describe(apis.APIS()...)
}

func (server *Server) ServerDidReady(ctx context.Context) {

}

func (server *Server) ServerWillClose(ctx context.Context) {

}

func (server *Server) ServerDidClose(ctx context.Context) {

}

func newServer() *Server {
	return &Server{
		Server: rest.NewServer(),
	}
}
		`,
	})
}

func registerAPIS() {
	register(sourceInfo{
		path: "src/apis/apis.go",
		content: `
package apis

import (
	"github.com/boxgo/rest"
)

var (
	apis = []rest.API{}
)

func register(api ...rest.API) {
	apis = append(apis, api...)
}

// APIS get register apis
func APIS() []rest.API {
	return apis
}
		`,
	})

	register(sourceInfo{
		path: "src/apis/version.go",
		content: `
package apis

import (
	"github.com/boxgo/rest"
	"github.com/gin-gonic/gin"
)

var (
	BoxVersion   string
	BoxBuildTime string
)

func init() {
	register(
		rest.API{
			Method:      "get",
			Path:        "/version",
			Tags:        []string{"System"},
			Summary:     "Get version",
			Description: "Get server version info",
			Handlers: gin.HandlersChain{func(ctx *gin.Context) {
				ctx.JSON(200, map[string]string{
					"version":   BoxVersion,
					"buildTime": BoxBuildTime,
				})
			}},
		},
	)
}
		`,
	})
}

func registerMiddleware() {
	register(sourceInfo{
		path: "src/middlewares/middleware.go",
		content: `
package middlewares

import (
	"github.com/boxgo/box/minibox"
)

var (
	boxes = []minibox.MiniBox{}
)

func register(miniBoxes ...minibox.MiniBox) {
	boxes = append(boxes, miniBoxes...)
}

// MiniBoxes get register miniboxes
func MiniBoxes() []minibox.MiniBox {
	return boxes
}
		`,
	})

	register(sourceInfo{
		path: "src/middlewares/default.go",
		content: `
package middlewares

import (
	"github.com/boxgo/kit/middlewares"
	"github.com/boxgo/kit/middlewares/session"
)

func init() {
	register(
		middlewares.DefaultLogger,
		session.Default,
	)
}
		`,
	})
}

func registerConfig() {
	register(sourceInfo{
		path: "src/config.go",
		content: `
package main

import (
	"os"

	"github.com/boxgo/box/config"
	"github.com/boxgo/box/config/loader"
)

const (
	envFlag = "APP_MODE"
	envDev  = "dev"
	envTest = "test"
	envUat  = "uat"
	envProd = "prod"
)

func newConfig() config.Config {
	ld := loader.NewFileEnvConfig(configPath())
	config := config.NewConfig(ld)

	return config
}

func configPath() string {
	configFileName := envDev

	if name, ok := os.LookupEnv(envFlag); ok {
		configFileName = name
	}

	return "./config/" + configFileName + ".yaml"
}

func envMode() string {
	return os.Getenv(envFlag)
}
		`,
	})

	register(sourceInfo{
		path: "src/config-hook.go",
		content: `
package main

import (
	"context"

	"{{.Repo}}/{{.ProjectName}}/src/schedules"
)

type (
	configHook struct{}
)

func (*configHook) ConfigWillLoad(ctx context.Context) {

}

func (*configHook) ConfigDidLoad(ctx context.Context) {
	schedules.Start()
}

func newConfigHook() *configHook {
	return &configHook{}
}
		`,
	})

	register(sourceInfo{
		path:    "src/config/test.yaml",
		content: "",
	})

	register(sourceInfo{
		path:    "src/config/uat.yaml",
		content: "",
	})

	register(sourceInfo{
		path:    "src/config/prod.yaml",
		content: "",
	})

	register(sourceInfo{
		path: "src/config/dev.yaml",
		content: `
name: {{.ProjectName}}

rest:
  port: 8080
  mode: release
  doc: true

logger:
  level: debug

request:
  timeout: 5000
  showLog: true

middleware:
  logger:
    requestQueryLimit: 2000
    requestBodyLimit: 2000
    responseBodyLimit: 2000
  session:
    mode: 1
    db: 0
    keyPair: box is awesome
    address:
      - 127.0.0.1:6379

redislock:
  db: 0
  address:
    - 127.0.0.1:6379
    `,
	})
}

func registerSchedule() {
	register(sourceInfo{
		path: "src/schedules/schedules.go",
		content: `
package schedules

import (
	"github.com/boxgo/box/minibox"
)

var (
	boxes = []minibox.MiniBox{}
)

func register(miniBoxes ...minibox.MiniBox) {
	boxes = append(boxes, miniBoxes...)
}

// MiniBoxes get register miniboxes
func MiniBoxes() []minibox.MiniBox {
	return boxes
}

// Start schedules
func Start() {

}
		`,
	})

	register(sourceInfo{
		path: "src/schedules/default.go",
		content: `
package schedules

import "github.com/boxgo/kit/schedule/lock"

func init() {
	register(lock.DefaultRedisLock)
}
		`,
	})
}

func registerService() {
	register(sourceInfo{
		path: "src/services/services.go",
		content: `
package services

import (
	"github.com/boxgo/box/minibox"
)

var (
	boxes = []minibox.MiniBox{}
)

func register(miniBoxes ...minibox.MiniBox) {
	boxes = append(boxes, miniBoxes...)
}

// MiniBoxes get register miniboxes
func MiniBoxes() []minibox.MiniBox {
	return boxes
}
		`,
	})
}

func registerProject() {
	register(sourceInfo{
		path: ".gitignore",
		content: `
.vscode
.debug
.DS_Store

box
vendor
		`,
	})

	register(sourceInfo{
		path: "makefile",
		content: `
#
# Box Makefile
#

GO          ?=  go
BASEDIR      =  $(shell pwd)
GOOS         =  $(shell go env GOOS)
GOARCH       =  $(shell go env GOARCH)
GOPATH       =  $(shell go env GOPATH)
GOFILE       =  box
pkgs         =  $(shell $(GO) list ./... | grep -v /vendor/)

export CONFIG_ROOTPATH=$(BASEDIR)

all: clean format vet


# go format
format:
	@echo ">> formatting code"
	$(GO) fmt $(pkgs)


# go vet
vet:
	@echo ">> vetting code"
	$(GO) vet $(pkgs)


# go clean
clean:
	@echo ">> cleaning project."
	rm -f ${BASEDIR}/${GOFILE}
	rm -f ${BASEDIR}/.debug
	$(GO) clean -i


build_info:
	$(eval BoxVersion := $(shell git rev-parse HEAD))
	$(eval BoxBuildTime := $(shell date '+%Y/%m/%d,%H:%M:%S'))
	@echo BoxVersion=$(BoxVersion)
	@echo BoxBuildTime=$(BoxBuildTime)


# go build linux
build_linux: build_info
	GOOS=linux GOARCH=amd64 go build -o ${GOFILE} -ldflags "-X {{.Repo}}/{{.ProjectName}}/src/apis.BoxVersion=${BoxVersion} -X {{.Repo}}/{{.ProjectName}}/src/apis.BoxBuildTime=${BoxBuildTime}" src/*.go


# go build mac
build_mac: build_info
	GOOS=darwin GOARCH=amd64 go build -o ${GOFILE} -ldflags "-X {{.Repo}}/{{.ProjectName}}/src/apis.BoxVersion=${BoxVersion} -X {{.Repo}}/{{.ProjectName}}/src/apis.BoxBuildTime=${BoxBuildTime}" src/*.go

# go build windows
build_windows: build_info
	GOOS=windows GOARCH=amd64 go build -o ${GOFILE} -ldflags "-X {{.Repo}}/{{.ProjectName}}/src/apis.BoxVersion=${BoxVersion} -X {{.Repo}}/{{.ProjectName}}/src/apis.BoxBuildTime=${BoxBuildTime}" src/*.go

# go run dev
dev:export APP_MODE=dev
dev:
	cd src && gowatch -o ../.debug
		`,
	})
}
