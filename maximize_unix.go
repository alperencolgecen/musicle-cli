//go:build !windows

package main

import "os"

func maximizeTerminal() {
	// dtterm/xterm pencere büyütme sekansı (9;2t = maximize). Çoğu modern
	// terminalde pencereyi ekranı kaplayacak şekilde büyütür; desteklemeyen
	// terminallerde yok sayılır. Alt ekran tamponu (tea.WithAltScreen) zaten
	// tüm terminali kapladığından uygulama içeriği yine de tam ekran görünür.
	os.Stdout.Write([]byte("\033[9;2t"))
}
