// Command cli processes a single local image and prints a JSON manifest.
//
// It must stay thin: parse flags, load config, build dependencies, call the
// shared application service, print the result, and set the exit code. No
// image-processing logic belongs here.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "cli: error:", err)
		os.Exit(1)
	}
}

// run holds the real (testable) CLI logic. Right now it does nothing.
//
// TODO (story 04): parse --input/--output/--profiles, validate the input path,
// open the file, and hand it to the application service.
func run() error {
	fmt.Println("cli: not implemented yet")
	return nil
}
