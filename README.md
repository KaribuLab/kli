# Karibu CLI

Karibu CLI (kli) es una herramienta de línea de comandos diseñada para ayudar a los desarrolladores con tareas comunes como la gestión de versiones semánticas y la creación de proyectos a partir de plantillas.

## Instalación

```bash
# Clonar el repositorio
git clone https://github.com/KaribuLab/kli.git

# Entrar al directorio
cd kli

# Compilar
go build -o bin/kli ./*.go

# Opcional: Mover el binario a un directorio en el PATH
sudo mv bin/kli /usr/local/bin/
```

## Comandos disponibles

### Comando `semver`

El comando `semver` analiza los mensajes de commit y genera una versión semántica basada en el [Versionado Semántico](https://semver.org/lang/es/).

#### Uso básico

```bash
kli semver
```

Este comando analizará el historial de commits y generará una versión semántica siguiendo las reglas:
- Commits con `fix:` incrementan el número de parche (0.0.X)
- Commits con `feat:` incrementan el número menor (0.X.0)
- Commits con `!` o `BREAKING CHANGE` incrementan el número mayor (X.0.0)

#### Opciones

```bash
kli semver [flags]

Flags:
  -d, --dryrun           Ejecutar sin crear tags reales
  -p, --pattern string   Patrón a utilizar para el tag (default "v{major}.{minor}.{patch}")
  -r, --remove           Eliminar tags
  -t, --tags             Crear todos los tags si no están presentes
  -v, --verbose          Salida detallada
```

#### Ejemplos

**Generar la versión actual basada en los commits:**
```bash
kli semver
# Salida: v0.1.2
```

**Crear tags automáticamente:**
```bash
kli semver -t
# Crea los tags necesarios en Git y los muestra
```

**Modo detallado:**
```bash
kli semver -v
# Muestra información detallada sobre los commits analizados
```

**Usar un patrón personalizado:**
```bash
kli semver -p "version-{major}.{minor}.{patch}"
# Salida: version-0.1.2
```

### Comando `project`

El comando `project` permite crear nuevos proyectos basados en plantillas alojadas en repositorios Git.

El sistema utiliza el motor de plantillas (templates) de Go para procesar los archivos de la plantilla y generar el código personalizado. Esto te permite crear estructuras de proyecto dinámicas basadas en los inputs proporcionados por el usuario.

#### Uso básico

```bash
kli project <URL-del-repositorio>
```

Este comando clonará el repositorio especificado y aplicará transformaciones según la configuración en el archivo `.kliproject.json`.

#### Opciones

```bash
kli project [URL-del-repositorio] [flags]

Flags:
  -b, --branch string   Rama a clonar (default "main")
  -w, --workdir string  Directorio de trabajo (default ".")
```

#### Ejemplo de uso

**Crear un nuevo proyecto basado en una plantilla:**
```bash
kli project https://github.com/usuario/plantilla-react
# Solicitará inputs según se definan en .kliproject.json
# Por ejemplo:
# > Nombre del proyecto:
# > Descripción del proyecto:
# > Autor:
```

**Especificar una rama diferente:**
```bash
kli project https://github.com/usuario/plantilla-node -b desarrollo
```

**Especificar un directorio de trabajo:**
```bash
kli project https://github.com/usuario/plantilla-go -w mi-proyecto-nuevo
```

## Estructura de archivos de configuración

### Configuración de plantillas (`.kliproject.json`)

Para crear tus propias plantillas de proyecto, debes incluir un archivo `.kliproject.json` en la raíz del repositorio:

```json
{
  "prompts": [
    {
      "name": "projectName",
      "description": "Nombre del proyecto:",
      "type": "string"
    }
  ],
  "templates": [
    {
      "rootDir": "templates",
      "delete": true,
      "files": [
        {
          "source": "templates/README.md.tmpl",
          "destination": "README.md"
        }
      ]
    }
  ],
  "posthooks": [
    {
      "name": "Instalar dependencias",
      "command": "npm install"
    }
  ]
}
```

### Sintaxis de plantillas

El comando `project` utiliza el [paquete text/template de Go](https://pkg.go.dev/text/template) para procesar las plantillas. Puedes utilizar esta sintaxis en tus archivos de plantilla:

```
# Proyecto: {{.Inputs.projectName}}

## Descripción
{{.Inputs.description}}

## Autor
{{.Inputs.author}}
```

Además de las variables directas, el sistema ofrece las siguientes funciones de transformación:

- `{{toLowerCase "MiTexto"}}` - Convierte el texto a minúsculas
- `{{toUpperCase "MiTexto"}}` - Convierte el texto a mayúsculas
- `{{toPascalCase "mi-texto"}}` - Convierte a formato PascalCase (MiTexto)
- `{{toCamelCase "mi-texto"}}` - Convierte a formato camelCase (miTexto)

## Contribuir

Las contribuciones son bienvenidas. Por favor, envía un pull request o abre un issue para discutir los cambios propuestos.