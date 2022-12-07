/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package iuf

// Stage
type Stage struct {
	Name       string       `yaml:"name" json:"name" binding:"required"`             // Name of the stage
	Type       string       `yaml:"type" json:"type" binding:"required"`             // Type of the stage
	Operations []Operations `yaml:"operations" json:"operations" binding:"required"` // operations
} // @name Stage

type Stages struct {
	Version string  `yaml:"version" json:"version" binding:"required"`
	Stages  []Stage `yaml:"stages" json:"stages" binding:"required"`
} // @name Stages

type Operations struct {
	Name             string                 `yaml:"name" json:"name" binding:"required"`             // Name of the operation
	LocalPath        string                 `yaml:"local-path" json:"local-path" binding:"required"` // Argo operation file path
	StaticParameters map[string]interface{} `yaml:"static-parameters" json:"static-parameters" binding:"required"`
} // @name Operations
