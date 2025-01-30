package level

import (
	_ "embed"
	"encoding/json"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/txt/level1_dialog.json
var level_1_dialog string

//go:embed assets/txt/level2_dialog.json
var level_2_dialog string

//go:embed assets/txt/level3_dialog.json
var level_3_dialog string

//go:embed assets/txt/level4_dialog.json
var level_4_dialog string

//go:embed assets/txt/level5_dialog.json
var level_5_dialog string

type DialogState struct {
	Intro    []string `json:"intro"`
	Entities []Entity `json:"entities"`
}

type EntityDialogOption struct {
	Text      string `json:"text"`
	NextState string `json:"nextState,omitempty"`
}

type EntityDialogState struct {
	ID        string               `json:"id"`
	Text      string               `json:"text"`
	NextState string               `json:"nextState,omitempty"`
	Options   []EntityDialogOption `json:"options,omitempty"`
}

type Entity struct {
	ID           string              `json:"id"`
	Symbol       rune                `json:"symbol"`
	X            int                 `json:"x"`
	Y            int                 `json:"y"`
	Text         string              `json:"text,omitempty"`
	DialogStates []EntityDialogState `json:"dialogStates,omitempty"`
	active       bool
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	type alias Entity
	var tmp struct {
		Symbol string `json:"symbol"`
		*alias
	}
	tmp.alias = (*alias)(e)

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	runes := []rune(tmp.Symbol)
	if len(runes) == 0 {
		e.Symbol = 0
	} else {
		e.Symbol = runes[0]
	}

	return nil
}

func (l *Level) loadDialog() {
	var dialogState DialogState
	err := json.Unmarshal([]byte(levelDialogs[l.level]), &dialogState)
	if err != nil {
		log.Fatal(err)
	}

	l.dialogState = nil

	l.dialogState = &dialogState
	l.levelIntroDialog = dialogState.Intro
	l.dialogActive = true
	l.selectedDialog = 0
	l.showingIntro = true

	for _, entity := range dialogState.Entities {
		key := strconv.Itoa(entity.Y) + strconv.Itoa(entity.X)
		l.entities[key] = entity
		Obstacles[entity.Symbol] = entity.Symbol
		l.MapGrid[entity.Y][entity.X] = entity.Symbol
		if entity.ID == "Player" {
			l.Player.SetPosition(entity.X, entity.Y)
		}
	}

	l.setDialog(l.levelIntroDialog, []string{"Press SPACE to continue"})

}

func (l *Level) updateDialog() {
	l.selectCooldown--
	if l.selectCooldown > 0 {
		return
	}
	l.selectCooldown = 10

	if l.showingIntro {
		l.showIntro()
		return
	}

	if l.currentEntity != nil {
		if !l.currentEntity.active {
			l.activateCurrentEntity()
			return
		}

		if l.currentEntity.active {
			l.processACtivatedEntity()
			return
		}
	}
}

func (l *Level) processACtivatedEntity() {
	state := l.getCurrentEntityState()
	if state == nil {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		if len(state.Options) == 0 {
			if state.NextState == "end" {
				l.currentEntity.active = false
				l.clearDialogArea()
			}
			l.currentDialogStateID = state.NextState
			l.selectedDialog = 0
			l.setDialogForCurrentState()
		} else {
			opt := state.Options[l.selectedDialog]
			if opt.NextState == "end" {
				l.currentEntity.active = false
				l.clearDialogArea()
			}

			if opt.NextState == "exit" {
				l.level++
				l.LoadLevel()
				return
			}

			if opt.NextState == "end_game" {
				l.level = 0
				l.LoadLevel()
			}

			l.currentDialogStateID = opt.NextState
			l.selectedDialog = 0
			l.setDialogForCurrentState()
		}
	}

	if len(state.Options) > 0 {
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			l.selectedDialog = (l.selectedDialog + 1) % len(state.Options)
			l.setDialogForCurrentState()
		}

		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			l.selectedDialog = (l.selectedDialog - 1 + len(state.Options)) % len(state.Options)
			l.setDialogForCurrentState()
		}
	}
}

func (l *Level) activateCurrentEntity() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) {
		l.currentEntity.active = true
		l.currentDialogStateID = "start"
		l.selectedDialog = 0
		l.setDialogForCurrentState()
		return
	}
}

func (l *Level) showIntro() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) {
		l.showingIntro = false
		l.clearDialogArea()
	}
}

func (l *Level) clearDialogArea() {
	for y := 0; y < ROW_DIVIDER; y++ {
		for x := COL_DIVIDER; x < COLS; x++ {
			l.TheGrid[y][x] = LEVEL_EMPTY
		}
	}
}

func (l *Level) setDialog(text []string, options []string) {
	dialRow := 0
	for _, textLine := range text {
		lineStart := 0
		for lineStart < len(textLine) {
			lineEnd := lineStart + (COLS - COL_DIVIDER - DIALOG_INDENT)
			if lineEnd > len(textLine) {
				lineEnd = len(textLine)
			}

			gridRow := dialRow + DIALOG_INDENT

			for dialCol := 0; dialCol < lineEnd-lineStart; dialCol++ {
				gridCol := COL_DIVIDER + DIALOG_INDENT + dialCol
				l.TheGrid[gridRow][gridCol] = rune(textLine[lineStart+dialCol])
			}

			dialRow++
			lineStart = lineEnd

			if gridRow+1 >= ROW_DIVIDER {
				return
			}
		}
	}

	for i := 0; i < DIALOG_INDENT; i++ {
		dialRow++
		if dialRow+DIALOG_INDENT >= ROW_DIVIDER {
			return
		}
	}

	sepRow := dialRow + DIALOG_INDENT
	if sepRow < ROW_DIVIDER {
		for col := COL_DIVIDER + DIALOG_INDENT; col < COLS; col++ {
			l.TheGrid[sepRow][col] = '-'
		}
		dialRow++
	}

	for i := 0; i < DIALOG_INDENT; i++ {
		dialRow++
		if dialRow+DIALOG_INDENT >= ROW_DIVIDER {
			return
		}
	}

	for _, opt := range options {
		optionRow := dialRow + DIALOG_INDENT
		if optionRow >= ROW_DIVIDER {
			return
		}
		startCol := COL_DIVIDER + DIALOG_INDENT
		for j, r := range opt {
			gridCol := startCol + j
			if gridCol < COLS {
				l.TheGrid[optionRow][gridCol] = r
			}
		}
		dialRow++
	}

}

func (l *Level) getCurrentEntityState() *EntityDialogState {
	if l.currentEntity == nil {
		return nil
	}
	for i := range l.currentEntity.DialogStates {
		if l.currentEntity.DialogStates[i].ID == l.currentDialogStateID {
			return &l.currentEntity.DialogStates[i]
		}
	}
	return nil
}

func (l *Level) setDialogForCurrentState() {
	l.clearDialogArea()
	state := l.getCurrentEntityState()
	if state == nil {
		return
	}

	text := []string{state.Text}

	var optionTexts []string
	for i, opt := range state.Options {
		prefix := "  "
		if i == l.selectedDialog {
			prefix = "> "
		}
		optionTexts = append(optionTexts, prefix+opt.Text)
	}

	l.setDialog(text, optionTexts)
}
