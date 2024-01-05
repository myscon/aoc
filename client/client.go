/*
Package client provides a simple http client for advent of code
*/
package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/k3a/html2text"
	"github.com/mitchellh/go-wordwrap"
)

const firstPuzzleYear = 2015
const firstPuzzleDay = 1
const lastPuzzleDay = 25
const baseURL = "https://adventofcode.com"

type AocClient struct {
	year      string
	day       string
	overwrite bool
	session   *string
	calendar  *string
	puzzle    *string
	input     *string
}

func (aoc *AocClient) ShowCalendar() error {
	if aoc.calendar == nil {
		if err := aoc.getCalendar(); err != nil {
			return fmt.Errorf("ShowCalendar: failed to get calendar: %w", err)
		}
	}
	// TODO: write a formatter to make a pretty tree
	fmt.Println(wordwrap.WrapString(html2text.HTML2Text(*aoc.calendar), 80))
	return nil
}

func (aoc *AocClient) ShowPuzzle() error {
	if aoc.puzzle == nil {
		if err := aoc.getPuzzle(); err != nil {
			return fmt.Errorf("ShowPuzzle: failed to get puzzle: %w", err)
		}
	}
	// // TODO: write a formatter to match the presentation online
	fmt.Println(wordwrap.WrapString(html2text.HTML2Text(*aoc.puzzle), 80))
	return nil
}

func (aoc *AocClient) SavePuzzle(filePath *string) error {
	if aoc.puzzle == nil {
		if err := aoc.getPuzzle(); err != nil {
			return fmt.Errorf("SavePuzzle: failed to get puzzle: %w", err)
		}
	}
	return aoc.saveFile(filePath, *aoc.puzzle)
}

func (aoc *AocClient) ShowInput() error {
	if aoc.puzzle == nil {
		if err := aoc.getInput(); err != nil {
			return fmt.Errorf("ShowInput: failed to get input: %w", err)
		}
	}
	fmt.Println(*aoc.input)
	return nil
}

func (aoc *AocClient) SaveInput(filePath *string) error {
	if aoc.puzzle == nil {
		if err := aoc.getInput(); err != nil {
			return fmt.Errorf("SaveInput: failed to get input: %w", err)
		}
	}
	return aoc.saveFile(filePath, *aoc.input)
}

func (aoc *AocClient) SubmitAnswer(part *string, answer *string) (string, error) {
	targetURL, _ := url.JoinPath(baseURL, aoc.year, "day", aoc.day)
	if !(*part == "1" || *part == "2") {
		return "", fmt.Errorf("SubmitAnswer: invalid part number {%s}. please enter 1 or 2", *part)
	}
	data := url.Values{}
	data.Set("part", *part)
	data.Set("answer", *answer)
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("SubmitAnswer: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := aoc.clientBuilder().Do(req)
	if err != nil {
		return "", fmt.Errorf("SubmitAnswer: failed to perform request: %w", err)
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		return parseBody(res.Body)
	default:
		return "", fmt.Errorf("SubmitAnswer: failed to post request with error code {%d}", res.StatusCode)
	}
}

func (aoc *AocClient) getCalendar() error {
	targetURL, _ := url.JoinPath(baseURL, aoc.year)
	if calendar, err := aoc.httpGet(targetURL); err == nil {
		aoc.calendar = &calendar
		return err
	} else {
		return fmt.Errorf("getCalendar: failed for url {%s} %w", targetURL, err)
	}
}

func (aoc *AocClient) getPuzzle() error {
	targetURL, _ := url.JoinPath(baseURL, aoc.year, "day", aoc.day)
	if puzzle, err := aoc.httpGet(targetURL); err == nil {
		aoc.puzzle = &puzzle
		return err
	} else {
		return fmt.Errorf("getPuzzle: failed for url {%s} %w", targetURL, err)
	}
}

func (aoc *AocClient) getInput() error {
	targetURL, _ := url.JoinPath(baseURL, aoc.year, "day", aoc.day, "input")
	if input, err := aoc.httpGet(targetURL); err == nil {
		aoc.input = &input
		return err
	} else {
		return fmt.Errorf("getInput: failed for url {%s}: %w", targetURL, err)
	}
}

func (aoc *AocClient) saveFile(filePath *string, content string) error {
	if aoc.overwrite {
		if err := writeToFile(*filePath, content); err != nil {
			return fmt.Errorf("saveFile: failed to overwrite to {%s} with error: %w", *filePath, err)
		}
	} else {
		if _, err := os.Stat(*filePath); os.IsNotExist(err) {
			err := writeToFile(*filePath, content)
			if err != nil {
				return fmt.Errorf("saveFile: failed to create/write to {%s} with error: %w", *filePath, err)
			}
		} else if err == nil {
			return fmt.Errorf("saveFile: file {%s} already exists", *filePath)
		} else {
			return fmt.Errorf("saveFile: unexpected error with filepath {%s}: %w", *filePath, err)
		}
	}
	return nil
}

func (aoc *AocClient) getSession(sessionPath *string) error {
	// TODO: perform more validation checks on the value (ie whitespace)
	if sessionPath != nil && *sessionPath != "" {
		println(*sessionPath, "2")
		session, err := os.ReadFile(*sessionPath)
		if err == nil {
			text := string(session)
			aoc.session = &text
		}
		return err
	}
	if session := os.Getenv("AOC_SESSION"); session != "" {
		aoc.session = &session
		return nil
	}
	homeDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	homePath := filepath.Join(homeDir, ".adventofcode.session")
	if _, err := os.Stat(homePath); err == nil {
		session, err := os.ReadFile(homePath)
		if err != nil {
			return fmt.Errorf("getSession: error reading filepath {%s}: %w", homePath, err)
		}
		text := string(session)
		aoc.session = &text
		return err
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "adventofcode.session")
	if _, err := os.Stat(configPath); err == nil {
		session, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("getSession: error reading filepath {%s}: %w", configPath, err)
		}
		text := string(session)
		aoc.session = &text
		return err
	}
	return errors.New("getSession: unable to find session cookie")
}

