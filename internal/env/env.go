package env

import (
	"os"
	"strconv"
)

func String(s string,fallback string) string {
	variable,ok:=os.LookupEnv(s)
	if !ok {
		return fallback
	}
	return variable
}

func Int (s string,fallback int) int {
	variable,ok:= os.LookupEnv(s)
	if !ok {
		return fallback
	}
	valueAsInt,err:= strconv.Atoi(variable)
	if err!=nil {
		return fallback
	}
	return valueAsInt
}
func Bool (s string,fallback bool) bool{
	variable,ok:= os.LookupEnv(s)
	if !ok {
		return fallback
	}
	valueAsBool,err:= strconv.ParseBool(variable)
	if err!=nil {
		return fallback
	}
	return valueAsBool
}