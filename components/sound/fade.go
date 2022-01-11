package sound

import (
	"sync"
	"time"
)

const fadeLoopTime = 20 * time.Millisecond

// MusicFader exposes a Fader for usage on Music.
var MusicFader = Fader{}

// Fader is a safe way to fade in our out Music.
type Fader struct {
	sync.Mutex

	fadeTo    float64
	fadedFrom float64
	running   bool
	runTil    time.Time
}

// FadeOut the Music over the given time.
func (f *Fader) FadeOut(duration time.Duration) {
	f.Lock()
	if f.fadeTo == musicVolume && f.fadeTo == 0 {
		f.Unlock()
		return
	}
	f.fadedFrom = musicVolume
	if f.fadeTo > musicVolume {
		f.fadedFrom = f.fadeTo
	}
	f.fadeTo = 0
	f.Unlock()
	f.startFade(duration)
}

// FadeIn the Music over the given time.
func (f *Fader) FadeIn(duration time.Duration) {
	f.Lock()
	if f.fadeTo == f.fadedFrom {
		f.Unlock()
		return
	}
	f.fadeTo = f.fadedFrom
	f.Unlock()
	f.startFade(duration)
}

func (f *Fader) startFade(duration time.Duration) {
	f.Lock()
	f.runTil = time.Now().Add(duration)
	f.Unlock()
	if f.running {
		return
	}
	go func() {
		f.Lock()
		f.running = true
		f.Unlock()

		t := time.NewTicker(fadeLoopTime)
		defer t.Stop()

		for {
			f.Lock()
			if time.Now().After(f.runTil) {
				SetMusicVolume(f.fadeTo)
				break
			}
			f.Unlock()

			loopsLeft := time.Until(f.runTil) / fadeLoopTime
			fadeTo := musicVolume + (f.fadeTo-musicVolume)/float64(loopsLeft)
			SetMusicVolume(fadeTo)
			<-t.C
		}
		f.running = false
		f.Unlock()
	}()
}
