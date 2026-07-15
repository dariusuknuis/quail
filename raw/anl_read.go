package raw

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type Anl struct {
	MetaFileName string
	Version      uint32

	Animations      []string
	DefaultAnimation string
}

func (anl *Anl) Identity() string {
	return "anl"
}

func (anl *Anl) String() string {
	out := ""
	out += fmt.Sprintf("metafilename %s\n", anl.MetaFileName)
	out += fmt.Sprintf("version %d\n", anl.Version)
	out += fmt.Sprintf("animations %d\n", len(anl.Animations))

	for _, anim := range anl.Animations {
		out += fmt.Sprintf("  %s\n", anim)
	}

	out += fmt.Sprintf("default %s", anl.DefaultAnimation)

	return out
}

func (anl *Anl) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQAL" {
		return fmt.Errorf("invalid header %q, wanted EQAL", header)
	}

	anl.Version = dec.Uint32()
	if anl.Version != 1 {
		return fmt.Errorf("unsupported EQAL version %d", anl.Version)
	}

	stringBlockLength := int(dec.Uint32())
	animationCount := int(dec.Uint32())

	// Read the string block
	stringBlock := dec.Bytes(stringBlockLength)

	// Read animation offsets
	offsets := make([]uint32, animationCount)
	for i := range offsets {
		offsets[i] = dec.Uint32()
	}

	// Read default animation offset
	defaultOffset := dec.Uint32()

	// Resolve animation names
	anl.Animations = make([]string, animationCount)
	for i, off := range offsets {
		if int(off) >= len(stringBlock) {
			return fmt.Errorf("animation offset %d beyond string block (%d bytes)", off, len(stringBlock))
		}

		end := int(off)
		for end < len(stringBlock) && stringBlock[end] != 0 {
			end++
		}

		anl.Animations[i] = string(stringBlock[off:end])
	}

	// Resolve default animation
	if int(defaultOffset) >= len(stringBlock) {
		return fmt.Errorf("default animation offset %d beyond string block (%d bytes)", defaultOffset, len(stringBlock))
	}

	end := int(defaultOffset)
	for end < len(stringBlock) && stringBlock[end] != 0 {
		end++
	}

	anl.DefaultAnimation = string(stringBlock[defaultOffset:end])

	// Verify we've consumed the file
	pos := dec.Pos()
	endPos, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}

	if pos < endPos {
		remaining := dec.Bytes(int(endPos - pos))
		if !bytes.Equal(remaining, []byte{0x00, 0x00, 0x00, 0x00}) {
			fmt.Printf("remaining bytes:\n%s\n", hex.Dump(remaining))
			return fmt.Errorf("%d bytes remaining (%d total)", endPos-pos, endPos)
		}
	}

	if pos > endPos {
		return fmt.Errorf("read past end of file")
	}

	if err := dec.Error(); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	return nil
}

func (anl *Anl) SetFileName(name string) {
	anl.MetaFileName = name
}

func (anl *Anl) FileName() string {
	return anl.MetaFileName
}