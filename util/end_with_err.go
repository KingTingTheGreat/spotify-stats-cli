package util

import (
	"os"
)

func EndWithErr(errorMessage string) {
	os.Stderr.WriteString("Error: " + errorMessage + "\n")

	outputFile := WriteOutputToFile("", "")

	os.Stdout.WriteString(outputFile)
	os.Exit(1)
}
