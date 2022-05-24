package arc_test

import (
	"bytes"
	"testing"

	"github.com/qba73/arc"
)

func TestCLI_PrintsVersion(t *testing.T) {
	t.Parallel()
	tc := struct {
		name string
		args []string
		want string
	}{
		name: "Print version",
		args: []string{"-version"},
		want: "Version: \nGitRef: \nBuild Time: \n",
	}

	out := &bytes.Buffer{}

	err := arc.CLI(
		arc.WithArgs(tc.args),
		arc.WithOutput(out),
	)
	if err != nil {
		t.Fatal(err)
	}
	got := out.String()

	if tc.want != got {
		t.Errorf("want %q, got %q", tc.want, got)
	}
}
