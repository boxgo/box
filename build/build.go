// Package build contains application build information
// ```sh
// go build -ldflags="-X 'github.com/boxgo/box/v2/build.ID={ID}' -X 'github.com/boxgo/box/v2/build.Name={Name}' -X 'github.com/boxgo/box/v2/build.Version={Version}' -X 'github.com/boxgo/box/v2/build.Namespace={Namespace}'"
// ```
package build

import (
	"os"
)

var (
	Namespace = ""
	Name      = ""
	Version   = ""
	ID        = ""
)

func init() {
	if ID == "" {
		ID, _ = os.Hostname()
	}
}

func SetNamespace(val string) {
	Namespace = val
}

func SetName(val string) {
	Name = val
}

func SetVersion(val string) {
	Version = val
}

func SetId(val string) {
	ID = val
}
