package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/arthurnavah/gocurlp/models"
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

//coloringText Coloria e indenta la entrada
func coloringText(scanner *bufio.Scanner) string {
	var output, jsonCURL string
	var prettyJSON bytes.Buffer

	line := 1

	for scanner.Scan() {
		if line > 1 {
			if scanner.Text() != "" {
				if scanner.Text()[0] == '{' {
					jsonCURL = scanner.Text()
				} else {
					header := strings.SplitN(scanner.Text(), ":", 2)

					output += colorBoldBlue + header[0] + ":" + colorDefault
					output += colorBoldWhite + header[1] + colorDefault
					output += "\n"
				}
			}
		} else {
			httpinfo := strings.Split(scanner.Text(), " ")
			output += "\n"

			output += colorBoldWhite + httpinfo[0] + " " + colorDefault
			if httpinfo[1][0] == '2' {
				output += colorBoldGreen + httpinfo[1] + " "
				output += colorBoldGreen + httpinfo[2] + colorDefault
			} else if httpinfo[1][0] == '3' {
				output += colorBoldBlue + httpinfo[1] + " "
				output += colorBoldBlue + httpinfo[2] + colorDefault
			} else {
				output += colorBoldRed + httpinfo[1] + " "
				output += colorBoldRed + httpinfo[2] + colorDefault
			}

			output += "\n\n"
		}

		line++
	}

	err := json.Indent(&prettyJSON, []byte(jsonCURL), "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	readerJSON := strings.NewReader(prettyJSON.String())
	scannerJSON := bufio.NewScanner(readerJSON)

	var newJSON string
	for scannerJSON.Scan() {
		var spaces string
		for _, v := range scannerJSON.Text() {
			if v == ' ' {
				spaces += " "
			} else {
				break
			}
		}

		jsonLine := strings.TrimSpace(scannerJSON.Text())
		jsonFields := strings.SplitN(jsonLine, ":", 2)

		if jsonFields[0][0] == '"' {
			newJSON += spaces + colorBoldYellow + jsonFields[0] + colorDefault + ":" + jsonFields[1] + "\n"
		} else {
			newJSON += spaces + jsonLine + "\n"
		}
	}
	output += "\n" + newJSON

	return output
}

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		log.Fatal("No se esta recibiendo entrada por pipes")
	}

	scanner := bufio.NewReader(os.Stdin)

	curl, err := models.NewCURLData(scanner)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(curl.Body))
}
