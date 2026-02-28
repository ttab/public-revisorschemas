package publicrevisorschemas_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testingT interface {
	Helper()
	Fatalf(format string, args ...any)
}

func must(t testingT, err error, format string, a ...any) {
	t.Helper()

	if err != nil {
		t.Fatalf("failed: %s: %v", fmt.Sprintf(format, a...), err)
	}
}

// testAgainstGolden compares a result against the contents of the file at the
// goldenPath. Run with regenerate set to true to create or update the file.
func testAgainstGolden[T any](
	t *testing.T,
	regenerate bool,
	got T,
	goldenPath string,
) {
	t.Helper()

	if regenerate {
		data, err := json.MarshalIndent(got, "", "  ")
		must(t, err, "marshal for storage in %q", goldenPath)

		// End all files with a newline
		data = append(data, '\n')

		err = os.WriteFile(goldenPath, data, 0o600)
		must(t, err, "write golden file %q", goldenPath)
	}

	wantData, err := os.ReadFile(goldenPath)
	must(t, err, "read from golden file %q", goldenPath)

	var wantValue T

	err = json.Unmarshal(wantData, &wantValue)
	must(t, err, "unmarshal data from golden file %q", goldenPath)

	diff := cmp.Diff(wantValue, got)
	if diff != "" {
		t.Fatalf("must match golden file %q: mismatch (-want +got):\n%s",
			goldenPath, diff)
	}
}
