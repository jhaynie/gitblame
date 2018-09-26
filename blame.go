package gitblame

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
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
	authorPrefix     = "author "
	authorMailPrefix = "author-mail "
	authorTimePrefix = "author-time "
)

// the maximum of one line of output. testing with 1K which seems OK
const maxLineSize = 1024

// GenerateOutput will take in a reader that's already in the correct blame output
// and return each line of the blame entry to callback
func GenerateOutput(ctx context.Context, r io.Reader, callback Callback, w io.Writer) error {
	lr := bufio.NewReaderSize(r, maxLineSize)
	s := bufio.NewScanner(lr)
	var writer *bufio.Writer
	if w != nil {
		writer = bufio.NewWriter(w)
		defer writer.Flush()
	}
	var current BlameLine
	for s.Scan() {
		err := s.Err()
		if err != nil {
			if strings.Contains(s.Err().Error(), "file already closed") {
				break
			}
			return err
		}
		// make sure our context isn't done
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		buf := s.Text()
		if writer != nil {
			_, err := writer.WriteString(buf)
			if err != nil {
				return fmt.Errorf("error writing buffer to output. %v", err)
			}
			_, err = writer.WriteString("\n")
			if err != nil {
				return fmt.Errorf("error writing buffer to output. %v", err)
			}
		}
		if strings.HasPrefix(buf, authorPrefix) {
			current.Name = string(buf[7:])
		} else if strings.HasPrefix(buf, authorMailPrefix) {
			current.Email = string(buf[13 : len(buf)-1])
		} else if strings.HasPrefix(buf, authorTimePrefix) {
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
	return generateWithRetry(context.Background(), dir, sha, filename, callback, writer, 1)
}

// GenerateWithContext will generate blame detail for a specific commit sha and filename
// writer is an optional (set to nil if not needed) io.Writer that will stream the output
// of blame to cache the result for later usage with GenerateOutput
func GenerateWithContext(ctx context.Context, dir string, sha string, filename string, callback Callback, writer io.Writer) error {
	return generateWithRetry(ctx, dir, sha, filename, callback, writer, 1)
}

func generateWithRetry(ctx context.Context, dir string, sha string, filename string, callback Callback, writer io.Writer, count int) error {
	if count > 3 {
		return fmt.Errorf("tried to run git blame %v on %v (%v) too many times and it failed on 3 retries", sha, filename, dir)
	}
	// fmt.Println("git", "blame", sha, "-e", "--root", "--line-porcelain", "--", filename)
	cmd := exec.CommandContext(ctx, "git", "blame", sha, "-e", "--root", "--line-porcelain", "--", filename)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Dir = dir
	r, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error running blame in %v for %v (%v). %v", dir, sha, filename, err)
	}
	defer r.Close()
	if err := cmd.Start(); err != nil {
		if strings.Contains(err.Error(), "exit status") {
			r.Close()
			cmd.Wait()
			if strings.Contains(stderr.String(), "no such path") {
				return errors.New(strings.TrimSpace(stderr.String()))
			}
			// sleep a bit to backoff in case too many other processes or something like that..
			time.Sleep(time.Millisecond * 100 * time.Duration(count))
			return generateWithRetry(ctx, dir, sha, filename, callback, writer, count+1)
		}
		return fmt.Errorf("git blame %s exited with %s. %v", sha, err.Error(), strings.TrimSpace(stderr.String()))
	}
	if err := GenerateOutput(ctx, r, callback, writer); err != nil {
		r.Close()
		cmd.Wait()
		return err
	}
	r.Close()
	if err := cmd.Wait(); err != nil {
		if strings.Contains(err.Error(), "exit status 128") {
			if strings.Contains(stderr.String(), "no such path") {
				return errors.New(strings.TrimSpace(stderr.String()))
			}
		}
		return fmt.Errorf("git blame %s exited with %v. %v", sha, err, strings.TrimSpace(stderr.String()))
	}
	return nil
}
