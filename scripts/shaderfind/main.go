package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 3 {
		fmt.Println("usage: shaderfind <path> <shader>")
		fmt.Println("example: shaderfind C:\\EverQuest Opaque_MPLBasic.fx")
		os.Exit(1)
	}

	path := os.Args[1]
	targetShader := os.Args[2]

	fmt.Println("Path:  ", path)
	fmt.Println("Shader:", targetShader)

	start := time.Now()

	err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".eqg":
			// supported
		default:
			return nil
		}

		archive, err := pfs.NewFile(path)
		if err != nil {
			fmt.Printf("pfs open %s: %v\n", filepath.Base(path), err)
			return nil
		}

		archiveName := filepath.Base(path)
		fmt.Println("Reading", archiveName)

		for _, file := range archive.Files() {
			fileExt := filepath.Ext(file.Name())
			switch fileExt {
			case ".mod", ".mds", ".ter":
			default:
				continue
			}

			r := bytes.NewReader(file.Data())

			rawRead, err := raw.Read(fileExt, r)
			if err != nil {
				fmt.Printf("raw read %s: %v\n", filepath.Base(file.Name()), err)
				continue
			}

			found := false

			switch dat := rawRead.(type) {
			case *raw.Mod:
				for _, mat := range dat.Materials {
					if mat.ShaderName == targetShader {
						found = true
						break
					}
				}

			case *raw.Mds:
				for _, mat := range dat.Materials {
					if mat.ShaderName == targetShader {
						found = true
						break
					}
				}

			case *raw.Ter:
				for _, mat := range dat.Materials {
					if mat.ShaderName == targetShader {
						found = true
						break
					}
				}
			}

			if found {
				fmt.Printf("%s : %s\n", archiveName, file.Name())
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}

	fmt.Printf("Finished in %.2f seconds\n", time.Since(start).Seconds())
	return nil
}
