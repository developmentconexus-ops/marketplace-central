package migrate_test

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/migrate"
	canonical "marketplace-central/apps/server_core/migrations"
)

func TestCanonicalSourceListsEveryMigrationByFullFilename(t *testing.T) {
	wantMatches, err := filepath.Glob(filepath.Join("..", "..", "..", "migrations", "*.sql"))
	if err != nil {
		t.Fatalf("glob canonical migrations: %v", err)
	}
	want := make([]string, 0, len(wantMatches))
	for _, path := range wantMatches {
		want = append(want, filepath.Base(path))
	}
	sort.Strings(want)
	if len(want) != 87 {
		t.Fatalf("fixture inventory drift: got %d canonical migrations, want 87", len(want))
	}

	got, err := migrate.Filenames(canonical.Source())
	if err != nil {
		t.Fatalf("list embedded migrations: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("embedded filenames differ from canonical files\ngot:  %v\nwant: %v", got, want)
	}
	if !sort.StringsAreSorted(got) {
		t.Fatal("embedded migrations are not lexicographically sorted")
	}
	var prefix0021 int
	for _, name := range got {
		if strings.HasPrefix(name, "0021_") {
			prefix0021++
		}
	}
	if prefix0021 != 2 {
		t.Fatalf("full filenames were collapsed by numeric prefix: got %d files with 0021_", prefix0021)
	}
}

func TestCanonicalSourceDoesNotDependOnCallerWorkingDirectory(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	got, err := migrate.Filenames(canonical.Source())
	if err != nil {
		t.Fatalf("list embedded migrations from foreign CWD: %v", err)
	}
	if len(got) != 87 {
		t.Fatalf("foreign CWD returned %d migrations, want 87", len(got))
	}
}
