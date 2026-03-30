package cli

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewType int

const (
	ViewMain ViewType = iota
	ViewSearch
	ViewFiles
	ViewExamples
	ViewStatus
	ViewVersion
	ViewAbout
	ViewFTPInfo
	ViewConfigMenu
	ViewConfigFTP
	ViewConfigHTTP
	ViewConfigFTPToggle
	ViewConfigMCPToggle
	ViewConfigPorts
	ViewConfigSave
	ViewConfigView
)

type Model struct {
	serverIP    string
	fileMap     map[string]string
	lang        string
	quitFlag    *bool
	width       int
	height      int
	view        ViewType
	prevView    ViewType
	input       textinput.Model
	inputActive bool
	message     string
	configStep  int
}

func NewModel(serverIP string, fileMap map[string]string, lang string, quitFlag *bool) Model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.PromptStyle = promptStyle
	ti.CharLimit = 200

	return Model{
		serverIP: serverIP,
		fileMap:  fileMap,
		lang:     lang,
		quitFlag: quitFlag,
		view:     ViewMain,
		prevView: ViewMain,
		input:    ti,
		width:    80,
		height:   24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) contentWidth() int {
	return minInt(m.width-4, 70)
}

func (m Model) panelWidth() int {
	return minInt(m.width-4, 70)
}

func (m *Model) enterView(view ViewType) {
	m.prevView = m.view
	m.view = view
	m.message = ""
	m.configStep = 0

	needsInput := view == ViewSearch ||
		view == ViewConfigFTP ||
		view == ViewConfigHTTP ||
		view == ViewConfigPorts

	if needsInput {
		m.inputActive = true
		m.input.SetValue("")
		m.input.Focus()
		switch view {
		case ViewSearch:
			m.input.Placeholder = t(m.lang, "search_prompt")
		case ViewConfigFTP:
			if m.configStep == 0 {
				m.input.Placeholder = "username"
			}
		case ViewConfigHTTP:
			m.input.Placeholder = "username"
		case ViewConfigPorts:
			m.input.Placeholder = "HTTP port"
		}
	} else {
		m.inputActive = false
		m.input.Blur()
	}
}

func (m *Model) enterConfigStep(step int) {
	m.configStep = step
	m.inputActive = true
	m.input.SetValue("")
	m.input.Focus()
}

func (m Model) goBack() (tea.Model, tea.Cmd) {
	m.inputActive = false
	m.input.Blur()
	m.input.SetValue("")
	m.view = m.prevView
	m.prevView = ViewMain
	m.message = ""
	return m, nil
}

var _ tea.Model = Model{}
