// Command backfill adds a default verification block (tier "t1", the current CI
// floor) to every registry entry that does not already have one.
//
// It is a one-shot, idempotent, line-based transform: it inserts the block as
// the last key of each entry and touches only those entries, preserving the
// existing formatting of every other line (no reordering, no re-escaping). Safe
// to re-run — entries that already carry a verification block are left untouched.
//
//	go run ./backfill ../registry.json
package main

import (
	"fmt"
	"os"
	"strings"
)

const defaultTier = "t1"

func main() {
	path := "../registry.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
		os.Exit(1)
	}

	out, added := backfill(string(data))

	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Printf("backfilled %d entr(y/ies) with verification tier %q in %s\n", added, defaultTier, path)
}

// backfill inserts a verification block after the last field of each entry that
// lacks one. An entry's last field is a 6-space-indented "key": value line with
// no trailing comma, immediately followed by the entry-closing brace (`    }`).
// It returns the transformed document and the number of entries changed.
func backfill(doc string) (string, int) {
	lines := strings.Split(doc, "\n")
	block := []string{
		`      "verification": {`,
		fmt.Sprintf(`        "tier": %q`, defaultTier),
		`      }`,
	}

	var b strings.Builder
	added := 0
	for i, line := range lines {
		if isLastField(line) && i+1 < len(lines) && isEntryClose(lines[i+1]) {
			b.WriteString(line)
			b.WriteString(",\n")
			b.WriteString(strings.Join(block, "\n"))
			b.WriteString("\n")
			added++
			continue
		}
		b.WriteString(line)
		if i < len(lines)-1 {
			b.WriteString("\n")
		}
	}
	return b.String(), added
}

// isLastField reports whether line is a 6-space-indented object field with no
// trailing comma (i.e. the final field of its object). An already-backfilled
// entry's min_mint_version line carries a trailing comma and is skipped, which
// is what makes the transform idempotent.
func isLastField(line string) bool {
	if !strings.HasPrefix(line, `      "`) {
		return false
	}
	return !strings.HasSuffix(strings.TrimRight(line, " "), ",")
}

// isEntryClose reports whether line closes a top-level entry object in the apis
// array (4-space indent), with or without a trailing comma.
func isEntryClose(line string) bool {
	t := strings.TrimRight(line, " ")
	return t == "    }" || t == "    },"
}
