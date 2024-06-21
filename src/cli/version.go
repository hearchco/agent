package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type versionFlag string

func (v versionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v versionFlag) IsBool() bool                         { return true }
func (v versionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

type version struct {
	ver       string
	timestamp string
	commit    string
}

func (v version) String() string {
	if v.ver == "" {
		return "dev"
	} else {
		return fmt.Sprintf("%v (%v@%v)", v.ver, v.commit, v.timestamp)
	}
}
