package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	//"testing/quick"
)

type Instruction struct {
	typeofInstruction string
	rawInstruction    string //binary input converted to string
	linevalue         uint64 //
	programCnt        int
	opcode            uint64 //needs to be 64 bit to apply mask and shift
	op                string //whether its B, I, BREAK, etcrd
	rd                uint8
	rn                uint8
	rm                uint8
	shamt             uint8 //6 bits, 0-63
}

type Register struct {
	registerValue int //value stored in register
}

type Snapshot struct {
	cycle    int
	insLabel string
	regis    [32]Register //snapshot of current register values, updates every instruction that changes them
}

var InputParsed []Instruction
var SnapshotArray []Snapshot

var Breaknow = false
var InputFileName *string
var OutputFileName *string

//var OutputFileName2 *string

func readInstruction(filePath string) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}

	// Close file after main runs
	defer file.Close()
	// program counter
	var pc = 96

	// Read in file with scanner
	//// InputParsed := []Instruction{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ns := Instruction{rawInstruction: scanner.Text(), programCnt: pc}
		InputParsed = append(InputParsed, ns)
		pc += 4

	}
	if err := scanner.Err(); err != nil {
		fmt.Print(err)
	}
}

func processInput(list []Instruction) {
	for i := 0; i < len(list); i++ {
		convertToInt(&list[i])
		opcodeMasking(&list[i])
		opcodeMatching(&list[i])
		RTypeFormat(&list[i])
	}
}

/*func writeInstruction(filePath string, list []Instruction) {
	f, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for i := 0; i < len(list); i++ {
		switch list[i].typeofInstruction {
		case "R":
			// Prints out bits
			_, err := fmt.Fprintf(f, "%s %s %s %s %s\t", list[i].rawInstruction[0:11],
				list[i].rawInstruction[11:16], list[i].rawInstruction[16:22],
				list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])
			if err != nil {
				log.Fatal(err)
			}
			// Prints out pc and opcode instruction
			_, err = fmt.Fprintf(f, "%d\t%s\t", list[i].programCnt, list[i].op)
			// Prints out rd and rn
			_, err = fmt.Fprintf(f, "R%d, R%d, ", list[i].rd, list[i].rn)
			_, err = fmt.Fprintf(f, "R%d\n", list[i].rm)
		}
	}
}
*/
/*
	func writeSimulator(filePath string, list []Register) {
		f, err := os.Create(filePath)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		_, err := fmt.Fprintf(f, "%s %s %s %s %s\t", list[i].rawInstruction[0:11],
			list[i].rawInstruction[11:16], list[i].rawInstruction[16:22],
			list[i].rawInstruction[22:27], list[i].rawInstruction[27:32])
		if err != nil {
			log.Fatal(err)
				}
	}
*/
func convertToInt(ins *Instruction) {
	i, err := strconv.ParseUint(ins.rawInstruction, 2, 64)
	//fmt.Println(i)
	if err == nil {
		ins.linevalue = i
	} else {
		fmt.Println(err)
	}
}

func opcodeMasking(ins *Instruction) {
	ins.opcode = (ins.linevalue & 4292870144) >> 21
}

func opcodeMatching(ins *Instruction) {
	if ins.opcode == 1112 {
		ins.op = "ADD"
		ins.typeofInstruction = "R"
	} else {
		fmt.Println("Error")
	}
}

func RTypeFormat(ins *Instruction) {
	ins.rm = uint8((ins.linevalue & 2031616) >> 16)
	ins.shamt = uint8((ins.linevalue & 64512) >> 10)
	ins.rn = uint8((ins.linevalue & 992) >> 5)
	ins.rd = uint8(ins.linevalue & 31)
}

func processSnapshot(list []Instruction, array []Snapshot) {
	for i := 0; i < len(list); i++ {
		initializeRegisters(&array[i])
	}
	for i := 0; i < len(list); i++ {
		updateRegisters(&list[i], &array[i])
	}
}

func initializeRegisters(array *Snapshot) {
	for i := 0; i < len(array.regis); i++ {
		array.regis[i].registerValue = 0
	}
}

func updateRegisters(ins *Instruction, array *Snapshot) {
	if ins.op == "ADD" {
		array.regis[ins.rd].registerValue = array.regis[ins.rd].registerValue + array.regis[ins.rd].registerValue
	}
}

func main() {
	InputFileName := flag.String("i", "input_.txt", "Gets the input file name")
	//OutputFileName := flag.String("o", "out_.dis", "Gets the output file name")
	//OutputFileName2 := flag.String("o2", "Ro_.sim", "Gets the output2 file name")

	flag.Parse()
	readInstruction(*InputFileName)
	processInput(InputParsed)
	processSnapshot(InputParsed, SnapshotArray)
	//writeInstruction(*OutputFileName, InputParsed)
	//writeSimulator(*OutputFileName2, RegisterArray)
}
