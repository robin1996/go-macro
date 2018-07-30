package macrofile

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Point struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type Step struct {
	Type     int           `yaml:"type"`
	Pos      Point         `yaml:"pos"`
	Colour   string        `yaml:"colour"`
	Duration time.Duration `yaml:"duration"`
}

const (
	LeftClick = iota
	RightClick
	Test
)

const macroFile = "C:\\Users\\robdo\\Desktop\\macro.yaml"

func WriteMacro(recording []Step) {
	content, err := yaml.Marshal(recording)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(macroFile, content, 0644)
	if err != nil {
		panic(err)
	}
}
