package level

import (
	"strconv"
)

func (l *Level) IsWalkable(x, y int) bool {
	if l.disableInput {
		return false
	}
	tile := l.MapGrid[y][x]

	if l.showingItem {
		l.showingItem = false
		l.currentEntity = nil
		l.clearDialogArea()
	}

	_, exists := Obstacles[tile]
	if exists {
		key := strconv.Itoa(y) + strconv.Itoa(x)
		item, interactionExists := l.entities[key]
		if interactionExists {
			l.currentEntity = &item
			l.showingItem = true
			l.showingIntro = false
			l.clearDialogArea()
			l.setDialog([]string{item.Text}, []string{"Press SPACE to interact"})
		} else {
			l.currentEntity = nil
			l.showingItem = true
			l.showingIntro = false
			l.clearDialogArea()
			l.setDialog([]string{"Its a wall, you probably", "shouldn't lick it", "It doesn't look very clean"}, []string{"Press SPACE to interact"})
		}

	}

	return !exists
}

func (l *Level) UpdateBoard(x, y, dx, dy int) {
	l.MapGrid[y][x] = LEVEL_EMPTY
	l.MapGrid[dy][dx] = l.dialogState.Entities[0].Symbol
}
