package ini

import "strings"

func trimComment(v string) string {
	rest, _, _ := strings.Cut(v, "#")
	rest, _, _ = strings.Cut(rest, ";")
	return rest
}

// assumes no surrounding comment
func splitProperty(s string) (string, string, bool) {
	equalsi := strings.Index(s, "=")
	coloni := strings.Index(s, ":") // LEGACY: also supported for property assignment
	sep := "="
	if equalsi == -1 || coloni != -1 && coloni < equalsi {
		sep = ":"
	}

	k, v, ok := strings.Cut(s, sep)
	if !ok {
		return "", "", false
	}
	return strings.TrimSpace(k), strings.TrimSpace(v), true
}

// assumes no surrounding comment, whitespace, or profile brackets
func splitProfile(s string) (string, string) {
	var first int
	for i, r := range s {
		if isLineSpace(r) {
			if first == 0 {
				first = i
			}
		} else {
			if first != 0 {
				return s[:first], s[i:]
			}
		}
	}
	if first == 0 {
		return "", s // type component is effectively blank
	}
	return "", ""
}

func isLineSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func unquote(s string) string {
	if isSingleQuoted(s) || isDoubleQuoted(s) {
		return s[1 : len(s)-1]
	}
	return s
}

// applies various legacy conversions to property values:
//   - remote wrapping single/doublequotes
//   - expand escaped quote and newline sequences
func legacyStrconv(s string) string {
	s = unquote(s)
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\'`, `'`)
	s = strings.ReplaceAll(s, `\n`, "\n")
	return s
}

func isSingleQuoted(s string) bool {
	return hasAffixes(s, "'", "'")
}

func isDoubleQuoted(s string) bool {
	return hasAffixes(s, `"`, `"`)
}

func isBracketed(s string) bool {
	return hasAffixes(s, "[", "]")
}

func hasAffixes(s, left, right string) bool {
	return strings.HasPrefix(s, left) && strings.HasSuffix(s, right)
}
