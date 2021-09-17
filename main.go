package main

// sed -n '4710,4730p;4730q' Pong.asm > PongPonggameVar.asm

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
	// first parsing regexp
	rj         = regexp.MustCompile(`^\s*?(.{1,3});([a-zA-Z]{1,4}).*$`)
	rd         = regexp.MustCompile(`^\s*?([a-zA-Z]{1,4})=(.{1,3}).*$`)
	binary     strings.Builder
	sb         strings.Builder
	dv         = regexp.MustCompile(`^\s*?\((.*)\).*$`)
	vv         = regexp.MustCompile(`^\s*?@([^SP|LCL|ARG|THAT|THIS][A-Za-z.]+[^$_0-9])*$`)
	vvcomplete = regexp.MustCompile(`^.*@([A-Za-z.]+\.[0-9]).*$`)
	le         = regexp.MustCompile(`^\s*?(\/\/.*|\s*)$`)

	// second parsing regexp
	rv    = regexp.MustCompile(`^\s*?@R([0-9]{1,2}).*$`)
	tthat = regexp.MustCompile(`^\s*?@(SP|LCL|ARG|THAT|THIS|ponggame.0|math.1|math.0|memory.0|output.6|output.5|output.4|output.3|output.2|output.1|output.0|screen.1|screen.2|screen.0).*$`)
	gv    = regexp.MustCompile(`^\s*?@(\D[A-Za-z._$]*[0-9]*).*`)
	ga    = regexp.MustCompile(`^\s*?@(\d*).*`)

	varVar           = make(map[string]string)
	countDestVar int = 0 // start at 0
	countVarVar  int = 0 // start at 0
	v            string
	b            string
)

func main() {
	cInstruction := tableParser("./tables/cInstructions.txt")
	cDestination := tableParser("./tables/cInstDestination.txt")
	cJump := tableParser("./tables/cJump.txt")
	rVar := tableParser("./tables/rVariable.txt")
	thisThat := tableParser("./tables/thisThatTable.txt")

	// parse progTest
	// prog, _ := filepath.Abs("./progTest/Add.asm")
	// prog, _ := filepath.Abs("./progTest/MaxL.asm")
	// prog, _ := filepath.Abs("./progTest/MaxCopy2.asm")
	// prog, _ := filepath.Abs("./progTest/Max.asm")
	// prog, _ := filepath.Abs("./progTest/PongL.asm")
	prog, _ := filepath.Abs("./progTest/Pong.asm")
	// prog, _ := filepath.Abs("./progTest/PongTest.asm")
	// prog, _ := filepath.Abs("./progTest/PongPonggameVar.asm")
	// prog, _ := filepath.Abs("./progTest/MaxCopy.asm")

	progData, _ := ioutil.ReadFile(prog)
	progText := string(progData)

	// first parse file to identify various var and addresses
	for _, line := range strings.Split(progText, "\n") {

		if len(line) > 0 {
			if dv.MatchString(line) {
				// get destination var
				d := dv.FindAllStringSubmatch(line, -1)[0][1]
				// add to table Var
				varVar[strings.TrimSpace(d)] = completeBitsFront(strconv.FormatInt(int64(countDestVar), 2), 16) // in binary
			} else if vv.MatchString(line) {
				// get named var
				v := vv.FindAllStringSubmatch(line, -1)[0][1]
				// check if exist in table Var
				if _, ok := varVar[strings.TrimSpace(v)]; !ok {
					varVar[strings.TrimSpace(v)] = completeBitsFront(strconv.FormatInt(int64(countVarVar), 2), 16) // in binary
					countVarVar++
				}
				countDestVar++

			} else if vvcomplete.MatchString(line) {
				// constant, do nothing with them during first parsing
				// feel kind of problem with course assembler
				// as the value of these constant change sometime...
				countDestVar++

			} else if le.MatchString(line) {
				// fmt.Println("line empty oor comm")

			} else {
				countDestVar++
			}
		}
	}

	// second parsing set the binary file
	for _, line := range strings.Split(progText, "\n") {
		if len(line) > 0 {
			if rv.MatchString(line) {
				// get variable @Rn
				v = rv.FindAllStringSubmatch(line, -1)[0][1]
				// parse table R
				b = rVar[strings.TrimSpace(v)]
				binary.WriteString(b)
				binary.WriteString("\n")

			} else if tthat.MatchString(line) {
				// get the constant SP/LCL/ARG etc...
				v = tthat.FindAllStringSubmatch(line, -1)[0][1]
				// parse table thisThat
				b = thisThat[strings.TrimSpace(v)]
				binary.WriteString(b)
				binary.WriteString("\n")

			} else if gv.MatchString(line) {
				// get var
				v = gv.FindAllStringSubmatch(line, -1)[0][1]
				// parse table Var
				b := varVar[strings.TrimSpace(v)]
				binary.WriteString(b)
				binary.WriteString("\n")

			} else if ga.MatchString(line) {
				// get address from @num
				v = ga.FindAllStringSubmatch(line, -1)[0][1]
				vi, _ := strconv.Atoi(strings.TrimSpace(v))
				b = completeBitsFront(strconv.FormatInt(int64(vi), 2), 16) // in binary
				binary.WriteString(b)
				binary.WriteString("\n")

			} else if rj.MatchString(line) {
				// get binary from jumping command (with ;)
				b = assembleJump(line, sb, cInstruction, cJump)
				binary.WriteString(b)

			} else if rd.MatchString(line) {
				// get binary from no Jumping command (with =)
				b = assembleNoJump(line, sb, cInstruction, cDestination)
				binary.WriteString(b)
			} else {
				// nothing to do on remainings line
			}
		}
	}
	fmt.Println("Binaries")
	codeStrr := binary.String()
	codeBin := []byte(codeStrr)
	if err := ioutil.WriteFile("Add.hack", codeBin, 0777); err != nil {
		log.Fatal(err)
	}
}

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

