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
package services

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestInitializeWorkflowTemplate(t *testing.T) {
	t.Skip("skip: TODO")
	t.Run("It should initialize workflow template", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should update workflow template", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestCreateWorkflow(t *testing.T) {
	t.Skip("skip: TODO")
	t.Run("It can create a new workflow", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should NOT create a new workflow when there is a running one", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should prevent jobs running on the ncn being rebuilt", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestGetWorkflows(t *testing.T) {
	t.Skip("skip: TODO")
	t.Run("It should get workflows", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should get workflows by (type,ncn, some other filters .....)", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}
