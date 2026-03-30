package cli

import (
	"fmt"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zy84338719/upftp/internal/conf"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			*m.quitFlag = true
			return m, tea.Quit
		}
		return m.handleViewKeys(msg)
	}

	if m.inputActive {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.view {
	case ViewMain:
		return m.handleMainKeys(msg)
	case ViewSearch:
		return m.handleSearchKeys(msg)
	case ViewConfigMenu:
		return m.handleConfigMenuKeys(msg)
	case ViewConfigFTP:
		return m.handleConfigFTPKeys(msg)
	case ViewConfigHTTP:
		return m.handleConfigHTTPKeys(msg)
	case ViewConfigFTPToggle:
		return m.handleConfigFTPToggleKeys(msg)
	case ViewConfigMCPToggle:
		return m.handleConfigMCPToggleKeys(msg)
	case ViewConfigPorts:
		return m.handleConfigPortsKeys(msg)
	case ViewConfigSave:
		return m.handleConfigSaveKeys(msg)
	default:
		return m.handleStaticKeys(msg)
	}
}

func (m Model) handleMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "1":
		m.enterView(ViewSearch)
		return m, nil
	case "2":
		m.enterView(ViewFiles)
		return m, nil
	case "3":
		m.enterView(ViewExamples)
		return m, nil
	case "4":
		m.fileMap = scanDirectory(conf.AppConfig.Root, m.serverIP)
		m.message = greenStyle.Render(fmt.Sprintf(t(m.lang, "refresh_done"), len(m.fileMap)))
		return m, nil
	case "5":
		m.enterView(ViewStatus)
		return m, nil
	case "6":
		m.prevView = ViewMain
		m.enterView(ViewConfigMenu)
		m.prevView = ViewMain
		return m, nil
	case "7":
		if conf.AppConfig.EnableFTP {
			m.enterView(ViewFTPInfo)
		} else {
			m.enterView(ViewAbout)
		}
		return m, nil
	case "8":
		m.enterView(ViewAbout)
		return m, nil
	case "v":
		m.enterView(ViewVersion)
		return m, nil
	case "l":
		if m.lang == "en" {
			m.lang = "zh"
		} else {
			m.lang = "en"
		}
		conf.AppConfig.Language = m.lang
		_ = conf.SaveConfig()
		return m, nil
	case "q":
		*m.quitFlag = true
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.inputActive {
		switch msg.String() {
		case "enter":
			m.inputActive = false
			m.input.Blur()
			return m, nil
		case "esc":
			return m.goBack()
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}

	switch strings.ToLower(msg.String()) {
	case "b", "esc":
		return m.goBack()
	}
	return m, nil
}

func (m Model) handleStaticKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "b", "esc":
		return m.goBack()
	}
	return m, nil
}

