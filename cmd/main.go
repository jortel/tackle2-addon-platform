package main

import (
	"os"
	"path"

	"github.com/konveyor/tackle2-addon/ssh"
	hub "github.com/konveyor/tackle2-hub/addon"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/nas"
)

var (
	addon       = hub.Addon
	TemplateDir = ""
	AssetDir    = ""
	Dir         = ""
)

func init() {
	Dir, _ = os.Getwd()
	TemplateDir = path.Join(Dir, "/templates")
	AssetDir = path.Join(Dir, "assets")
}

// Data Addon data passed in the secret.
type Data struct {
	// Action
	// - import
	// - fetch
	// - generate
	Action string `json:"action"`
	// Filter applications.
	Filter api.Map `json:"filter"`
}

// main
func main() {
	addon.Run(func() (err error) {
		addon.Activity("TemplateDir: %s", TemplateDir)
		addon.Activity("AssetDir: %s", AssetDir)
		//
		// Get the addon data associated with the task.
		d := &Data{}
		err = addon.DataWith(d)
		if err != nil {
			return
		}
		//
		// Create directories.
		for _, dir := range []string{TemplateDir, AssetDir} {
			err = nas.MkDir(dir, 0755)
			if err != nil {
				return
			}
		}
		//
		// SSH
		agent := ssh.Agent{}
		err = agent.Start()
		if err != nil {
			return
		}
		//
		// action
		action, err := NewAction(d)
		if err != nil {
			return
		}
		//
		// Run action
		err = action.Run(d)
		if err != nil {
			return
		}

		addon.Activity("Done.")

		return
	})
}
