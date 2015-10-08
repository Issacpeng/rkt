// Copyright 2015 The rkt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/coreos/rkt/Godeps/_workspace/src/github.com/steveeJ/gexpect"
)

var expectedResults = []string{
	"prestart OK",
	"main OK",
	"sidekick OK",
	"poststop OK",
}

func TestAceValidator(t *testing.T) {
	ctx := newRktRunCtx()
	defer ctx.cleanup()

	if err := ctx.launchMDS(); err != nil {
		t.Fatalf("Cannot launch metadata service: %v", err)
	}

	aceMain := os.Getenv("RKT_ACE_MAIN_IMAGE")
	if aceMain == "" {
		panic("empty RKT_ACE_MAIN_IMAGE env var")
	}
	aceSidekick := os.Getenv("RKT_ACE_SIDEKICK_IMAGE")
	if aceSidekick == "" {
		panic("empty RKT_ACE_SIDEKICK_IMAGE env var")
	}

	rktArgs := fmt.Sprintf("--debug --insecure-skip-verify run --volume database,kind=empty %s %s",
		aceMain, aceSidekick)
	rktCmd := fmt.Sprintf("%s %s", ctx.cmd(), rktArgs)

	t.Logf("Command: %v", rktCmd)
	child, err := gexpect.Spawn(rktCmd)
	if err != nil {
		t.Fatalf("Cannot exec rkt: %v", err)
	}

	for _, e := range expectedResults {
		err = expectWithOutput(child, e)
		if err != nil {
			t.Fatalf("Expected %q but not found: %v", e, err)
		}
	}

	err = child.Wait()
	if err != nil {
		t.Fatalf("rkt didn't terminate correctly: %v", err)
	}
}