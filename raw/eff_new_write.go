package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

func (eff *EffNew) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)

	writeStr64 := func(s string) error {
		if len(s) >= effNewNameLen {
			return fmt.Errorf(
				"effect name %q is too long: maximum is %d bytes",
				s,
				effNewNameLen-1,
			)
		}

		var b [effNewNameLen]byte
		copy(b[:], []byte(s))
		enc.Bytes(b[:])
		return nil
	}

	for i, rec := range eff.Records {
		if rec == nil {
			return fmt.Errorf("record %d is nil", i)
		}

		err := writeStr64(rec.Name)
		if err != nil {
			return fmt.Errorf("record %d: %w", i, err)
		}

		for j := 0; j < 4; j++ {
			emitter := &rec.FirstEmitters[j]
			enc.Int32(emitter.UnknownA)
			enc.Int32(emitter.EmitterID)
			enc.Int32(emitter.UnknownB)
			enc.Int32(emitter.UnknownC)
		}

		for j := 0; j < 19; j++ {
			enc.Int32(rec.Unknown[j])
		}

		for j := 0; j < 4; j++ {
			emitter := &rec.SecondEmitters[j]
			enc.Int32(emitter.UnknownA)
			enc.Int32(emitter.EmitterID)
			enc.Int32(emitter.UnknownB)
			enc.Int32(emitter.UnknownC)
		}
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
