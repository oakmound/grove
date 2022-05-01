package sound4

import "github.com/oakmound/oak/v4/event"

// EventVolumeChange for any given volume type which should be one of the 3 supported.
var EventVolumeChange = event.RegisterEvent[VolumeChangePayload]()

// VolumeChangePayload encodes that type of volume to manipulate and the volume to set to.
type VolumeChangePayload struct {
	Kind      VolumeKind
	NewVolume float64
}
