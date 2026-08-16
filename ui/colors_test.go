package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// Uygulama artık arka planı zorlamaz; terminalin kendi teması kullanılır.
// AppStyle, tüm ekranı siyaha boyayan bir zemin özelliğine sahip olmamalı.
func TestAppStyleDoesNotForceBackground(t *testing.T) {
	InitStyles()
	empty := lipgloss.NoColor{}
	if got := AppStyle.GetBackground(); got != empty {
		t.Fatalf("AppStyle arka planı = %v, beklenen boş (terminal teması)", got)
	}
}
