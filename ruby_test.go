package cli

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestFindRubyGems(t *testing.T) {
	withTestDir(t, func(directory string) {
		configDirectory := path.Join(directory, "config")
		err := os.MkdirAll(configDirectory, 0700)
		if err != nil {
			t.Fatal(err)
		}

		projectDirectory := path.Join(directory, "project")
		srcDirectory := path.Join(projectDirectory, "src")
		err = os.MkdirAll(srcDirectory, 0700)
		if err != nil {
			t.Fatal(err)
		}

		err = ioutil.WriteFile(
			path.Join(projectDirectory, "Gemfile"),
			[]byte(`
source "https://rubygems.org"
gem "licensezero_rubygem", git: 'https://github.com/licensezero/licensezero_rubygem'
			`),
			0700,
		)
		if err != nil {
			t.Fatal(err)
		}

		err = ioutil.WriteFile(
			path.Join(projectDirectory, "Gemfile.lock"),
			[]byte(`
GIT
  remote: https://github.com/licensezero/licensezero_rubygem
  revision: 68cc71897a4c094e967245456c5e12adc8e2b641
  specs:
    licensezero_rubygem (0.1.0)

GEM
  remote: https://rubygems.org/
  specs:

PLATFORMS
  ruby

DEPENDENCIES
  licensezero_rubygem!

BUNDLED WITH
   2.0.2
			`),
			0700,
		)
		if err != nil {
			t.Fatal(err)
		}

		err = ioutil.WriteFile(
			path.Join(srcDirectory, "main.rb"),
			[]byte(`
require 'licensezero_rubygem'

puts LicenseZeroRubyGem::MESSAGE
			`),
			0700,
		)
		if err != nil {
			t.Fatal(err)
		}

		install := exec.Command("bundle", "install")
		install.Dir = projectDirectory
		err = install.Run()
		if err != nil {
			t.Fatal(err)
		}

		findings, err := findRubyGems(projectDirectory)
		if err != nil {
			t.Fatal("read error")
		}

		if len(findings) != 1 {
			t.Fatal("did not find one offer")
		}
		finding := findings[0]
		if finding.Type != "rubygem" {
			t.Error("did not set RubyGem type")
		}
		if finding.API != "https://api.licensezero.com" {
			t.Error("did not set API")
		}
	})
}
