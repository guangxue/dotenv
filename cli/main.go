package main

import (
	"flag"
	"fmt"
	"slices"

	"github.com/guangxue/dotenv"
)

func main() {
	// create sub-command `encrypt`
	encrypt := flag.NewFlagSet("encrypt", flag.ExitOnError)

	// encrypt flag -f for file to encrypt into .env.vault
	fileToEncrypt := encrypt.String("f", ".env", "Specify .env file location")

	// encrypt flage -i for extra file to ignore in .gitignore
	defaultIgnoredFiles := []string{".env", ".env.keys"}
	ignoredFile := encrypt.String("i", "", "List file names for .gitignore")

	// encrypt flage -o for name customized output encrypted file name
	outputFile := encrypt.String("o", ".env.vault", "Set customized output file name")

	// create sub-command `ignore`
	ignore := flag.NewFlagSet("ignore", flag.ExitOnError)
	fileToIgnore := ignore.String("f", ".env.keys", "Specify file append to .gitignore")

	flag.Parse()

	// sub-command name: flag.Args[0]
	subCmdName := flag.Args()[0]
	restArgsOfSubCmd := flag.Args()[1:]
	switch subCmdName {

	case "encrypt":
		encrypt.Parse(restArgsOfSubCmd)

		// Rest of args remaining after encrypt
		restArgsOfEncrypt := encrypt.Args()

		flagSymbol := ""
		for v := range slices.Values(restArgsOfEncrypt) {
			found := slices.Index(restArgsOfEncrypt, v)
			if found > 0 {
				flagSymbol = v
				break
			}
		}

		if flagSymbol != "" {
			panic("illegal command line input:\n\t before flag:" + flagSymbol + " (only accept one file per flag)")
		}

		ignoredFileList := append(defaultIgnoredFiles, *ignoredFile)
		fmt.Printf("dotenv-cli: encrypt file(%s) into (%s) and ignore file %s\n", *fileToEncrypt, *outputFile, ignoredFileList)

		dotenv.Encryption(*fileToEncrypt, *outputFile)
		dotenv.GitIgnore(ignoredFileList...)
	case "ignore":
		ignore.Parse(restArgsOfSubCmd)

		// start to ignore file
		fmt.Println("dotenv cli: ignored file:", *fileToIgnore)
		dotenv.GitIgnore(*fileToIgnore)

	}

}
