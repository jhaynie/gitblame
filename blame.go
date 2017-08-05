package gitblame

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"
)

// BlameLine is a structure for a blame result for a specific user
type BlameLine struct {
	Name  string
	Email string
	Date  time.Time
	Line  string
}

// Callback is called for each line
type Callback func(line BlameLine) error

var (
	authorPrefix     = []byte("author ")
	authorMailPrefix = []byte("author-mail ")
	authorTimePrefix = []byte("author-time ")
	eol              = []byte("\n")
)

// the maximum of one line of output. testing with 1K which seems OK
const maxLineSize = 1024

// GenerateOutput will take in a reader that's already in the correct blame output
// and return each line of the blame entry to callback
func GenerateOutput(r io.Reader, callback Callback, writer io.Writer) error {
	lr := bufio.NewReaderSize(r, maxLineSize)
	var current BlameLine
	for {
		buf, _, err := lr.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
		if writer != nil {
			_, err := writer.Write(buf)
			if err != nil {
				return fmt.Errorf("error writing buffer to output. %v", err)
			}
			_, err = writer.Write(eol)
			if err != nil {
				return fmt.Errorf("error writing buffer to output. %v", err)
			}
		}
		if bytes.HasPrefix(buf, authorPrefix) {
			current.Name = string(buf[7:])
		} else if bytes.HasPrefix(buf, authorMailPrefix) {
			current.Email = string(buf[13 : len(buf)-1])
		} else if bytes.HasPrefix(buf, authorTimePrefix) {
			i, err := strconv.ParseInt(string(buf[12:]), 10, 64)
			if err != nil {
				return err
			}
			current.Date = time.Unix(i, 0)
		} else if buf[0] == 9 {
			current.Line = string(buf[1:])
			if err := callback(current); err != nil {
				return err
			}
		}
	}
	return nil
}

// Generate will generate blame detail for a specific commit sha and filename
// writer is an optional (set to nil if not needed) io.Writer that will stream the output
// of blame to cache the result for later usage with GenerateOutput
func Generate(dir string, sha string, filename string, callback Callback, writer io.Writer) error {
	cmd := exec.Command("git", "blame", sha, filename, "-e", "--root", "-w", "-t", "--line-porcelain")
	cmd.Dir = dir
	r, err := cmd.StdoutPipe()
	if err != nil {
		cmd.Process.Kill()
		return err
	}
	defer r.Close()
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := GenerateOutput(r, callback, writer); err != nil {
		cmd.Process.Kill()
		return err
	}
	r.Close()
	return cmd.Wait()
}
