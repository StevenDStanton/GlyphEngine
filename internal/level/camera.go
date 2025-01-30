package level

var (
	cameraX int
	cameraY int
)

func (l *Level) UpdateCamera(playerX, playerY int) {
	halfWidth := COL_DIVIDER / 2
	halfHeight := ROW_DIVIDER / 2
	desiredX := playerX - halfWidth
	desiredY := playerY - halfHeight

	mapHeight := len(l.MapGrid)
	if mapHeight == 0 {
		return
	}
	mapWidth := len(l.MapGrid[0])

	viewW := COL_DIVIDER
	viewH := ROW_DIVIDER

	if mapWidth <= viewW {
		cameraX = 0
	} else {
		if desiredX < 0 {
			desiredX = 0
		}
		maxCameraX := mapWidth - viewW
		if desiredX > maxCameraX {
			desiredX = maxCameraX
		}
		cameraX = desiredX
	}

	if mapHeight <= viewH {
		cameraY = 0
	} else {
		if desiredY < 0 {
			desiredY = 0
		}
		maxCameraY := mapHeight - viewH
		if desiredY > maxCameraY {
			desiredY = maxCameraY
		}
		cameraY = desiredY
	}
	l.updateGridFromCamera()
}

func (l *Level) updateGridFromCamera() {
	viewW := COL_DIVIDER
	viewH := ROW_DIVIDER

	for screenY := 0; screenY < viewH; screenY++ {
		for screenX := 0; screenX < viewW; screenX++ {
			mapY := cameraY + screenY
			mapX := cameraX + screenX

			if mapY < 0 || mapY >= len(l.MapGrid) ||
				mapX < 0 || mapX >= len(l.MapGrid[mapY]) {
				l.TheGrid[screenY][screenX] = ' '
			} else {
				l.TheGrid[screenY][screenX] = l.MapGrid[mapY][mapX]
			}
		}
	}
}
