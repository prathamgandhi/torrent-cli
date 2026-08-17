package torrentfile

import "os"

func Read(fpath string) {
	os.ReadFile(fpath)
}
