package db_reader

import "fmt"

// makeLink возвращает строку, содержащую html-ссылку. Параметры - текст ссылки, якорь и стили вида
// "style1:value1;style2:value2;"
func makeLink(text, link, styles string) string {
	if len(styles) == 0 {
		return fmt.Sprintf(`<a href="%s">%s</a>`, link, text)
	}
	return fmt.Sprintf(`<a href="%s" style="%s">%s</a>`, link, styles, text)
}

// makeStyles возвращает строку стилей вида "style1:value1;style2:value2;" полученные из мапы вида "map[стиль]значение"
func makeStyles(styles map[string]string) string {
	var result string

	for style, value := range styles {
		result += fmt.Sprintf("%s:%s;", style, value)
	}

	return result
}

// makeTableHeader возвращает строку, содержащую html-код заголовка таблицы
func makeTableHeader(header []string) string {
	th := tags("th",
		makeStyles(map[string]string{"border": "solid 1px black", "padding": "5px", "background": "#E0DBD9"}),
		header)
	row := tag("tr", th, "")
	return tag("thead", row, "")
}

// makeTableBody возвращает строку, содержащую html-код тела таблицы
func makeTableBody(data [][]string) string {
	var rows string
	for _, row := range data {
		rows += tag("tr",
			tags("td", makeStyles(map[string]string{"border": "solid 1px black", "padding": "5px"}), row),
			"")
	}
	return tag("tbody", rows, "")
}

// makeTable возвращает строку, содержащую html-код таблицы. В качестве параметров массив заголовков таблицы и матрица
// со значениями ячеек
func makeTable(caption []string, data [][]string) string {
	return tag("div",
		tag("table",
			makeTableHeader(caption)+makeTableBody(data),
			makeStyles(map[string]string{"border": "solid", "display": "inline-block", "border-collapse": "collapse"})),
		makeStyles(map[string]string{"text-align": "center"}))
}

// tags возвращает строку, содержащую html-код повторяющихся тегов, содержащих переданные в массиве строки и оформленные
// соответственно переданным стилям
func tags(htmlTag, styles string, data []string) string {
	var result string
	for _, text := range data {
		result += tag(htmlTag, text, styles)
	}

	return result
}

// tag возвращает строку, где в обрамлении указанного тега находится переданный текст. К тегу применяются переданные
// стили
//
// Например:
// tag("span", "go", "color:red;")
// вернет строку '<span style="color:red;">go</span>'
func tag(htmlTag, text, styles string) string {
	if len(styles) == 0 {
		return fmt.Sprintf(`<%s>%s</%s>`, htmlTag, text, htmlTag)
	}
	return fmt.Sprintf(`<%s style="%s">%s</%s>`, htmlTag, styles, text, htmlTag)
}
