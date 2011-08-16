// persistence.go
//
package main

import (
	"fmt"
	"path"
	"os"
	o "orchestra"
)

// changing this will result in fire.  you have been warned.
const bucketDepth = 2

var spoolDirectory = ""

func SetSpoolDirectory(spooldir string) {
	if spoolDirectory == "" {
		spoolDirectory = spooldir
	} else {
		if spooldir != spoolDirectory {
			o.Warn("Spool Directory Not Changed.")
		}
	}
}

func GetSpoolDirectory() string {
	if spoolDirectory == "" {
		o.Assert("GetSpoolDirectory() called before set")
	}
	return spoolDirectory
}

func FilenameForJobId(jobid uint64) (fullpath string) {
	fnbits := make([]string, bucketDepth+1)
	for i := 0; i < bucketDepth; i++ {
		fnbits[i] = fmt.Sprintf("%01X", (jobid >> uint(i*4)) & 0xF)
	}
	fnbits[bucketDepth] = fmt.Sprintf("%0*X", 16-bucketDepth, uint64(jobid >> uint(bucketDepth*4)))

	return path.Join(fnbits...)
}

func makeSpoolDirInner(prefix string, depth int) {
	for i := 0; i < 16; i++ {
		dirname := path.Join(prefix, fmt.Sprintf("%01X", i))
		if (depth == 1) {
			err := os.MkdirAll(dirname, 0700)
			o.MightFail(err, "Couldn't make directory building spool tree")
		} else {
			makeSpoolDirInner(dirname, depth-1)
		}
	}
}

func MakeSpoolDir() {
	makeSpoolDirInner(path.Join(spoolDirectory, "active"), bucketDepth)
	makeSpoolDirInner(path.Join(spoolDirectory, "finished"), bucketDepth)
	os.MkdirAll(path.Join(spoolDirectory, "corrupt"), 0700)
}

// blat
func UnlinkNodesForJobId(jobid uint64) {
	suffix := FilenameForJobId(jobid)

	os.Remove(path.Join(spoolDirectory, "active", suffix))
	os.Remove(path.Join(spoolDirectory, "finished", suffix))
}