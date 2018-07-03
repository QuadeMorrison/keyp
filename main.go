package main

import (
   "fmt"
   "log"
   "os"

   "github.com/urfave/cli"
)

func main() {
   app := setup_gui()

   app.Action = func(c *cli.Context) error {
      if len(c.Args()) == 1 {
         fmt.Println("hello, it worked")
      } else {
         cli.ShowAppHelpAndExit(c, 0)
      }

      return nil
   }

   err := app.Run(os.Args)
   if err != nil {
      log.Fatal(err)
   }
}

func setup_gui() *cli.App {
   app := cli.NewApp()

   app.Name = "keyp"
   app.Usage = "A cli password manager using encryption."
   app.HideHelp = true // no dumb help command.
   app.Version = "0.5"

   app.Authors = []cli.Author{
      {Name: "Quade Morrison", Email: "quademorrison@gmail.com"},
      {Name: "Alan Morgan",    Email: "alanxoc3@gmail.com"},
   }

   app.Flags = []cli.Flag {
      cli.BoolFlag   { Name: "a, add",            Usage: "Add an account." },
      cli.BoolFlag   { Name: "c, change",         Usage: "Change an account." },
      cli.BoolFlag   { Name: "r, remove",         Usage: "Remove an account." },
      cli.BoolFlag   { Name: "l, list",           Usage: "Lists all account titles." },
      cli.StringFlag { Name: "u, username",       Usage: "Account username." },
      cli.StringFlag { Name: "p, password",       Usage: "Account password." },
      cli.StringFlag { Name: "d, description",    Usage: "Account description." },
      cli.BoolFlag   { Name: "U, no-username",    Usage: "Don't prompt username." },
      cli.BoolFlag   { Name: "P, no-password",    Usage: "Don't prompt password." },
      cli.BoolFlag   { Name: "D, no-description", Usage: "Don't prompt description." },
   }

   app.HelpName = "keyp"

   // Gotta override the template in order to have a better usage.
   // EXAMPLE: Override a template
   cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}title{{end}}
{{if .Version}}
VERSION:
   {{.Version}}{{end}}
{{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}
   {{end}}{{end}}{{if .VisibleFlags}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

   app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
      cli.ShowAppHelpAndExit(c, 0)
      return nil
   }

   return app
}
