package model

import (
	"testing"
	"time"
)

func TestCase(t *testing.T) {
	snap := PlantSnapshot{
		UnitID: "EXP-1",
		Slabskid: SlabskidReading{LevelPercent: 55},
		Alarms: []AlarmEvent{
			{Code: "HI", Severity: "warning", Message: "high", RaisedAt: time.Now(), Active: true},
		},
	}
	clone := CloneSnapshot(snap)
	if len(clone.Alarms) != 1 {
		t.Fatal("expected one alarm in clone")
	}
	clone.Alarms[0].Code = "MUTATED"
	if snap.Alarms[0].Code == "MUTATED" {
		t.Fatal("mutating clone alarms should not affect source snapshot")
	}
}
