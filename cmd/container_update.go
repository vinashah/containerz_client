// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"context"
	"google.golang.org/grpc/metadata"

	"github.com/openconfig/containerz/client"
	"github.com/spf13/cobra"
)

var (
	updateAsync bool
)

var cntUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a container with the specified image and tag",
	RunE: func(command *cobra.Command, args []string) error {
		if image == "" {
			return fmt.Errorf("--image must be specified")
		}
		if instance == "" {
			return fmt.Errorf("--instance must be specified")
		}

		opts := []client.StartOption{}
		if len(ports) > 0 {
			opts = append(opts, client.WithPorts(ports))
		}
		if len(envs) > 0 {
			opts = append(opts, client.WithEnv(envs))
		}
		if len(volumes) > 0 {
			opts = append(opts, client.WithVolumes(volumes))
		}
		if len(devices) > 0 {
			opts = append(opts, client.WithDevices(devices))
		}
		if network != "" {
			opts = append(opts, client.WithNetwork(network))
		}
		if runAs != "" {
			opts = append(opts, client.WithRunAs(runAs))
		}
		if restartPolicy != "" {
			opts = append(opts, client.WithRestartPolicy(restartPolicy))
		}
		if len(addCaps) > 0 || len(delCaps) > 0 {
			opts = append(opts, client.WithCapabilities(addCaps, delCaps))
		}
		ctx, cancel := context.WithCancel(command.Context())
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "username", "cisco", "password", "cisco123")

		id, err := containerzClient.UpdateContainer(ctx, instance, image, tag, cntCommand, updateAsync, client.WithEnv(envs), client.WithPorts(ports), client.WithVolumes(volumes))
		if err != nil {
			return err
		}

		if updateAsync {
			fmt.Printf("Container with id %s started updating asynchronously.\n", id)
		} else {
			fmt.Printf("Container with id %s updated successfully.\n", id)
		}
		return nil
	},
}

func init() {
	containerCmd.AddCommand(cntUpdateCmd)

	cntUpdateCmd.PersistentFlags().StringVar(&cntCommand, "command", "/bin/bash", "command to run.")
	cntUpdateCmd.PersistentFlags().StringVar(&instance, "instance", "", "Container to update.")
	cntUpdateCmd.PersistentFlags().StringVar(&network, "network", "", "Network to attach container to.")
	cntUpdateCmd.PersistentFlags().StringVar(&runAs, "runas", "", "User to use (format: <user>[:<group>]")
	cntUpdateCmd.PersistentFlags().StringVar(&restartPolicy, "restart_policy", "", "Restart policy to use. "+
		"Valid policies are \"always\", \"on-failure\", \"unless-stopped\", and \"none\". "+
		"Some policies (e.g., \"on-failure\") optionally accept a maximum number of restart attempts. "+
		"(format: <policy>[:<max_attempts>])")
	cntUpdateCmd.PersistentFlags().StringArrayVar(&ports, "port", []string{}, "Ports to expose (format: <internal_port>:<external_port>")
	cntUpdateCmd.PersistentFlags().StringArrayVar(&envs, "env", []string{}, "Environment vars to set (format: <VAR_NAMEt>=<VAR_VALUE>")
	cntUpdateCmd.PersistentFlags().StringArrayVarP(&volumes, "volume", "v", []string{}, "Volumes to attach to the container (format: <volume-name>:<mountpoint>[:ro])")
	cntUpdateCmd.PersistentFlags().StringArrayVarP(&devices, "device", "d", []string{}, "Devices to attach to the container (format: <src-path>[:<dst-path>[:<permissions>]])")
	cntUpdateCmd.PersistentFlags().StringArrayVar(&addCaps, "add_caps", []string{}, "Capabilities to add.")
	cntUpdateCmd.PersistentFlags().StringArrayVar(&delCaps, "del_caps", []string{}, "Capabilities to remove.")
	cntUpdateCmd.PersistentFlags().StringVar(&cntCommand, "command", "/bin/bash", "command to run.")
	cntUpdateCmd.PersistentFlags().StringVar(&instance, "instance", "", "Name of the container to update.")
	cntUpdateCmd.PersistentFlags().StringArrayVar(&ports, "port", []string{}, "Ports to expose (format: <internal_port>:<external_port>")
	cntUpdateCmd.PersistentFlags().StringArrayVar(&envs, "env", []string{}, "Environment vars to set (format: <VAR_NAME>=<VAR_VALUE>")
	cntUpdateCmd.PersistentFlags().StringArrayVarP(&volumes, "volume", "v", []string{}, "Volumes to attach to the container (format: <volume-name>:<mountpoint>[:ro])")
	cntUpdateCmd.PersistentFlags().BoolVar(&updateAsync, "async", false, "Run the update operation asynchronously")
}
