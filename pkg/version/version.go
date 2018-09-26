package version

import (
	"fmt"
	"os"
	"runtime"
)

var (
	VERSION = "Not provided."
	COMMIT  = "Not provided."
)

// PrintInfoAndExit prints versions from the array returned by Info() and exit
func PrintInfoAndExit() {
	for _, i := range Info() {
		fmt.Printf("%v\n", i)
	}
	os.Exit(0)
}

// Info returns an array of various service versions
func Info() []string {
	return []string{
		fmt.Sprintf("Version: %s", VERSION),
		fmt.Sprintf("Git SHA: %s", COMMIT),
		fmt.Sprintf("Go Version: %s", runtime.Version()),
		fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
