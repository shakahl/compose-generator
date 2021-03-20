package model

// TemplateConfig represents the JSON structure of predefined template configuration file
type TemplateConfig struct {
	Label     string
	Dir       string
	Files     []File
	Questions []Question
	Volumes   []Volume
	Secrets   []Secret
}

// File represents an important file and holds the path and the type of this file
type File struct {
	Path string
	Type    string
}

// Question represents the JSON structure of a question of a predefined template
type Question struct {
	Text           string
	Type           int // 1 = Yes/No; 2 = Text
	DefaultValue   string
	Validator      string
	EnvVar         string
	Advanced       bool
	WithDockerfile bool
}

// Volume represents the JSON structure of a volume of a predefined template
type Volume struct {
	Text           string
	DefaultValue   string
	EnvVar         string
	Advanced       bool
	WithDockerfile bool
}

// Secret represents the JSON structure of a secret of a predefined template
type Secret struct {
	Name   string
	Var    string
	Length int
}
