package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	brainfact "github.com/alkemist-17/brain-fact"
)

const (
	REPLWelcome = "Welcome to the Brainfact 🧠 REPL!"
	Welcome     = "Brainfact! 🧠"
	Success     = "Your script was successfully compiled and run."
	Description = "A brainfuck interpreter written in Go."
	Help        = "Type .help for assistance"
	scriptIndex = 1
)

func main() {
	switch len(os.Args) {
	case 1:
		repl()
	case 2:
		runScript()
	case 3:
		if os.Args[1] == "c" {
			compileScript()
		}
	}
}

func repl() {
	clearScreen()
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
				fmt.Println()
				if err := brainfact.Run(code); err != nil {
					fmt.Printf("\n\n\n Error: %s\n\n\n", err)
				}
				fmt.Printf("\n\n")
			}
		}
	}
	clearScreen()
	fmt.Printf("\n\n\n   Leaving the Brainfact REPL! 🧠\n\n\n")
}

func runScript() {
	if strings.HasSuffix(os.Args[scriptIndex], brainfact.FileExt) {
		if bin, err := os.ReadFile(os.Args[scriptIndex]); err == nil && utf8.Valid(bin) {
			code := string(bin)
			if err := brainfact.Run(code); err != nil {
				fmt.Printf("\n\n\n Error: %s\n\n\n", err)
			} else {
				fmt.Println()
			}
		} else if err != nil {
			fmt.Printf("\n\nError reading file '%v'\n%v\n\n\n", os.Args[scriptIndex], err)
		} else {
			fmt.Printf("\n\nThe file '%v' does not contain a valid utf-8 sequence of bytes\n\n\n", os.Args[scriptIndex])
		}
	} else if strings.HasSuffix(os.Args[scriptIndex], brainfact.CompiledFileExt) {
		if bin, err := os.ReadFile(os.Args[scriptIndex]); err == nil && utf8.Valid(bin) {
			if err := brainfact.RunBytecode(bin); err != nil {
				fmt.Printf("\n\n\n Error: %s\n\n\n", err)
			}
		} else if err != nil {
			fmt.Printf("\n\nError reading file '%v'\n%v\n\n\n", os.Args[scriptIndex], err)
		} else {
			fmt.Printf("\n\nThe file '%v' does not contain a valid utf-8 sequence of bytes\n\n\n", os.Args[scriptIndex])
		}
	}
}

func compileScript() error {
	if strings.HasSuffix(os.Args[scriptIndex+1], brainfact.FileExt) {
		if bin, err := os.ReadFile(os.Args[scriptIndex+1]); err == nil && utf8.Valid(bin) {
			code := string(bin)
			filename, _ := strings.CutSuffix(os.Args[scriptIndex+1], brainfact.FileExt)
			if err := brainfact.Compile(filename, code); err != nil {
				fmt.Printf("\n\n\n Error: %s\n\n\n", err)
			}
		} else if err != nil {
			fmt.Printf("\n\nError reading file '%v'\n%v\n\n\n", os.Args[scriptIndex+1], err)
		} else {
			fmt.Printf("\n\nThe file '%v' does not contain a valid utf-8 sequence of bytes\n\n\n", os.Args[scriptIndex+1])
		}
	}
	return nil
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
