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
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/moby/moby/v/v24/client"
	"github.com/moby/moby/v/v24/docker"
	"github.com/openconfig/containerz/server"
)

var (
	dockerHost string
	chunkSize  int
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the containerz server",
	RunE: func(command *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(command.Context())
		defer cancel()

		cli, err := client.NewClientWithOpts(client.WithHost(dockerHost), client.WithAPIVersionNegotiation())
		if err != nil {
			return err
		}

		opts := []server.Option{
			server.WithAddr(addr),
			server.WithChunkSize(chunkSize),
		}
		mgr := docker.New(cli)
		s := server.New(docker.New(cli), opts...)

		// listen for ctrl-c
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		go func() {
			<-interrupt // wait for signal
			cancel()
			s.Halt(ctx)
			mgr.Stop()
		}()

		return s.Serve(ctx)
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVar(&dockerHost, "docker_host", "unix:///var/run/docker.sock", "Docker host to connect to.")
	startCmd.PersistentFlags().IntVar(&chunkSize, "chunk_size", 64000, "the size of the chunks supported by this server")
}
