## To Dos

- Some github workflows and testing suites.

- Provide an overwrite flag and prompt user if the flag is not provided and a file is found.

- Shorthand flags. Apparently golang doesn't like them.

- Handle potential responses from advent of code answer submittions

- Improve puzzle display. Currently it constrains the read to 80 characters but code snippets and input examples do not display in the same format as the website.

- Cleaning up error handling. Perhaps define them in a separate package though that might not be appropriate for the language.

- Provide a subcommand to initialize, parse, and manipulate a config file so workspace setup can be automated

- Add timing functionality to set up puzzle and input files at puzzle release time

- Potential Subcommands:
  - `help`                  provide a help dialog for subcommands and flags
  - `calendar`              partialy implemented. The display for html2text does not work well with aoc.
  - `read`                  partially implemented. See above
  - `personal-stats`        show personal stats(https://adventofcode.com/2023/leaderboard/self).
  - `private-leaderboard`   show a leaderboard given an id number