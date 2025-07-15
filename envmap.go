package dotenv

import (
	"os"
	"strings"

	"github.com/guangxue/ciphertext"
)

type envMap map[string]string

func (evMap envMap) from(dotenvlist []string) envMap {
	entries := make(envMap)
	for _, val := range dotenvlist {
		kname, kval := strings.SplitN(val, "=", 2)[0], strings.SplitN(val, "=", 2)[1]
		kvalue := normalizeVal(kval)
		entries[kname] = kvalue
	}
	return entries
}

func (evMap envMap) setEnvkey(keystr string) {
	evMap["DOTENV_PRIVATE_KEY"] = keystr
}

func (evMap envMap) get(keyname string) string {
	if value, ok := evMap[keyname]; ok {
		return value
	} else {
		panic("envMap: key name not found")
	}
}

func (evMap envMap) decrypt(privateKey string) envMap {
	decrptedEnvMap := make(envMap)
	for key, val := range evMap {
		decrptedEnvMap[key] = ciphertext.Decrypt(val, privateKey)
	}
	return decrptedEnvMap
}

func (evMap envMap) setEnv() error {
	for key, val := range evMap {
		if err := os.Setenv(key, val); err != nil {
			return err
		}
	}
	return nil
}

func (evMap envMap) str() string {
	var str strings.Builder
	for key, val := range evMap {
		str.WriteString(key)
		str.WriteString("=\"")
		str.WriteString(val)
		str.WriteString("\"\n")
	}
	return str.String()
}

func (env envMap) encrypt() (envMap, string) {
	encrypted := make(envMap)
	pkey := generateSecretKeyString()
	for key, val := range env {
		encrypted[key] = ciphertext.Create(val, pkey)
	}
	return encrypted, pkey
}
