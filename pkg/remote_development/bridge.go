package remote_development

import "bunnyshell.com/dev/pkg/ssh"

func GetSSHConfigFilePath() (string, error) {
	return ssh.GetConfigFilePath()
}
