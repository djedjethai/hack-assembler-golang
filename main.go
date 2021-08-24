package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	r = regexp.MustCompile(`^(.{1}).*`)
)

// A instruction: @value
// 0 vvv vvvv vvvv vvvv

// C instruction
//  1 1 1 c0(a) c1 c2 c3 c4 c5 c6 d1 d2 d3 j1 j2 j3
// dest = comp; jump

func main() {
	// cInstruction := tableParser("./tables/cInstructions.txt")
	// cDestination := tableParser("./tables/cInstDestination.txt")
	// cJump := tableParser("./tables/cJump.txt")
	// fmt.Println(cInstruction)
	// fmt.Println(cDestination)
	// fmt.Println(cJump)

	// parse progTest
	prog, _ := filepath.Abs("./progTest/Add.asm")
	progData, _ := ioutil.ReadFile(prog)
	progText := string(progData)

	// fmt.Println(progText)
	for _, line := range strings.Split(progText, "\n") {
		if len(line) > 0 {
			firstChar := r.FindAllStringSubmatch(line, -1)[0][1]
			// fmt.Println(firstChar)

			// will need to add case "(" and "0" later on
			switch firstChar {
			case "@":
				fmt.Println("@: ", line)
			case "M", "A", "D":
				fmt.Println("A or M or D: ", line)
			default:
				fmt.Println("empty or comm: ", line)
			}
		}
	}

}

func tableParser(filename string) map[string]string {
	cAbsPath, _ := filepath.Abs(filename)
	cDat, _ := ioutil.ReadFile(cAbsPath)
	cText := string(cDat)

	// set cInstructions into a map
	cData := make(map[string]string)
	for _, line := range strings.Split(cText, "\n") {
		var ref string
		for i, value := range strings.Split(line, ":") {
			if i == 0 {
				ref = value
			}
			if i == 1 {
				cData[ref] = value
			}
		}
	}
	return cData
}
