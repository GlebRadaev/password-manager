package utils

import (
	"strings"
)

func StringsJoin(strs ...string) string {
	var strb strings.Builder

	for _, str := range strs {
		strb.WriteString(str)
	}

	return strb.String()
}
