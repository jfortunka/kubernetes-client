/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package strategy

import (
	"strconv"
	"strings"

	kapi "k8s.io/kubernetes/pkg/apis/core"
)

// DeploymentStrategy knows how to make a deployment active.
type DeploymentStrategy interface {
	// Deploy transitions an old deployment to a new one.
	Deploy(from *kapi.ReplicationController, to *kapi.ReplicationController, desiredReplicas int) error
}

// UpdateAcceptor is given a chance to accept or reject the new controller
// during a deployment each time the controller is scaled up.
//
// After the successful scale-up of the controller, the controller is given to
// the UpdateAcceptor. If the UpdateAcceptor rejects the controller, the
// deployment is stopped with an error.
//
// DEPRECATED: Acceptance checking has been incorporated into the rolling
// strategy, but we still need this around to support Recreate.
type UpdateAcceptor interface {
	// Accept returns nil if the controller is okay, otherwise returns an error.
	Accept(*kapi.ReplicationController) error
}

type errConditionReached struct {
	msg string
}

func NewConditionReachedErr(msg string) error {
	return &errConditionReached{msg: msg}
}

func (e *errConditionReached) Error() string {
	return e.msg
}

func IsConditionReached(err error) bool {
	value, ok := err.(*errConditionReached)
	return ok && value != nil
}

func PercentageBetween(until string, min, max int) bool {
	if !strings.HasSuffix(until, "%") {
		return false
	}
	until = until[:len(until)-1]
	i, err := strconv.Atoi(until)
	if err != nil {
		return false
	}
	return i >= min && i <= max
}

func Percentage(until string) (int, bool) {
	if !strings.HasSuffix(until, "%") {
		return 0, false
	}
	until = until[:len(until)-1]
	i, err := strconv.Atoi(until)
	if err != nil {
		return 0, false
	}
	return i, true
}