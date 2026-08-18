package main

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"MusicLeCLI/internal/browser"
	"MusicLeCLI/ui"
)

// ConnectModel is the platform-selection screen of the browser connector.
// It shows two brand-colored cards (Spotify / YouTube Music) between which the
// user navigates with F1 or the arrow keys.
type ConnectModel struct {
	width    int
	height   int
	focus    int // 0 = Spotify, 1 = YouTube Music
	chosen   browser.Platform
	confirmed bool
}

func NewConnectModel() *ConnectModel {
	return &ConnectModel{focus: 0}
}

func (m *ConnectModel) Init() tea.Cmd { return nil }

func (m *ConnectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "left", "up", "shift+tab":
			m.focus = 0
			return m, nil
		case "right", "down", "tab":
			m.focus = 1
			return m, nil
		case "enter", " ":
			if m.focus == 0 {
				m.chosen = browser.PlatformSpotify
			} else {
				m.chosen = browser.PlatformYouTube
			}
			m.confirmed = true
			return m, nil
		}
	}
	return m, nil
}

func (m *ConnectModel) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	title := lipgloss.NewStyle().Foreground(ui.ColorPrimary).Bold(true).
		Render("Tarayıcı Bağlayıcı — bir platform seçin   •   F1 ile geçiş   •   Enter ile onayla")

	cardW := (m.width - 10) / 2
	if cardW < 24 {
		cardW = 24
	}
	cardH := m.height - 8
	if cardH < 14 {
		cardH = 14
	}

	spotify := m.renderCard("Spotify", "assets/Spotify_logo.png", ui.ColorSpotify, ui.ColorSpotifyFocus, m.focus == 0, cardW, cardH)
	yt := m.renderCard("YouTube Music", "assets/Youtube_logo.png", ui.ColorYouTube, ui.ColorYouTubeFocus, m.focus == 1, cardW, cardH)

	row := lipgloss.JoinHorizontal(lipgloss.Top, spotify, "    ", yt)

	var footer string
	if m.confirmed {
		footer = lipgloss.NewStyle().Foreground(ui.ColorAccent).Render(
			fmt.Sprintf("Seçildi: %s — tarayıcı taraması sonraki adımda başlatılacak.", m.chosen))
	} else {
		footer = lipgloss.NewStyle().Foreground(ui.ColorSecondary).Render(
			"Spotify veya YouTube Music sekmesini tarayıcıda açık tutun.")
	}

	return lipgloss.JoinVertical(lipgloss.Center, title, "", row, "", footer)
}

func (m *ConnectModel) renderCard(name, logoPath string, base, focus lipgloss.Color, focused bool, w, h int) string {
	border := base
	if focused {
		border = focus
	}
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Width(w).
		Height(h).
		Align(lipgloss.Center, lipgloss.Center)
	if focused {
		style = style.Bold(true)
	}

	logo := renderLogo(logoPath, w-8, h-8)
	labelColor := base
	if focused {
		labelColor = focus
	}
	label := lipgloss.NewStyle().Foreground(labelColor).Bold(true).Render(name)
	content := lipgloss.JoinVertical(lipgloss.Center, logo, "", label)
	return style.Render(content)
}

// renderLogo renders a PNG as half-block ANSI art to fit within cols x rows,
// reusing the package-level scaleImage helper used for playlist art.
func renderLogo(path string, cols, rows int) string {
	if cols < 4 || rows < 3 {
		return ""
	}
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return ""
	}
	resized := scaleImage(img, cols, rows*2)
	var out strings.Builder
	for cy := 0; cy < rows; cy++ {
		for cx := 0; cx < cols; cx++ {
			r1, g1, b1, a1 := resized.At(cx, cy*2).RGBA()
			r2, g2, b2, a2 := resized.At(cx, cy*2+1).RGBA()
			if a1 < 128 && a2 < 128 {
				out.WriteByte(' ')
				continue
			}
			if a2 < 128 {
				out.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm▀\033[0m", r1>>8, g1>>8, b1>>8))
				continue
			}
			if a1 < 128 {
				out.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm▄\033[0m", r2>>8, g2>>8, b2>>8))
				continue
			}
			out.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%d;48;2;%d;%d;%dm▄\033[0m", r2>>8, g2>>8, b2>>8, r1>>8, g1>>8, b1>>8))
		}
		out.WriteByte('\n')
	}
	return out.String()
}
