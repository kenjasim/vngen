package printing

import (
	"fmt"
)

const (
	InfoColour    = "\033[1;34m%s\033[0m\n"
	SuccessColour = "\033[1;32m%s\033[0m\n"
	NoticeColour  = "\033[1;36m%s\033[0m\n"
	WarningColour = "\033[1;33m%s\033[0m\n"
	ErrorColour   = "\033[1;31m%s\033[0m\n"
	DebugColour   = "\033[0;36m%s\033[0m\n"
)

//PrintError - Prints a Formatted Error Message
func PrintError(str string) {
	fmt.Printf(ErrorColour, "[!] "+str)
}

//PrintSuccess - Prints a formatted success message
func PrintSuccess(str string) {
	fmt.Printf(SuccessColour, "[✓] "+str)
}

//PrintInfo - Prints a formatted info message
func PrintInfo(str string) {
	fmt.Printf(InfoColour, "[i] "+str)
}

//PrintWarning - Prints a formatted warning message
func PrintWarning(str string) {
	fmt.Printf(WarningColour, "[?] "+str)
}

//SprintError - Prints a Formatted Error Message
func SprintError(str string) (fstr string) {
	return fmt.Sprintf(ErrorColour, "[!] "+str)
}

//SprintSuccess - Prints a formatted success message
func SprintSuccess(str string) (fstr string) {
	return fmt.Sprintf(SuccessColour, "[✓] "+str)
}

//SprintInfo - Prints a formatted info message
func SprintInfo(str string) (fstr string) {
	return fmt.Sprintf(InfoColour, "[i] "+str)
}

//SprintWarning - Prints a formatted warning message
func SprintWarning(str string) (fstr string) {
	return fmt.Sprintf(WarningColour, "[?] "+str)
}
