package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitorTimes = 5
const delayInSeconds = 5

var logFile *os.File

func main() {
	showIntro()

	var err error
	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	for {
		showMenu()
		command := readCommand()

		switch command {
			case 1:
				startMonitoring()
			case 2:
				fmt.Println("Showing logs")
				printLogs()
			case 0:
				fmt.Println("Exiting program...")
				os.Exit(0)
			default:
				fmt.Println("Command not found")
				os.Exit(-1)
		}
	}
}

func showMenu() {
	fmt.Println("1- Start monitoring")
	fmt.Println("2- Show logs")
	fmt.Println("0- Exit program")
}

func showIntro() {
	name := "Hélio"
	version := 1.1

	fmt.Println("Hello,", name)
	fmt.Println("This program it's in version:", version)
}

func readCommand() int {
	var command int
	fmt.Scan(&command)
	return command
}

func startMonitoring() {
	fmt.Println("Monitoring")
	sites := readFileSites()
	for i := 0; i < monitorTimes; i++{
		for _, site := range sites {		
			checkSite(site)
		}
		time.Sleep(delayInSeconds * time.Second)
	}
}

func checkSite(site string) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if resp.StatusCode == 200 {
		fmt.Println(site, "it's running ok")
		registerLog(site, true)
	} else {
		fmt.Println(site, "it's off air")
		registerLog(site, false)
	}
}

func readFileSites() []string {
	var sites []string
	file, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("Error opening sites file:", err)
		return sites
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		sites = append(sites, line)
		if err == io.EOF {
			break
		}
	}
	file.Close()
	return sites
}

func registerLog(site string, status bool) {
	logFile.WriteString(time.Now().Format("02/01/2006 15:04:05") + " " + site + " - online: " + strconv.FormatBool(status) + "\n")
}

func printLogs() {
	file, err := os.ReadFile("log.txt")

	if err != nil { 
		fmt.Println("Error reading log file:", err)
		return
	}

	fmt.Println(string(file))
}