package pkg

import (
	"errors"
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
	// cmd := exec.Command("mongo", primaryAddr)
	// cmd.Stdin = strings.NewReader("rs.initiate(); rs.reconfig(" + rsConf + ")\n")
	// return cmd.Run()

	initCmd := exec.Command("mongo", primaryAddr)
	initCmd.Stdin = strings.NewReader("rs.initiate()\n")
	if err := initCmd.Run(); err != nil {
		return err
	}

	confCmd := exec.Command("mongo", primaryAddr)
	confCmd.Stdin = strings.NewReader("rs.reconfig(" + rsConf + ")\n")
	out, err := confCmd.CombinedOutput()
	if err != nil {
		return err
	}

	outstr := string(out)
	if strings.Contains(outstr, `{ "ok" : 1 }`) {
		return errors.New("rs.reconfig() did not work. full output: " + outstr)
	}

	return nil
}
