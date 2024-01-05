## AOC

`aoc` is a cli to conveniently interact with [Advent of Code](https://adventofcode.com). This is a barebones project thus far written in a lazy winter afternoon and will need some much needed refinement.

## Features

- Load Advent of Code via session cookie from file or env variable
- Read puzzle description and optionally save to a given file path
- Read puzzle input and optionally save to a given file path
- Submit puzzle answers

## Getting Started

### Download and Install

- Install go version > 1.2 from [go.dev](https://go.dev/dl/) then run the install command

```
# go install github.com/myscon/aoc
```

### Session Cookie

In order to access your unique input for the day's puzzle, you must provide a session cookie value. Visit the [Advent of Code](https://adventofcode.com), log in, and inspect the value for the `session` cookie that is stored in your browser.

You may provide the cookie by one of several ways. The client will attempt to retrieve the session cookie in the same order as below and will use the first one found.

1. In a file specified via the `-session-file` command line option.
2. In an `AOC_SESSION` environment variable.
3. In a file called `.adventofcode.session` (note the dot) in your home
   directory (`/home/username` on Linux, `C:\Users\Username` on Windows,
   `/Users/Username` on macOS).
4. In a file called `adventofcode.session` (no dot) in your user's config
   directory (`/home/username/.config` on Linux, `C:\Users\Username\AppData\Roaming`
   on Windows, `/Users/Username/Library/Application Support` on macOS).

## Usage

### Flags and Subcommands

```
usage: aoc <Command> [-h|--help] [-y|--year <integer>] [-d|--day <integer>]
           [-s|--session-path "<value>"]

           Command line interface for Advent of Code.

Commands:

  read      Read the information for a given puzzle.
  download  Download the information for a given puzzle.
  download  Submit answers for a given puzzle.

Arguments:

  -h  --help          Print help information
  -y  --year          Provide an optional year for puzzle. Default: 2023
  -d  --day           Provide an option day for puzzle. Default: 25
  -s  --session-path  Provide a session cookie file path

Read Arguments:
  -i  --info          Choose what part of AOC to display. Default: puzzle
                      (puzzle|input|calendar)

Download Arguments:
  -o  --output        Output file path for download
  -z  --puzzle        Download puzzle description instead of puzzle input.
                      Default: false
  -w  --overwrite     Overwrite file at path. Default: false

Submit Arguments:
  -p  --part          Specify part to submit answer (1|2). Default: 1
  -a  --answer        Answer to puzzle
```

### Show puzzle

```
# aoc read
```

### Show input

```
# aoc read --info input
```

### Save input to path

```
# aoc download -o ./day_2_input
```

### Save input from previous year and day

```
# aoc download --year 2022 --day 22
```

## Credits

If you haven't noticed already, there are some similarities to [aoc-cli](https://github.com/scarvalhojr/aoc-cli/). The reason I started this project was because I couldn't find a comparable package in go and thought it would be a fun thing to do over winter break.

## Contributing

There are quite a few holes in the project that I'd like to complete and some feedback would be much appreciated. Please see [CONTRIBUTING](CONTRIBUTING.md) for more information.

## Support Advent of Code

[Advent of Code](https://adventofcode.com) is an Advent calendar of small programming puzzles for a variety of skill sets and skill levels that can be solved in any programming language you like. People use them as interview prep, company training, university coursework, practice problems, a speed contest, or to challenge each other. The project is maintained by [Eric Wastl](http://was.tl/) and some volunteers outside normal day jobs and family time. Please [support their work](https://adventofcode.com/support).