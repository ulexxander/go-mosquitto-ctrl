package mosquittoctrl

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Logger interface {
	Cmd(cmd string, stdin, stdout, stderr *bytes.Buffer)
}

type LoggerStd struct {
	Logger *log.Logger
}

func (ls *LoggerStd) Cmd(cmd string, stdin, stdout, stderr *bytes.Buffer) {
	if ls.Logger == nil {
		ls.Logger = log.New(os.Stderr, "mosquitto_ctrl: ", log.LstdFlags)
	}
	ls.Logger.Printf("executed command: %s", cmd)
	ls.Logger.Printf("stdout: %s", stdout)
	ls.Logger.Printf("stderr: %s", stderr)
}

type Dynsec struct {
	client        *ssh.Client
	adminUsername string
	adminPassword string
	Logger        Logger
}

func NewDynsec(client *ssh.Client, adminUsername, adminPassword string) *Dynsec {
	return &Dynsec{
		client:        client,
		adminUsername: adminUsername,
		adminPassword: adminPassword,
	}
}

const DefaultClientConfigFile = "/mosquitto/config/dynamic-security.json"

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

func (d *Dynsec) run(cmd string, stdinLines ...string) error {
	session, err := d.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdinLines = append(stdinLines, d.adminPassword)
	stdin := bytes.NewBufferString(strings.Join(stdinLines, "\n"))
	session.Stdin = stdin
	var stdout bytes.Buffer
	session.Stdout = &stdout
	var stderr bytes.Buffer
	session.Stderr = &stderr

	err = session.Run(cmd)
	if d.Logger != nil {
		d.Logger.Cmd(cmd, stdin, &stdout, &stderr)
	}
	if err != nil {
		return err
	}
	return seekOutputErrors(stderr.String())
}

// ConnectionError represents MQTT errors printed to stderr
// See https://github.com/eclipse/mosquitto/blob/master/lib/strings_mosq.c for possible error messages
// Example:
// Connection error: Not authorized
type ConnectionError struct {
	Reason string
}

func (ce *ConnectionError) Error() string {
	return "Connection error: " + ce.Reason
}

func seekOutputErrors(out string) error {
	prefix := "Connection error: "
	crIdx := strings.Index(out, prefix)
	if crIdx == -1 {
		return nil
	}
	out = out[crIdx+len(prefix):]
	newlineIdx := strings.Index(out, "\n")
	if newlineIdx != -1 {
		out = out[:newlineIdx]
	}
	return &ConnectionError{Reason: out}
}
