package box

import (
	"fmt"
	"runtime/debug"
	"strings"
)

const banner = `
____
|  _ \
| |_) | _____  __
|  _ < / _ \ \/ /
| |_) | (_) >  <
|____/ \___/_/\_\
`

func init() {
	fmt.Print(banner)

	splitter := strings.Repeat("-", 80)
	fmt.Println(splitter)
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		fmt.Printf("go version: %s\n", buildInfo.GoVersion)

		for _, setting := range buildInfo.Settings {
			fmt.Printf("%s=%s\n", setting.Key, setting.Value)
		}
	}
	fmt.Println(splitter)
}
