package player

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"time"
)

var (
	_IdToOfflineTimer = cmap.New[*time.Timer]()
)
