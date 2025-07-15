package dotenv

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var vaultFile = ".env.vault"

func quotedString(str string) string {
	return str[1 : len(str)-1]
}

func isQuoted(str string) bool {
	return (str[0] == '\'' && str[len(str)-1] == '\'') || (str[0] == '"' && str[len(str)-1] == '"')
}

func getCmd(str string) []string {
	cmd := cmdExpr.FindAllString(str, -1)
	if len(cmd) > 0 {
		return cmd
	}
	return nil
}

func cmdopt(cmd string) string {
	cmd = cmd[2 : len(cmd)-1]
	shellCmd := exec.Command(cmd)
	opt, err := shellCmd.Output()

	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSpace(string(opt))
}

func normalizeVal(rawKval string) string {
	if isQuoted(rawKval) {
		rawKval = quotedString(rawKval)
	}

	if cmds := getCmd(rawKval); len(cmds) > 0 {
		for _, cmd := range cmds {
			rawKval = strings.ReplaceAll(rawKval, cmd, cmdopt(cmd))
		}
	}
	return rawKval
}

func generateSecretKeyString() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func Encryption(dotenvfile string, outputEncryptedFile string) {
	dotenvMap := envFile(dotenvfile).envMap
	encryptedMap, pkey := dotenvMap.encrypt()
	envFile(outputEncryptedFile).write(encryptedMap.str())

	pkeyMap := make(envMap)
	pkeyMap.setEnvkey(pkey)
	envFile(".env.keys").write(pkeyMap.str())
}

func GitIgnore(filenames ...string) {
	envFile(".gitignore").append(filenames...)
}

// Parse only parse encrypted .env.vault or any encrypted files with different names
// into envMap and save decrypted into environs. if no existing .env.vault found,
// use `./dotenv encrypt` to generate .env.vault or `./dotenv encrypt -o outputfile`
// to load different encrypted filename
func Parse(filename ...string) {

	if len(filename) > 1 {
		fmt.Println("dotenv:too many files: only accept one file")
		os.Exit(0)
	}

	if len(filename) == 0 && !fileExists(vaultFile) {
		fmt.Println("dotenv: no such file .env found")
		os.Exit(0)
	}

	if len(filename) == 1 {
		vaultFile = filename[0]
	}

	encryptedEnvMap := envFile(vaultFile).envMap
	privateKey := envFile(".env.keys").get("DOTENV_PRIVATE_KEY")

	decryptedMap := encryptedEnvMap.decrypt(privateKey)
	if err := decryptedMap.setEnv(); err != nil {
		fmt.Println(err)
	}
}
