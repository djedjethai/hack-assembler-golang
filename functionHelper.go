package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

//
func assembleJump(line string, sb strings.Builder, cInstruction map[string]string, cJump map[string]string) string {
	sb.WriteString("111")

	ci := rj.FindAllStringSubmatch(line, -1)[0][1]

	sb.WriteString(string(cInstruction[strings.TrimSpace(ci)]))

	sb.WriteString(string("000"))

	cj := rj.FindAllStringSubmatch(line, -1)[0][2]

	sb.WriteString(string(cJump[strings.TrimSpace(cj)]))
	sb.WriteString("\n")

	return sb.String()
}

//
func assembleNoJump(line string, sb strings.Builder, cInstruction map[string]string, cDestination map[string]string) string {
	sb.WriteString("111")

	c := rd.FindAllStringSubmatch(line, -1)[0][2]
	sb.WriteString(string(cInstruction[strings.TrimSpace(c)]))

	d := rd.FindAllStringSubmatch(line, -1)[0][1]
	sb.WriteString(string(cDestination[strings.TrimSpace(d)]))

	final := completeBitsBack(sb.String(), 16)

	return final
}

//
func completeBitsFront(s string, lgt int) string {
	if len(s) < lgt {
		for i := len(s); i < lgt; i++ {
			s = "0" + s
		}
		return s
	} else {
		return s
	}
}

func completeBitsBack(s string, lgt int) string {
	if len(s) < lgt {
		for i := len(s); i < lgt; i++ {
			s = s + "0"
		}
		return s + "\n"
	} else {
		return s + "\n"
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