func (m Model) handleConfigMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "1":
		m.prevView = ViewConfigMenu
		m.view = ViewConfigFTP
		m.message = ""
		m.configStep = 0
		m.inputActive = true
		m.input.SetValue("")
		m.input.Placeholder = "username"
		m.input.Focus()
		return m, nil
	case "2":
		m.prevView = ViewConfigMenu
		m.view = ViewConfigHTTP
		m.message = ""
		m.configStep = 0
		m.inputActive = false
		m.input.Blur()
		return m, nil
	case "3":
		m.prevView = ViewConfigMenu
		m.view = ViewConfigFTPToggle
		m.message = ""
		m.inputActive = false
		m.input.Blur()
		return m, nil
	case "4":
		m.prevView = ViewConfigMenu
		m.view = ViewConfigMCPToggle
		m.message = ""
		m.inputActive = false
		m.input.Blur()
		return m, nil
	case "5":
		m.prevView = ViewConfigMenu
		m.view = ViewConfigPorts
		m.message = ""
		m.configStep = 0
		m.inputActive = true
		m.input.SetValue("")
		m.input.Placeholder = "HTTP port"
		m.input.Focus()
		return m, nil
	case "6":
		m.prevView = ViewConfigMenu
		m.view = ViewConfigSave
		m.message = ""
		m.inputActive = false
		m.input.Blur()
		return m, nil
	case "7":
		m.prevView = ViewConfigMenu
		m.view = ViewConfigView
		m.message = ""
		m.inputActive = false
		m.input.Blur()
		return m, nil
	case "b", "esc":
		m.view = ViewMain
		m.prevView = ViewMain
		m.message = ""
		return m, nil
	case "q":
		*m.quitFlag = true
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleConfigFTPKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.inputActive {
		return m.goBack()
	}

	switch msg.String() {
	case "enter":
		val := strings.TrimSpace(m.input.Value())
		if m.configStep == 0 {
			if val != "" {
				conf.AppConfig.Username = val
				m.message = greenStyle.Render("username → " + val)
			}
			m.enterConfigStep(1)
			m.input.Placeholder = "password"
			return m, nil
		}
		if val != "" {
			conf.AppConfig.Password = val
			m.message = greenStyle.Render("password updated")
		}
		m.inputActive = false
		m.input.Blur()
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		m.configStep = 0
		return m, nil
	case "esc":
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		m.inputActive = false
		m.input.Blur()
		m.input.SetValue("")
		m.message = ""
		m.configStep = 0
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m Model) handleConfigHTTPKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.configStep == 0 {
		switch strings.ToLower(msg.String()) {
		case "y":
			conf.AppConfig.HTTPAuth.Enabled = true
			m.message = greenStyle.Render("http auth enabled")
			m.enterConfigStep(1)
			m.input.Placeholder = "username"
			return m, nil
		case "n":
			conf.AppConfig.HTTPAuth.Enabled = false
			m.message = greenStyle.Render("http auth disabled")
			m.view = ViewConfigMenu
			m.prevView = ViewMain
			m.configStep = 0
			return m, nil
		case "esc", "b":
			m.view = ViewConfigMenu
			m.prevView = ViewMain
			m.message = ""
			return m, nil
		}
		return m, nil
	}

	if !m.inputActive {
		return m.goBack()
	}

	switch msg.String() {
	case "enter":
		val := strings.TrimSpace(m.input.Value())
		if m.configStep == 1 {
			if val != "" {
				conf.AppConfig.HTTPAuth.Username = val
				m.message = greenStyle.Render("username → " + val)
			}
			m.enterConfigStep(2)
			m.input.Placeholder = "password"
			return m, nil
		}
		if val != "" {
			conf.AppConfig.HTTPAuth.Password = val
			m.message = greenStyle.Render("password updated")
		}
		m.inputActive = false
		m.input.Blur()
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		m.configStep = 0
		return m, nil
	case "esc":
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		m.inputActive = false
		m.input.Blur()
		m.input.SetValue("")
		m.message = ""
		m.configStep = 0
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m Model) handleConfigFTPToggleKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "y":
		conf.AppConfig.EnableFTP = !conf.AppConfig.EnableFTP
		m.message = greenStyle.Render("ftp server → " + fmt.Sprintf("%v", conf.AppConfig.EnableFTP))
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		return m, nil
	case "n", "esc", "b":
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		m.message = ""
		return m, nil
	}
	return m, nil
}

func (m Model) handleConfigMCPToggleKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "y":
		conf.AppConfig.EnableMCP = !conf.AppConfig.EnableMCP
		m.message = greenStyle.Render("mcp server → " + fmt.Sprintf("%v", conf.AppConfig.EnableMCP))
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		return m, nil
	case "n", "esc", "b":
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		m.message = ""
		return m, nil
	}
	return m, nil
}

func (m Model) handleConfigPortsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.inputActive {
		return m.goBack()
	}

	switch msg.String() {
	case "enter":
		val := strings.TrimSpace(m.input.Value())
		if m.configStep == 0 {
			if val != "" {
				conf.AppConfig.Port = ":" + val
				m.message = greenStyle.Render("http port → " + val)
			}
			m.enterConfigStep(1)
			m.input.Placeholder = "FTP port"
			return m, nil
		}
		if val != "" {
			conf.AppConfig.FTPPort = ":" + val
			m.message = greenStyle.Render("ftp port → " + val)
		}
		m.inputActive = false
		m.input.Blur()
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		m.configStep = 0
		return m, nil
	case "esc":
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		m.inputActive = false
		m.input.Blur()
		m.input.SetValue("")
		m.message = ""
		m.configStep = 0
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m Model) handleConfigSaveKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "y":
		if err := conf.SaveConfig(); err != nil {
			m.message = redStyle.Render("failed: " + err.Error())
		} else {
			currentPath := conf.GetConfigPath()
			if currentPath == "defaults" {
				currentPath = conf.GetDefaultConfigPath()
			}
			m.message = greenStyle.Render("saved → " + currentPath)
		}
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		return m, nil
	case "n", "esc", "b":
		m.view = ViewConfigMenu
		m.prevView = ViewMain
		m.message = ""
		return m, nil
	}
	return m, nil
}

var _ = syscall.SIGQUIT
