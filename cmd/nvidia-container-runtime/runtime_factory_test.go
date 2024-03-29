/*
# Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/NVIDIA/nvidia-container-toolkit/internal/config"
	"github.com/opencontainers/runtime-spec/specs-go"
	testlog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestFactoryMethod(t *testing.T) {
	logger, _ := testlog.NewNullLogger()

	testCases := []struct {
		description   string
		cfg           *config.Config
		spec          *specs.Spec
		expectedError bool
	}{
		{
			description: "empty config no error",
			cfg: &config.Config{
				NVIDIAContainerRuntimeConfig: config.RuntimeConfig{},
			},
		},
		{
			description: "experimental flag supported",
			cfg: &config.Config{
				NVIDIAContainerRuntimeConfig: config.RuntimeConfig{
					Experimental: true,
					DiscoverMode: "legacy",
				},
			},
			spec: &specs.Spec{
				Process: &specs.Process{
					Env: []string{
						"NVIDIA_VISIBLE_DEVICES=all",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			bundleDir := t.TempDir()

			specFile, err := os.Create(filepath.Join(bundleDir, "config.json"))
			require.NoError(t, err)
			require.NoError(t, json.NewEncoder(specFile).Encode(tc.spec))

			argv := []string{"--bundle", bundleDir}

			_, err = newNVIDIAContainerRuntime(logger, tc.cfg, argv)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
