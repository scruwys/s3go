package s3go

import (
	"fmt"
    "testing"
)

func TestHumanizeBytes(t *testing.T) {
    var tests = []struct {
        b int64
        want string
    }{
    	{1024, "1.0 KiB"},
    	{0, "0 B"},
    	{161661570, "154.2 MiB"},
    	{87020, "85.0 KiB"},
    	{648799, "633.6 KiB"},
    	{1689536476, "1.6 GiB"},
    }

    for _, tt := range tests {
        testname := fmt.Sprintf("%d,%s", tt.b, tt.want)
        t.Run(testname, func(t *testing.T) {
            ans := HumanizeBytes(tt.b)
            if ans != tt.want {
                t.Errorf("got %s, want %s", ans, tt.want)
            }
        })
    }
}

func TestByteSizeToString(t *testing.T) {
	s := ByteSizeToString(1024, false)

	if s != "1024" {
		t.Errorf("%v is not a valid string", s)
	}

	h := ByteSizeToString(1024, true)

	if h != "1.0 KiB" {
		t.Errorf("%v is not a human readable", h)
	}
}

func TestRegexpCompile(t *testing.T) {
    var tests = []struct {
        a string
        want bool
    }{
        {"test.go", true},
        {"should/fs.go", true},
        {"should/compiled.go/fcn.go", true},
        {"should/other_file.go", true},
        {"should/other_file.go.txt", false},
        {"should/compiled.go/README", false},
    }

	s, _ := RegexpCompile("^(.*).go$", "$^")

    for _, tt := range tests {
        testname := fmt.Sprintf("%s", tt.a)
        t.Run(testname, func(t *testing.T) {
        	ans := s.MatchString(tt.a)
            if ans != tt.want {
                t.Errorf("%v evaluted to %v, should be %v", tt.a, ans, tt.want)
            }
        })
    }

    // Test with default pattern instead of main pattern.
	s, _ = RegexpCompile("", "^(.*).go$")

    for _, tt := range tests {
        testname := fmt.Sprintf("%s", tt.a)
        t.Run(testname, func(t *testing.T) {
        	ans := s.MatchString(tt.a)
            if ans != tt.want {
                t.Errorf("%v evaluted to %v, should be %v", tt.a, ans, tt.want)
            }
        })
    }
}

func TestIntMin(t *testing.T) {
    var tests = []struct {
        a, b int
        want int
    }{
        {0, 1, 0},
        {1, 0, 0},
        {2, -2, -2},
        {0, -1, -1},
        {-1, 0, -1},
    }

    for _, tt := range tests {
        testname := fmt.Sprintf("%d,%d", tt.a, tt.b)
        t.Run(testname, func(t *testing.T) {
            ans := IntMin(tt.a, tt.b)
            if ans != tt.want {
                t.Errorf("got %d, want %d", ans, tt.want)
            }
        })
    }
}
