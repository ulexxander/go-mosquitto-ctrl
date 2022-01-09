# Go Mosquitto Ctrl

Manage Mosquitto users backed by [Dynamic Security Plugin](https://mosquitto.org/documentation/dynamic-security) over SSH.

Official [eclipse-mosquitto](https://github.com/eclipse/mosquitto) Docker image does not have SSH daemon set up.
I recommend using my [eclipse-mosquitto-ssh](https://github.com/ulexxander/eclipse-mosquitto-ssh) Docker image which has SSH and Dynsec initialization ready-to-go.

## Usage

```go
// Connect to Mosquitto container via SSH
sshClient, err := ssh.Dial("tcp", "localhost:1882", &ssh.ClientConfig{
  User: "admin",
  Auth: []ssh.AuthMethod{
    ssh.Password("123"),
  },
  // Example only, use safer option in production
  HostKeyCallback: ssh.InsecureIgnoreHostKey(),
})
if err != nil {
  return fmt.Errorf("dialing ssh: %w", err)
}
defer sshClient.Close()

// Initialize client
ds := mosquittoctrl.NewDynsec(sshClient, "admin", "123")

roleName := "time"
clientUsername := "time_publisher"
clientPassword := "123"

// Manage users and permissions
err = ds.CreateRole(roleName)
if err != nil {
  return fmt.Errorf("creating role: %w", err)
}
err = ds.AddRoleACL(roleName, "publishClientSend", "time_current", "allow", 1)
if err != nil {
  return fmt.Errorf("adding role ACL (publish): %w", err)
}
err = ds.CreateClient(clientUsername, clientPassword)
if err != nil {
  return fmt.Errorf("creating client: %w", err)
}
err = ds.AddClientRole(clientUsername, roleName)
if err != nil {
  return fmt.Errorf("adding client role: %w", err)
}
```
