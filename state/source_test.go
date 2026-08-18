package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSongSourceRoundTrip_WriteRead(t *testing.T) {
	dir := t.TempDir()
	listPath := filepath.Join(dir, "song_list.txt")

	songs := []Song{
		{Filename: "abc123.mp3", Title: "Song A", Artist: "Artist A", DateAdded: "2024-01-01", Duration: "03:00", Source: "https://music.youtube.com/watch?v=abc123"},
		{Filename: "def456.mp3", Title: "Song B", Artist: "Artist B", DateAdded: "2024-01-02", Duration: "04:00", Source: "Artist B - Song B"},
	}
	if err := WriteSongs(listPath, songs); err != nil {
		t.Fatalf("WriteSongs: %v", err)
	}

	got, err := ReadSongs(listPath)
	if err != nil {
		t.Fatalf("ReadSongs: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d songs, want 2", len(got))
	}
	if got[0].Source != songs[0].Source {
		t.Errorf("Source[0] = %q, want %q", got[0].Source, songs[0].Source)
	}
	if got[1].Source != songs[1].Source {
		t.Errorf("Source[1] = %q, want %q", got[1].Source, songs[1].Source)
	}
}

func TestParseSongList_SixField(t *testing.T) {
	dir := t.TempDir()
	listPath := filepath.Join(dir, "song_list.txt")
	content := "vid.mp3|Title|Artist|2024-01-01|03:30|https://music.youtube.com/watch?v=vid\n"
	if err := os.WriteFile(listPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	songs := parseSongList(listPath, dir)
	if len(songs) != 1 {
		t.Fatalf("got %d songs, want 1", len(songs))
	}
	if songs[0].Source != "https://music.youtube.com/watch?v=vid" {
		t.Errorf("Source = %q", songs[0].Source)
	}
	if songs[0].FilePath != filepath.Join(dir, "vid.mp3") {
		t.Errorf("FilePath = %q", songs[0].FilePath)
	}
}

func TestParseSongList_FiveFieldBackwardCompat(t *testing.T) {
	dir := t.TempDir()
	listPath := filepath.Join(dir, "song_list.txt")
	content := "a.mp3|Title|Artist|2024-01-01|03:30\n"
	if err := os.WriteFile(listPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	songs := parseSongList(listPath, dir)
	if len(songs) != 1 {
		t.Fatalf("got %d songs, want 1", len(songs))
	}
	if songs[0].Source != "" {
		t.Errorf("Source = %q, want empty (backward compat)", songs[0].Source)
	}
}
