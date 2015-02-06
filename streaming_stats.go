package anaconda

type StreamStats struct {
	Listens     uint
	Disconnects uint
	Backoffs    uint
	Requests    uint
}

func (bss *StreamStats) Listened()     { bss.Listens++ }
func (bss *StreamStats) Disconnected() { bss.Disconnects++ }
func (bss *StreamStats) Backedoff()    { bss.Backoffs++ }
func (bss *StreamStats) Requested()    { bss.Requests++ }

func (bss *StreamStats) Minus(previous *StreamStats) StreamStats {
	return StreamStats{
		Listens:     bss.Listens - previous.Listens,
		Disconnects: bss.Disconnects - previous.Disconnects,
		Backoffs:    bss.Backoffs - previous.Backoffs,
		Requests:    bss.Requests - previous.Requests,
	}
}
