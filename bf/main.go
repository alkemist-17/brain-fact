package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	brainfact "github.com/alkemist-17/brain-fact"
)

const BFFileExtension string = ".bf"
const REPLWelcome string = "Welcome to Brainfact 🧠 REPL!"
const Welcome string = "Brainfact! 🧠"
const Success string = "Your script was successfully compiled and run."
const Description string = "A brainfuck interpreter written in Go."
const Help string = "Type .help for assistance"

func main() {
	clearScreen()
	if len(os.Args) > 1 {
		scriptIndex := 1
		if strings.HasSuffix(os.Args[scriptIndex], BFFileExtension) {
			if bin, err := os.ReadFile(os.Args[scriptIndex]); err == nil && utf8.Valid(bin) {
				code := string(bin)
				if err := brainfact.Run(code); err != nil {
					fmt.Printf("\n\n\n Error: %s\n\n\n", err)
				} else {
					fmt.Printf("\n\n\n   %s\n   %s\n\n\n\n", Welcome, Success)
				}
			} else if err != nil {
				fmt.Printf("\n\nError reading file '%v'\n%v\n\n\n", os.Args[scriptIndex], err)
			} else {
				fmt.Printf("\n\nThe file '%v' does not contain a valid utf-8 sequence of bytes\n\n\n", os.Args[scriptIndex])
			}
		}
	} else {
		var code string
		running := true
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("\n\n\n   %s\n   %s\n   %s\n\n\n", REPLWelcome, Description, Help)
		for running {
			fmt.Printf("%v", brainfact.Prompt)
			if scanner.Scan() {
				code = scanner.Text()
				switch code {
				case ".help":
					printHelp()
				case ".clear":
					clearScreen()
					fmt.Printf("\n\n\n   %s\n   %s\n   %s\n\n\n", Welcome, Description, Help)
				case ".exit":
					running = false
				default:
					println()
					if err := brainfact.Run(code); err != nil {
						fmt.Printf("\n\n\n Error: %s\n\n\n", err)
					}
					println()
				}
			}
		}
		clearScreen()
		fmt.Printf("\n\n\n   Leaving the Brainfact REPL! 🧠\n\n\n")
	}
}

func clearScreen() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}

func printHelp() {
	clearScreen()
	fmt.Printf("\n\n\n   Brainfact REPL! 🧠\n\n\n")
	fmt.Printf("   Command list:\n")
	fmt.Println("   .help  - Show this message")
	fmt.Println("   .exit  - Quit the REPL")
	fmt.Println("   .clear - Clear the terminal")
	fmt.Printf("\n\n")
}
