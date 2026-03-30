package cli

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func panel(w int, content string) string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(w).
		Padding(0, 1).
		Render(content)
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(m.renderHeader())
	b.WriteString("\n")
	b.WriteString(m.renderContent())

	if m.message != "" {
		b.WriteString("\n  ")
		b.WriteString(m.message)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.renderFooter())

	return b.String()
}

func (m Model) renderHeader() string {
	if m.view == ViewMain {
		return m.renderBanner()
	}
	return m.renderCompactHeader()
}

func (m Model) renderFooter() string {
	lang := m.lang

	if m.view == ViewMain {
		enLabel := "EN"
		zhLabel := "中文"
		if lang == "en" {
			enLabel = greenStyle.Bold(true).Render("[EN]")
			zhLabel = dimStyle.Render("中文")
		} else {
			enLabel = dimStyle.Render("EN")
			zhLabel = greenStyle.Bold(true).Render("[中文]")
		}
		return "  " + dimStyle.Render(t(lang, "lang_hint")) + " " + enLabel + " " + zhLabel + "  " + dimStyle.Render(t(lang, "lang_switch")) + "\n"
	}

	if m.inputActive {
		return "  " + dimStyle.Render(t(lang, "press_esc")) + "\n"
	}

	return "  " + dimStyle.Render(t(lang, "press_back")) + "\n"
}

func (m Model) renderContent() string {
	switch m.view {
	case ViewMain:
		return m.renderMenu()
	case ViewSearch:
		return m.renderSearch()
	case ViewFiles:
		return m.renderFiles()
	case ViewExamples:
		return m.renderExamples()
	case ViewStatus:
		return m.renderStatus()
	case ViewVersion:
		return m.renderVersion()
	case ViewAbout:
		return m.renderAbout()
	case ViewFTPInfo:
		return m.renderFTPInfo()
	case ViewConfigMenu:
		return m.renderConfigMenu()
	case ViewConfigFTP:
		return m.renderConfigFTP()
	case ViewConfigHTTP:
		return m.renderConfigHTTP()
	case ViewConfigFTPToggle:
		return m.renderConfigFTPToggle()
	case ViewConfigMCPToggle:
		return m.renderConfigMCPToggle()
	case ViewConfigPorts:
		return m.renderConfigPorts()
	case ViewConfigSave:
		return m.renderConfigSave()
	case ViewConfigView:
		return m.renderConfigView()
	default:
		return m.renderMenu()
	}
}
