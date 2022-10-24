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

import v1 "github.com/Cray-HPE/cray-nls/src/api/models/v1"

type CreateOrPatchActivityRequest struct {
	Products []Product                `json:"products" validate:"required"`
	MediaDir string                   `json:"media-dir"  validate:"required"`
	Inputs   v1.IufSessionInputParams `json:"inputs"  validate:"required"`
} // @name Iuf.CreateOrPatchActivityRequest

type Product struct {
	Name    string `json:"name" validate:"required"`
	Version string `json:"version" validate:"required"`
} // @name Iuf.Product

type Activity struct {
	Name string `json:"name" validate:"required"`
	CreateOrPatchActivityRequest
	Sessions []struct {
		Name string
		v1.IufSession
	} `json:"sessions,omitempty"`
} // @name Iuf.Activity
