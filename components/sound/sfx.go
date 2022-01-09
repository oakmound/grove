package sound

import (
	"fmt"
	"sync"
	"math/rand"

	"github.com/oakmound/oak/v3/audio"
	"github.com/oakmound/oak/v3/audio/font"
	"github.com/oakmound/oak/v3/event"
)

var (
	SFXs = map[string]*audio.Audio{}

	sfxFileLock sync.RWMutex
	sfxFiles    = map[string]*font.Font{}
)

func RegisterSFX(file string, f *font.Font) {
	sfxFileLock.Lock()
	sfxFiles[file] = f
	sfxFileLock.Unlock()
}

func ReloadSFXAssets() error {
	sfxFileLock.Lock()
	defer sfxFileLock.Unlock()
	for s, f := range sfxFiles {
		a, err := audio.Get(s)
		if err != nil {
			return err 
		}
		SFXs[s] = audio.New(f, a)
	}
	return nil
}

func PlayOneOfSFX(files ...string) {
	i := rand.Intn(len(files))
	PlaySFX(files[i])
}

func PlaySFX(s string) error {
	audOrigin, ok := SFXs[s]
	if !ok {
		return fmt.Errorf("Tried to play unloaded Audio %q", s)
	}
	aud, err := audOrigin.Copy()
	if err != nil {
		return err
	}
	a := aud.(*audio.Audio)

	v := volume * sfxVolume
	if v > 1 || v < 0 {
		v = 0
	}
	scalar := convertVolumeScale(v)
	a.SetVolume(scalar)
	a.Play()
	return nil
}

func SetSFXVolume(newVolume float64) {
	sfxVolume = newVolume
	updateSFXVolume(volume, sfxVolume)
	event.Trigger(EventSFXVolumeChanged, newVolume)
}

func updateSFXVolume(volume, sfxVolume float64) {
	newVolume := volume * sfxVolume
	if newVolume > 1 || newVolume < 0 {
		newVolume = 0
	}
	scalar := convertVolumeScale(newVolume)

	for _, s := range SFXs {
		s.SetVolume(scalar)
	}
}
