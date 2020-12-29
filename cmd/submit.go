package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/xalanq/cf-tool/client"
	"github.com/xalanq/cf-tool/config"
)

// Submit command
func Submit() (err error) {
	// confirm
	color.Cyan("Are you sure to submit? (y/n)")
	var ans string
	fmt.Scan(&ans)
	if ans != "y" {
		color.Magenta("Submit canceled.")
		return
	}

	cln := client.Instance
	cfg := config.Instance
	info := Args.Info
	filename, index, err := getOneCode(Args.File, cfg.Template)
	if err != nil {
		return
	}

	bundledname := withoutext(filename) + ".bdl"

	cmd := exec.Command("bash", "-c", "/bin/python3.8 ~/my-bundle "+filename+" "+bundledname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if ret := cmd.Run(); ret != nil {
		return
	}

	bytes, err := ioutil.ReadFile(bundledname)
	if err != nil {
		return
	}
	source := string(bytes)

	lang := cfg.Template[index].Lang
	if err = cln.Submit(info, lang, source); err != nil {
		if err = loginAgain(cln, err); err == nil {
			err = cln.Submit(info, lang, source)
		}
	}
	return
}

func withoutext(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
