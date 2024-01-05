package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/myscon/aoc/client"
)

type BaseArgs struct {
	year        *int
	day         *int
	sessionPath *string
}

type ReadArgs struct {
	page *string
}

type DownloadArgs struct {
	out       *string
	puzzle    *bool
	overwrite *bool
}

type SubmitArgs struct {
	part   *string
	answer *string
}

func defineBaseArgs() (*argparse.Parser, *BaseArgs) {
	baseArgs := &BaseArgs{}
	est, _ := time.LoadLocation("America/New_York")
	currentTime := time.Now().In(est)

	parser := argparse.NewParser("aoc", "Command line interface for Advent of Code.")
	baseArgs.year = parser.Int("y", "year", &argparse.Options{Help: "Provide an optional year for puzzle", Default: getYear(currentTime)})
	baseArgs.day = parser.Int("d", "day", &argparse.Options{Help: "Provide an option day for puzzle", Default: getDayOfDecember(currentTime)})
	baseArgs.sessionPath = parser.String("s", "session-path", &argparse.Options{Help: "Provide a session cookie file path"})
	return parser, baseArgs
}

func defineReadArgs(parser *argparse.Parser) (*argparse.Command, *ReadArgs) {
	readArgs := &ReadArgs{}
	readCmd := parser.NewCommand("read", "Read the information for a given puzzle.")
	readArgs.page = readCmd.Selector("i", "info", []string{"puzzle", "input", "calendar"}, &argparse.Options{Help: "Choose what part of AOC to display", Default: "puzzle"})

	return readCmd, readArgs
}

func defineDownloadArgs(parser *argparse.Parser) (*argparse.Command, *DownloadArgs) {
	downloadArgs := &DownloadArgs{}
	downloadCmd := parser.NewCommand("download", "Download the information for a given puzzle.")
	downloadArgs.out = downloadCmd.String("o", "output", &argparse.Options{Help: "Output file path for download"})
	downloadArgs.puzzle = downloadCmd.Flag("z", "puzzle", &argparse.Options{Help: "Download puzzle description instead of puzzle input", Default: false})
	downloadArgs.overwrite = downloadCmd.Flag("w", "overwrite", &argparse.Options{Help: "Overwrite file at path", Default: false})
	return downloadCmd, downloadArgs
}

func defineSubmitArgs(parser *argparse.Parser) (*argparse.Command, *SubmitArgs) {
	submitArgs := &SubmitArgs{}
	submitCmd := parser.NewCommand("submit", "Submit answers for a given puzzle.")
	submitArgs.part = submitCmd.Selector("p", "part", []string{"1", "2"}, &argparse.Options{Help: "Specify part to submit answer", Default: "1"})
	submitArgs.answer = submitCmd.String("a", "answer", &argparse.Options{Help: "Answer to puzzle", Required: true})
	return submitCmd, submitArgs
}

func getYear(currentTime time.Time) int {
	if currentTime.Month() == time.December {
		return currentTime.Year()
	} else {
		return currentTime.Year() - 1
	}
}

func getDayOfDecember(currentTime time.Time) int {
	if currentTime.Month() == time.December && currentTime.Day() <= 25 {
		return currentTime.Day()
	} else {
		return 25
	}
}

func main() {
	parser, baseArgs := defineBaseArgs()

	readCmd, readArgs := defineReadArgs(parser)
	downloadCmd, downloadArgs := defineDownloadArgs(parser)
	submitCmd, submitArgs := defineSubmitArgs(parser)
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}
	client, err := client.NewAocClient(*baseArgs.year, *baseArgs.day, baseArgs.sessionPath, false)
	if err != nil {
		log.Fatalln(err)
	}
	if readCmd.Happened() {
		switch *readArgs.page {
		case "input":
			err = client.ShowInput()
		case "puzzle":
			err = client.ShowPuzzle()
		case "calendar":
			err = client.ShowCalendar()
		}
		if err != nil {
			log.Println(err)
		}
	} else if downloadCmd.Happened() {
		if *downloadArgs.puzzle {
			err = client.SavePuzzle(downloadArgs.out)
		} else {
			err = client.SaveInput(downloadArgs.out)
		}
		if err != nil {
			log.Println(err)
		}
	} else if submitCmd.Happened() {
		if submitArgs.answer == nil {
			log.Fatalln(errors.New("submit: please provide an answer"))
		}
		res, err := client.SubmitAnswer(submitArgs.part, submitArgs.answer)
		if err == nil {
			// TODO: graciously handle potential responses
			fmt.Println(res)
		}
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Fatalln("Command parsing fatal error.")
	}
}
