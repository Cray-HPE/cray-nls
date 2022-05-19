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
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

# HMS Build changed charts container image
HMS_BUILD_IMAGE ?= hms-build-changed-charts-action:local 

# Helm Chart
TARGET_BRANCH ?= main
UNSTABLE_BUILD_SUFFIX ?= "" # If this variable is the empty string, then this is a stable build
							# Otherwise, if this variable is non-empty then this is an unstable build

all-charts:
	docker run --rm -it -v $(shell pwd):/workspace ${HMS_BUILD_IMAGE} build_all_charts.sh ./charts

changed-charts: ct-config
	# If the repo was cloned with SSH, then the docker container needs those credentails to interact with the 
	# locally checkouted repo. TODO for now this is broken for not macOS. 
	# The following works on macOS, assuming you have ran "ssh-add" to add your SSH identity to the SSH agent.
	docker run --rm -it -v $(shell pwd):/workspace \
		-v /run/host-services/ssh-auth.sock:/ssh-agent -e SSH_AUTH_SOCK=/ssh-agent \
		${HMS_BUILD_IMAGE} build_changed_charts.sh ./charts ${TARGET_BRANCH}

ct-config:
	git checkout -- ct.yaml
	docker run --rm -v $(shell pwd):/workspace ${HMS_BUILD_IMAGE} update-ct-config-with-chart-dirs.sh charts

lint: ct-config
	docker run --rm -it -v $(shell pwd):/workspace ${HMS_BUILD_IMAGE} ct lint --config ct.yaml

clean:
	git checkout -- ct.yaml
	docker run --rm -it -v $(shell pwd):/workspace ${HMS_BUILD_IMAGE} clean_all_charts.sh ./charts
	rm -rf .packaged
