package main

import "github.com/l-ILINDAN-l/WorkerPool/cmd"

// main is the entry point for the application.
// Its sole responsibility is to execute the root command from the cmd package,
// which in turn handles all command-line logic.
func main() {
	cmd.Execute()
}