func (aoc *AocClient) clientBuilder() *http.Client {
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse("https://adventofcode.com")
	cookies := []*http.Cookie{
		{Name: "session", Value: *aoc.session},
	}
	jar.SetCookies(u, cookies)
	return &http.Client{
		Jar: jar,
	}
}

func (aoc *AocClient) httpGet(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	res, err := aoc.clientBuilder().Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		return parseBody(res.Body)
	default:
		return "", fmt.Errorf("httpGet: get request failed with error code %d", res.StatusCode)
	}
}

func writeToFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err == nil {
		defer file.Close()
		_, err = file.WriteString(content)
		return err
	}
	return err
}

func parseBody(body io.ReadCloser) (string, error) {
	re := regexp.MustCompile(`(?i)(?s)<main>(?P<main>.*)</main>`)
	if text, err := io.ReadAll(body); err == nil {
		match := re.Find(text)
		if len(match) < 1 {
			return string(text), err
		}
		return string(match), err
	} else {
		return "", err
	}
}

func verifyDay(year int, day int) bool {
	est, _ := time.LoadLocation("America/New_York")
	ninePM := time.Date(year, 12, day, 0, 0, 0, 0, est)
	return !(year < firstPuzzleYear ||
		year > time.Now().In(est).Year() ||
		day > lastPuzzleDay ||
		day < firstPuzzleDay ||
		time.Now().In(est).Before(ninePM))
}

func NewAocClient(year int, day int, sessionPath *string, greedy bool) (AocClient, error) {
	if verifyDay(year, day) {
		newClient := AocClient{
			year:      fmt.Sprint(year),
			day:       fmt.Sprint(day),
			overwrite: false,
		}
		if err := newClient.getSession(sessionPath); err != nil {
			return AocClient{}, fmt.Errorf("NewAocClient: failed to get session cookie with path {%s}: %w", *sessionPath, err)
		}
		if greedy {
			if err := newClient.getPuzzle(); err != nil {
				return AocClient{}, fmt.Errorf("NewAocClient: failed to get puzzle description: %w", err)
			}

			if err := newClient.getInput(); err != nil {
				return AocClient{}, fmt.Errorf("NewAocClient: failed to get puzzle input: %w", err)
			}
		}
		return newClient, nil
	} else {
		return AocClient{}, fmt.Errorf("NewAocClient: puzzle for year {%d} and day {%d} is not available", year, day)
	}
}
