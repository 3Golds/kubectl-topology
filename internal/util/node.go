// Copyright 2020 bmcustodio
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

package util

import (
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	Name   string
	Region string
	Zone   string
	Taint  string
	Age    string
	Label  string
}

func NewNode(node corev1.Node, l string) Node {
	r := Node{
		Name: node.Name,
	}
	var strSlice []string
	var labelSlice []string
	if v, exists := node.Labels[RegionLabel]; exists && v != "" {
		r.Region = v
	} else {
		r.Region = node.Labels[Pre117RegionLabel]
	}
	if v, exists := node.Labels[ZoneLabel]; exists && v != "" {
		r.Zone = v
	} else {
		r.Zone = node.Labels[Pre117ZoneLabel]
	}
	if node.Spec.Taints != nil {
		for i := 0; i < len(node.Spec.Taints); i++ {
			str := node.Spec.Taints[i].Key + "=" + node.Spec.Taints[i].Value + ":" + string(node.Spec.Taints[i].Effect)
			strSlice = append(strSlice, str)
		}
	} else {
		r.Taint = "<none>"
	}
	r.Age = strconv.FormatFloat(time.Now().Sub(node.CreationTimestamp.Time).Hours(), 'f', 1, 64)+"h"
	flagLabelSlice := strings.Split(l, ",")
	for i := 0; i < len(flagLabelSlice); i++ {
		if flagLabelSlice[i] == "" {
			labelSlice = append(labelSlice, "<none>")
			continue
		}
		if node.Labels[flagLabelSlice[i]] == "" {
			labelSlice = append(labelSlice, "")
			continue
		}
		str := flagLabelSlice[i] + "=" + node.Labels[flagLabelSlice[i]]
		labelSlice = append(labelSlice, str)
	}
	r.Label = strings.Join(labelSlice, ",")
	//r.AppLabel = node.Labels["app"]
	r.Taint = strings.Join(strSlice, ",")
	return r
}

type NodeList []Node

func (l NodeList) Headers() string {
	return "NAME\tREGION\tZONE\tTAINTS\tAge\t\tLABEL\n"
}

func (l NodeList) Items() []string {
	r := make([]string, 0, len(l))
	for _, ll := range l {
		r = append(r, ll.Name+"\t"+ll.Region+"\t"+ll.Zone+"\t"+ll.Taint+"\t"+ll.Age+"\t\t"+ll.Label+"\n")
		//r = append(r, ll.Name+"\t"+ll.Region+"\t"+ll.Zone+"\t"+ll.Taint+"\t"+ll.Age+"\t"+ll.Label+"\n")
	}
	return r
}

func (l NodeList) Length() int {
	return len(l)
}