package mosquittoctrl

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Dynsec struct {
	client        *ssh.Client
	adminUsername string
	adminPassword string
}

func NewDynsec(client *ssh.Client, adminUsername, adminPassword string) *Dynsec {
	return &Dynsec{
		client:        client,
		adminUsername: adminUsername,
		adminPassword: adminPassword,
	}
}

const DefaultConfigFile = "/mosquitto/config/dynamic-security.json"

func (d *Dynsec) Init(configFile string) error {
	return d.run(
		fmt.Sprintf(
			"mosquitto_ctrl dynsec init %s %s",
			configFile, d.adminUsername,
		),
		d.adminPassword,
		d.adminPassword,
	)
}

func (d *Dynsec) CreateRole(name string) error {
	return d.run(
		fmt.Sprintf(
			"mosquitto_ctrl -u %s dynsec createRole %s",
			d.adminUsername, name,
		),
	)
}

func (d *Dynsec) DeleteRole(name string) error {
	return d.run(
		fmt.Sprintf(
			"mosquitto_ctrl -u %s dynsec deleteRole %s",
			d.adminUsername, name,
		),
	)
}

func (d *Dynsec) AddRoleACL(role, aclType, topicFilter, allowOrDeny string, priority int) error {
	return d.run(
		fmt.Sprintf(
			"mosquitto_ctrl -u %s dynsec addRoleACL %s %s %s %s %d",
			d.adminUsername, role, aclType, topicFilter, allowOrDeny, priority,
		),
	)
}

func (d *Dynsec) CreateClient(name, password string) error {
	return d.run(
		fmt.Sprintf(
			"mosquitto_ctrl -u %s dynsec createClient %s",
			d.adminUsername, name,
		),
		password,
		password,
	)
}

func (d *Dynsec) DeleteClient(name string) error {
	return d.run(
		fmt.Sprintf(
			"mosquitto_ctrl -u %s dynsec deleteClient %s",
			d.adminUsername, name,
		),
	)
}

func (d *Dynsec) AddClientRole(client string, role string) error {
	return d.run(
		fmt.Sprintf(
			"mosquitto_ctrl -u %s dynsec addClientRole %s %s",
			d.adminUsername, client, role,
		),
	)
}

func (d *Dynsec) run(cmd string, stdin ...string) error {
	session, err := d.client.NewSession()
	if err != nil {
		return err
	}
	stdin = append(stdin, d.adminPassword)
	session.Stdin = strings.NewReader(strings.Join(stdin, "\n"))
	defer session.Close()
	return session.Run(cmd)
}
