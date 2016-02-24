package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"regexp"

	"io"

	"github.com/bborbe/log"
	"github.com/bborbe/stringutil"
)

type entry struct {
	date  string
	from  string
	until string
}

var logger = log.DefaultLogger
var PARAMETER_LOGLEVEL = "loglevel"

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	writer := os.Stdout
	err := do(writer)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	entries, err := read(fmt.Sprintf("%s/Library/Logs/ping.log", u.HomeDir))
	if err != nil {
		return err
	}
	for _, i := range entries {
		fmt.Fprintf(writer, "%s: %s - %s\n", i.date, i.from, i.until)
	}
	return nil
}

func read(filename string) ([]*entry, error) {
	entries := make([]*entry, 0)
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2})\\s(.*)")
	var e *entry
	reader := bufio.NewReader(fi)
	for {
		read, err := reader.ReadString('\n')
		if err != nil {
			return entries, nil
		}
		line := stringutil.Trim(read)
		result := re.FindStringSubmatch(line)
		if len(result) > 0 {
			date := result[1]
			time := result[2]
			place := result[3]
			if place == "seibert-media" {
				if e == nil || e.date != date {
					e = &entry{date: date, from: time, until: time}
					entries = append(entries, e)
				} else {
					e.until = time
				}
			}
		}
	}
	return nil, nil
}
