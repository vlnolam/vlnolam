package format

import (
	"sync"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/vlnolam/vlnolam/orchestrator/internal/errors"
)

var (
	doneStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#10b504"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4757")).Bold(true)
)

const (
	iconCheck = "✔"
	iconError = "✘"
)

// BubbleFormatter provides a TUI using the Bubble Tea library
type BubbleFormatter struct {
	waitGroup sync.WaitGroup
	program   *tea.Program
	isInit    bool
}

func (f *BubbleFormatter) EnsureInit() {
	if f.isInit {
		return
	}

	m := &bubbleModel{}
	f.program = tea.NewProgram(m)
	f.waitGroup.Go(f.RunAsync)
	f.isInit = true
}

func (f *BubbleFormatter) Finalize() {
	if !f.isInit {
		return
	}

	f.program.Quit()
	f.waitGroup.Wait()
}

func (f *BubbleFormatter) Done() {
	f.EnsureInit()
	f.program.Println(doneStyle.Render(iconCheck + " Done"))
}

func (f *BubbleFormatter) Error(e *errors.Error) {
	f.EnsureInit()

	f.program.Println(errorStyle.Render(iconError + " Error: " + e.Data["key"].(string)))
	for key, value := range e.Data {
		if key == "key" {
			continue
		}

		f.program.Printf("  %s: %s", key, value)
	}
}

type bubbleModel struct {
}

func (m *bubbleModel) Init() tea.Cmd {
	return nil
}

func (m *bubbleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *bubbleModel) View() tea.View {
	v := tea.NewView("Test")
	return v
}

func (f *BubbleFormatter) RunAsync() {
	f.program.Run()
}
