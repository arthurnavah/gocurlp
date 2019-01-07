package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	for scanner.Scan() {
		if scanner.Text() != "" {
			if scanner.Text()[0] == '{' {
				jsonCURL = scanner.Text()
			} else {
				fmt.Println(scanner.Text())
			}
		}
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(jsonCURL), "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(prettyJSON.String())

}
