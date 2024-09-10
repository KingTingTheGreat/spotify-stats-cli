package util

import (
	"image"
	"net/http"
	"os"
	"path/filepath"
	"spotify-stats-cli/cnsts"
	"spotify-stats-cli/conv"
	"strings"
)

func BasePath() string {
	execPath, err := os.Executable()
	if err != nil {
		EndWithErr("cannot get executable path")
	}

	execDir := filepath.Dir(execPath)
	cwd, err := os.Getwd()
	if err != nil {
		EndWithErr("cannot get current working directory")
	}

	if filepath.Base(execPath) == "go.exe" || filepath.Base(execPath) == "go" && strings.HasPrefix(cwd, os.TempDir()) {
		return execDir
	}

	return cwd
}

func AnsiImage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		EndWithErr("cannot get image from URL")
	}

	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		EndWithErr("cannot decode image")
	}

	ansiImage := conv.Convert(img)

	return ansiImage
}

func formatTrackText(trackText string) string {
	if len(trackText) < cnsts.DIM {
		totalPadding := cnsts.DIM - len(trackText)
		leftPadding := totalPadding / 2
		trackText = strings.Repeat(" ", leftPadding) + trackText + strings.Repeat(" ", totalPadding-leftPadding)
	}
	return trackText
}

func WriteOutputToFile(ansiImage string, trackText string) string {
	outputFile := cnsts.OUTPUT_FILE_NAME
	basePath := BasePath()

	file, err := os.Create(basePath + "\\" + outputFile)
	if err != nil {
		file, err = os.Create(basePath + "/" + outputFile)
		if err != nil {
			os.Stderr.WriteString("Cannot create output file\n")
			os.Exit(1)
		} else {
			outputFile = basePath + "/" + outputFile
		}
	} else {
		outputFile = basePath + "\\" + outputFile
	}

	defer file.Close()

	if ansiImage == "" || trackText == "" {
		_, err = file.WriteString("")
		if err != nil {
			os.Stderr.WriteString("Cannot write to file\n")
			os.Exit(1)
		}
	} else {
		_, err = file.WriteString(ansiImage + "\n" + formatTrackText(trackText) + "\n\n")
		if err != nil {
			os.Stderr.WriteString("Cannot write to file\n")
			os.Exit(1)
		}
	}

	return outputFile
}
