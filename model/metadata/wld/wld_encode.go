package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// Encode writes wld.Fragments to a .wld writer. Use quail.WldMarshal to convert a common.Wld to wld.Fragments
func Encode(wld *common.Wld, w io.Writer) error {
	if wld.Fragments == nil {
		wld.Fragments = make(map[int]common.FragmentReader)
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Bytes([]byte{0x02, 0x3D, 0x50, 0x54})
	enc.Uint32(uint32(wld.Header.Version))
	enc.Uint32(uint32(len(wld.Fragments)))
	enc.Uint32(0) //unk1
	enc.Uint32(0) //unk2

	fragBuf := bytes.NewBuffer(nil)
	nameData, err := writeFragments(wld, fragBuf)
	if err != nil {
		return fmt.Errorf("write fragments: %w", err)
	}
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(0) //unk3
	enc.Bytes(nameData)
	enc.Bytes(fragBuf.Bytes())

	if enc.Error() != nil {
		return fmt.Errorf("encode: %w", enc.Error())
	}
	return nil
}

// writeFragments converts fragment structs to bytes
func writeFragments(wld *common.Wld, w io.Writer) ([]byte, error) {
	nameBuf := bytes.NewBuffer(nil)
	for i, frag := range wld.Fragments {
		err := frag.Encode(w)
		if err != nil {
			return nil, fmt.Errorf("fragment id %d 0x%x (%s): %w", i, frag.FragCode(), common.FragName(frag.FragCode()), err)
		}
		// Name builder?
		/*
			pos, ok := names[name]
			if !ok {
				tmpNames = append(tmpNames, model.Header.Name)
				names[name] = int32(nameBuf.Len())
				nameBuf.Write([]byte(name))
				nameBuf.Write([]byte{0})
				pos = names[name]
			}
		*/
	}
	return nameBuf.Bytes(), nil
}