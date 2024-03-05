package main

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println("Cannot find append only file")
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
		mu:   sync.Mutex{},
	}

	go func() {
		for {
			aof.mu.Lock()
			err := aof.file.Sync()
			if err != nil {
				log.Fatal(err)
			}
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *Aof) Write(val Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(val.Marshal())
	if err != nil {
		return err
	}

	return nil
}
