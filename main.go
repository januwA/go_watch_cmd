package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var clear map[string]func()

// https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go
func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func event(args []string) {
	CallClear()
	c := exec.Command(args[0], args[1:]...)
	out, _ := c.Output()
	fmt.Printf("%s", out)
}

func main() {
	interval := flag.Int("n", 2, "seconds")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)

	// watch ctrl+c
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	done := make(chan bool, 1)
	go func() {
		<-sigs
		done <- true
	}()

	event(args)
	for {
		select {
		case <-done:
			ticker.Stop()
			CallClear()
			return
		case <-ticker.C:
			event(args)
		}
	}
}
