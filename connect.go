package main

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"MusicLeCLI/bridge"
	"MusicLeCLI/internal/browser"
	"MusicLeCLI/state"
	"MusicLeCLI/ui"
)

// ConnectModel is the platform-selection screen of the browser connector.
// It shows two brand-colored cards (Spotify / YouTube Music) between which the
// user navigates with F1 or the arrow keys. After a choice it runs a browser
// scan behind a loading modal and then lists the discovered playlists.
type ConnectModel struct {
	width  int
	height int
	focus  int // 0 = Spotify, 1 = YouTube Music
	chosen browser.Platform

	scanning  bool
	scanStart time.Time
	scanDone  bool
	playlists []browser.Playlist
	scanErr   error

	plFocus   int
	confirmed map[int]bool
	saved     bool
	saveErr   error
}

func NewConnectModel() *ConnectModel {
	return &ConnectModel{focus: 0}
}

func (m *ConnectModel) Init() tea.Cmd { return nil }

func (m *ConnectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch {
		case m.scanning:
			return m, nil

		case m.scanDone && m.scanErr == nil:
			switch k.String() {
			case "up", "k":
				if m.plFocus > 0 {
					m.plFocus--
				}
				return m, nil
			case "down", "j":
				if m.plFocus < len(m.playlists)-1 {
					m.plFocus++
				}
				return m, nil
			case "enter", " ":
				if !m.confirmed[m.plFocus] {
					if err := m.confirm(m.plFocus); err != nil {
						m.saveErr = err
					}
				}
				if m.saved {
					return m, func() tea.Msg { return connectDoneMsg{} }
				}
				return m, nil
			case "esc":
				if m.saved {
					return m, func() tea.Msg { return connectDoneMsg{} }
				}
				m.scanDone = false
				m.playlists = nil
				m.plFocus = 0
				m.confirmed = nil
				return m, nil
			}

		default:
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
				m.scanning = true
				m.scanStart = time.Now()
				return m, scanConnectCmd(m.chosen)
			}
		}
	}
	return m, nil
}

// confirm saves the playlist at index i: it derives an auto-named profile for
// the chosen platform and imports the tracks via bridge.ImportFromBrowser.
func (m *ConnectModel) confirm(i int) error {
	if m.confirmed == nil {
		m.confirmed = map[int]bool{}
	}
	if m.confirmed[i] {
		return nil
	}
	pl := m.playlists[i]
	displayName, folderName := browser.NextProfileName(state.Current.Profiles, m.chosen, state.Current.ProfilesDir())
	if err := state.Current.CreateProfileStructure(folderName, displayName, "", "", state.Current.Language); err != nil {
		return err
	}
	plFolder := uniqueFolder(browser.Slugify(pl.Name), folderName)
	if err := bridge.ImportFromBrowser(folderName, plFolder, pl.Name, pl.Tracks); err != nil {
		return err
	}
	if err := state.Current.ScanProfiles(); err != nil {
		return err
	}
	for idx := range state.Current.Profiles {
		if state.Current.Profiles[idx].FolderName == folderName {
			state.Current.CurrentProfile = &state.Current.Profiles[idx]
			if len(state.Current.CurrentProfile.Playlists) > 0 {
				state.Current.CurrentPlaylist = &state.Current.CurrentProfile.Playlists[0]
			}
			break
		}
	}
	m.confirmed[i] = true
	m.saved = true
	return nil
}

// uniqueFolder returns a playlist folder name that does not yet exist under the
// given profile directory.
func uniqueFolder(name, profileFolder string) string {
	base := name
	i := 2
	for {
		if _, err := os.Stat(state.Current.PlaylistDir(profileFolder, name)); os.IsNotExist(err) {
			break
		}
		name = fmt.Sprintf("%s_%d", base, i)
		i++
	}
	return name
}

// finishScan stores the result of a browser scan and dismisses the modal.
func (m *ConnectModel) finishScan(pls []browser.Playlist, err error) {
	m.scanning = false
	m.scanDone = true
	m.playlists = pls
	m.scanErr = err
}

// scanConnectCmd runs the real CDP scan but guarantees the loading modal is
// visible for at least 2 seconds (simulated API latency).
func scanConnectCmd(platform browser.Platform) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		pls, err := browser.Connect(platform)
		if d := time.Since(start); d < 2*time.Second {
			time.Sleep(2*time.Second - d)
		}
		return ConnectResultMsg{Platform: platform, Playlists: pls, Err: err}
	}
}

