package main

import (
	"fmt"
	"os"
)

func shellExec(cmd string, apiCommand *API_Command) {

}

func pwdExec(cmdArgs []string) string {
	out, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)
	return out
}
