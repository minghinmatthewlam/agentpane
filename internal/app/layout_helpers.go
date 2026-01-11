package app

func layoutForPaneCount(count int) string {
	if count == 2 {
		return "even-horizontal"
	}
	return "tiled"
}
