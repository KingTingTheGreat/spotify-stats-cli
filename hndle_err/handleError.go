package hndle_err

import (
	"os"
	"spotify-stats-cli/util"
)

func EndWithErr(errorMessage string) {
	os.Stderr.WriteString("Error: " + errorMessage + "\n")

	outputFile, err := util.WriteOutputToFile("", "")
	if err != nil {
		os.Stderr.WriteString("Cannot write to file\n")
		return
	}

	os.Stdout.WriteString(outputFile)
	os.Exit(1)
}
