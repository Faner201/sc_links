package shorten

import (
	"strings"

	"github.com/Faner201/sc_links/internal/utils"
)

const alphabet = "kfsPcd2aYGAVEnDzD14WjdQzqgPGcTi2L2uuuhkr"

var alphabetLen = uint32(len(alphabet))

func Shorten(id uint32) string {
	var (
		digits  []uint32
		num     = id
		builder strings.Builder
	)

	for num > 0 {
		digits = append(digits, num%alphabetLen)
		num /= alphabetLen
	}
	if len(digits) > 1 {
		utils.Reverse(digits)
	}

	for _, digit := range digits {
		builder.WriteString(string(alphabet[digit]))
	}

	return builder.String()
}
