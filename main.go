package main

import "glookbs.github.com/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
