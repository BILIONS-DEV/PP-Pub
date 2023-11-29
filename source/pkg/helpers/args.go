package helpers

import "os"

func issetArgsConfig(arg string) bool {
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == arg {
			return true
		}
	}
	return false
}

func IsHaiMode() bool {
	return issetArgsConfig("--hai")
}
