// Package skills carries the agent skills (setup, deploy, debug, rollback)
// inside the CLI binary. The files in this directory are the ONE source: the
// Claude Code plugin ships them as they are, and `nesly mcp setup` installs
// them for the other agent tools from this same embedded copy
// (internal/agentclients/skills.go).
//
// Each skill is <name>/SKILL.md plus optional <name>/references/*.md, one
// level deep (progressive disclosure: SKILL.md stays short, detail is read
// on demand).
package skills

import "embed"

// FS holds <name>/SKILL.md and <name>/references/*.md for every skill.
//
//go:embed */SKILL.md */references/*.md
var FS embed.FS
