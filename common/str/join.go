package str

import (
	"fmt"
	"github.com/jack0829/letsgo/common/types"
	"strings"
)

func JoinInteger[T types.Integer](slice []T, sep string) string {

	l := len(slice)
	if l < 1 {
		return ""
	}

	s := make([]string, 0, l)
	for _, i := range slice {
		s = append(s, fmt.Sprintf("%d", i))
	}

	return strings.Join(s, sep)
}
