package main

import (
	"os"
	"strings"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/encoder"
	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/proto"
)

type AvgComputeListener struct {
	cumulativePower     uint64
	cumulativeHeartrate uint64
	cumulativeCadence   uint64
	records             uint64
}

func (l *AvgComputeListener) OnMesg(msg proto.Message) {
	if msg.Num == typedef.MesgNumRecord {
		record := mesgdef.NewRecord(&msg)
		l.cumulativePower += uint64(record.Power)
		l.cumulativeHeartrate += uint64(record.HeartRate)
		l.cumulativeCadence += uint64(record.Cadence)
		l.records += 1
	}
}

func (l *AvgComputeListener) GetAveragePower() uint64 {
	return l.cumulativePower / l.records
}
func (l *AvgComputeListener) GetAverageHeartRate() uint64 {
	return l.cumulativeHeartrate / l.records
}
func (l *AvgComputeListener) GetAverageCadence() uint64 {
	return l.cumulativeCadence / l.records
}

func CreatePatchFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	l := &AvgComputeListener{}
	dec := decoder.New(f, decoder.WithMesgListener(l))

	fit, err := dec.Decode()
	if err != nil {
		panic(err)
	}

	activity := filedef.NewActivity(fit.Messages...)

	activity.Sessions[0].SetAvgCadence(uint8(l.GetAverageCadence()))
	activity.Sessions[0].SetAvgHeartRate(uint8(l.GetAverageHeartRate()))
	activity.Sessions[0].SetAvgPower(uint16(l.GetAveragePower()))

	// Convert back to FIT protocol messages
	patchedFit := activity.ToFIT(nil)

	f_patched, err := os.OpenFile("patched_"+fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer f_patched.Close()

	enc := encoder.New(f_patched, encoder.WithProtocolVersion(fit.FileHeader.ProtocolVersion))
	if err := enc.Encode(&patchedFit); err != nil {
		panic(err)
	}
}

func main() {
	current_dir, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, entry := range current_dir {
		if entry.IsDir() {
			continue
		}
		if strings.Contains(entry.Name(), ".fit") && !strings.Contains(entry.Name(), "patched") {
			CreatePatchFile(entry.Name())
		}
	}
}
