package game

import (
	"bytes"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

// AudioPlayer manages all sound-related operations.
type AudioPlayer struct {
	context *audio.Context
	sounds  map[string]*audio.Player
}

// NewAudioPlayer initializes the audio system.
func NewAudioPlayer(sampleRate int) *AudioPlayer {
	return &AudioPlayer{
		context: audio.NewContext(sampleRate),
		sounds:  make(map[string]*audio.Player),
	}
}

// LoadSound loads an MP3 file from disk and stores it in the map.
func (ap *AudioPlayer) LoadSound(name, filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read sound file: %v", err)
	}

	decoded, err := mp3.DecodeWithSampleRate(ap.context.SampleRate(), bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Failed to decode sound file: %v", err)
	}

	player, err := ap.context.NewPlayer(decoded)
	if err != nil {
		log.Fatalf("Failed to create audio player: %v", err)
	}

	ap.sounds[name] = player
}

// PlaySound plays a sound by name.
func (ap *AudioPlayer) PlaySound(name string) {
	if player, exists := ap.sounds[name]; exists {
		player.Rewind() // Restart the sound from the beginning
		player.Play()
	} else {
		log.Println("Sound not found:", name)
	}
}

// StopSound stops a sound if it's playing.
func (ap *AudioPlayer) StopSound(name string) {
	if player, exists := ap.sounds[name]; exists && player.IsPlaying() {
		player.Pause()
	}
}

// IsPlaying checks if a sound is currently playing.
func (ap *AudioPlayer) IsPlaying(name string) bool {
	if player, exists := ap.sounds[name]; exists {
		return player.IsPlaying()
	}
	return false
}
