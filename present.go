package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"io"
	"runtime"
	"github.com/bborbe/stringutil"
	"github.com/golang/glog"
)

type entry struct {
	date  string
	from  string
	until string
}

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout
	err := do(writer)
	if err != nil {
		glog.Exit(err)
	}
}

func do(writer io.Writer) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	entries, err := readFromFile(fmt.Sprintf("%s/Library/Logs/ping.log", u.HomeDir))
	if err != nil {
		return err
	}
	for _, i := range entries {
		fmt.Fprintf(writer, "%s: %s - %s\n", i.date, i.from, i.until)
	}
	return nil
}

func readFromFile(filename string) ([]*entry, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return read(fi)
}

func read(rd io.Reader) ([]*entry, error) {
	entries := make([]*entry, 0)
	re := regexp.MustCompile("(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2})\\s(.*)")
	var e *entry
	reader := bufio.NewReader(rd)
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
}
