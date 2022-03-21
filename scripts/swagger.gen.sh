#!/bin/bash
swag fmt
swag init --md docs/
swagger-markdown -i  docs/swagger.yaml