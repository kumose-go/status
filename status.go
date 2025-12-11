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

import (
	"runtime"

	"github.com/kumose-go/clog"
)

type StatusCode int

const (
	Ok                 StatusCode = 0
	Cancelled          StatusCode = 1
	Unknown            StatusCode = 2
	InvalidArgument    StatusCode = 3
	DeadlineExceeded   StatusCode = 4
	NotFound           StatusCode = 5
	AlreadyExists      StatusCode = 6
	PermissionDenied   StatusCode = 7
	ResourceExhausted  StatusCode = 8
	FailedPrecondition StatusCode = 9
	Aborted            StatusCode = 10
	OutOfRange         StatusCode = 11
	Unimplemented      StatusCode = 12
	Internal           StatusCode = 13
	Unavailable        StatusCode = 14
	DataLoss           StatusCode = 15
	Unauthenticated    StatusCode = 16
	IoError            StatusCode = 17
)

type Level clog.Level

const (
	InvalidLevel Level = Level(clog.InvalidLevel)
	DebugLevel   Level = Level(clog.DebugLevel)
	InfoLevel    Level = Level(clog.InfoLevel)
	WarnLevel    Level = Level(clog.WarnLevel)
	ErrorLevel   Level = Level(clog.ErrorLevel)
	FatalLevel   Level = Level(clog.FatalLevel)
)

func (l Level) String() string {
	return clog.Level(l).String()
}

type Frame struct {
	Pc   uintptr
	File string
	Line int
	Ok   bool
}

type Status interface {
	Code() StatusCode

	Level() Level

	Message() string

	Frame() []Frame
	// if not exist, return empty string
	GetPayload(uri string) string
	//  return true means continue, false stop
	WalkPayload(func(string, string) bool)
	/// builders

	SetPayload(uri string, data string) Status

	// capture frame
	Capture() Status

	Ok() bool
}

type status struct {
	code    StatusCode
	level   Level
	message string
	payload map[string]string
	frame   []Frame
}

func NewStatus() Status {
	return &status{
		payload: make(map[string]string),
		level:   InfoLevel,
		code:    Ok,
		frame:   make([]Frame, 0),
	}
}

func NewStatusImpl() *status {
	return &status{
		payload: make(map[string]string),
		level:   InfoLevel,
		code:    Ok,
		frame:   make([]Frame, 0),
	}
}

func (s *status) Code() StatusCode {
	return s.code
}

func (s *status) Level() Level {
	return s.level
}

func (s *status) Frame() []Frame {
	return s.frame
}

func (s *status) Message() string {

	return s.message
}

func (s *status) SetPayload(uri, data string) Status {
	s.payload[uri] = data
	return s
}

func (s *status) Ok() bool {
	return s.Code() == Ok
}

func (s *status) GetPayload(uri string) string {
	return s.payload[uri]
}

func (s *status) WalkPayload(fn func(string, string) bool) {
	for k, v := range s.payload {
		if !fn(k, v) {
			break
		}
	}
}

func (s *status) Capture() Status {
	var f Frame
	f.Pc, f.File, f.Line, f.Ok = runtime.Caller(0)
	s.frame = append(s.frame, f)
	return s
}

func Debug(code StatusCode, msg string) Status {
	s := NewStatusImpl()
	s.code = code
	s.level = DebugLevel
	s.message = msg
	return s
}

func Info(code StatusCode, msg string) Status {
	s := NewStatusImpl()
	s.code = code
	s.level = InfoLevel
	s.message = msg
	return s
}

func Warn(code StatusCode, msg string) Status {
	s := NewStatusImpl()
	s.code = code
	s.level = WarnLevel
	s.message = msg
	return s
}

func Error(code StatusCode, msg string) Status {
	s := NewStatusImpl()
	s.code = code
	s.level = ErrorLevel
	s.message = msg
	return s
}

func Fatal(code StatusCode, msg string) Status {
	s := NewStatusImpl()
	s.code = code
	s.level = FatalLevel
	s.message = msg
	return s
}

func NewOk(msg string) Status {
	return Info(Ok, msg)
}
