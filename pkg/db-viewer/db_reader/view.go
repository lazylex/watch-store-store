package db_reader

import "fmt"

// makeLink возвращает строку, содержащую html-ссылку. Параметры - текст ссылки, якорь и стили вида
// "style1:value1;style2:value2;"
func makeLink(text, link, styles string) string {
	return fmt.Sprintf(`<a href="%s" style="%s">%s</a>`, link, styles, text)
}

// maleStyles возвращает строку стилей вида "style1:value1;style2:value2;" полученные из мапы вида "map[стиль]значение"
func makeStyles(styles map[string]string) string {
	var result string

	for style, value := range styles {
		result += fmt.Sprintf("%s:%s;", style, value)
	}

	return result
}

// makePath объединяет переданные строки в путь вида "строка1/строка2/.../строкаN"
func makePath(pathParts []string) string {
	var result string
	for idx, part := range pathParts {
		result += part
		if idx != len(pathParts) {
			result += "/"
		}
	}
	return result
}
