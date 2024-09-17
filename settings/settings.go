package settings

func Color(names []string) []string {
	colorForGraf := []string{"rgba(255, 99, 132, 0.5)",
		"rgba(255, 159, 64, 0.5)",
		"rgba(255, 205, 86, 0.5)",
		"rgba(75, 192, 192, 0.5)",
		"rgba(54, 162, 235, 0.5)",
		"rgba(153, 102, 255, 0.5)",
		"rgba(201, 203, 207, 0.5)"}
	colors := []string{}
	for i := 0; i < len(names); i++ {
		j := i
		if j == 6 {
			j = 0
		}
		colors = append(colors, colorForGraf[j])

	}
	return colors

}
