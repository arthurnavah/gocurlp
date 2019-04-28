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

	PrintDataCURL(curl)
}

//PrintDataCURL ...
func PrintDataCURL(curl models.CURLData) (err error) {
	httpinfo := curl.HTTPInfo

	var colorVersion string

	switch fmt.Sprintf("%d", httpinfo.StatusCode)[0] {
	case '2':
		colorVersion = colorBoldGreen
	case '3':
		colorVersion = colorBoldBlue
	default:
		colorVersion = colorBoldRed
	}

	fmt.Printf("%s%s%s %s%d%s %s%s%s\n",
		colorBoldWhite, httpinfo.Version, colorDefault,
		colorVersion, httpinfo.StatusCode, colorDefault,
		colorVersion, httpinfo.Status, colorDefault,
	)

	for k, v := range curl.Headers {
		fmt.Printf("%s%s%s:%s%s%s\n",
			colorBoldBlue, k, colorDefault,
			colorBoldWhite, v, colorDefault,
		)
	}

	fmt.Println()

	if curl.BodyType == "json" {
		var bufJSON bytes.Buffer
		err = json.Indent(&bufJSON, curl.Body, "", "    ")
		if err != nil {
			return
		}

		readerJSON := strings.NewReader(bufJSON.String())
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

			if jsonFields[0][0] == '"' && len(jsonFields) >= 2 {
				if len(jsonFields) >= 2 {
					newJSON += spaces + colorBoldYellow + jsonFields[0] + colorDefault + ":" + jsonFields[1] + "\n"
				}
			} else {
				newJSON += spaces + jsonLine + "\n"
			}
		}
		fmt.Println(newJSON)
	} else {
		fmt.Println(string(curl.Body))
	}

	return
}
