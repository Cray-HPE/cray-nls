//
//  MIT License
//
//  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
//
//  Permission is hereby granted, free of charge, to any person obtaining a
//  copy of this software and associated documentation files (the "Software"),
//  to deal in the Software without restriction, including without limitation
//  the rights to use, copy, modify, merge, publish, distribute, sublicense,
//  and/or sell copies of the Software, and to permit persons to whom the
//  Software is furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included
//  in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
//  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
//  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
//  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//  OTHER DEALINGS IN THE SOFTWARE.
//
package utils

import (
	"fmt"
	"go.uber.org/fx"
	"math/rand"
	"regexp"
)

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewRequestHandler),
	fx.Provide(NewEnv),
	fx.Provide(GetLogger),
	fx.Provide(NewValidator),
)

type GenericError struct {
	Message string
}

func (e GenericError) Error() string {
	return e.Message
}

const (
	maxNameLength          = 63
	randomLength           = 5
	maxGeneratedNameLength = maxNameLength - randomLength - 1 // -1 because we want to put a dash after the prefix
)

var (
	letterRunes          = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)
)

// RandomString A random string of the given length.
func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// GenerateName generates a name with the given prefix. Makes sure result doesn't go beyond 63 characters
func GenerateName(prefix string) string {
	prefix = nonAlphanumericRegex.ReplaceAllString(prefix, "-")
	if len(prefix) > maxGeneratedNameLength {
		prefix = prefix[:maxGeneratedNameLength]
	}
	return fmt.Sprintf("%s-%s", prefix, RandomString(randomLength))
}