func (m *ConnectModel) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	if m.scanning {
		return m.renderModal()
	}
	if m.scanDone {
		if m.scanErr != nil {
			return lipgloss.PlaceVertical(m.height, lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.ColorError).Render(
					fmt.Sprintf("Bağlantı hatası: %v", m.scanErr)))
		}
		return m.renderPlaylistList()
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

	footer := lipgloss.NewStyle().Foreground(ui.ColorSecondary).Render(
		"Spotify veya YouTube Music sekmesini tarayıcıda açık tutun.")

	return lipgloss.JoinVertical(lipgloss.Center, title, "", row, "", footer)
}

// renderModal shows the "reading browser info" loader while scanning.
func (m *ConnectModel) renderModal() string {
	frames := []string{"|", "/", "-", "\\"}
	spin := frames[(time.Now().Nanosecond()/100000000)%4]
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorAccent).
		Padding(1, 4).
		Align(lipgloss.Center)
	content := lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Foreground(ui.ColorPrimary).Bold(true).Render("Tarayıcı bilgisi alınıyor..."),
		"",
		lipgloss.NewStyle().Foreground(ui.ColorAccent).Render(spin+" "+string(m.chosen)+" taranıyor..."),
	)
	modal := box.Render(content)
	return lipgloss.PlaceVertical(m.height, lipgloss.Center, modal)
}

// renderPlaylistList shows the discovered playlists, each with an "Onayla"
// button. The focused row is confirmed with Enter, which advances to the next
// unconfirmed row so playlists can be approved sequentially.
func (m *ConnectModel) renderPlaylistList() string {
	header := lipgloss.NewStyle().Foreground(ui.ColorPrimary).Bold(true).
		Render(fmt.Sprintf("%s — %d playlist bulundu", m.chosen, len(m.playlists)))

	var rows []string
	for i, pl := range m.playlists {
		cursor := "  "
		if i == m.plFocus {
			cursor = "▶ "
		}
		name := cursor + pl.Name
		count := lipgloss.NewStyle().Foreground(ui.ColorSecondary).
			Render(fmt.Sprintf("(%d şarkı)", pl.TrackCount))

		var action string
		switch {
		case m.confirmed[i]:
			action = lipgloss.NewStyle().Foreground(ui.ColorSuccess).Render("[Onaylandı]")
		case i == m.plFocus:
			action = ui.AccentButtonStyle.Render(" Onayla ")
		default:
			action = ui.ButtonStyle.Render(" Onayla ")
		}
		row := lipgloss.JoinHorizontal(lipgloss.Center, name, "   ", count, "   ", action)
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		rows = append(rows, lipgloss.NewStyle().Foreground(ui.ColorSecondary).Render("(playlist bulunamadı)"))
	}

	body := strings.Join(rows, "\n")

	var footer string
	switch {
	case m.saveErr != nil:
		footer = lipgloss.NewStyle().Foreground(ui.ColorError).
			Render(fmt.Sprintf("Kayıt hatası: %v", m.saveErr))
	case m.saved:
		footer = lipgloss.NewStyle().Foreground(ui.ColorSuccess).
			Render("Kaydedildi ✓  •  Enter/Esc ile ana ekrana dön")
	default:
		footer = lipgloss.NewStyle().Foreground(ui.ColorSecondary).
			Render("↑/↓ ile gezin  •  Enter ile onayla  •  Esc ile geri")
	}

	return lipgloss.JoinVertical(lipgloss.Center, header, "", body, "", footer)
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

	logo := renderPNGLogo(logoPath, w-8, h-8)
	labelColor := base
	if focused {
		labelColor = focus
	}
	label := lipgloss.NewStyle().Foreground(labelColor).Bold(true).Render(name)
	content := lipgloss.JoinVertical(lipgloss.Center, logo, "", label)
	return style.Render(content)
}

// renderPNGLogo renders a PNG as half-block ANSI art to fit within cols x rows,
// reusing the package-level scaleImage helper used for playlist art.
func renderPNGLogo(path string, cols, rows int) string {
	if cols < 4 || rows < 3 {
		return ""
	}
	data, err := logoFS.ReadFile(path)
	if err != nil {
		// Fallback to the filesystem (e.g. when running from the repo root).
		data, err = os.ReadFile(path)
		if err != nil {
			return ""
		}
	}
	img, _, err := image.Decode(bytes.NewReader(data))
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
