package mds

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/eq/bxi.eqg"
	inFile := "bxi.mds"

	a, err := eqg.New(path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	ra, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer ra.Close()
	err = a.Load(ra)
	if err != nil {
		t.Fatalf("archive load: %s", err)
	}

	d, err := dump.New(inFile)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save(fmt.Sprintf("test/eq/%s.png", inFile))

	e, err := NewEQG(inFile, a)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	data, err := a.File(inFile)
	if err != nil {
		t.Fatalf("file: %s", err)
	}

	err = e.Load(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("load mds: %s", err)
	}
}

func TestLoadSaveLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/"
	inFile := "test/obj_gears.mod"
	outFile := "test/obj_gears_loadsaveload.mod"
	f, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	d.Save(fmt.Sprintf("%s.png", outFile))
	dump.Close()

	r, err := os.Open(outFile)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
}

func TestLoadSaveGLTF(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/"
	inFile := "test/obj_gears.mod"
	outFile := "test/obj_gears_loadsavegtlf.gltf"

	f, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(inFile)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create gltf: %s", err)
	}
	defer w.Close()

	doc, err := gltf.New()
	if err != nil {
		t.Fatalf("gltf.New: %s", err)
	}
	err = e.GLTFExport(doc)
	if err != nil {
		t.Fatalf("gltf: %s", err)
	}

	err = doc.Export(w)
	if err != nil {
		t.Fatalf("export: %s", err)
	}
	d.Save(fmt.Sprintf("%s.png", outFile))
	dump.Close()
}