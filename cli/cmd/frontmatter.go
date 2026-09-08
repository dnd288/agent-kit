package cmd

import "strings"

// parseFrontmatter returns the YAML frontmatter block between the opening and
// closing `---` fences, and whether a well-formed block was found.
func parseFrontmatter(content string) (string, bool) {
	if !strings.HasPrefix(content, "---") {
		return "", false
	}
	rest := content[3:]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return "", false
	}
	return rest[:end], true
}

// frontmatterField returns the trimmed, unquoted value of a top-level scalar
// field in a frontmatter block, or "" if the field is absent or empty.
func frontmatterField(fm, key string) string {
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key+":") {
			v := strings.TrimSpace(strings.TrimPrefix(line, key+":"))
			return strings.Trim(v, `"'`)
		}
	}
	return ""
}
