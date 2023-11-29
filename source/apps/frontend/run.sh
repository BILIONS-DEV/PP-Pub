#!/usr/bin/env bash
go build
echo "go build"
supervisorctl restart frontend
echo "supervisorctl restart"
$SHELL
