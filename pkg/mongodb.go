package pkg

import (
	"bytes"
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

	cmd := exec.Command("mongo", primaryAddr, "--eval", `'rs.initiate(); rs.reconfig(`+rsConf+`)'`)

	fmt.Println(cmd.Path, strings.Join(cmd.Args, " "))

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	fmt.Println("OUTPUT")
	fmt.Println(out.String())

	fmt.Println("ERR")
	fmt.Println(stderr.String())

	return err
}
