package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		host           string
		port           uint
		framesFilePath string
		linesPerFrame  uint
		fps            uint
	)
	flag.StringVar(&host, "host", "0.0.0.0", "listen host")
	flag.UintVar(&port, "port", 5010, "listen port")
	flag.StringVar(&framesFilePath, "path", "./data/frames.txt", "frames file path")
	flag.UintVar(&linesPerFrame, "lines", 8, "lines per frame")
	flag.UintVar(&fps, "fps", 25, "frames per second")
	flag.Parse()

	if linesPerFrame < 1 {
		fmt.Fprintln(os.Stderr, "lines argument should be greater than 0")
		os.Exit(1)
	}
	if fps < 1 {
		fmt.Fprintln(os.Stderr, "fps argument should be greater than 0")
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen %s: %v\n", addr, err)
		os.Exit(1)
	}
	fmt.Printf("listening on %s\n", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept: %v\n", err)
			continue
		}
		go serveConn(conn, framesFilePath, int(linesPerFrame), fps)
	}
}

func serveConn(conn net.Conn, framesFilePath string, linesPerFrame int, fps uint) {
	defer func() {
		_ = conn.Close()
	}()

	file, err := os.Open(framesFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open frames: %v\n", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1024), 64<<20)

	writer := bufio.NewWriter(conn)
	defer func() {
		_ = writer.Flush()
	}()

	frame := make([]string, 0, linesPerFrame)
	frameDelay := time.Second / time.Duration(fps)

	for scanner.Scan() {
		frame = append(frame, scanner.Text())
		if len(frame) < linesPerFrame {
			continue
		}
		if err := writeFrame(writer, frame); err != nil {
			return
		}
		frame = frame[:0]
		if frameDelay > 0 {
			time.Sleep(frameDelay)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scan frames: %v\n", err)
		return
	}
	if len(frame) > 0 {
		_ = writeFrame(writer, frame)
	}
}

func writeFrame(w *bufio.Writer, lines []string) error {
	if _, err := w.WriteString("\033[2J\033[H"); err != nil {
		return err
	}
	if _, err := w.WriteString(strings.Join(lines, "\n")); err != nil {
		return err
	}
	if _, err := w.WriteString("\n"); err != nil {
		return err
	}
	return w.Flush()
}
