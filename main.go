package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	colorBoldBlue   = "\033[1;34m"
	colorBoldGreen  = "\033[1;32m"
	colorBoldWhite  = "\033[1;37m"
	colorBoldRed    = "\033[1;31m"
	colorBoldCyan   = "\033[1;36m"
	colorBoldYellow = "\033[1;33m"
	colorDefault    = "\033[0m"
)

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		log.Fatal("No se esta recibiendo entrada por pipes")
	}

	scanner := bufio.NewScanner(os.Stdin)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var jsonCURL string

	line := 1
	for scanner.Scan() {
		if line > 1 {
			if scanner.Text() != "" {
				if scanner.Text()[0] == '{' {
					jsonCURL = scanner.Text()
				} else {
					fmt.Println(scanner.Text())
				}
			}
		} else {
			httpinfo := strings.Split(scanner.Text(), " ")
			fmt.Print(colorBoldBlue + httpinfo[0] + " ")
			fmt.Print(colorBoldGreen + httpinfo[1] + " ")
			fmt.Println(colorBoldCyan + httpinfo[2] + colorDefault)
		}

		line++
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(jsonCURL), "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(prettyJSON.String())

}
