package handle

import (
	"log"
	"os"
	"os/user"
	"runtime/debug"

	"nenvoy.com/pkg/utils/printing"
)

// Error - handles the error by placing the stack trace into a log file and printing the error
func Error(inerr error) {
	if inerr != nil {
		//If the directiory doesnt exit, make it
		err := os.MkdirAll("/tmp/nenvoy/", 0777)
		if err != nil {
			log.Fatal(err)
		}

		//create log file file with desired read/write permissions if it doesnt already exist
		f, err := os.OpenFile("/tmp/nenvoy/nenvn.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			log.Fatal(err)
		}

		//defer the file closure
		defer f.Close()

		//set output of logs to f
		log.SetOutput(f)

		// Get the user who ran the program
		userAccount, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		// Generate the log entry
		logEntry := "ERROR, " + userAccount.Username + ", " + inerr.Error() + "\n" + string(debug.Stack()) + "\n#########################################################\n"

		//Write the stack trace
		log.Println(logEntry)

		printing.PrintError(inerr.Error())
	}
}
