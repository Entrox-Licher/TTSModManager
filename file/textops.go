package file

import (
	"fmt"
	"io"
	"os"
	"path"
)

// TextOps allows for arbitrary reads and writes of text files
type TextOps struct {
	basepaths        []string
	writeBasepath    string
	readFileToBytes  func(string) ([]byte, error)
	writeBytesToFile func(string, []byte) error
}

// TextReader serves to describe all ways to read luascripts
type TextReader interface {
	EncodeFromFile(string) (string, error)
}

// TextWriter serves to describe all ways to write luascripts
type TextWriter interface {
	EncodeToFile(script, file string) error
}

// NewTextOps initializes our object on a directory
func NewTextOps(base string) *TextOps {
	return NewTextOpsMulti([]string{base}, base)
}

// NewTextOpsMulti allows for luascript to be read from multiple directories
func NewTextOpsMulti(readDirs []string, writeDir string) *TextOps {
	return &TextOps{
		basepaths:     readDirs,
		writeBasepath: writeDir,
		readFileToBytes: func(s string) ([]byte, error) {
			sFile, err := os.Open(s)
			if err != nil {
				return nil, fmt.Errorf("os.Open(%s): %v", s, err)
			}
			defer sFile.Close()

			b, err := io.ReadAll(sFile)
			if err != nil {
				return nil, fmt.Errorf("ReadAll(%s): %v", s, err)
			}
			if l := len(b); l > 0 {
				if b[l-1] == '\n' {
					b = b[0 : l-1]
				}
			}
			return b, nil
		},
		writeBytesToFile: func(p string, b []byte) error {
			b = append(b, '\n')
			err := os.MkdirAll(path.Dir(p), 0750)
			if err != nil && !os.IsExist(err) {
				return fmt.Errorf("MkdirAll(%s): %v", path.Dir(p), err)
			}
			return os.WriteFile(p, b, 0644)
		},
	}
}

// EncodeFromFile pulls a file from configs and encodes it as a string.
func (l *TextOps) EncodeFromFile(filename string) (string, error) {
	for _, base := range l.basepaths {
		p := path.Join(base, filename)
		b, err := l.readFileToBytes(p)
		if err != nil {
			continue
		}
		s := string(b)
		return s, nil
	}
	return "", fmt.Errorf("%s not found among any known paths", filename)
}

// EncodeToFile takes a single string and decodes escape characters; writes it.
func (l *TextOps) EncodeToFile(script, file string) error {
	p := path.Join(l.writeBasepath, file)
	return l.writeBytesToFile(p, []byte(script))
}
