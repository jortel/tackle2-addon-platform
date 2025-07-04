package main

import (
	"os"
	"path"

	"github.com/konveyor/tackle2-addon/ssh"
	hub "github.com/konveyor/tackle2-hub/addon"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/nas"
	"k8s.io/utils/env"
)

var (
	addon     = hub.Addon
	SharedDir = ""
	CacheDir  = ""
	SourceDir = ""
	Dir       = ""
	OptDir    = ""
)

func init() {
	Dir, _ = os.Getwd()
	OptDir = path.Join(Dir, "opt")
	SharedDir = env.GetString(hub.EnvSharedDir, "/tmp/shared")
	CacheDir = env.GetString(hub.EnvCacheDir, "/tmp/cache")
	SourceDir = path.Join(SharedDir, "source")
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
		addon.Activity("OptDir:    %s", OptDir)
		addon.Activity("SharedDir: %s", SharedDir)
		addon.Activity("CacheDir:  %s", CacheDir)
		addon.Activity("SourceDir: %s", SourceDir)
		//
		// Get the addon data associated with the task.
		d := &Data{}
		err = addon.DataWith(d)
		if err != nil {
			return
		}
		//
		// Create directories.
		for _, dir := range []string{OptDir} {
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
