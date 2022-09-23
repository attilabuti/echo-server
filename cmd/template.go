package cmd

var helpTemplate = `{{$v := offset .Name 6}}{{wrap .Name 3}}

Usage:
   {{.HelpName}} [options] {{if .Description}}

Description:
   {{wrap .Description 3}}{{end}}

Options:{{range .VisibleFlagCategories}}
   {{if .Name}}{{.Name}}
   {{end}}{{range .Flags}}{{.}}
   {{end}}{{end}}
`
