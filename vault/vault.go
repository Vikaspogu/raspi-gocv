package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"os"
)

var c *api.Logical

func init() {
	config := &api.Config{
		Address: os.Getenv("VAULT_ADDR"),
	}
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	client.SetToken(os.Getenv("TOKEN"))
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
