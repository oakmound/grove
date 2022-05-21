package sound

import (
	"math/rand"
	"sync"

	"github.com/oakmound/oak/v4/event"
)

var (
	// SFXs = map[string]*audio.Audio{}

	sfxFileLock sync.RWMutex
	// sfxFiles    = map[string]*font.Font{}
)

// // RegisterSFX
// func RegisterSFX(file string, f *font.Font) {
// 	sfxFileLock.Lock()
// 	sfxFiles[file] = f
// 	sfxFileLock.Unlock()
// }

// // ReloadSFXAssets
// func ReloadSFXAssets() error {
// 	sfxFileLock.Lock()
// 	defer sfxFileLock.Unlock()
// 	for s, f := range sfxFiles {
// 		a, err := audio.Get(s)
// 		if err != nil {
// 			return err
// 		}
// 		SFXs[s] = audio.New(f, a)
// 	}
// 	return nil
// }

// PlayOneOfSFX playes a random entry from the provided files
func PlayOneOfSFX(files ...string) {
	i := rand.Intn(len(files))
	PlaySFX(files[i])
}

func PlaySFX(s string) error {
	// audOrigin, ok := SFXs[s]
	// if !ok {
	// 	return fmt.Errorf("Tried to play unloaded Audio %q", s)
	// }
	// aud, err := audOrigin.Copy()
	// if err != nil {
	// 	return err
	// }
	// a := aud.(*audio.Audio)

	// v := volume * sfxVolume
	// if v > 1 || v < 0 {
	// 	v = 0
	// }
	// scalar := convertVolumeScale(v)
	// a.SetVolume(scalar)
	// a.Play()
	return nil
}

// SetSFXVolume for grove's sound component
func SetSFXVolume(eh event.Handler, newVolume float64) {
	sfxVolume = newVolume
	updateSFXVolume(volume, sfxVolume)
	event.TriggerOn(eh, EventVolumeChange, VolumeChangePayload{Kind: KindSFX, NewVolume: sfxVolume})
}

func updateSFXVolume(volume, sfxVolume float64) {
	// newVolume := volume * sfxVolume
	// if newVolume > 1 || newVolume < 0 {
	// 	newVolume = 0
	// }
	// scalar := convertVolumeScale(newVolume)

	// for _, s := range SFXs {
	// 	s.SetVolume(scalar)
	// }
}
