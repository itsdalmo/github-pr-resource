package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	dir   string
	f     *os.File
	debug bool
)

func init() {
	dir = defaultDir
	if os.Getenv("LOG_DIRECTORY") != "" {
		dir = os.Getenv("LOG_DIRECTORY")
	}

	if os.Getenv("LOG_DEBUG") != "" || true {
		debug = true
	}

	cleanStaleLogs(dir)

	f, err := os.OpenFile(fmt.Sprintf("%s/resource-%s.log", dir, time.Now().Format("2006-01-02")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(f)
}

// WriteStdin will write the contents of stdin to the log, then return the contents if env var `LOG_DEBUG` != ""
func WriteStdin() []byte {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if debug {
		Write(string(input))
	}

	return input
}

// Write writes a message to a log file in `/tmp`
func Write(msg string) {
	log.Println(msg)
}

// Close the *os.File connection for the logger
func Close() {
	f.Close()
}
