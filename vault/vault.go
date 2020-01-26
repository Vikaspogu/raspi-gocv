package vault

import (
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var c *api.Logical

func init() {
	//Read service account token
	content, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Fatal(err)
	}
	config := &api.Config{
		Address: os.Getenv("VAULT_ADDR"),
	}
	client, err := api.NewClient(config)
	if err != nil {
		log.Warn(err)
		return
	}
	//Lookup VAULT_LOGIN_PATH Environment variable
	mountPath, set := os.LookupEnv("VAULT_LOGIN_PATH")
	if !set {
		mountPath = "auth/kubernetes/login"
	}
	//Attempt Vault login
	s, err := client.Logical().Write(mountPath, map[string]interface{}{
		"jwt":  string(content[:]),
		"role": "vault-demo-role",
	})
	if err != nil {
		log.Warn(err)
		return
	}
	client.SetToken(s.Auth.ClientToken)
	c = client.Logical()
}

func ReadSecret(path string) (string, string){
	secret, err := c.Read(path)
	if err != nil {
		log.Warn(err)
		return "", ""
	}
	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Warn("%T %#v\n", secret.Data["data"], secret.Data["data"])
		return "", ""
	}
	for key, value := range m {
		return key, value.(string)
	}
	return "", ""
}
