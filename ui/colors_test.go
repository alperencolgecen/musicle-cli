package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// Terminal arka planı (uygulamanın kendi penceresi) saf siyah olmalı; böylece
// 3 işletim sisteminde de kullanıcının terminal temasından bağımsız olarak
// uygulama siyah arka planla açılır.
func TestColorBackgroundIsBlack(t *testing.T) {
	if ColorBackground != lipgloss.Color("#000000") {
		t.Fatalf("ColorBackground = %q, beklenen #000000", ColorBackground)
	}
	InitStyles()
	if got := AppStyle.GetBackground(); got != lipgloss.Color("#000000") {
		t.Fatalf("AppStyle arka planı = %v, beklenen #000000", got)
	}
}
