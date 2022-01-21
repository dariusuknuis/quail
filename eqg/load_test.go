package eqg

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoad(t *testing.T) {
	path := "test/eqzip-test.eqg"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save("test/out.png")
	e := &EQG{}
	err = e.Load(f)
	if err != nil {
		d.Save("test/out.png")
		t.Fatalf("load: %s", err)
	}

}

func TestLoadSaveLoad(t *testing.T) {
	path := "test/soldungb.eqg"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	e := &EQG{}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	w, err := os.Create("test/out.eqg")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	d.Save("test/out.png")
	dump.Close()

	path = "test/out.eqg"
	d, err = dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save("test/out2.png")
	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

}