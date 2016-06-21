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

	primaryAddr := instances[0].Addresses.Internal[0].Address

	var rsConfMems []string
	for _, instance := range instances {
		rsConfMems = append(rsConfMems, `{_id: `+*instance.ID+`, host: "`+instance.Addresses.Internal[0].Address+`"}`)
	}
	rsConf := `{_id: "rs0", members: [` + strings.Join(rsConfMems, ", ") + `]}` // TODO TODO TODO rs name....... param

	fmt.Println(strings.Join([]string{"mongo", primaryAddr, "--eval", `'rs.initiate()'`}, " "))
	out, err := exec.Command("mongo", primaryAddr, "--eval", `'rs.initiate()'`).Output()
	fmt.Println(string(out))
	if err != nil {
		return err
	}

	fmt.Println(strings.Join([]string{"mongo", primaryAddr, "--eval", `'rs.reconfig(` + rsConf + `)'`}, " "))
	out, err = exec.Command("mongo", primaryAddr, "--eval", `'rs.reconfig(`+rsConf+`)'`).Output()
	fmt.Println(string(out))
	if err != nil {
		return err
	}

	return nil
}
