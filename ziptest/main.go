package main

import (
	"archive/zip"
	"io"
	"log"
	"time"
)

func main() {
	filename := "randomfile.jpg"

	readIO, writeIO := io.Pipe()
	zipWriter := zip.NewWriter(writeIO)
	f, err := zipWriter.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer zipWriter.Close()
	defer writeIO.Close()
	defer readIO.Close()
	go func() {
		for {
			log.Printf("reading")
			buf := make([]byte, 1024)
			n, err := readIO.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				break
			}
			log.Printf("read %d bytes", n)
		}
	}()
	timer := time.NewTimer(time.Second * 1)
	defer timer.Stop()
	for {
		<-timer.C
		_, err := f.Write([]byte("hello world"))
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("wrote %d bytes", len("hello world"))
		timer.Reset(time.Microsecond * 1)

	}

}
