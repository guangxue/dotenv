package dotenv

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

var (
	quoted    = regexp.MustCompile(`([a-zA-Z_0-9]+=\"[^\"]*\")`)
	nonQuoted = regexp.MustCompile(`([a-zA-Z_0-9]+=[a-zA-Z_0-9]+[\s|\n])`)
	cmdExpr   = regexp.MustCompile(`(\$\(\w+\))`)
)

type dotEnvFile struct {
	envfile *os.File
	envMap  envMap
	info    os.FileInfo
	content string
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

// read whole file content into string
func readfile(filename string) string {
	if buf, err := os.ReadFile(filename); err != nil {
		panic("envfile: readfile: can not read file:" + err.Error())
	} else {
		return string(buf)
	}
}

// open .env.vault for writing encrypted data, if doesn't exist then create it.
//
// open .env.keys for writing encrypted keys, if doesn't exist then create it.
func openfile(filname string) *os.File {
	if file, err := os.OpenFile(filname, os.O_RDWR|os.O_CREATE, 0644); err != nil {
		panic("envfile: open file error:" + err.Error())
	} else {
		return file
	}
}

// open existing specified file to read / write / append otherwise create it.
func envFile(filename string) *dotEnvFile {

	// .env file must present otherwise panic
	envfileInfo, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) && filename == ".env" {
		panic("envfile: get fileInfo error: " + err.Error())
	}
	file := openfile(filename)
	content := readfile(file.Name())
	quotedExpr := quoted.FindAllString(content, -1)
	nonQuotedExpr := nonQuoted.FindAllString(content, -1)
	exprList := append(quotedExpr, nonQuotedExpr...)
	exprEnvMap := make(envMap)

	return &dotEnvFile{envfile: file, envMap: exprEnvMap.from(exprList), info: envfileInfo, content: content}
}

func (file *dotEnvFile) write(content string) {
	defer file.envfile.Close()
	_, err := file.envfile.WriteString(content)
	if err != nil {
		fmt.Println("Writing to the file error:", err)
	}
}

func (file *dotEnvFile) append(content ...string) {
	defer file.envfile.Close()
	sliceToAppend := []string{}
	for val := range slices.Values(content) {
		if !strings.Contains(file.content, val) {
			sliceToAppend = append(sliceToAppend, val)
		}
	}
	readyContent := "\n" + strings.Join(sliceToAppend, "\n")
	fileSize := file.info.Size()
	file.envfile.WriteAt([]byte(readyContent), fileSize)
}

func (file *dotEnvFile) get(keyname string) string {
	defer file.envfile.Close()
	if val := file.envMap.get(keyname); val != "" {
		return val
	} else {
		panic("envfile:Get error: key not found")
	}
}
