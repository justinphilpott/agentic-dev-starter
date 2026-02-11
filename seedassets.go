package seedassets

import "embed"

// FS embeds canonical Seed assets into the CLI binary so seeded repos are self-contained.
//
//go:embed seed-contract/manifest.json skills/seed-validate/SKILL.md
var FS embed.FS
