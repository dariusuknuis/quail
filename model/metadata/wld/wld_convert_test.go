package wld

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
)

func TestConvert(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		path      string
		file      string
		fragIndex int
		want      common.FragmentReader
		wantErr   bool
	}{
		{"btp_chr.s3d", "btp_chr.wld", 0, &Mesh{NameRef: 1414544642}, false},
		{"bac_chr.s3d", "bac_chr.wld", 0, &Mesh{NameRef: 1414544642}, false},
		{"avi_chr.s3d", "avi_chr.wld", 0, &Mesh{NameRef: 1414544642}, false},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.path))
			if err != nil {
				t.Fatalf("failed to open s3d %s: %s", tt.file, err.Error())
			}
			defer pfs.Close()
			data, err := pfs.File(tt.file)
			if err != nil {
				t.Fatalf("failed to open wld %s: %s", tt.file, err.Error())
			}
			world := &common.Wld{}
			err = Decode(world, bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to decode wld %s: %s", tt.file, err.Error())
			}

			models, err := Convert(world)
			if err != nil {
				t.Fatalf("failed to convert wld %s: %s", tt.file, err.Error())
			}

			if len(models) == 0 {
				t.Fatalf("failed to convert wld %s: no models", tt.file)
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("decodeMesh() = %v, want %v", got, tt.want)
			//}
		})
	}
}
