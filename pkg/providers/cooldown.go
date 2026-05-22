package providers

import "time"

// cooldownCutoff returns the latest publish-time a release may have to be
// considered installable when a cooldown of `days` is configured. Releases
// published after this instant must be skipped.
func cooldownCutoff(days int) time.Time {
	return time.Now().UTC().AddDate(0, 0, -days)
}
