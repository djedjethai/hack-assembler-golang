package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	// "regexp"
	"strings"
)

var (
// c     = regexp.MustCompile(`^([01]{7})*`)
// d =
// j =
)

// A instruction: @value
// 0 vvv vvvv vvvv vvvv

// C instruction
//  1 1 1 c0(a) c1 c2 c3 c4 c5 c6 d1 d2 d3 j1 j2 j3
// dest = comp; jump

func main() {
	cAbsPath, _ := filepath.Abs("./tables/cInstructions.txt")
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
	fmt.Println(cData)
	fmt.Printf("cool: %v\n", cData["M-1"])
}
