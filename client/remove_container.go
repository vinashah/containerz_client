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

package client

import (
	"context"

	cpb "github.com/openconfig/gnoi/containerz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (

	// ErrNotFound indicates that the specified image was not found on the target system.
	ErrNotFound = status.Error(codes.NotFound, "resource was not found")

	// ErrRunning indicates that there is a container running this image.
	ErrRunning = status.Error(codes.FailedPrecondition, "resource is running")
	ErrUnspecified = status.Error(codes.FailedPrecondition, "unspecified")
)

// RemoveContainer removes an image from the target system. It returns nil upon success. Otherwise it
// returns an error indicating whether the image was not found or is associated to running
// container.
//
// Deprecated -- Use ImageRemove instead.
func (c *Client) RemoveContainer(ctx context.Context, cnt string, forceopt ...bool) error {
	force := false
	if len(forceopt) > 0 {
		force = forceopt[0]
	}

	if _, err := c.cli.RemoveContainer(ctx, &cpb.RemoveContainerRequest{
		Name:  cnt,
		Force: force,
	}); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		switch st.Code() {
		case codes.FailedPrecondition:
			return ErrRunning
		case codes.NotFound:
			return ErrNotFound
		default:
			return status.Errorf(codes.Unknown, "unknown error: %v", st.Message())
		}
	}

	return nil
}
