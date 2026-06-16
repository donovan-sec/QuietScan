package storage

import "testing"

func TestSanitizeCSVCell(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"plain hostname", "printer-01", "printer-01"},
		{"plain vendor", "Apple, Inc.", "Apple, Inc."},
		{"equals formula", "=cmd|'/c calc'!A1", "'=cmd|'/c calc'!A1"},
		{"plus formula", "+1+1", "'+1+1"},
		{"minus formula", "-2+3", "'-2+3"},
		{"at formula", "@SUM(A1:A9)", "'@SUM(A1:A9)"},
		{"leading tab", "\t=evil", "'\t=evil"},
		{"leading cr", "\r=evil", "'\r=evil"},
		{"safe ip", "192.168.1.10", "192.168.1.10"},
		{"safe mac", "00:1a:2b:3c:4d:5e", "00:1a:2b:3c:4d:5e"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := sanitizeCSVCell(c.in); got != c.want {
				t.Errorf("sanitizeCSVCell(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestSanitizeCSVRecord(t *testing.T) {
	in := []string{"192.168.1.5", "00:11:22:33:44:55", "Acme", "=HYPERLINK(\"http://evil\")"}
	got := sanitizeCSVRecord(in)
	if got[3] != "'=HYPERLINK(\"http://evil\")" {
		t.Errorf("malicious hostname not neutralized: got %q", got[3])
	}
	// Benign fields must be untouched.
	for i := 0; i < 3; i++ {
		if got[i] != in[i] {
			t.Errorf("benign field %d altered: got %q want %q", i, got[i], in[i])
		}
	}
}
