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
	adminUsername := "admin"
	adminPassword := "admin"
	clientUsername := "time_publisher"
	clientPassword := "123"
	roleName := "time"
	ds := setupDynsec(t, adminUsername, adminPassword)
	ds.Logger = &mosquittoctrl.LoggerStd{}

	err := ds.Init(mosquittoctrl.DefaultClientConfigFile)
	if err != nil {
		t.Fatalf("error initializing: %s", err)
	}

	err = ds.DeleteClient(clientUsername)
	if err != nil {
		t.Fatalf("error cleaning up client: %s", err)
	}
	err = ds.DeleteRole(roleName)
	if err != nil {
		t.Fatalf("error cleaning up role: %s", err)
	}

	err = ds.CreateRole(roleName)
	if err != nil {
		t.Fatalf("error creating role: %s", err)
	}
	err = ds.AddRoleACL(roleName, "publishClientSend", "time_current", "allow", 1)
	if err != nil {
		t.Fatalf("error allowing publishClientSend: %s", err)
	}
	err = ds.AddRoleACL(roleName, "subscribeLiteral", "time_current", "allow", 1)
	if err != nil {
		t.Fatalf("error allowing subscribeLiteral: %s", err)
	}
	err = ds.CreateClient(clientUsername, clientPassword)
	if err != nil {
		t.Fatalf("error creating client: %s", err)
	}
	err = ds.AddClientRole(clientUsername, roleName)
	if err != nil {
		t.Fatalf("error adding client role: %s", err)
	}
}

func TestDynsecConnectionError(t *testing.T) {
	adminUsername := "guythatdoesnotexit"
	adminPassword := "123"
	ds := setupDynsec(t, adminUsername, adminPassword)

	err := ds.CreateClient("someclient", "clientpass")
	connErr, ok := err.(*mosquittoctrl.ConnectionError)
	if !ok {
		t.Fatalf("expected connection error, got: %T", err)
	}
	reason := "Not authorized"
	if connErr.Reason != reason {
		t.Fatalf("expected reason %s, got: %s", reason, connErr.Reason)
	}
}

func setupDynsec(t *testing.T, adminUsername, adminPassword string) *mosquittoctrl.Dynsec {
	sshClient, err := ssh.Dial("tcp", sshServer, sshConfig)
	if err != nil {
		t.Fatalf("error dialing ssh: %s", err)
	}
	t.Cleanup(func() {
		sshClient.Close()
	})
	return mosquittoctrl.NewDynsec(sshClient, adminUsername, adminPassword)
}
