package main

import (
	"os"
)
func main() {

	
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "0":
			generate_input()
		case "1":
			generate_output()
		case "2":
			push_to_db()
		}
	}
}