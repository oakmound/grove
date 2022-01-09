package sound

import (
	"fmt"
	"sync"

	"github.com/oakmound/oak/v3/audio"
	"github.com/oakmound/oak/v3/audio/font"
	"github.com/oakmound/oak/v3/event"
)

var (
	Musics = map[string]*audio.Audio{}

	musicFileLock sync.RWMutex
	musicFiles    = map[string]*font.Font{}
)

func RegisterMusic(file string, f *font.Font) {
	musicFileLock.Lock()
	musicFiles[file] = f
	musicFileLock.Unlock()
}

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

func PlayMusic(s string) error {
	StopMusic()
	if PlayingMusicLabel == s {
		return nil
	}
	PlayingMusicLabel = s
	audOrigin, ok := Musics[s]
	if !ok {
		return fmt.Errorf("Tried to play unloaded Audio %q", s)
	}
	audOrigin.Play()
	PlayingMusic = audOrigin
	return nil
}

func StopMusic() {
	if PlayingMusic != nil {
		PlayingMusic.Stop()
	}
}

// only one music can be playing
var PlayingMusic *audio.Audio
var PlayingMusicLabel string

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
