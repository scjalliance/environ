package main

// go build -ldflags -H=windowsgui

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	env := make(Environment)
	args := env.Parse(os.Args[1:])
	if len(args) == 0 {
		log.Print("No executable provided")
		os.Exit(2)
	}
	for name, value := range env {
		if err := os.Setenv(name, value); err != nil {
			log.Printf("Failed to set environment variable \"%s\" to \"%s\": %v", name, value, err)
			os.Exit(2)
		}
		log.Printf("%s=%s", name, value)
	}
	log.Printf("Args: %+v", args)
	cmd := exec.Command(args[0], args[1:]...)
	if os.Stdout == nil {
		log.Print("nil out")
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to run \"%s\": %v", args[0], err)
		os.Exit(2)
	}
	cmd.Wait()
}
