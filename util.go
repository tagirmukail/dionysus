package gotemplconstr

import (
	"strings"
)

func xmlReplaceSymbols(s string) string {
	for _, repS := range replaceXMLSymbols {
		s = strings.ReplaceAll(s, repS.from, repS.to)
	}

	return s
}
