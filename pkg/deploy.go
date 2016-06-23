package pkg

import supergiant "github.com/supergiant/supergiant/client"

func Deploy(appName *string, componentName *string) error {

	// supergiant.Log.SetLevel("debug")

	sg := supergiant.New("http://supergiant-api.supergiant.svc.cluster.local/v0", "", "", true)

	app, err := sg.Apps().Get(appName)
	if err != nil {
		return err
	}

	component, err := app.Components().Get(componentName)
	if err != nil {
		return err
	}

	var currentRelease *supergiant.ReleaseResource
	if component.CurrentReleaseTimestamp != nil {
		currentRelease, err = component.CurrentRelease()
		if err != nil {
			return err
		}
	}

	targetRelease, err := component.TargetRelease()
	if err != nil {
		return err
	}

	targetList, err := targetRelease.Instances().List()
	if err != nil {
		return err
	}
	targetInstances := targetList.Items

	if currentRelease == nil { // first release
		for _, instance := range targetInstances {
			if err = instance.Start(); err != nil {
				return err
			}
		}
		for _, instance := range targetInstances {
			if err = instance.WaitForStarted(); err != nil {
				return err
			}
		}

		// Initiate replica set
		return configureReplicaSet(targetInstances)
	}

	currentList, err := currentRelease.Instances().List()
	if err != nil {
		return err
	}
	currentInstances := currentList.Items

	// remove instances
	if currentRelease.InstanceCount > targetRelease.InstanceCount {
		instancesRemoving := currentRelease.InstanceCount - targetRelease.InstanceCount

		for _, instance := range currentInstances[len(currentInstances)-instancesRemoving:] {
			if err := instance.Stop(); err != nil {
				return err
			}
		}

		// add new instances
	} else if currentRelease.InstanceCount < targetRelease.InstanceCount {
		instancesAdding := targetRelease.InstanceCount - currentRelease.InstanceCount
		newInstances := targetInstances[len(targetInstances)-instancesAdding:]
		for _, instance := range newInstances {
			if err := instance.Start(); err != nil {
				return err
			}
		}
		for _, instance := range newInstances {
			if err := instance.WaitForStarted(); err != nil {
				return err
			}
		}
	}

	// update instances

	if *currentRelease.InstanceGroup == *targetRelease.InstanceGroup {
		// no need to update restart instances
		return configureReplicaSet(targetInstances)
	}

	var instancesRestarting int
	if currentRelease.InstanceCount < targetRelease.InstanceCount {
		instancesRestarting = currentRelease.InstanceCount
	} else {
		instancesRestarting = targetRelease.InstanceCount
	}

	for i := 0; i < instancesRestarting; i++ {
		currentInstance := currentInstances[i]
		targetInstance := targetInstances[i]

		currentInstance.Stop()
		currentInstance.WaitForStopped()

		targetInstance.Start()
		targetInstance.WaitForStarted()
	}

	for i := 0; i < instancesRestarting; i++ {
		currentInstance := currentInstances[i]
		targetInstance := targetInstances[i]

		currentInstance.Stop()
		currentInstance.WaitForStopped()

		targetInstance.Start()
		targetInstance.WaitForStarted()
	}

	return configureReplicaSet(targetInstances)
}
