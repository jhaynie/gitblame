package gitblame

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestBlame(t *testing.T) {
	f, err := os.Open("./testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	lines := make([]BlameLine, 0)
	callback := func(line BlameLine) error {
		lines = append(lines, line)
		return nil
	}
	if err := GenerateOutput(f, callback); err != nil {
		t.Fatal(err)
	}
	expected := []string{
		"0 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"1 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"2 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"3 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"4 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"5 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"6 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"7 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"8 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"9 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"10 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"11 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"12 2017-04-11 21:20:16 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"13 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"14 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"15 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"16 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"17 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"18 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"19 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"20 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"21 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"22 2017-07-07 21:16:40 +0000 UTC Robindiddams robindiddams@gmail.com",
		"23 2017-07-07 21:16:40 +0000 UTC Robindiddams robindiddams@gmail.com",
		"24 2017-06-21 08:48:44 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"25 2017-06-21 08:47:45 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"26 2017-07-02 17:33:19 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"27 2017-07-02 17:33:19 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"28 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"29 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"30 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"31 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"32 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"33 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"34 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"35 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"36 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"37 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"38 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"39 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"40 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"41 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"42 2017-06-21 08:45:22 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"43 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"44 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"45 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"46 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"47 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"48 2017-04-28 03:42:07 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"49 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"50 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"51 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"52 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"53 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"54 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
		"55 2017-04-11 21:17:47 +0000 UTC Jeff Haynie jhaynie@gmail.com",
	}
	if len(lines) != 56 {
		t.Fatal("didn't return 56 lines as expected")
	}
	var out bytes.Buffer
	for i, line := range lines {
		tline := fmt.Sprintf("%d %v %s %s", i, line.Date.UTC(), line.Name, line.Email)
		if expected[i] != tline {
			t.Fatalf("expected %s but was %s", expected[i], tline)
		}
		out.WriteString(line.Line)
		out.WriteRune('\n')
	}
	tf, err := os.Open("./testdata/test.out")
	if err != nil {
		t.Fatal(err)
	}
	defer tf.Close()
	buf, err := ioutil.ReadAll(tf)
	if err != nil {
		t.Fatal(err)
	}
	if out.String() != string(buf) {
		t.Fatal("output not expected")
	}
}

func BenchmarkBlame(b *testing.B) {
	f, err := os.Open("./testdata/test.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	callback := func(line BlameLine) error {
		return nil
	}
	for n := 0; n < b.N; n++ {
		if err := GenerateOutput(f, callback); err != nil {
			b.Fatal(err)
		}
	}
}
