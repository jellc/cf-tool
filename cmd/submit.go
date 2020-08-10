package cmd

import (
	"io/ioutil"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/xalanq/cf-tool/client"
	"github.com/xalanq/cf-tool/config"
)

// Submit command
func Submit() (err error) {
	cln := client.Instance
	cfg := config.Instance
	info := Args.Info
	filename, index, err := getOneCode(Args.File, cfg.Template)
	if err != nil {
		return
	}

	bundledname := withoutext(filename) + "_bundled"

	cmd := exec.Command("my-bundle",filename,bundledname)
	if ret := cmd.Run(); ret != nil { return }
	fmt.Println("INFO: successfully bundled.")

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
