package mosquittoctrl_test

import (
	"testing"

	mosquittoctrl "github.com/ulexxander/go-mosquitto-ctrl"
	"golang.org/x/crypto/ssh"
)

const sshServer = "localhost:1882"
const sshUsername = "admin"
const sshPassword = "admin"

var sshConfig = &ssh.ClientConfig{
	User: sshUsername,
	Auth: []ssh.AuthMethod{
		ssh.Password(sshPassword),
	},
	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
}

func TestDynsec(t *testing.T) {
	sshClient, err := ssh.Dial("tcp", sshServer, sshConfig)
	if err != nil {
		t.Fatalf("error dialing ssh: %s", err)
	}
	defer sshClient.Close()

	adminUsername := "admin"
	adminPassword := "admin"
	clientUsername := "time_publisher"
	clientPassword := "123"
	roleName := "time"

	mcd := mosquittoctrl.NewDynsec(sshClient, adminUsername, adminPassword)

	err = mcd.DeleteClient(clientUsername)
	if err != nil {
		t.Fatalf("error cleaning up client: %s", err)
	}
	err = mcd.DeleteRole(roleName)
	if err != nil {
		t.Fatalf("error cleaning up role: %s", err)
	}

	err = mcd.CreateRole(roleName)
	if err != nil {
		t.Fatalf("error creating role: %s", err)
	}
	err = mcd.AddRoleACL(roleName, "publishClientSend", "time_current", "allow", 1)
	if err != nil {
		t.Fatalf("error allowing publishClientSend: %s", err)
	}
	err = mcd.AddRoleACL(roleName, "subscribeLiteral", "time_current", "allow", 1)
	if err != nil {
		t.Fatalf("error allowing subscribeLiteral: %s", err)
	}
	err = mcd.CreateClient(clientUsername, clientPassword)
	if err != nil {
		t.Fatalf("error creating client: %s", err)
	}
	err = mcd.AddClientRole(clientUsername, roleName)
	if err != nil {
		t.Fatalf("error adding client role: %s", err)
	}
}
