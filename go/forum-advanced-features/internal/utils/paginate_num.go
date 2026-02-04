package utils

import "html/template"

// TemplateFuncs returns a map of template functions
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"subtract": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sequence": func(start, end int) []int {
			var sequence []int
			for i := start; i <= end; i++ {
				sequence = append(sequence, i)
			}
			return sequence
		},
	}
}
