package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

//HTTPinfo ...
type HTTPinfo struct {
	Version       string
	VersionNumber float64
	StatusCode    int
	Status        string
}

//NewHTTInfo ...
func NewHTTInfo(line []byte) (info HTTPinfo, err error) {
	lineSplit := bytes.SplitN(line, []byte{' '}, 3)

	info.Version = string(lineSplit[0])
	info.VersionNumber, err = strconv.ParseFloat(string(bytes.SplitN(lineSplit[0], []byte{'/'}, 2)[1]), 64)
	if err != nil {
		return
	}

	info.StatusCode, err = strconv.Atoi(string(lineSplit[1]))
	if err != nil {
		return
	}
	info.Status = string(lineSplit[2])

	return
}

//CURLData ...
type CURLData struct {
	HTTPInfo HTTPinfo
	Headers  map[string]string
	Body     []byte
}

//NewCURLData ...
func NewCURLData(read *bufio.Reader) (curl CURLData, err error) {
	curl.Headers = make(map[string]string)

	var line []byte
	var passedEnter bool
	for lineNumber := 1; err == nil; lineNumber++ {
		line, _, err = read.ReadLine()

		//HTTPInfo
		if lineNumber == 1 {
			if !bytes.ContainsAny(line, "{}[]()<>-&$:=") {
				curl.HTTPInfo, err = NewHTTInfo(line)
				if err != nil {
					return
				}
			}
			continue
		}

		if len(line) == 0 {
			passedEnter = true
			continue
		}

		//Headers
		if !passedEnter {
			header := bytes.SplitN(line, []byte(":"), 2)
			curl.Headers[string(header[0])] = string(header[1])
		}

		//Body
		if passedEnter {
			if line[0] == '{' || line[0] == '[' {
				curl.Body, err = json.Marshal(string(line))
				if err != nil {
					return
				}
			} else {
				err = errors.New("El body esta codificado en un formato no soportado")
			}
			break
		}
	}

	return
}
