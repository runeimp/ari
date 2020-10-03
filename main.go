package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
 * CONSTANTS
 */
const (
	AppName                = "Ari"
	AppVersion             = "0.1.0-alpha"
	DataFormatISO8601RegEx = `\d\d\d\d-?\d\d-?\d\d\D\d\d:?\d\d:?\d\d(\.\d*)?(GMT|z|Z)`
	FileSizeKiloByte       = 1024
	FileSizeMegaByte       = 1048576
	FileSizeGigaByte       = 1073741824
	TimeFormatIso8601      = "2006-01-02 15:04:05 UTC"
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
	AppLabel       = fmt.Sprintf("%s v%s", AppName, AppVersion)
	PriorityFields = make(map[int]priorityField)
)

func init() {
	priorityFields := []string{"timestamp", "thread", "level", "logger", "message", "msg", "error", "exception"}

	for i, key := range priorityFields {
		formated := strings.Title(key)
		switch key {
		case "lvl":
			formated = "Level"
		case "msg":
			formated = "Message"
		case "timestamp":
			formated = "TimeStamp"
		}
		PriorityFields[i] = priorityField{
			Key:      key,
			Formated: formated,
			Order:    i,
		}
	}
}

/*
 * DATA TYPES
 */

// DateTime represents a string that is for data time display
type DateTime string

// JSONDatum represents a key/value field in a JSON object
type JSONDatum struct {
	OriginalKey string
	viewKey     string
	Value       interface{}
	valueType   string
}

// Key returns the current key for the datum
func (jd JSONDatum) Key() string {
	if len(jd.viewKey) > 0 {
		return jd.viewKey
	}
	return jd.OriginalKey
}

// SetKey sets the current key for the datum
func (jd JSONDatum) SetKey(k string) {
	jd.viewKey = k
}

// ValueType returns the value type for the datum
func (jd JSONDatum) ValueType() string {
	return jd.valueType
}

// NewJSONDatum returns a new JSONDatum instance
func NewJSONDatum(key string, value interface{}) JSONDatum {
	return JSONDatum{
		OriginalKey: key,
		Value:       value,
		valueType:   fmt.Sprintf("%T", value),
	}
}

