package main

import (
   "fmt"
   "log"
   "os"
   "os/user"
   "bufio"
   "strings"
   "path/filepath"
   "io/ioutil"

   "github.com/pelletier/go-toml"
   "github.com/urfave/cli"
)

type Account struct {
   Title string
   Username string
   Password string
   Description string
}

func (a *Account) Print() {
   fmt.Println("Title:       " + a.Title)
   fmt.Println("Username:    " + a.Username)
   fmt.Println("Password:    " + a.Password)
   fmt.Println("Description: " + a.Description)
}

func open_config_file() {
   usr, err := user.Current()
   if err != nil {
      log.Fatal( err )
   }

   rc_path := filepath.Join(usr.HomeDir, ".keyp")
   rc_bytes, err := ioutil.ReadFile(rc_path)
   if err != nil {
      log.Fatal( err )
   }

   // decrypt the file here

   rc_string := string(rc_bytes)
   fmt.Println(rc_string)
   config, err2 := toml.Load(rc_string)
   if err2 != nil {
      log.Fatal( err )
   }

   // retrieve data directly
   m := config.ToMap()
   for k, v := range m {
      fmt.Printf("key[%s] value[%s]\n", k, v)
   }
}

func executeParameters(c *cli.Context) {
   if c.Bool("add") {
      a := Account {}
      a.Title = c.Args()[0]
      prompt_info(&a)
      fmt.Println() // for debugging
      a.Print()
   } else if c.Bool("change") {
      a := Account {} // should open from toml file instead
      prompt_info(&a)
      fmt.Println() // for debugging
   }
   // remove is next
}

func main() {
   open_config_file()
   app := setup_gui()

   app.Action = func(c *cli.Context) error {
      if c.NArg() == 1 {
         executeParameters(c)
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

func prompt(s string) string {
   reader := bufio.NewReader(os.Stdin)
   fmt.Print(s)
   text, _ := reader.ReadString('\n')
   return strings.TrimSpace(text)
}

func prompt_info(a *Account) {
   if len(a.Username)    == 0 { a.Username    = prompt("Username: ") }
   if len(a.Password)    == 0 { a.Password    = prompt("Password: ") }
   if len(a.Description) == 0 { a.Description = prompt("Description: ") }
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
