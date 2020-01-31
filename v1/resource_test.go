package acl

import (
	"testing"
)

// Test path
func TestPath(t *testing.T) {

	// no options
	checkPath(t, "*", "Anything!", 0, true)
	checkPath(t, "*", "/", 0, true)
	checkPath(t, "*", "", 0, true)

	checkPath(t, "/*", "/a", 0, true)

	checkPath(t, "/a", "/a", 0, true)
	checkPath(t, "/a", "/b", 0, false)

	checkPath(t, "/a/", "/a", 0, true)
	checkPath(t, "/a", "/a/", 0, true)
	checkPath(t, "/a/b/", "/a/b", 0, true)
	checkPath(t, "/a/b", "/a/b/", 0, true)

	checkPath(t, "/a/b", "/a/b", 0, true)
	checkPath(t, "/a/b", "/a", 0, false)
	checkPath(t, "/a/b", "/a/a", 0, false)

	checkPath(t, "/a/*", "/a/a", 0, true)
	checkPath(t, "/a/*", "/a/Z", 0, true)
	checkPath(t, "/a/*", "/a/*", 0, true)

	checkPath(t, "/a/*", "/a/b/c", 0, true)
	checkPath(t, "/a/*", "/a/XYZ", 0, true)

	checkPath(t, "/a/*/z", "/a/b", 0, false)
	checkPath(t, "/a/*/z", "/a/z", 0, false)
	checkPath(t, "/a/*/z", "/a/l/z", 0, true)

	// encompass
	checkPath(t, "/a", "/a", MatchOptionEncompass, true)
	checkPath(t, "/a", "/b", MatchOptionEncompass, false)
	checkPath(t, "/a", "/b/c", MatchOptionEncompass, false)

	checkPath(t, "/a/b", "/a/b", MatchOptionEncompass, true)
	checkPath(t, "/a/b", "/a", MatchOptionEncompass, false)
	checkPath(t, "/a/b", "/a/Z", MatchOptionEncompass, false)

	checkPath(t, "/a", "/a/a", MatchOptionEncompass, true)
	checkPath(t, "/a", "/a/Z", MatchOptionEncompass, true)
	checkPath(t, "/a", "/a/*", MatchOptionEncompass, true)

	checkPath(t, "/a", "/a/b/c", MatchOptionEncompass, true)
	checkPath(t, "/a", "/a/XYZ", MatchOptionEncompass, true)
	checkPath(t, "/a/b", "/a/b/c", MatchOptionEncompass, true)
	checkPath(t, "/a/b/c/d", "/a/b/c", MatchOptionEncompass, false)

	checkPath(t, "/a/*/z", "/a/b", MatchOptionEncompass, false)
	checkPath(t, "/a/*/z", "/a/z", MatchOptionEncompass, false)
	checkPath(t, "/a/*/z", "/a/l/z", MatchOptionEncompass, true)

}

func checkPath(t *testing.T, l, r string, opts int, expect bool) {
	if Path(l).MatchesStringWithOptions(r, opts) != expect {
		if expect {
			t.Errorf("<%v> == <%v>", l, r)
		} else {
			t.Errorf("<%v> != <%v>", l, r)
		}
	}
}