type jsonMessage struct {
	TimeStamp int64  `json:"timestamp"`
	Thread    int    `json:"thread"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Error     string `json:"error"`
	Exception string `json:"exception"`
	Logger    string `json:"logger"`
}

type priorityField struct {
	Key      string
	Formated string
	Order    int
}

/*
 * VARIABLES
 */
var (
	humanNumbers       = true
	verbosity    uint8 = 1

	logFiles []string
	maxDepth uint64
)

/*
 * FUNCTIONS
 */

func limitLength(s string, length int) string {
	if len(s) < length {
		return s
	}
	return s[:length]
}

func printMessage(msgValues []JSONDatum) {
	var fieldsUsed []string

	i := 0

	// for i, field := range PriorityFields {
	for i < len(PriorityFields) {
		field := PriorityFields[i]
		// fmt.Printf("%s | %s | %d\n", field.Key, PriorityFields[field.Key].Formated, PriorityFields[field.Key].Order)
		for _, jd := range msgValues {
			if jd.OriginalKey == field.Key {
				msg := fmt.Sprintf("| %s: %s ", field.Formated, jd.Value)
				switch field.Key {
				case "level":
					fmt.Printf("%-17s", msg)
				case "modtime", "timestamp":
					fmt.Printf("%-37s", msg)
				default:
					fmt.Printf("%-30s", msg)
				}
				fieldsUsed = append(fieldsUsed, field.Key)
			}
		}
		i++
	}

	jsonOutput := ""
	for _, jd := range msgValues {
		matchFound := false
		for _, key := range fieldsUsed {
			if key == jd.OriginalKey {
				matchFound = true
				break
			}
		}
		if matchFound == false {
			// fmt.Printf("| %s: %v ", jd.Key(), jd.Value)
			if jd.ValueType() == "string" {
				jsonOutput += fmt.Sprintf("%q: %q, ", jd.Key(), jd.Value)
			} else {
				jsonOutput += fmt.Sprintf("%q: %v, ", jd.Key(), jd.Value)
			}
		}
	}
	if len(jsonOutput) > 0 {
		end := len(jsonOutput) - 2
		fmt.Printf("| Extra: {%s}", jsonOutput[:end])
	}
	fmt.Println()
}

func processLog(log string) {
	f, err := os.Open(log)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	defer f.Close()

	// fi, err := f.Stat()
	// if fi.Size() < (FileSizeMegaByte * 50) {
	// 	dat, err := ioutil.ReadFile(log)
	// }

	reader := bufio.NewReader(f)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		// Process the line here.
		// fmt.Printf(" > Read %d characters\n", len(line))
		// fmt.Printf(" > > %s\n", limitLength(line, 120))

		processMessage(line)

		if err != nil {
			break
		}
	}
	if err != io.EOF {
		fmt.Printf(" > Failed with error: %v\n", err)
		return
	}
}

func processMessage(msg string) {
	// var m jsonMessage
	if len(msg) > 0 {
		var m map[string]interface{}
		var msgValues []JSONDatum

		err := json.Unmarshal([]byte(msg), &m)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		// fmt.Printf("Message: %#v\n", m)

		levelSet := false
		errorSet := false
		for k, v := range m {
			// fmt.Printf("%s: %v\n", k, v)
			switch k {
			case "error":
				if v != nil && len(v.(string)) > 0 {
					// fmt.Printf("Error: %q", v)
					jd := NewJSONDatum(k, v)
					jd.SetKey("Error")
					errorSet = true
					msgValues = append(msgValues, jd)
				}
			case "level":
				levelSet = true
				l := strings.ToUpper(v.(string))
				// fmt.Printf("Level: %s", l)

				jd := NewJSONDatum(k, l)
				jd.SetKey("Level")
				msgValues = append(msgValues, jd)
			case "message", "msg":
				// fmt.Printf("Message: %q", v)

				jd := NewJSONDatum(k, v)
				jd.SetKey("Message")
				msgValues = append(msgValues, jd)
			case "modtime":
				t := processTime(v)
				// fmt.Printf("ModTime: %s", t)

				jd := NewJSONDatum(k, t)
				jd.SetKey("ModTime")
				msgValues = append(msgValues, jd)
			case "timestamp":
				t := processTime(v)
				// fmt.Printf("TimeStamp: %s", t)

				jd := NewJSONDatum(k, t)
				jd.SetKey("TimeStamp")
				msgValues = append(msgValues, jd)
			default:
				// fmt.Printf("%s: %v", k, v)
				jd := NewJSONDatum(k, v)
				// jd.SetKey("TimeStamp")
				msgValues = append(msgValues, jd)
			}
			// fmt.Println()
		}
		if levelSet == false {
			jd := NewJSONDatum("level", "INFO")
			if errorSet {
				// fmt.Println("Level: ERROR")
				jd = NewJSONDatum("level", "ERROR")
			} else {
				// fmt.Println("Level: INFO")
			}
			jd.SetKey("Level")
			msgValues = append(msgValues, jd)
		}
		// fmt.Println()
		printMessage(msgValues)
	}
}

func processTime(v interface{}) (result DateTime) {
	switch v.(type) {
	case string:
		matched, err := regexp.MatchString(DataFormatISO8601RegEx, v.(string))
		fmt.Printf("TimeStamp | matched = %t | err = %q | value = %v\n", matched, err)
	case float64:
		i := int64(v.(float64))
		t := time.Unix(i, 0).UTC().Format(TimeFormatIso8601)
		// fmt.Printf("TimeStamp | float64 = %f | int = %d | time = %s\n", v, i, t)
		// fmt.Printf("TimeStamp: %s", t)
		result = DateTime(fmt.Sprintf("%s", t))
	default:
		fmt.Printf("TimeStamp | type = %T | value = %v\n", v)
	}
	return result
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
		processLog(log)
	}
}