// =============================================================
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
// 	r        = regexp.MustCompile(`^(.{1}).*`)
// 	i        = regexp.MustCompile(`^(.*)=(.*)\s$`)
// 	ii       = regexp.MustCompile(`^(.*);(.*)\s$`)
// 	rj       = regexp.MustCompile(`^.*(;).*$`)
// 	rd       = regexp.MustCompile(`^.*(=).*$`)
// 	binary   strings.Builder
// 	sb       strings.Builder
// 	codeStrr string
// )
//
// func main() {
// 	cInstruction := tableParser("./tables/cInstructions.txt")
// 	cDestination := tableParser("./tables/cInstDestination.txt")
// 	cJump := tableParser("./tables/cJump.txt")
// 	// fmt.Println(cJump)
//
// 	// parse progTest
// 	// prog, _ := filepath.Abs("./progTest/Add.asm")
// 	// prog, _ := filepath.Abs("./progTest/MaxL.asm")
// 	// prog, _ := filepath.Abs("./progTest/PongL.asm")
// 	prog, _ := filepath.Abs("./progTest/Max.asm")
//
// 	progData, _ := ioutil.ReadFile(prog)
// 	progText := string(progData)
//
// 	for _, line := range strings.Split(progText, "\n") {
// 		if len(line) > 0 {
// 			firstChar := r.FindAllStringSubmatch(line, -1)[0][1]
//
// 			if firstChar == "@" {
// 				codeStrr = assembleAddress(line, sb)
// 				binary.WriteString(codeStrr)
//
// 			} else if rj.MatchString(line) {
// 				codeStrr = assembleJump(line, sb, cInstruction, cJump)
// 				binary.WriteString(codeStrr)
//
// 			} else if rd.MatchString(line) {
// 				codeStrr = assembleNoJump(line, sb, cInstruction, cDestination)
// 				binary.WriteString(codeStrr)
//
// 			} else {
// 				fmt.Println("other cases: ", line)
// 			}
// 		}
// 	}
// 	fmt.Println("Binaries")
// 	codeStrr = binary.String()
// 	fmt.Println(codeStrr)
// 	codeBin := []byte(codeStrr)
// 	if err := ioutil.WriteFile("Add.hack", codeBin, 0777); err != nil {
// 		log.Fatal(err)
// 	}
// }
//
// func assembleAddress(line string, sb strings.Builder) string {
// 	i, err := strconv.Atoi(strings.TrimSpace(line[1:]))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	b := completeBitsFront(strconv.FormatInt(int64(i), 2), 16)
// 	// fmt.Println(reflect.TypeOf(b))
// 	sb.WriteString(b)
// 	sb.WriteString("\n")
//
// 	return sb.String()
// }
//
// func assembleJump(line string, sb strings.Builder, cInstruction map[string]string, cJump map[string]string) string {
// 	sb.WriteString("111")
//
// 	c := ii.FindAllStringSubmatch(line, -1)[0][1]
//
// 	sb.WriteString(string(cInstruction[c]))
//
// 	sb.WriteString(string("000"))
//
// 	j := ii.FindAllStringSubmatch(line, -1)[0][2]
//
// 	sb.WriteString(string(cJump[j]))
// 	sb.WriteString("\n")
//
// 	return sb.String()
// }
//
// func assembleNoJump(line string, sb strings.Builder, cInstruction map[string]string, cDestination map[string]string) string {
// 	sb.WriteString("111")
//
// 	c := i.FindAllStringSubmatch(line, -1)[0][2]
// 	sb.WriteString(string(cInstruction[c]))
//
// 	d := i.FindAllStringSubmatch(line, -1)[0][1]
// 	sb.WriteString(string(cDestination[d]))
//
// 	final := completeBitsBack(sb.String(), 16)
//
// 	return final
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
// 		return s + "\n"
// 	} else {
// 		return s + "\n"
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
