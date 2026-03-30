package cli

import (
	"fmt"
	"strings"

	"github.com/zy84338719/upftp/pkg/conf"
)

func (m Model) renderStatus() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "  "+labelStyle.Render(" "+t(lang, "server_config")+":"))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "http_port")+":"))+" "+fmt.Sprintf("%d", conf.AppConfig.GetHTTPPort()))
	if conf.AppConfig.EnableFTP {
		lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_port")+":"))+" "+fmt.Sprintf("%d", conf.AppConfig.GetFTPPort()))
	}
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "root_directory")+":"))+" "+truncatePath(conf.AppConfig.Root, 36))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render(" "+t(lang, "features_status")+":"))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "http_auth")+":"))+" "+renderEnabled(conf.AppConfig.HTTPAuth.Enabled))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "file_upload_feat")+":"))+" "+renderEnabled(conf.AppConfig.Upload.Enabled))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_server_feat")+":"))+" "+renderEnabled(conf.AppConfig.EnableFTP))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "https")+":"))+" "+renderEnabled(conf.AppConfig.HTTPS.Enabled))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render(" "+t(lang, "statistics")+":"))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "files_available")+":"))+" "+fmt.Sprintf("%d", len(m.fileMap)))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "upload_max_size")+":"))+" "+formatSize(conf.AppConfig.Upload.MaxSize))

	return "  " + dimStyle.Render(t(lang, "server_status_box")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderVersion() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "version_label")+":"))+" "+conf.AppConfig.Version)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "git_commit")+":"))+" "+conf.AppConfig.LastCommit)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "build_date")+":"))+" "+conf.AppConfig.BuildDate)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "go_version")+":"))+" "+conf.AppConfig.GoVersion)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "platform")+":"))+" "+conf.AppConfig.Platform)
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "project_homepage")+":"))+" "+greenStyle.Render(truncateString(conf.AppConfig.ProjectURL, 40)))
	lines = append(lines, "")

	return "  " + dimStyle.Render(t(lang, "version_info_box")) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderAbout() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "")
	lines = append(lines, "  "+greenStyle.Render(conf.AppConfig.ProjectName)+" "+t(lang, "about_desc"))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "key_features")+":"))
	lines = append(lines, "    ● "+t(lang, "feat_web_ui"))
	lines = append(lines, "    ● "+t(lang, "feat_mcp"))
	lines = append(lines, "    ● "+t(lang, "feat_upload"))
	lines = append(lines, "    ● "+t(lang, "feat_http_auth"))
	lines = append(lines, "    ● "+t(lang, "feat_ftp"))
	lines = append(lines, "    ● "+t(lang, "feat_qr"))
	lines = append(lines, "    ● "+t(lang, "feat_tree"))
	lines = append(lines, "    ● "+t(lang, "feat_lang"))
	lines = append(lines, "    ● "+t(lang, "feat_yaml"))
	lines = append(lines, "    ● "+t(lang, "feat_https"))
	lines = append(lines, "    ● "+t(lang, "feat_cli"))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "links")+":"))
	lines = append(lines, "    Homepage:  "+truncateString(conf.AppConfig.ProjectURL, 47))
	lines = append(lines, "    GitHub:    https://github.com/zy84338719/upftp")
	lines = append(lines, "")
	lines = append(lines, "  "+dimStyle.Render(t(lang, "license")))
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(t(lang, "build_info")+":"))
	lines = append(lines, "    "+t(lang, "version_label")+": "+conf.AppConfig.Version)
	lines = append(lines, "    "+t(lang, "git_commit")+": "+conf.AppConfig.LastCommit)
	lines = append(lines, "    "+t(lang, "build_date")+": "+conf.AppConfig.BuildDate)
	lines = append(lines, "    "+t(lang, "platform")+": "+conf.AppConfig.Platform)
	lines = append(lines, "")

	return "  " + dimStyle.Render(t(lang, "about_box")+" v"+conf.AppConfig.Version) + "\n" + panel(w, strings.Join(lines, "\n"))
}

func (m Model) renderFTPInfo() string {
	lang := m.lang
	w := m.panelWidth()

	var lines []string
	lines = append(lines, "")
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_server")+":"))+" "+m.serverIP)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "ftp_port")+":"))+" "+fmt.Sprintf("%d", conf.AppConfig.GetFTPPort()))
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "username")+":"))+" "+conf.AppConfig.Username)
	lines = append(lines, "  "+labelStyle.Render(fmt.Sprintf(" %-14s", t(lang, "password")+":"))+" "+conf.AppConfig.Password)
	lines = append(lines, "")
	lines = append(lines, "  $ ftp "+m.serverIP+":"+fmt.Sprintf("%d", conf.AppConfig.GetFTPPort()))
	lines = append(lines, "  Name: "+conf.AppConfig.Username)
	lines = append(lines, "  Password: "+conf.AppConfig.Password)
	lines = append(lines, "  ftp> ls")
	lines = append(lines, "  ftp> get filename")
	lines = append(lines, "  ftp> put localfile")
	lines = append(lines, "")

	return "  " + dimStyle.Render(t(lang, "ftp_connection_info")) + "\n" + panel(w, strings.Join(lines, "\n"))
}
