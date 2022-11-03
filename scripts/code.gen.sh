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
swag init --md  docs/ --outputTypes go,yaml \
    --exclude src/api/controllers/v1/misc,src/api/controllers/v1/iuf \
    --instanceName NLS

# update iuf swagger doc yaml
swag init --md docs/ --outputTypes go,yaml \
    --exclude src/api/controllers/v1/misc,src/api/controllers/v1/nls \
    --instanceName IUF --parseDependency --parseDepth 1

# fix copyright headers
docker run -it --rm -v $(pwd):/github/workspace artifactory.algol60.net/csm-docker/stable/license-checker --fix docs

go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.9.0
~/go/bin/controller-gen crd webhook paths="./src/api/models/nls/v1/..." output:crd:artifacts:config="src/api/services/nls"


# mockgen
~/go/bin/mockgen -destination=src/api/mocks/services/workflow.go -package=mocks -source=src/api/services/shared/workflow.go
~/go/bin/mockgen -destination=src/api/mocks/services/iuf.go -package=mocks -source=src/api/services/iuf/iuf.go
~/go/bin/mockgen -destination=src/api/mocks/services/ncn.go -package=mocks -source=src/api/services/nls/ncn.go