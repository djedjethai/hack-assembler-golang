package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	// "reflect"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	r        = regexp.MustCompile(`^(.{1}).*`)
	i        = regexp.MustCompile(`^(.*)=(.*)\s$`)
	ii       = regexp.MustCompile(`^(.*);(.*)\s$`)
	rj       = regexp.MustCompile(`^.*(;).*$`)
	rd       = regexp.MustCompile(`^.*(=).*$`)
	binary   strings.Builder
	sb       strings.Builder
	codeStrr string
)

func main() {
	cInstruction := tableParser("./tables/cInstructions.txt")
	cDestination := tableParser("./tables/cInstDestination.txt")
	cJump := tableParser("./tables/cJump.txt")
	// fmt.Println(cJump)

	// parse progTest
	// prog, _ := filepath.Abs("./progTest/Add.asm")
	prog, _ := filepath.Abs("./progTest/MaxL.asm")
	progData, _ := ioutil.ReadFile(prog)
	progText := string(progData)

	for _, line := range strings.Split(progText, "\n") {
		if len(line) > 0 {
			firstChar := r.FindAllStringSubmatch(line, -1)[0][1]

			if firstChar == "@" {
				codeStrr = assembleAddress(line, sb)
				binary.WriteString(codeStrr)

			} else if rj.MatchString(line) {
				codeStrr = assembleJump(line, sb, cInstruction, cJump)
				binary.WriteString(codeStrr)

			} else if rd.MatchString(line) {
				codeStrr = assembleNoJump(line, sb, cInstruction, cDestination)
				binary.WriteString(codeStrr)

			} else {
				fmt.Println("other cases: ", line)
			}
		}
	}
	fmt.Println("Binaries")
	codeStrr = binary.String()
	fmt.Println(codeStrr)
	codeBin := []byte(codeStrr)
	if err := ioutil.WriteFile("Add.hack", codeBin, 0777); err != nil {
		log.Fatal(err)
	}
}

func assembleAddress(line string, sb strings.Builder) string {
	i, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		fmt.Println(err)
	}
	b := completeBitsFront(strconv.FormatInt(int64(i), 2), 16)
	// fmt.Println(reflect.TypeOf(b))
	sb.WriteString(b)
	sb.WriteString("\n")

	return sb.String()
}

func assembleJump(line string, sb strings.Builder, cInstruction map[string]string, cJump map[string]string) string {
	sb.WriteString("111")

	c := ii.FindAllStringSubmatch(line, -1)[0][1]

	sb.WriteString(string(cInstruction[c]))

	sb.WriteString(string("000"))

	j := ii.FindAllStringSubmatch(line, -1)[0][2]

	sb.WriteString(string(cJump[j]))
	sb.WriteString("\n")

	return sb.String()
}

func assembleNoJump(line string, sb strings.Builder, cInstruction map[string]string, cDestination map[string]string) string {
	sb.WriteString("111")

	c := i.FindAllStringSubmatch(line, -1)[0][2]
	sb.WriteString(string(cInstruction[c]))

	d := i.FindAllStringSubmatch(line, -1)[0][1]
	sb.WriteString(string(cDestination[d]))

	final := completeBitsBack(sb.String(), 16)

	return final
}

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

// =============================================================

// func assembleJump(r *regexp.Regexp, line string, sb strings.Builder, cInstruction map[string]string, cJump map[string]string) string {
// 	firstChar := r.FindAllStringSubmatch(line, -1)[0][1]
// 	switch firstChar {
// 	case "@":
// 		// fmt.Println("ahhh: ", line)
// 		i, err := strconv.Atoi(strings.TrimSpace(line[1:]))
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		b := completeBitsFront(strconv.FormatInt(int64(i), 2), 16)
// 		// fmt.Println(reflect.TypeOf(b))
// 		sb.WriteString(b)
// 		sb.WriteString("\n")
//
// 		return sb.String()
// 	case "M", "A", "D", "0":
// 		sb.WriteString("111")
//
// 		// same block(no jump)
// 		c := i.FindAllStringSubmatch(line, -1)[0][1]
// 		sb.WriteString(string(cInstruction[c]))
// 		sb.WriteString(string("000"))
// 		j := i.FindAllStringSubmatch(line, -1)[0][2]
// 		sb.WriteString(string(cJump[j]))
//
// 		return sb.String()
// 	default:
// 		// fmt.Println("empty or comm: ", line)
// 		return ""
// 	}
//
// }
//
// func assembleNoJump(r *regexp.Regexp, line string, sb strings.Builder, cInstruction map[string]string, cDestination map[string]string) string {
// 	firstChar := r.FindAllStringSubmatch(line, -1)[0][1]
// 	// will need to add case "(" and "0" later on
// 	switch firstChar {
// 	case "@":
// 		// fmt.Println("ahhh: ", line)
// 		i, err := strconv.Atoi(strings.TrimSpace(line[1:]))
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		b := completeBitsFront(strconv.FormatInt(int64(i), 2), 16)
// 		// fmt.Println(reflect.TypeOf(b))
// 		sb.WriteString(b)
// 		sb.WriteString("\n")
//
// 		return sb.String()
// 	case "M", "A", "D":
// 		sb.WriteString("111")
//
// 		// same block(no jump)
// 		c := i.FindAllStringSubmatch(line, -1)[0][2]
// 		sb.WriteString(string(cInstruction[c]))
//
// 		d := i.FindAllStringSubmatch(line, -1)[0][1]
// 		sb.WriteString(string(cDestination[d]))
//
// 		final := completeBitsBack(sb.String(), 16)
//
// 		return final
// 	default:
// 		// fmt.Println("empty or comm: ", line)
// 		return ""
// 	}
// }

