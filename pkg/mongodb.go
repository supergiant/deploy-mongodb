package pkg

import (
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
	cmd := exec.Command("mongo", primaryAddr)
	cmd.Stdin = strings.NewReader("rs.initiate()\nrs.reconfig(" + rsConf + ")\n")
	return cmd.Run()
}
