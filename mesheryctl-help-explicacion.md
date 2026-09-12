# `mesheryctl help`: flujo, variables y jerarquía

## Flujo principal

```text
mesheryctl/cmd/mesheryctl/main.go
        |
        v
root.Execute()
        |
        v
RootCmd.Execute()
        |
        +--> mesheryctl help [command]
        |       |
        |       +--> newHelpCommand()
        |       +--> Root().Find(args)
        |       +--> cmd.HelpFunc()
        |
        +--> mesheryctl --help
        |       |
        |       +--> ayuda automática de Cobra
        |
        +--> mesheryctl
                |
                +--> RootCmd.RunE()
                +--> cmd.Help()
```

En `mesheryctl` no existe una única función `help`; el comportamiento lo controla **Cobra** mediante un árbol de comandos.

## Archivo de entrada

<u>[mesheryctl/cmd/mesheryctl/main.go](mesheryctl/cmd/mesheryctl/main.go#L20-L28)</u>

```go
func main() {
    err := root.Execute()
    if err != nil {
        os.Exit(1)
    }
}
```

Su función es iniciar el controlador raíz y terminar con código `1` si ocurre un error.

## Controlador raíz

<u>[mesheryctl/internal/cli/root/root.go](mesheryctl/internal/cli/root/root.go#L56-L88)</u>

Aquí está definido `RootCmd`, el controlador principal:

```go
var RootCmd = &cobra.Command{
    Use:   "mesheryctl",
    Short: "Meshery Command Line tool",
    Long:  "...",
    Example: "...",

    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        mesheryctlflags.InitValidators(cmd)
        return nil
    },

    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            return cmd.Help()
        }

        ...
    },
}
```

### Variables principales

<u>[mesheryctl/internal/cli/root/root.go](mesheryctl/internal/cli/root/root.go#L46-L53)</u>

| Variable | Tipo | Uso |
|---|---|---|
| `cfgFile` | `string` | Ruta del archivo de configuración |
| `verbose` | `bool` | Activa logs detallados con `-v` o `--verbose` |
| `availableSubcommands` | `[]*cobra.Command` | Lista de comandos registrados |
| `RootCmd` | `*cobra.Command` | Controlador raíz de toda la CLI |

También existen variables locales importantes:

| Variable | Lugar | Uso |
|---|---|---|
| `cmd` | `RunE`, `newHelpCommand` | Comando encontrado o ejecutado |
| `args` | Cobra | Argumentos recibidos |
| `c` | `newHelpCommand` | Comando `help` que recibe la ejecución |
| `err` | Varias funciones | Errores de búsqueda, configuración o ejecución |
| `utils.CfgFile` | `initConfig` | Ruta de configuración usada globalmente |
| `viper.ConfigFileUsed()` | `initConfig` | Archivo de configuración finalmente leído |

## Implementación de `mesheryctl help`

<u>[mesheryctl/internal/cli/root/root.go](mesheryctl/internal/cli/root/root.go#L141-L155)</u>

```go
func newHelpCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "help [command]",
        Short: "Show help for any command",
        Long:  "Show help for any command.",
        Run: func(c *cobra.Command, args []string) {
            cmd, _, err := c.Root().Find(args)

            if cmd == nil || err != nil {
                c.Println(c.UsageString())
                return
            }

            cmd.InitDefaultHelpFlag()
            cmd.HelpFunc()(cmd, args)
        },
    }
}
```

El flujo es:

1. Cobra recibe `help`.
2. `newHelpCommand()` crea el comando especial.
3. `c.Root().Find(args)` busca, por ejemplo, `system start`.
4. `cmd.InitDefaultHelpFlag()` agrega `-h` y `--help`.
5. `cmd.HelpFunc()` genera la ayuda del comando encontrado.
6. Cobra imprime el resultado.

Ejemplos:

```bash
mesheryctl help
mesheryctl help system
mesheryctl help system start
```

## Diferencia con `--help`

### `mesheryctl --help`

No pasa por `newHelpCommand()`. Cobra agrega automáticamente el flag de ayuda a los comandos. La salida se genera con la configuración de:

```go
RootCmd.Use
RootCmd.Short
RootCmd.Long
RootCmd.Example
RootCmd.Commands()
RootCmd.Flags()
```

### `mesheryctl`

Cuando no se proporcionan argumentos, se ejecuta:

```go
RunE: func(cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
        return cmd.Help()
    }
}
```

## Registro de comandos

<u>[mesheryctl/internal/cli/root/root.go](mesheryctl/internal/cli/root/root.go#L101-L138)</u>

Los comandos se agregan mediante:

```go
availableSubcommands = []*cobra.Command{
    completionCmd,
    versionCmd,
    system.SystemCmd,
    design.DesignCmd,
    perf.PerfCmd,
    adapter.AdapterCmd,
    experimental.ExpCmd,
    filter.FilterCmd,
    registry.RegistryCmd,
    components.ComponentCmd,
    model.ModelCmd,
    environments.EnvironmentCmd,
    connections.ConnectionsCmd,
    organizations.OrgCmd,
    relationships.RelationshipCmd,
    workspaces.WorkSpaceCmd,
}

RootCmd.AddCommand(availableSubcommands...)
RootCmd.SetHelpCommand(newHelpCommand())
```

La jerarquía principal es:

```text
mesheryctl
├── completion
├── version
├── system
├── design
├── perf
├── adapter
├── exp
├── filter
├── registry
├── component
├── model
├── environment
├── connection
├── organization
├── relationship
└── workspace
```

## Archivos de los comandos principales

<u>[mesheryctl/internal/cli/root/system/system.go](mesheryctl/internal/cli/root/system/system.go#L57-L102)</u>

Registra:

```text
system
├── reset
├── logs
├── start
├── stop
├── restart
├── status
├── update
├── config
├── context
├── channel
├── provider
├── check
├── login
├── logout
├── token
├── dashboard
└── delete
```

Los demás controladores principales son:

- <u>[mesheryctl/internal/cli/root/adapter/adapter.go](mesheryctl/internal/cli/root/adapter/adapter.go#L42-L88)</u>
- <u>[mesheryctl/internal/cli/root/design/design.go](mesheryctl/internal/cli/root/design/design.go#L38-L80)</u>
- <u>[mesheryctl/internal/cli/root/perf/perf.go](mesheryctl/internal/cli/root/perf/perf.go#L34-L77)</u>
- <u>[mesheryctl/internal/cli/root/filter/filter.go](mesheryctl/internal/cli/root/filter/filter.go#L34-L76)</u>
- <u>[mesheryctl/internal/cli/root/registry/registry.go](mesheryctl/internal/cli/root/registry/registry.go#L28-L64)</u>
- <u>[mesheryctl/internal/cli/root/components/component.go](mesheryctl/internal/cli/root/components/component.go#L34-L98)</u>
- <u>[mesheryctl/internal/cli/root/model/model.go](mesheryctl/internal/cli/root/model/model.go#L39-L128)</u>
- <u>[mesheryctl/internal/cli/root/environments/environment.go](mesheryctl/internal/cli/root/environments/environment.go#L29-L78)</u>
- <u>[mesheryctl/internal/cli/root/connections/connection.go](mesheryctl/internal/cli/root/connections/connection.go#L18-L80)</u>
- <u>[mesheryctl/internal/cli/root/organizations/organization.go](mesheryctl/internal/cli/root/organizations/organization.go#L14-L62)</u>
- <u>[mesheryctl/internal/cli/root/relationships/relationship.go](mesheryctl/internal/cli/root/relationships/relationship.go#L31-L103)</u>
- <u>[mesheryctl/internal/cli/root/workspaces/workspace.go](mesheryctl/internal/cli/root/workspaces/workspace.go#L26-L66)</u>

Cada archivo normalmente contiene un `*cobra.Command` y una llamada a:

```go
ParentCmd.AddCommand(availableSubcommands...)
```

Los comandos hoja están en archivos como:

```text
internal/cli/root/system/start.go
internal/cli/root/system/stop.go
internal/cli/root/system/status.go
internal/cli/root/model/list.go
internal/cli/root/model/view.go
internal/cli/root/model/search.go
internal/cli/root/design/apply.go
internal/cli/root/design/delete.go
internal/cli/root/filter/list.go
internal/cli/root/filter/view.go
```

Cada comando hoja aporta información que luego aparece en `help`:

```go
Use
Short
Long
Example
Args
Flags
Run
RunE
```

## Configuración y funciones auxiliares usadas

Durante el arranque se ejecutan callbacks registrados con Cobra:

<u>[mesheryctl/internal/cli/root/root.go](mesheryctl/internal/cli/root/root.go#L101-L116)</u>

```go
cobra.OnInitialize(setupLogger)
cobra.OnInitialize(initConfig)
```

### Configuración

<u>[mesheryctl/internal/cli/root/root.go](mesheryctl/internal/cli/root/root.go#L162-L218)</u>

`initConfig()` usa:

- `cfgFile`
- `utils.CfgFile`
- `utils.DefaultConfigPath`
- `utils.MesheryFolder`
- `viper.SetConfigFile`
- `viper.AutomaticEnv`
- `viper.ReadInConfig`
- `utils.SetKubeConfig`
- `config.AddTokenToConfig`
- `config.AddContextToConfig`

### Logs

<u>[mesheryctl/internal/cli/root/root.go](mesheryctl/internal/cli/root/root.go#L220-L235)</u>

`setupLogger()` usa:

- `verbose`
- `logrus.InfoLevel`
- `logrus.DebugLevel`
- `mesheryctllogger.GetMeshkitLogger`
- `utils.Log`

## Modelo arquitectónico

El modelo es un **árbol de comandos basado en Cobra**, utilizando el patrón compuesto:

```text
RootCmd
└── Command
    ├── Command
    │   ├── Command
    │   └── Command
    └── Command
```

| Rol | Implementación |
|---|---|
| Controlador raíz | `RootCmd` |
| Controladores secundarios | `SystemCmd`, `ModelCmd`, `DesignCmd`, etc. |
| Modelo de ayuda | Campos `Use`, `Short`, `Long`, `Example`, flags y subcomandos |
| Vista | Templates internos de Cobra y salida de `HelpFunc()` |
| Registro | `AddCommand()` |
| Ejecución | `RootCmd.Execute()` |
| Documentación | `GenMarkdownCustom()` |

No existe un modelo de base de datos para `help`. El “modelo” de ayuda es la estructura `cobra.Command`.

## Documentación generada

<u>[mesheryctl/doc/doc.go](mesheryctl/doc/doc.go#L158-L190)</u>

`GenMarkdownCustom()` transforma el árbol de Cobra en documentación Markdown.

Usa:

```go
cmd.InitDefaultHelpCmd()
cmd.InitDefaultHelpFlag()
cmd.CommandPath()
cmd.Short
cmd.Long
cmd.Example
cmd.UseLine()
cmd.NonInheritedFlags()
cmd.InheritedFlags()
```

La función `printOptions()` imprime los flags:

<u>[mesheryctl/doc/doc.go](mesheryctl/doc/doc.go#L143-L156)</u>

La documentación se genera desde:

<u>[mesheryctl/Makefile](mesheryctl/Makefile#L120-L124)</u>

```bash
cd doc
go run doc.go
```

El índice generado está en:

<u>[docs/content/en/reference/references/mesheryctl/_index.md](docs/content/en/reference/references/mesheryctl/_index.md#L1-L35)</u>

## Tests relacionados

Archivos que prueban este flujo:

- <u>[mesheryctl/internal/cli/root/root_test.go](mesheryctl/internal/cli/root/root_test.go#L24-L88)</u>
- <u>[mesheryctl/doc/doc_test.go](mesheryctl/doc/doc_test.go#L25-L210)</u>
- <u>[mesheryctl/internal/cli/root/filter/filter_test.go](mesheryctl/internal/cli/root/filter/filter_test.go#L14-L39)</u>
- <u>[mesheryctl/internal/cli/root/filter/testdata/filter.help.output.golden](mesheryctl/internal/cli/root/filter/testdata/filter.help.output.golden)</u>
- <u>[mesheryctl/tests/e2e/001-system/00-system.bats](mesheryctl/tests/e2e/001-system/00-system.bats#L7-L12)</u>
- <u>[mesheryctl/tests/e2e/005-component/01-component-list.bats](mesheryctl/tests/e2e/005-component/01-component-list.bats#L11-L15)</u>
- <u>[mesheryctl/tests/e2e/007-connection/04-connection-delete.bats](mesheryctl/tests/e2e/007-connection/04-connection-delete.bats#L50-L53)</u>
- <u>[mesheryctl/tests/e2e/008-workspace/00-workspace.bats](mesheryctl/tests/e2e/008-workspace/00-workspace.bats#L9-L13)</u>

## Archivos mínimos para crear un controlador nuevo

Para crear, por ejemplo, `mesheryctl sample`, se necesitarían:

```text
mesheryctl/cmd/mesheryctl/main.go
mesheryctl/internal/cli/root/root.go
mesheryctl/internal/cli/root/sample/sample.go
mesheryctl/internal/cli/root/sample/list.go
mesheryctl/internal/cli/root/sample/view.go
mesheryctl/internal/cli/root/sample/sample_test.go
mesheryctl/doc/doc.go
mesheryctl/doc/doc_test.go
```

El archivo principal sería:

```go
var SampleCmd = &cobra.Command{
    Use:   "sample",
    Short: "Manage samples",
    Long:  "Manage Meshery samples.",
    RunE: func(cmd *cobra.Command, args []string) error {
        return cmd.Help()
    },
}
```

Después se registra en `root.go`:

```go
availableSubcommands = []*cobra.Command{
    ...
    sample.SampleCmd,
}
```

Y sus subcomandos se registran dentro de `sample.go`:

```go
SampleCmd.AddCommand(listCmd, viewCmd)
```

La dependencia que proporciona el motor de comandos y ayuda es:

<u>[go.mod](go.mod#L1-L20)</u>

```text
github.com/spf13/cobra
```

## Resumen

Para cambiar el comportamiento de `mesheryctl help`, el archivo central es <u>[mesheryctl/internal/cli/root/root.go](mesheryctl/internal/cli/root/root.go#L138-L155)</u>. Para cambiar el texto mostrado por un comando concreto, hay que modificar el archivo que define ese `cobra.Command`, porque sus campos `Use`, `Short`, `Long`, `Example` y flags alimentan automáticamente la ayuda.
