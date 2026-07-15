package raw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestAnlWrite(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}

	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	testDir := helper.DirTest()

	tests := []struct {
		eqg     string
		file    string
		wantErr bool
	}{
		// Add more as desired
		{eqg: "omensequip.eqg", file: "it10739_animlist.anl"},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {

			archive, err := pfs.NewFile(filepath.Join(eqPath, tt.eqg))
			if err != nil {
				t.Fatalf("pfs.NewFile() error = %v", err)
			}

			data, err := archive.File(tt.file)
			if err != nil {
				t.Fatalf("archive.File() error = %v", err)
			}

			anl := &Anl{}
			err = anl.Read(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}

			err = os.WriteFile(filepath.Join(testDir, tt.file+".src.anl"), data, 0644)
			if err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			buf := bytes.NewBuffer(nil)
			err = anl.Write(buf)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			anl2 := &Anl{}
			err = anl2.Read(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("Read(roundtrip) error = %v", err)
			}

			buf2 := bytes.NewBuffer(nil)
			err = anl2.Write(buf2)
			if err != nil {
				t.Fatalf("Write(roundtrip) error = %v", err)
			}

			err = os.WriteFile(filepath.Join(testDir, tt.file+".dst.anl"), buf2.Bytes(), 0644)
			if err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			// Validate parsed data
			if anl.Version != anl2.Version {
				t.Fatalf("Version mismatch: %d != %d", anl.Version, anl2.Version)
			}

			if len(anl.Animations) != len(anl2.Animations) {
				t.Fatalf("Animation count mismatch: %d != %d",
					len(anl.Animations), len(anl2.Animations))
			}

			for i := range anl.Animations {
				if anl.Animations[i] != anl2.Animations[i] {
					t.Fatalf("Animation %d mismatch: %q != %q",
						i,
						anl.Animations[i],
						anl2.Animations[i],
					)
				}
			}

			if anl.DefaultAnimation != anl2.DefaultAnimation {
				t.Fatalf("Default animation mismatch: %q != %q",
					anl.DefaultAnimation,
					anl2.DefaultAnimation,
				)
			}

			// Verify deterministic encoding
			err = helper.ByteCompareTest(buf.Bytes(), buf2.Bytes())
			if err != nil {
				t.Fatalf("Round-trip mismatch: %v", err)
			}

			// Since the writer preserves the original ordering and
			// duplicates the default animation exactly like the source,
			// this should also match the original file byte-for-byte.
			err = helper.ByteCompareTest(data, buf.Bytes())
			if err != nil {
				t.Fatalf("Original vs encoded mismatch: %v", err)
			}
		})
	}
}