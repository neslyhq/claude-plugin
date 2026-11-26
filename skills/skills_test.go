package skills

import (
	"io/fs"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var wantSkills = []string{"debug", "deploy", "rollback", "setup"}

func lines(b []byte) int { return strings.Count(strings.TrimRight(string(b), "\n"), "\n") + 1 }

func TestTheFourSkillsAreEmbedded(t *testing.T) {
	entries, err := fs.ReadDir(FS, ".")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range entries {
		if e.IsDir() {
			got = append(got, e.Name())
		}
	}
	sort.Strings(got)
	if strings.Join(got, ",") != strings.Join(wantSkills, ",") {
		t.Fatalf("skills = %v, want %v", got, wantSkills)
	}
}

// Every SKILL.md is short: detail lives in references/, read on demand.
func TestEverySkillMDIsAtMost120Lines(t *testing.T) {
	for _, n := range wantSkills {
		b, err := FS.ReadFile(n + "/SKILL.md")
		if err != nil {
			t.Fatal(err)
		}
		if got := lines(b); got > 120 {
			t.Errorf("%s/SKILL.md has %d lines, the limit is 120: move detail to references/", n, got)
		}
	}
}

var (
	frontName  = regexp.MustCompile(`(?m)\A---\nname: (\S+)\ndescription: (.+)\n---\n`)
	refMention = regexp.MustCompile("`references/([a-z0-9-]+\\.md)`")
)

// The front matter is the name (the directory) and a third-person description
// with the trigger words in both languages.
func TestEverySkillHasFrontMatterAndTriggers(t *testing.T) {
	for _, n := range wantSkills {
		b, _ := FS.ReadFile(n + "/SKILL.md")
		m := frontName.FindSubmatch(b)
		if m == nil {
			t.Errorf("%s: front matter is not name + one-line description", n)
			continue
		}
		if string(m[1]) != n {
			t.Errorf("%s: name is %q", n, m[1])
		}
		d := string(m[2])
		if len(d) > 1024 {
			t.Errorf("%s: description is %d characters, the limit is 1024", n, len(d))
		}
		for _, first := range []string{"I ", "You ", "Use this", "Your "} {
			if strings.HasPrefix(d, first) {
				t.Errorf("%s: the description must be third person, it starts %q", n, first)
			}
		}
		if !strings.Contains(d, "Use when") {
			t.Errorf("%s: the description has no 'Use when' triggers", n)
		}
	}
	// the two languages the person may write in
	for n, words := range map[string][]string{
		"deploy":   {"deploy", "publish", "yayınla", "canlıya al"},
		"rollback": {"rollback", "geri al"},
		"debug":    {"logs"},
		"setup":    {"login"},
	} {
		b, _ := FS.ReadFile(n + "/SKILL.md")
		head := strings.ToLower(strings.SplitN(string(b), "\n---\n", 2)[0])
		for _, w := range words {
			if !strings.Contains(head, w) {
				t.Errorf("%s: the description lacks the trigger %q", n, w)
			}
		}
	}
}

// References sit one level deep, every one a SKILL.md mentions exists, and a
// reference over 100 lines opens with a table of contents.
func TestReferencesAreOneLevelDeepAndHaveContents(t *testing.T) {
	err := fs.WalkDir(FS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.Contains(p, "/references/") {
			return err
		}
		if strings.Count(p, "/") != 2 {
			t.Errorf("%s: references are one level deep", p)
		}
		b, _ := FS.ReadFile(p)
		if lines(b) > 100 && !strings.Contains(string(b), "Contents") {
			t.Errorf("%s has %d lines and no table of contents", p, lines(b))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range wantSkills {
		b, _ := FS.ReadFile(n + "/SKILL.md")
		for _, m := range refMention.FindAllSubmatch(b, -1) {
			found := false
			// a skill may point at another skill's reference by that skill's name
			for _, o := range wantSkills {
				if _, err := FS.ReadFile(o + "/references/" + string(m[1])); err == nil {
					found = true
				}
			}
			if !found {
				t.Errorf("%s/SKILL.md mentions references/%s, which does not exist", n, m[1])
			}
		}
	}
}

// The auth paragraph lives once, in setup/references/auth.md; the other
// skills point at it.
func TestAuthIsNotRepeated(t *testing.T) {
	for _, n := range []string{"deploy", "debug", "rollback"} {
		b, _ := FS.ReadFile(n + "/SKILL.md")
		s := string(b)
		if strings.Contains(s, "NESLY_PASSWORD") || strings.Contains(s, "nesly token create") {
			t.Errorf("%s repeats the auth detail: point at the setup skill's references/auth.md", n)
		}
		if !strings.Contains(s, "**setup** skill") {
			t.Errorf("%s does not point at the setup skill", n)
		}
	}
}
