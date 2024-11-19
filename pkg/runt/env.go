package runt

import "os"

// CfgEnv returns value of an environment variable called `name`, if there is no such, it returns def
func CfgEnv(name, def string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return v
}
