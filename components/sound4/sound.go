package sound4

import (
	"sync"

	"github.com/oakmound/oak/v4/event"
)

var (
	volume      = 0.0
	musicVolume = 0.0
	sfxVolume   = 0.0
)

var initOnce sync.Once

// Init loads assets and initializes volume levels. It will do nothing
// after it has been called once.
func Init(masterVolume, mVolume, sVolume float64) {
	initOnce.Do(func() {
		// ReloadMusicAssets()
		// ReloadSFXAssets()

		volume = masterVolume
		musicVolume = mVolume
		sfxVolume = sVolume

		updateSFXVolume(volume, sfxVolume)
		updateMusicVolume(volume, musicVolume)
	})
}

// convert a volume into the args for the windows api.
// Windows api is from 0 to -10000 but we see that -5000 and down is inaudible.
func convertVolumeScale(volumeScale float64) int32 {
	if volumeScale <= .1 {
		//  map 0 -> .1 to -10000 -> -5000
		volumeScale *= 5
		volumeScale--
		volumeScale *= 10000
	} else {
		// map .1 -> 1.0 to -5000 -> 0
		volumeScale--
		volumeScale *= 5555
	}
	return int32(volumeScale)
}

// SetMasterVolume value and update all volume values based of the given value.
func SetMasterVolume(eh event.Handler, masterVolume float64) {
	volume = masterVolume
	// updateSFXVolume(volume, sfxVolume)
	// updateMusicVolume(volume, musicVolume)
	// event.Trigger(EventMasterVolumeChanged, volume)
	event.TriggerOn(eh, EventVolumeChange, VolumeChangePayload{Kind: KindMaster, NewVolume: volume})
}
