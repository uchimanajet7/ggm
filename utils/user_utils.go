package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type userConfig struct {
	LastDate      int64
	LastTotal     int64
	UserEmail     string
	SpeakCommands [][]string
	UsbCommands   [][]string
	Filters       []selectFilter
}

type selectFilter struct {
	From     string
	Subjects []string
}

// GetUserConfig is get user config data
func GetUserConfig() (*userConfig, error) {
	userConf, err := loadUserConfig()
	if err == nil {
		// return saved config
		last := getTimeFromEpoch(userConf.LastDate)
		def := getDefaultLastDate()

		if def.After(last) {
			userConf.LastDate = getEpochFromTime(def)
		}
		return userConf, err
	}

	// create new config
	prof, err := getGmailProfile()
	if err != nil {
		return nil, err
	}

	userConf = &userConfig{}
	userConf.LastDate = getEpochFromTime(getDefaultLastDate())
	userConf.LastTotal = prof.MessagesTotal
	userConf.UserEmail = prof.EmailAddress

	return userConf, nil
}

func (c *userConfig) IsTargetData(data *gmailData) bool {
	if c.Filters == nil {
		// Non filters
		return true
	}

	for _, v := range c.Filters {
		if strings.Index(data.From, v.From) >= 0 {
			if v.Subjects == nil {
				return true
			}
			for _, s := range v.Subjects {
				if strings.Index(data.Subject, s) >= 0 {
					return true
				}
			}
		}
	}

	return false
}

func (c *userConfig) UpdateUserConfig(last int64) error {
	prof, err := getGmailProfile()
	if err != nil {
		return err
	}

	c.LastDate = last
	c.LastTotal = prof.MessagesTotal
	c.UserEmail = prof.EmailAddress
	err = saveUserConfig(c)

	return err
}

func (c *userConfig) RunSpeakCommand(text string) error {
	if len(c.SpeakCommands) <= 0 {
		fmt.Print("\nSpeak command not specified.\n\n")
		return nil
	}

	cmds := make([][]string, len(c.SpeakCommands))
	for i, v := range c.SpeakCommands {
		cmd := make([]string, len(v))
		for j, m := range v {
			if m == "%s" {
				m = fmt.Sprintf(m, text)
			}
			cmd[j] = m
		}
		cmds[i] = cmd
	}

	_, err := runPipeline(cmds...)

	return err
}

func (c *userConfig) RunUsbCommand(enabled bool) error {
	if len(c.UsbCommands) <= 0 {
		fmt.Print("\nUsb commands not specified.\n\n")
		return nil
	}

	var power int
	if enabled {
		power = 1
	}
	cmds := make([][]string, len(c.UsbCommands))
	for i, v := range c.UsbCommands {
		cmd := make([]string, len(v))
		for j, m := range v {
			if m == "%d" {
				m = fmt.Sprintf(m, power)
			}
			cmd[j] = m
		}
		cmds[i] = cmd
	}

	_, err := runPipeline(cmds...)

	return err
}

func runPipeline(commands ...[]string) ([]byte, error) {
	cmdList := make([]*exec.Cmd, len(commands))

	for i, v := range commands {
		cmdList[i] = exec.Command(v[0], v[1:]...)
		if i > 0 {
			r, err := cmdList[i-1].StdoutPipe()
			if err != nil {
				return nil, err
			}
			cmdList[i].Stdin = r
		}
		cmdList[i].Stderr = os.Stderr
	}

	var b bytes.Buffer
	cmdList[len(cmdList)-1].Stdout = &b

	for _, v := range cmdList {
		if err := v.Start(); err != nil {
			return nil, err
		}
	}
	for _, v := range cmdList {
		if err := v.Wait(); err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

func getUserConfigFilePath() (string, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dataDir, "user_config.json"), err
}

func loadUserConfig() (*userConfig, error) {
	path, err := getUserConfigFilePath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &userConfig{}
	err = json.NewDecoder(f).Decode(conf)

	return conf, err
}

func saveUserConfig(config *userConfig) error {
	path, err := getUserConfigFilePath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// write file
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(config)
	if err != nil {
		return err
	}

	fmt.Printf("\nSaving user config file to: %s\n\n", path)

	return err
}
