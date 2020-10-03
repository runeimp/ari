package main

import (
	"fmt"
	"os"
	"strconv"
)

/*
 * CONSTANTS
 */
const (
	AppName                = "Ari"
	AppVersion             = "0.1.0-alpha"
	DataFormatISO8601RegEx = `\d\d\d\d-?\d\d-?\d\d\D\d\d:?\d\d:?\d\d(\.\d*)?(GMT|z|Z)`
	PriorityFields         = []string{"timestamp", "thread", "level", "logger", "message", "msg", "error", "exception"}
)

/*
 * DERIVED CONSTANTS
 */
var (
	AppHelp = AppLabel + `

USAGE: ari [OPTIONS] JSONLOGFILE [JSONLOGFILE JSONLOGFILE ...]

OPTIONS:
  -d | -depth   | --depth           Max directory depth to decend
                                    (Default: 0 [unlimited])
  -h | -help                        Help info for this program
  -r | -raw     | -raw-numbers      Do not convert numbers to be more
                                    human friendly
  -Q | -quiet   | --quiet           Verbosity 0, errors only
  -T | -terse   | --terse           Verbosity 1, limited field output
  -V | -verbose | --verbose         Verbosity 2, all visual output
  -v | -ver     | -version          Version info for this program

Note that the longest options can be prefixed with a single hyphen per the
Multics standard or double hyphen per the GNU standard.
`
	AppLabel = fmt.Sprintf("%s v%s", AppName, AppVersion)
)

/*
 * VARIABLES
 */
var (
	humanNumbers       = true
	verbosity    uint8 = 1

	logFiles []string
	maxDepth uint64
)

func processLog(log string) {
	//
}

/*
 * MAIN ENTRYPOINT
 */
func main() {
	// fmt.Println(AppLabel)
	// userPass = os.Getenv("SFTPCMP_REMOTEPASS")

	// fmt.Printf("os.Args: %q\n", os.Args)
	skip := 0
	for i, arg := range os.Args {
		if i != skip {
			switch arg {
			case "-d", "-depth", "--depth":
				j := i + 1
				maxDepth, _ = strconv.ParseUint(os.Args[j], 10, 32)
				// fmt.Printf("arg[%d]: %s\n", i, arg)
				// fmt.Printf("Port set to %d\n", hostPort)
				skip = j
			case "-h", "-help", "--help":
				fmt.Println(AppHelp)
				os.Exit(0)
			case "-r", "-raw", "-raw-numbers", "--raw-numbers":
				humanNumbers = false
			case "-v", "-ver", "-version", "--version":
				fmt.Println(AppLabel)
				os.Exit(0)
			case "-verbosity", "--verbosity":
				j := i + 1
				verbosityInt, verbosityError := strconv.ParseInt(os.Args[j], 10, 32)
				if verbosityError != nil {
					fmt.Printf("Verbosity Error: %s\n", verbosityError.Error())
				} else {
					if verbosityInt > 2 {
						fmt.Fprintln(os.Stderr, "Verbosity is max 2")
						verbosity = 2
					} else if verbosityInt < 0 {
						fmt.Fprintln(os.Stderr, "Verbosity is min 0")
						verbosity = 0
					} else {
						verbosity = uint8(verbosityInt)
					}
				}
			case "-Q", "-quiet", "--quiet":
				verbosity = 0
			case "-T", "-terse", "--terse":
				verbosity = 1
			case "-V", "-verbose", "--verbose":
				verbosity = 2
			default:
				// fmt.Printf("arg[%d]: %s\n", i, arg)
				if len(arg) > 0 {
					logFiles = append(logFiles, arg)
				} else {
					fmt.Printf("What's this? %q\n", arg)
				}
			}
		}
	} // end range os.Args

	if len(logFiles) == 0 {
		fmt.Fprintln(os.Stderr, "No log files provided to display. Use the -h option for help.")
		os.Exit(10)
	}
	for _, log := range logFiles {
		fmt.Printf("Processing: %q\n", log)
	}
}
