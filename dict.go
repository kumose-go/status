// Copyright (C) Kumo inc. and its affiliates.
// Author: Jeff.li lijippy@163.com
// All rights reserved.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package status

type DictCombiner interface {
	Combine(children map[string]Status) Status
}

// DictStatus aggregates multiple Status objects, including a base Status and child Statuses.
type DictStatus struct {
	base     Status
	children map[string]Status
	combiner DictCombiner
}

// NewDictStatus creates a DictStatus with a given base Status and default combiner.
func NewDictStatus() *DictStatus {
	return &DictStatus{
		base:     nil,
		children: make(map[string]Status),
		combiner: &DefaultCombiner{},
	}
}

// SetCombiner sets a custom Combiner for aggregation.
func (d *DictStatus) ensureCombined() {
	if d.base == nil {
		d.base = d.combiner.Combine(d.children)
	}
}

// SetCombiner sets a custom Combiner for aggregation.
func (d *DictStatus) SetCombiner(c DictCombiner) *DictStatus {
	d.combiner = c
	return d
}

// AddChild adds a child Status to the DictStatus.
func (d *DictStatus) AddChild(key string, s Status) *DictStatus {
	if d.base != nil {
		panic("")
	}
	d.children[key] = s
	return d
}

// Code aggregates the base and child Status codes using the Combiner.
func (d *DictStatus) Code() StatusCode {
	d.ensureCombined()
	return d.base.Code()
}

// Level aggregates the base and child Levels using the Combiner.
func (d *DictStatus) Level() Level {
	d.ensureCombined()
	return d.base.Level()
}

// Message returns the base Status message.
func (d *DictStatus) Message() string {
	d.ensureCombined()
	return d.base.Message()
}

// Frame aggregates the call frames from base and child Statuses.
func (d *DictStatus) Frame() []Frame {
	d.ensureCombined()
	return d.base.Frame()
}

// GetPayload retrieves a payload value, checking the base first then children.
func (d *DictStatus) GetPayload(uri string) string {
	d.ensureCombined()
	return d.base.GetPayload(uri)
}

// WalkPayload iterates over all payloads in base and children.
func (d *DictStatus) WalkPayload(fn func(string, string) bool) {
	d.base.WalkPayload(fn)
	for _, s := range d.children {
		s.WalkPayload(fn)
	}
}

// SetPayload sets a payload on the base Status.
func (d *DictStatus) SetPayload(uri, data string) Status {
	d.ensureCombined()
	d.base.SetPayload(uri, data)
	return d
}

// Capture captures the current call frame for base and all child Statuses.
func (d *DictStatus) Capture() Status {
	d.ensureCombined()
	d.base.Capture()
	for _, s := range d.children {
		s.Capture()
	}
	return d
}

// DefaultCombiner provides default aggregation logic for DictStatus.
// base is already the aggregated result so far.
type DefaultCombiner struct{}

// CombineCode updates the base StatusCode based on child Status codes.
func (d DefaultCombiner) Combine(s map[string]Status) Status {
	for _, c := range s {
		if c.Code() != Ok {
			return c
		}
	}
	return Info(Ok, "")
}
