package sound

import (
	"fmt"
	"sync"

	"github.com/oakmound/oak/v3/audio"
	"github.com/oakmound/oak/v3/audio/font"
	"github.com/oakmound/oak/v3/event"
)

var (
	// Musics is a mapping of music files with an Audio to play.
	Musics = map[string]*audio.Audio{}

	musicFileLock sync.RWMutex
	musicFiles    = map[string]*font.Font{}
)

// RegisterMusic and make it ready for playback.
// Specify the file to load it from and the music font to use with it.
func RegisterMusic(file string, f *font.Font) {
	musicFileLock.Lock()
	musicFiles[file] = f
	musicFileLock.Unlock()
}

// ReloadMusicAssets and error if the reload from files fails somehow.
func ReloadMusicAssets() error {
	musicFileLock.Lock()
	defer musicFileLock.Unlock()
	for s, f := range musicFiles {
		a, err := audio.Get(s)
		if err != nil {
			return err
		}
		Musics[s] = audio.New(f, a)
	}
	return nil
}

// PlayMusic stops any playing music and then tries to play the specified music file.
func PlayMusic(s string) error {
	if PlayingMusicLabel == s {
		return nil
	}
	StopMusic()
	PlayingMusicLabel = s
	audOrigin, ok := Musics[s]
	if !ok {
		return fmt.Errorf("Tried to play unloaded Audio %q", s)
	}
	audOrigin.Play()
	PlayingMusic = audOrigin
	return nil
}

// StopMusic if any music is playing.
func StopMusic() {
	if PlayingMusic != nil {
		PlayingMusic.Stop()
	}
}

// only one music can be playing

// PlayingMusic is the currently playing audio.
var PlayingMusic *audio.Audio

// PlayingMusicLabel is the playing music's string name.
var PlayingMusicLabel string

// SetMusicVolume to the provided value.
func SetMusicVolume(newVolume float64) {
	musicVolume = newVolume
	updateMusicVolume(volume, musicVolume)
	event.Trigger(EventMusicVolumeChanged, newVolume)
}

func updateMusicVolume(volume, musicVolume float64) {
	newVolume := volume * musicVolume
	if newVolume > 1 || newVolume < 0 {
		newVolume = 0
	}
	scalar := convertVolumeScale(newVolume)
	for _, m := range Musics {
		m.SetVolume(scalar)
	}
}
