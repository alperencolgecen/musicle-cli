// Package audio enumerates the system's audio output devices so the Settings
// → Sound tab can present them with a human-friendly name and a connection
// type (bluetooth vs wired). Detection is best-effort and degrades gracefully
// when the underlying tools (PulseAudio/PipeWire's pactl, or ALSA's aplay) are
// unavailable.
package audio

import (
	"os/exec"
	"strings"
)

// Device describes one output sink the user can pick in the Sound tab.
type Device struct {
	// Name is the backend identifier (PulseAudio sink name or ALSA PCM), kept
	// so a future routing layer can target it.
	Name string
	// Description is the friendly label shown in the UI.
	Description string
	// Type is "bluetooth" or "wired".
	Type string
	// Card is the ALSA card index when derivable, used for best-effort routing.
	Card string
}

// ListOutputDevices returns the available audio output devices. It prefers
// PulseAudio/PipeWire (pactl) and falls back to ALSA (aplay). An empty slice
// means no backend could be queried; callers should show a sensible default.
func ListOutputDevices() []Device {
	if devs := listPulseAudio(); len(devs) > 0 {
		return devs
	}
	if devs := listALSA(); len(devs) > 0 {
		return devs
	}
	return nil
}

// bluetoothGuess reports whether the given identifiers look like a Bluetooth
// sink based on well-known naming patterns.
func bluetoothGuess(name, description string) bool {
	hay := strings.ToLower(name + " " + description)
	for _, kw := range []string{"bluez", "bluetooth", "headset", "hands-free", "hfp", "a2dp"} {
		if strings.Contains(hay, kw) {
			return true
		}
	}
	return false
}

// listPulseAudio parses `pactl list sinks` (PulseAudio / PipeWire).
func listPulseAudio() []Device {
	path, err := exec.LookPath("pactl")
	if err != nil {
		return nil
	}
	out, err := exec.Command(path, "list", "sinks").Output()
	if err != nil {
		return nil
	}
	return parsePulseSinks(string(out))
}

func parsePulseSinks(raw string) []Device {
	var devs []Device
	var cur Device
	flush := func() {
		if cur.Name != "" {
			if cur.Description == "" {
				cur.Description = cur.Name
			}
			if cur.Type == "" {
				if bluetoothGuess(cur.Name, cur.Description) {
					cur.Type = "bluetooth"
				} else {
					cur.Type = "wired"
				}
			}
			devs = append(devs, cur)
		}
		cur = Device{}
	}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "Name: "):
			cur.Name = strings.TrimPrefix(line, "Name: ")
		case strings.HasPrefix(line, "Description: "):
			cur.Description = strings.TrimPrefix(line, "Description: ")
		case strings.HasPrefix(line, "device.description = "):
			d := strings.Trim(strings.TrimPrefix(line, "device.description = "), `"`)
			if cur.Description == "" && d != "" {
				cur.Description = d
			}
		case strings.HasPrefix(line, "device.class = "):
			c := strings.Trim(strings.TrimPrefix(line, "device.class = "), `"`)
			if strings.EqualFold(c, "sound") {
				// sink modules are sound; bluetooth made explicit below
			}
		case strings.HasPrefix(line, "device.product.name = "):
			p := strings.Trim(strings.TrimPrefix(line, "device.product.name = "), `"`)
			if bluetoothGuess(cur.Name, p) {
				cur.Type = "bluetooth"
			}
		case strings.HasPrefix(line, "device.api = "):
			api := strings.Trim(strings.TrimPrefix(line, "device.api = "), `"`)
			if strings.Contains(strings.ToLower(api), "bluez") {
				cur.Type = "bluetooth"
			}
		case strings.HasPrefix(line, "alsa.card = "):
			cur.Card = strings.Trim(strings.TrimPrefix(line, "alsa.card = "), `"`)
		}
	}
	flush()
	return devs
}

// listALSA parses `aplay -L` (ALSA). Bluetooth sinks usually appear as
// bluez_sink_* or bluealsa entries.
func listALSA() []Device {
	path, err := exec.LookPath("aplay")
	if err != nil {
		return nil
	}
	out, err := exec.Command(path, "-L").Output()
	if err != nil {
		return nil
	}
	var devs []Device
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, " ") {
			continue // skip indented hints
		}
		if strings.HasPrefix(line, "null") {
			continue
		}
		name := line
		typ := "wired"
		if bluetoothGuess(name, "") {
			typ = "bluetooth"
		}
		desc := name
		if idx := strings.Index(name, ","); idx >= 0 {
			desc = name[:idx]
		}
		devs = append(devs, Device{Name: name, Description: desc, Type: typ})
	}
	return devs
}

// RoutingEnv returns environment variables that, when set before the audio
// backend initialises, best-effort route playback to the chosen device. This
// is deliberately permissive: an unknown device simply falls back to the
// system default. Returns nil when no routing hint applies.
func RoutingEnv(name, card string) []string {
	var env []string
	if name != "" {
		lower := strings.ToLower(name)
		if strings.Contains(lower, "bluez") || strings.Contains(lower, "pulse") || strings.HasPrefix(lower, "alsa_output") {
			env = append(env, "PULSE_SINK="+name)
		}
	}
	if card != "" {
		env = append(env, "ALSA_PCM_CARD="+card)
	}
	return env
}
