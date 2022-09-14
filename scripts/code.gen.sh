#!/bin/bash
#
# MIT License
#
# (C) Copyright 2022 Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.
#

# format swagger doc
swag fmt

# update swagger doc yaml
swag init --md  docs/ --outputTypes go,yaml --exclude api/controllers/v1/misc

# fix copyright headers
docker run -it --rm -v $(pwd):/github/workspace artifactory.algol60.net/csm-docker/stable/license-checker --fix docs

# update swagger.md
if ! command -v swagger-markdown &> /dev/null
then
    npx swagger-markdown -i  docs/swagger.yaml || true
else 
    swagger-markdown -i  docs/swagger.yaml || true
fi

# 
go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.9.0
~/go/bin/controller-gen crd webhook paths="api/models/v1/hooks.go" output:crd:artifacts:config="api/metacontroller"