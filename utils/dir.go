package utils

import (
	"strings"

	"github.com/spf13/viper"
)

func GetCurrentOpDir(args []string, index int) string {
	currentDir := viper.GetString("currentDir")
	if len(currentDir) == 0 {
		currentDir = "/"
	}
	opPath := currentDir + "/" + args[index]
	if strings.HasPrefix(args[index], "/") {
		opPath = args[index]
	}
	return opPath
}

func GetCurrentDir() string {
	currentDir := viper.GetString("currentDir")
	if len(currentDir) == 0 {
		currentDir = "/"
		viper.Set("currentDir", currentDir)
	}
	return currentDir
}

func SetCurrentDir(currentDir string) {
	viper.Set("currentDir", currentDir)
	viper.WriteConfig()
}
