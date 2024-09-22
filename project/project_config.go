package project

type projectConfig struct {
	Prompts   []input    `json:"prompts"`
	Posthooks []postHook `json:"posthooks"`
	Templates []template `json:"templates"`
}

type template struct {
	RootDir string `json:"rootDir"`
	Delete  bool   `json:"delete"`
	Files   []file `json:"files"`
}

type file struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type input struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type postHook struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

type projectPrompt struct {
	Inputs map[string]interface{}
}
