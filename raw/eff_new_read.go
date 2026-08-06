package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

const (
	EffNewRecordSize = 0x10C // 268
	effNewNameLen    = 0x40  // 64
)

// EffNew models spellsnew.eff.
type EffNew struct {
	MetaFileName string
	Records      []*EffNewRecord
}

type EffNewRecord struct {
	Name           string
	FirstEmitters  [4]EffNewEmitterRef
	Unknown        [19]int32
	SecondEmitters [4]EffNewEmitterRef
}

type EffNewEmitterRef struct {
	UnknownA  int32
	EmitterID int32
	UnknownB  int32
	UnknownC  int32
}

func (eff *EffNew) Identity() string {
	return "eff_new"
}

func (eff *EffNew) String() string {
	out := fmt.Sprintf("EffNew: %s,", eff.MetaFileName)
	out += fmt.Sprintf("records: %d,", len(eff.Records))
	return out
}

func (eff *EffNew) Read(r io.ReadSeeker) error {
	fileSize, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}
	if fileSize%EffNewRecordSize != 0 {
		return fmt.Errorf(
			"invalid spellsnew.eff size %d: not divisible by record size %d",
			fileSize,
			EffNewRecordSize,
		)
	}
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek start: %w", err)
	}

	recordCount := int(fileSize / EffNewRecordSize)
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	eff.Records = nil

	for i := 0; i < recordCount; i++ {
		effRec := &EffNewRecord{}

		nameData := dec.Bytes(effNewNameLen)
		if n := bytes.IndexByte(nameData, 0); n >= 0 {
			nameData = nameData[:n]
		}
		effRec.Name = string(nameData)

		for j := 0; j < 4; j++ {
			emitter := &effRec.FirstEmitters[j]
			emitter.UnknownA = dec.Int32()
			emitter.EmitterID = dec.Int32()
			emitter.UnknownB = dec.Int32()
			emitter.UnknownC = dec.Int32()
		}

		for j := 0; j < 19; j++ {
			effRec.Unknown[j] = dec.Int32()
		}

		for j := 0; j < 4; j++ {
			emitter := &effRec.SecondEmitters[j]
			emitter.UnknownA = dec.Int32()
			emitter.EmitterID = dec.Int32()
			emitter.UnknownB = dec.Int32()
			emitter.UnknownC = dec.Int32()
		}

		eff.Records = append(eff.Records, effRec)
	}

	err = dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if dec.Pos() != fileSize {
		return fmt.Errorf(
			"read ended at %d, expected %d",
			dec.Pos(),
			fileSize,
		)
	}

	return nil
}

// SetFileName sets the file name
func (eff *EffNew) SetFileName(name string) {
	eff.MetaFileName = name
}

// FileName returns the file name
func (eff *EffNew) FileName() string {
	return eff.MetaFileName
}
