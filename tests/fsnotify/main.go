package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

var verboseFlag = flag.Bool("v", false, "verbose output")

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Printf("usage: %s <file or directory> ...\n", os.Args[0])
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, file := range args {
		stat, err := os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}
		if stat.IsDir() {
			log.Print("watching directory: ", file)
		} else {
			log.Print("watching file: ", file)
		}
		if err := watcher.Add(file); err != nil {
			log.Fatal(err)
		}
	}

	done := make(chan int)
	go func() {
		for err := range watcher.Errors {
			log.Print("error: ", err)
		}
		done <- 1
	}()

	go func() {
		for evt := range watcher.Events {
			if *verboseFlag {
				log.Print(evt.String())
			}
		}
		done <- 1
	}()

	<-done
	<-done
}