// =========================================================

// package main
//
// import (
// 	"fmt"
// 	"io/ioutil"
// 	"path/filepath"
// 	// "reflect"
// 	"log"
// 	"regexp"
// 	"strconv"
// 	"strings"
// )
//
// var (
// 	r = regexp.MustCompile(`^(.{1}).*`)
// 	i = regexp.MustCompile(`^(.*)=(.*)\s$`)
// )
//
// // A instruction: @value
// // 0 vvv vvvv vvvv vvvv
//
// // C instruction
// //  1 1 1 c0(a) c1 c2 c3 c4 c5 c6 d1 d2 d3 j1 j2 j3
// // dest = comp; jump
//
// func main() {
// 	cInstruction := tableParser("./tables/cInstructions.txt")
// 	cDestination := tableParser("./tables/cInstDestination.txt")
// 	// cJump := tableParser("./tables/cJump.txt")
// 	// fmt.Println(cJump)
//
// 	// parse progTest
// 	prog, _ := filepath.Abs("./progTest/Add.asm")
// 	progData, _ := ioutil.ReadFile(prog)
// 	progText := string(progData)
//
// 	var binary strings.Builder
// 	for _, line := range strings.Split(progText, "\n") {
// 		if len(line) > 0 {
// 			firstChar := r.FindAllStringSubmatch(line, -1)[0][1]
// 			// fmt.Println(firstChar)
//
// 			// will need to add case "(" and "0" later on
// 			switch firstChar {
// 			case "@":
// 				i, err := strconv.Atoi(strings.TrimSpace(line[1:]))
// 				if err != nil {
// 					fmt.Println(err)
// 				}
// 				b := completeBitsFront(strconv.FormatInt(int64(i), 2), 16)
// 				// fmt.Println(reflect.TypeOf(b))
// 				// fmt.Println(b)
// 				binary.WriteString(b)
// 				binary.WriteString("\n")
//
// 			case "M", "A", "D":
// 				var sb strings.Builder
// 				sb.WriteString("111")
//
// 				// check if jump means have ";" as 2nd char
//
//
//
// 				// same block(no jump)
// 				c := i.FindAllStringSubmatch(line, -1)[0][2]
// 				sb.WriteString(string(cInstruction[c]))
//
// 				d := i.FindAllStringSubmatch(line, -1)[0][1]
// 				sb.WriteString(string(cDestination[d]))
//
// 				final := completeBitsBack(sb.String(), 16)
// 				// ===========
//
// 				// fmt.Println(final)
// 				binary.WriteString(final)
// 				binary.WriteString("\n")
//
// 			default:
// 				fmt.Println("empty or comm: ", line)
// 			}
// 		}
// 	}
// 	fmt.Println("Binaries")
// 	codeStr := binary.String()
// 	fmt.Println(codeStr)
// 	codeBin := []byte(codeStr)
// 	if err := ioutil.WriteFile("Add.hack", codeBin, 0777); err != nil {
// 		log.Fatal(err)
// 	}
// }
//
// func completeBitsFront(s string, lgt int) string {
// 	if len(s) < lgt {
// 		for i := len(s); i < lgt; i++ {
// 			s = "0" + s
// 		}
// 		return s
// 	} else {
// 		return s
// 	}
// }
//
// func completeBitsBack(s string, lgt int) string {
// 	if len(s) < lgt {
// 		for i := len(s); i < lgt; i++ {
// 			s = s + "0"
// 		}
// 		return s
// 	} else {
// 		return s
// 	}
// }
//
// func tableParser(filename string) map[string]string {
// 	cAbsPath, _ := filepath.Abs(filename)
// 	cDat, _ := ioutil.ReadFile(cAbsPath)
// 	cText := string(cDat)
//
// 	// set cInstructions into a map
// 	cData := make(map[string]string)
// 	for _, line := range strings.Split(cText, "\n") {
// 		var ref string
// 		for i, value := range strings.Split(line, ":") {
// 			if i == 0 {
// 				ref = value
// 			}
// 			if i == 1 {
// 				cData[ref] = value
// 			}
// 		}
// 	}
// 	return cData
// }
