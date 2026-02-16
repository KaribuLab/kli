package project

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	texttemplate "text/template"
	"unicode"

	"github.com/KaribuLab/kli/git"
	"github.com/spf13/cobra"
)

var templateFunctions map[string]any = texttemplate.FuncMap{
	"toLowerCase":  strings.ToLower,
	"toUpperCase":  strings.ToUpper,
	"toPascalCase": ToPascalCase,
	"toCamelCase":  ToCamelCase,
}

func ToPascalCase(s string) string {
	// Utilizamos un Builder para construir la cadena de manera eficiente.
	var builder strings.Builder

	// Función para determinar si un carácter es un separador.
	isSeparator := func(r rune) bool {
		return r == ' ' || r == '-' || r == '_' || r == '.' || r == '/' || r == '\\'
	}

	// Separa la cadena en palabras basadas en los separadores definidos.
	words := strings.FieldsFunc(s, isSeparator)

	for _, word := range words {
		if len(word) == 0 {
			continue
		}
		// Convierte el primer carácter a mayúscula y el resto a minúscula.
		// Esto maneja correctamente las letras Unicode.
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		for i := 1; i < len(runes); i++ {
			runes[i] = unicode.ToLower(runes[i])
		}
		builder.WriteString(string(runes))
	}

	return builder.String()
}

func ToCamelCase(s string) string {
	// Utilizamos un Builder para construir la cadena de manera eficiente.
	var builder strings.Builder

	// Función para determinar si un carácter es un separador.
	isSeparator := func(r rune) bool {
		return r == ' ' || r == '-' || r == '_' || r == '.' || r == '/' || r == '\\'
	}

	// Separa la cadena en palabras basadas en los separadores definidos.
	words := strings.FieldsFunc(s, isSeparator)

	for i, word := range words {
		if len(word) == 0 {
			continue
		}
		// Convertir la primera palabra a minúsculas y las siguientes con la primera letra en mayúscula.
		runes := []rune(word)
		if i == 0 {
			runes[0] = unicode.ToLower(runes[0])
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
		builder.WriteString(string(runes))
	}

	return builder.String()
}

func copyAll(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func buildDestinationPath(workdir string, destination string, propmt projectPrompt) (string, error) {
	tmpl, err := texttemplate.New("base").Parse(destination)
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	err = tmpl.Execute(&builder, propmt)
	if err != nil {
		return "", err
	}
	return path.Join(workdir, builder.String()), nil
}

func writeTemplates(template []template, workdir string, prompt projectPrompt) error {
	for _, t := range template {
		var templateFileList []string
		for _, f := range t.Files {
			sourceFilePath := path.Join(workdir, f.Source)
			templateFileList = append(templateFileList, sourceFilePath)
		}
		tmpl := texttemplate.New("base").Funcs(templateFunctions)
		tmpl, err := tmpl.ParseFiles(templateFileList...)
		if err != nil {
			return err
		}
		for _, f := range t.Files {
			destinationFilePath, err := buildDestinationPath(workdir, f.Destination, prompt)
			if err != nil {
				return err
			}
			sourceFilePath := filepath.Base(path.Join(workdir, f.Source))
			err = os.MkdirAll(path.Dir(destinationFilePath), os.ModePerm)
			if err != nil {
				return err
			}
			file, err := os.Create(destinationFilePath)
			if err != nil {
				return err
			}
			err = tmpl.ExecuteTemplate(file, sourceFilePath, prompt)
			if err != nil {
				return err
			}
			file.Close()
		}
		if t.Delete {
			os.RemoveAll(path.Join(workdir, t.RootDir))
		}
	}
	return nil
}

func NewProjectCommand(gitCmd git.Cmd) *cobra.Command {
	projectCommand := &cobra.Command{
		Use:   "project",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repository := args[0]
			branch, err := cmd.Flags().GetString("branch")
			if err != nil {
				return err
			}
			workdir, err := cmd.Flags().GetString("workdir")
			if err != nil {
				return err
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			tempPath, err := os.MkdirTemp("", "kli_*")
			if err != nil {
				return err
			}
			defer os.RemoveAll(tempPath)
			tempWorkingDirPath := path.Join(tempPath, workdir)
			err = gitCmd.Clone(repository, branch, tempWorkingDirPath)
			if err != nil {
				return err
			}
			os.RemoveAll(path.Join(tempWorkingDirPath, ".git"))
			payload, err := os.ReadFile(path.Join(tempWorkingDirPath, ".kliproject.json"))
			if err != nil {
				return err
			}
			var projectConfig projectConfig
			err = json.Unmarshal(payload, &projectConfig)
			if err != nil {
				return err
			}
			inputs := make(map[string]any)
			reader := bufio.NewReader(os.Stdin)
			for _, prompt := range projectConfig.Prompts {
				fmt.Println(prompt.Description)
				input, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				inputs[prompt.Name] = strings.TrimSpace(input)
			}
			projectPrompt := projectPrompt{
				Inputs: inputs,
			}
			err = writeTemplates(projectConfig.Templates, tempWorkingDirPath, projectPrompt)
			if err != nil {
				return err
			}
			err = copyAll(tempWorkingDirPath, path.Join(cwd, workdir))
			if err != nil {
				return err
			}
			err = os.Remove(path.Join(cwd, workdir, ".kliproject.json"))
			if err != nil {
				return err
			}
			for _, hook := range projectConfig.Posthooks {
				fmt.Printf("Running posthook '%s': %s\n", hook.Name, hook.Command)
				err = runHook(path.Join(cwd, workdir), hook.Command)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	projectCommand.Flags().StringP("branch", "b", "main", "Branch to clone")
	projectCommand.Flags().StringP("workdir", "w", ".", "Working directory")
	return projectCommand
}
