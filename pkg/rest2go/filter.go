package rest2go

import (
	"fmt"
	"strings"
)

var directions = [2]string{"asc", "desc"}

type Filter struct {
	Prefix string
	Params map[string]string
	Sort   string
}

func NewFilter(prefix, defaultSort, query string, availableSort []string) Filter {
	sort := parseQuery(prefix, query, availableSort)

	if sort == "" {
		sort = defaultSort
	}

	return Filter{
		Prefix: prefix,
		Params: make(map[string]string),
		Sort:   sort,
	}
}

func parseQuery(prefix, sort string, available []string) string {
	var builder strings.Builder
	fields := strings.SplitSeq(sort, ",")

	for field := range fields {
		sortChunks := strings.Split(field, ":")

		if len(sortChunks) == 2 {
			column := ""
			direction := ""

			for _, d := range directions {
				if sortChunks[1] == d {
					direction = d
					break
				}
			}

			for _, c := range available {
				if sortChunks[0] == c {
					column = c
					break
				}
			}

			if column != "" && direction != "" {
				if prefix == "" {
					builder.WriteString(fmt.Sprintf(",%s %s", column, direction))
				} else {
					builder.WriteString(fmt.Sprintf(",%s.%s %s", prefix, column, direction))
				}
			}
		}
	}

	result := builder.String()

	if result != "" {
		result = result[1:]
	}

	return result
}
