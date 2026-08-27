package store

import "github.com/lacsar712/slabheat/internal/model"

type SlabskidSnapshotView struct {
	UnitID   string
	Slabskid     model.SlabskidReading
	Alarms   []model.AlarmEvent
	Revision uint64
}

func CloneSlabskidSnapshot(s model.PlantSnapshot) SlabskidSnapshotView {
	out := SlabskidSnapshotView{
		UnitID:   s.UnitID,
		Slabskid:     s.Slabskid,
		Revision: s.Revision,
	}
	out.Alarms = make([]model.AlarmEvent, len(s.Alarms))
	copy(out.Alarms, s.Alarms)
	return out
}
