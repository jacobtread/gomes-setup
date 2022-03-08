package main

import (
	"bufio"
	"fmt"
	"github.com/jacobtread/gelv"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	if !gelv.IsElevated() { // This program must be run as admin
		gelv.Elevate()
		return
	}
	const intro = `

   ____           __  __   _____   ____  
  / ___|   ___   |  \/  | | ____| / ___|  
 | |  _   / _ \  | |\/| | |  _|   \___ \ 
 | |_| | | (_) | | |  | | | |___   ___) |
  \____|  \___/  |_|  |_| |_____| |____/ 
 
 Go Mass Effect (3) Server redirector.

 Start with -r or --remove to remove redirects

`
	println(intro)

	// The ip addresses that need to be redirected to the GoMES server
	addresses := []string{
		"gosredirector.ea.com",
		"383933-gosprapp396.ea.com",
		"gosgvaprod-qos01.ea.com",
		"gosiadprod-qos01.ea.com",
		"gossjcprod-qos01.ea.com",
		"reports.tools.gos.ea.com",
		"waleu2.tools.gos.ea.com",
		"me3.goscontent.ea.com",
	}

	nl := regexp.MustCompile("\r?\n")
	remove := IsRemove()

	suffix := "#gomes-redirect"
	hostsPath := os.Getenv("SystemRoot") + `\System32\drivers\etc\hosts`
	byteContents, err := ioutil.ReadFile(hostsPath)

	contents := nl.Split(string(byteContents), -1)
	lines := make([]string, 0)
	for _, content := range contents {
		if !strings.HasSuffix(content, suffix) {
			lines = append(lines, content)
		}
	}

	if !remove {
		fmt.Println("Please enter the IP address of the GoMES server:")
		reader := bufio.NewReader(os.Stdin)

		address, err := reader.ReadString('\n')
		if err != nil {
			_, _ = os.Stderr.WriteString("Err failed to take ip from input")
			return
		}
		address = nl.ReplaceAllString(address, "")

		println("Using the following address:")
		println(address)

		for _, s := range addresses {
			a := fmt.Sprintf("%s %s %s", address, s, suffix)
			lines = append(lines, a)
		}
	}

	out := strings.Join(lines, "\r\n")

	err = ioutil.WriteFile(hostsPath, []byte(out), 0644)
	if err != nil {
		_, _ = os.Stderr.WriteString("Failed to update entries in hosts file")
		return
	}

	if remove {
		println("Successfully removed entries")
	} else {
		println("Successfully added redirects")
	}
}

func IsRemove() bool {
	for _, arg := range os.Args {
		if strings.EqualFold("-r", arg) || strings.EqualFold("--remove", arg) {
			return true
		}
	}
	return false
}
