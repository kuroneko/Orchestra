// scores.go
//
// Score handling
//
// In here, we have the probing code that learns about scores, reads
// their configuration files, and does the heavy lifting for launching
// them, doing the privilege drop, etc.

package main

import (
	"os"
	"io"
	"strings"
	o "orchestra"
	"path"
	"github.com/kuroneko/configureit"
)

type ScoreInfo struct {
	Name		string
	Executable	string
	InitialPwd	string
	InitialEnv	map[string]string

	Interface	string

	Config		*configureit.Config
}

type ScoreExecution struct {
	Score	*ScoreInfo
	Task	*TaskRequest
}
	

func NewScoreInfo() (si *ScoreInfo) {
	si = new (ScoreInfo)
	si.InitialEnv = make(map[string]string)

	config := NewScoreInfoConfig()
	si.updateFromConfig(config)

	return si
}

func NewScoreInfoConfig() (config *configureit.Config) {
	config = configureit.New()

	config.Add("interface", configureit.NewStringOption("env"))
	config.Add("dir", configureit.NewStringOption(""))
	config.Add("path", configureit.NewStringOption("/usr/bin:/bin"))
	config.Add("user", configureit.NewUserOption(""))

	return config
}

func (si *ScoreInfo) updateFromConfig(config *configureit.Config) {
	// propogate PATH overrides.
	opt := config.Get("dir")
	sopt, _ := opt.(*configureit.StringOption)
	si.InitialEnv["PATH"] = sopt.Value

	// set the interface type.
	opt = config.Get("interface")
	sopt, _ = opt.(*configureit.StringOption)
	si.Interface = sopt.Value

	// propogate initial Pwd
	opt = config.Get("dir")
	sopt, _ = opt.(*configureit.StringOption)
	si.InitialPwd = sopt.Value	
}

var (
	Scores		map[string]*ScoreInfo
)

func ScoreConfigure(si *ScoreInfo, r io.Reader) {
	config := NewScoreInfoConfig()
	err := config.Read(r, 1)
	o.MightFail(err, "Error Parsing Score Configuration for %s", si.Name)
	si.updateFromConfig(config)
}

func LoadScores() {
	scoreDirectory := GetStringOpt("score directory")

	dir, err := os.Open(scoreDirectory)
	o.MightFail(err, "Couldn't open Score directory")
	defer dir.Close()

	Scores = make(map[string]*ScoreInfo)
	
	files, err := dir.Readdir(-1)
	for i := range files {
		name := files[i].Name()
		// skip ., .. and other dotfiles.
		if strings.HasPrefix(name, ".") {
			continue
		}
		// emacs backup files.  ignore these.
		if strings.HasSuffix(name, "~") || strings.HasPrefix(name, "#") {
			continue
		}
		// .conf is reserved for score configurations.
		if strings.HasSuffix(name, ".conf") {
			continue
		}

		// check to see if it's a file or symlink
		ftype := files[i].Mode() & os.ModeType
		if ftype != 0 && ftype != os.ModeSymlink {
			continue
		}

		// check for the executionable bit
		if files[i].Mode() & 0111 == 0 {
			continue
		}

		fullpath := path.Join(scoreDirectory, name)
		conffile := fullpath+".conf"
		o.Warn("Considering %s as score", name)

		si := NewScoreInfo()
		si.Name = name
		si.Executable = fullpath
	
		conf, err := os.Open(conffile)
		if err == nil {
			o.Warn("Parsing configuration for %s", fullpath)
			ScoreConfigure(si, conf)
			conf.Close()
		} else {
			o.Warn("Couldn't open config file for %s, assuming defaults: %s", name, err)
		}
		Scores[name] = si
	}
}
