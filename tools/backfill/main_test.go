package main

import (
	"strings"
	"testing"
)

const oneEntry = `{
  "schema_version": 1,
  "apis": [
    {
      "name": "ably-control",
      "description": "Control API for Ably applications",
      "tags": [
        "communication",
        "messaging"
      ],
      "spec_url": "https://example.com/ably.json",
      "auth_type": "api_key",
      "auth_env_var": "ABLY_API_KEY",
      "min_mint_version": "0.2.0"
    }
  ]
}
`

func TestBackfill(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		wantAdded int
		wantHas   string
	}{
		{
			name:      "adds block to entry ending in min_mint_version",
			in:        oneEntry,
			wantAdded: 1,
			wantHas:   "      \"verification\": {\n        \"tier\": \"t1\"\n      }",
		},
		{
			name: "adds block to entry ending in skip_url_check",
			in: `{
  "schema_version": 1,
  "apis": [
    {
      "name": "x",
      "description": "desc long enough",
      "tags": ["misc"],
      "spec_url": "https://example.com/x.json",
      "auth_type": "none",
      "auth_env_var": "X_TOKEN",
      "min_mint_version": "0.2.0",
      "skip_url_check": true
    }
  ]
}
`,
			wantAdded: 1,
			wantHas:   "\"verification\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, added := backfill(tt.in)
			if added != tt.wantAdded {
				t.Errorf("added = %d, want %d", added, tt.wantAdded)
			}
			if !strings.Contains(out, tt.wantHas) {
				t.Errorf("output missing %q\n---\n%s", tt.wantHas, out)
			}
		})
	}
}

func TestBackfill_Idempotent(t *testing.T) {
	once, added1 := backfill(oneEntry)
	if added1 != 1 {
		t.Fatalf("first pass added = %d, want 1", added1)
	}
	twice, added2 := backfill(once)
	if added2 != 0 {
		t.Errorf("second pass added = %d, want 0 (not idempotent)", added2)
	}
	if twice != once {
		t.Errorf("second pass mutated an already-backfilled document")
	}
}

func TestBackfill_PreservesOtherLines(t *testing.T) {
	out, _ := backfill(oneEntry)
	// Every original line must still be present verbatim; the only change to an
	// existing line is the comma appended to the previous last field.
	for _, line := range strings.Split(oneEntry, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line == `      "min_mint_version": "0.2.0"` {
			if !strings.Contains(out, line+",") {
				t.Errorf("last field did not gain a trailing comma")
			}
			continue
		}
		if !strings.Contains(out, line) {
			t.Errorf("original line not preserved: %q", line)
		}
	}
}
