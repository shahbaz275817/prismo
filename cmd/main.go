package main

import "fmt"

var (
	version = "0.0.1"
	commit  = "n/a"
)

func main() {
	cli := newCLI()
	cli.Version = fmt.Sprintf("%s (commit: %s)", version, commit)
	err := cli.Execute()
	if err != nil {
		panic(err)
	}
}
