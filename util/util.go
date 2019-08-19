package util

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

func Go(function func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf(
					"Panic: %s\n%s",
					err, debug.Stack(),
				)
				errf, err := os.OpenFile("posam.err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
				if err != nil {
					log.Println("Error when catching panic:", err)
					return
				}
				errf.WriteString(msg)
				errf.Sync()
				errf.Close()
				log.Fatal(msg)
			}
		}()
		function()
	}()
}
