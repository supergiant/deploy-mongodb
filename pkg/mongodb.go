package pkg

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/supergiant/supergiant/client"
)

func configureReplicaSet(instances []*client.InstanceResource) error {
	if len(instances) < 3 {
		return nil
	}

	var rsConfMems []string
	for _, instance := range instances {
		rsConfMems = append(rsConfMems, `{_id: `+*instance.ID+`, host: "`+instance.Addresses.Internal[0].Address+`"}`)
	}
	rsConf := `{_id: "rs0", members: [` + strings.Join(rsConfMems, ", ") + `]}` // TODO TODO TODO rs name....... param

	primaryAddr := instances[0].Addresses.Internal[0].Address

	initCmd := exec.Command("mongo", primaryAddr)
	initCmd.Stdin = strings.NewReader("rs.initiate()\n")
	initOut, err := initCmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(initOut))

	confCmd := exec.Command("mongo", primaryAddr)
	confCmd.Stdin = strings.NewReader("rs.reconfig(" + rsConf + ")\n")
	confOut, err := confCmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(confOut))

	return nil
}
