// Package cli 提供终端交互界面(TUI),让人能直观地浏览共享文件、
// 查看访问地址与二维码,用键盘选中文件即可在手机上扫码下载。
package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/skip2/go-qrcode"

	"github.com/zy84338719/upftp/internal/core"
)

// 视图状态
const (
	viewList = iota
	viewQR
)

// fileItem 是文件列表中的一项。
type fileItem struct {
	name  string
	isDir bool
	url   string // 选中后展示的下载链接
}

func (f fileItem) Title() string {
	if f.isDir {
		return "📁 " + f.name
	}
	return "📄 " + f.name
}
func (f fileItem) Description() string {
	if f.isDir {
		return "目录"
	}
	if f.url != "" {
		return f.url
	}
	return "文件"
}
func (f fileItem) FilterValue() string { return f.name }

// Model 是 TUI 的状态。
type Model struct {
	files    *core.Files
	list     list.Model
	qrItem   fileItem
	qrCode   string
	view     int
	host     string
	httpPort int
	ftpURL   string
	width    int
	height   int
}

// NewModel 构造初始 Model。
func NewModel(files *core.Files, host string, httpPort int, ftpURL string) Model {
	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{}, delegate, 80, 20)
	l.Title = "upftp · 共享文件"
	l.SetShowHelp(true)
	return Model{
		files: files, list: l,
		host: host, httpPort: httpPort, ftpURL: ftpURL,
	}
}

// loadDir 加载一个目录到列表。
func (m *Model) loadDir(rel string) {
	entries, err := m.files.List(rel)
	if err != nil {
		m.list.SetItems([]list.Item{fileItem{name: "无法读取目录: " + err.Error()}})
		return
	}
	items := make([]list.Item, 0, len(entries))
	for _, e := range entries {
		dlURL := ""
		if !e.IsDir {
			dlURL = fmt.Sprintf("http://%s:%d/api/download?path=%s/%s",
				m.host, m.httpPort, strings.TrimPrefix(rel, "/"), e.Name)
		}
		items = append(items, fileItem{name: e.Name, isDir: e.IsDir, url: dlURL})
	}
	m.list.SetItems(items)
	m.list.Title = "upftp · " + rel
}

// Init 启动钩子。
func (m Model) Init() tea.Cmd { return nil }

// Update 处理按键。
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.view == viewQR {
				m.view = viewList
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			if m.view == viewList {
				item, ok := m.list.SelectedItem().(fileItem)
				if !ok {
					return m, nil
				}
				if item.isDir {
					// 进入子目录(简化:仅支持一层,实际可扩展)
					m.loadDir("/" + item.name)
				} else if item.url != "" {
					m.qrItem = item
					m.qrCode = renderQR(item.url)
					m.view = viewQR
				}
			}
		case "esc":
			m.view = viewList
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View 渲染。
func (m Model) View() string {
	if m.view == viewQR {
		return m.viewQRCode()
	}
	return m.viewList()
}

func (m Model) viewList() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5b9dff")).Render("upftp")
	info := fmt.Sprintf("FTP: %s   |   Web: http://%s:%d", m.ftpURL, m.host, m.httpPort)
	infoStyled := lipgloss.NewStyle().Faint(true).Render(info)
	return header + "\n" + infoStyled + "\n" + m.list.View()
}

func (m Model) viewQRCode() string {
	title := lipgloss.NewStyle().Bold(true).Render("扫码下载 · " + m.qrItem.name)
	urlStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("#5b9dff")).Render(m.qrItem.url)
	hint := lipgloss.NewStyle().Faint(true).Render("按 q 或 esc 返回")
	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n\n%s\n", title, m.qrCode, urlStyled, hint)
}

// renderQR 把 URL 渲染成终端可显示的二维码(Unicode 半字符块)。
func renderQR(url string) string {
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return "(二维码生成失败: " + err.Error() + ")"
	}
	return qr.ToString(false)
}

// Run 在当前终端启动 TUI;非 TTY 环境静默跳过。
func Run(ctx context.Context, files *core.Files, host string, httpPort int, ftpURL string) {
	if !isTTY() {
		return
	}
	m := NewModel(files, host, httpPort, ftpURL)
	m.loadDir("/")
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "TUI 启动失败:", err)
	}
}

func isTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
