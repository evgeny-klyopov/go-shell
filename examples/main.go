package main

import (
	"fmt"
	"github.com/evgeny-klyopov/go-shell"
	"time"
)

func main() {
	s := shell.New("/bin/sh", "example", true, true)

	s.Run()

	go func(s shell.Sheller) {
		i := 0
		for {
			time.Sleep(time.Second * 1)
			i++
			if i > 4 {
				s.Stop()
			}
		}
	}(s)

	for {
		val, ok := <-s.GetChannel(shell.OutTypeChannel)
		if ok == false {
			break
		} else {
			fmt.Println(val)
		}
	}
}
