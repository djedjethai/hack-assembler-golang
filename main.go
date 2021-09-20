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
)

func main() {

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
	firstParseSetVariousVar(progText)

	// second parsing set the binary file
	secondParseSetBinaries(progText)

	fmt.Println("Binaries")
	codeStrr := binary.String()
	codeBin := []byte(codeStrr)
	if err := ioutil.WriteFile("Add.hack", codeBin, 0777); err != nil {
		log.Fatal(err)
	}
}

func firstParseSetVariousVar(progText string) {
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
}

func secondParseSetBinaries(progText string) {
	cInstruction := tableParser("./tables/cInstructions.txt")
	cDestination := tableParser("./tables/cInstDestination.txt")
	cJump := tableParser("./tables/cJump.txt")
	rVar := tableParser("./tables/rVariable.txt")
	thisThat := tableParser("./tables/thisThatTable.txt")

	for _, line := range strings.Split(progText, "\n") {
		var v string
		var b string

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
}
