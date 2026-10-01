// Command cli processes a single local image and prints a JSON manifest.
//
// It must stay thin: parse flags, load config, build dependencies, call the
// shared application service, print the result, and set the exit code. No
// image-processing logic belongs here.
package main

import (
	"flag"
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
	input := flag.String("input", "", "path to source image")
	output := flag.String("output", "./data/generated", "output directory")
	profiles := flag.String("profiles", "", "profiles to generate (ignored til later)")
	flag.Parse()

	// Just a way to pass the compiler as a "usage"
	_ = output
	_ = profiles

	if *input == "" {
		return fmt.Errorf("Input (--input) is required.")
	}

	info, err := os.Stat(*input)
	if err != nil {
		return fmt.Errorf("input %q: %w", *input, err)
	}
	if info.IsDir() {
		return fmt.Errorf("input %q is a directory, not a file", *input)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("inptu %q is not a regular file: %w", *input, err)
	}

	f, err := os.Open(*input)
	if err != nil {
		return fmt.Errorf("open %q: %w", *input, err)
	}
	defer f.Close()

	fmt.Printf("file: %q", *input)

	return nil
}
