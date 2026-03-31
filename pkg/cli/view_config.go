package cli

import (
	"fmt"
	"strings"

	"github.com/zy84338719/upftp/pkg/conf"
)

func (m Model) renderConfigMenu() string {
	lang := m.lang
	w := m.panelWidth()

	var rows []string
	rows = append(rows, "  "+keyStyle.Render("[1]")+" "+t(lang, "credentials_menu"))
	rows = append(rows, "  "+keyStyle.Render("[2]")+" "+t(lang, "http_auth_menu")+" "+yellowStyle.Render(renderEnabled(conf.AppConfig.HTTPAuth.Enabled)))
	rows = append(rows, "  "+keyStyle.Render("[3]")+" "+t(lang, "ftp_server_menu")+" "+yellowStyle.Render(renderEnabled(conf.AppConfig.EnableFTP)))
	rows = append(rows, "  "+keyStyle.Render("[4]")+" "+t(lang, "mcp_server_menu")+" "+yellowStyle.Render(renderEnabled(conf.AppConfig.EnableMCP)))
	rows = append(rows, "  "+keyStyle.Render("[5]")+" "+t(lang, "server_ports_menu"))
	rows = append(rows, "  "+keyStyle.Render("[6]")+" "+t(lang, "save_config_menu"))
	rows = append(rows, "  "+keyStyle.Render("[7]")+" "+t(lang, "view_config_menu"))
	rows = append(rows, "  "+greenStyle.Render("[b]")+" "+t(lang, "back_main"))
	rows = append(rows, "  "+redStyle.Render("[q]")+" "+t(lang, "quit"))

	return "  " + dimStyle.Render(t(lang, "config_menu")) + "\n" + panel(w, strings.Join(rows, "\n"))
}

func (m Model) renderConfigFTP() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "username")+":"))+" "+conf.AppConfig.Username)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "password")+":"))+" "+maskPassword(conf.AppConfig.Password))
	lines = append(lines, "")
	lines = append(lines, "  "+dimStyle.Render(t(lang, "credentials_hint")))
	lines = append(lines, "")

	if m.configStep == 0 {
		lines = append(lines, "  "+labelStyle.Render(t(lang, "username")+":")+" "+m.input.View())
	} else {
		lines = append(lines, "  "+labelStyle.Render(t(lang, "password")+":")+" "+m.input.View())
	}

	return "  " + dimStyle.Render(t(lang, "credentials_menu")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderConfigHTTP() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "http_auth")+":"))+" "+renderEnabled(conf.AppConfig.HTTPAuth.Enabled))
	if conf.AppConfig.HTTPAuth.Enabled {
		lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "username")+":"))+" "+conf.AppConfig.HTTPAuth.Username)
		lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "password")+":"))+" "+maskPassword(conf.AppConfig.HTTPAuth.Password))
	}
	lines = append(lines, "")

	if m.configStep == 0 {
		lines = append(lines, "  "+labelStyle.Render("Enable HTTP auth?")+" "+yellowStyle.Render("[y/n]"))
	} else if m.configStep == 1 {
		lines = append(lines, "  "+labelStyle.Render(t(lang, "username")+":")+" "+m.input.View())
	} else {
		lines = append(lines, "  "+labelStyle.Render(t(lang, "password")+":")+" "+m.input.View())
	}

	return "  " + dimStyle.Render(t(lang, "http_auth_menu")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderConfigFTPToggle() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_server_feat")+":"))+" "+renderEnabled(conf.AppConfig.EnableFTP))
	lines = append(lines, "  "+dimStyle.Render(t(lang, "restart_required")))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render("Toggle FTP server?")+" "+yellowStyle.Render("[y/n]"))

	return "  " + dimStyle.Render(t(lang, "ftp_server_menu")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderConfigMCPToggle() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "mcp_server")+":"))+" "+renderEnabled(conf.AppConfig.EnableMCP))
	lines = append(lines, "  "+dimStyle.Render(t(lang, "restart_required")))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render("Toggle MCP server?")+" "+yellowStyle.Render("[y/n]"))

	return "  " + dimStyle.Render(t(lang, "mcp_server_menu")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderConfigPorts() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "http_port")+":"))+" "+fmt.Sprintf("%d", conf.AppConfig.GetHTTPPort()))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_port")+":"))+" "+fmt.Sprintf("%d", conf.AppConfig.GetFTPPort()))
	lines = append(lines, "  "+dimStyle.Render(t(lang, "restart_required")))
	lines = append(lines, "")

	if m.configStep == 0 {
		lines = append(lines, "  "+labelStyle.Render(t(lang, "http_port")+":")+" "+m.input.View())
	} else {
		lines = append(lines, "  "+labelStyle.Render(t(lang, "ftp_port")+":")+" "+m.input.View())
	}

	return "  " + dimStyle.Render(t(lang, "server_ports_menu")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderConfigSave() string {
	lang := m.lang
	w := m.panelWidth()

	currentPath := conf.GetConfigPath()
	if currentPath == "defaults" {
		currentPath = conf.GetDefaultConfigPath()
	}

	var lines []string
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "config_file")+":"))+" "+truncateString(currentPath, 41))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render("Save to "+currentPath+"?")+" "+yellowStyle.Render("[y/n]"))

	return "  " + dimStyle.Render(t(lang, "save_config_menu")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderConfigView() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "  "+labelStyle.Render(t(lang, "server_settings")+":"))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "http_port")+":"))+" "+fmt.Sprintf("%d", conf.AppConfig.GetHTTPPort()))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_port")+":"))+" "+fmt.Sprintf("%d", conf.AppConfig.GetFTPPort()))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "root_directory")+":"))+" "+truncatePath(conf.AppConfig.Root, 36))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "auto_select_ip")+":"))+" "+fmt.Sprintf("%v", conf.AppConfig.AutoSelect))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "credentials")+":"))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "username")+":"))+" "+conf.AppConfig.Username)
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "password")+":"))+" "+maskPassword(conf.AppConfig.Password))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "ftp_server")+":"))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "enabled")+":"))+" "+renderEnabled(conf.AppConfig.EnableFTP))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "mcp_server")+":"))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "enabled")+":"))+" "+renderEnabled(conf.AppConfig.EnableMCP))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "http_authentication")+":"))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "enabled")+":"))+" "+renderEnabled(conf.AppConfig.HTTPAuth.Enabled))
	if conf.AppConfig.HTTPAuth.Enabled {
		lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "username")+":"))+" "+conf.AppConfig.HTTPAuth.Username)
		lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "password")+":"))+" "+maskPassword(conf.AppConfig.HTTPAuth.Password))
	}
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "upload_settings")+":"))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "enabled")+":"))+" "+renderEnabled(conf.AppConfig.Upload.Enabled))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "upload_max_size")+":"))+" "+formatSize(conf.AppConfig.Upload.MaxSize))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "logging")+":"))
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "level")+":"))+" "+conf.AppConfig.Logging.Level)
	lines = append(lines, "    "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "format")+":"))+" "+conf.AppConfig.Logging.Format)
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "config_file")+":"))+" "+truncateString(conf.GetConfigPath(), 36))

	return "  " + dimStyle.Render(t(lang, "current_config")) + "\n" + panel(w, strings.Join(lines, "\n")) + "\n  " + dimStyle.Render(t(lang, "note_save"))
}
