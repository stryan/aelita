package main

import (
	"fmt"
	"strings"
)

func Get(ael *Controller, args ...string) string {
	fmt.Println(args)
	if len(args) == 0 {
		return "No arguments provided"
	}
	var output []string
	for _, a := range args {
		get_command := ael.FindOutput(a)
		if get_command.GetName() == "nil" {
			return "N/A"
		}
		output = append(output, strings.TrimSpace(get_command.Run(ael, args...)))
	}
	fmt.Println(output)
	result := strings.Join(output, "; ")
	fmt.Println("Result",result)
	return result
}
