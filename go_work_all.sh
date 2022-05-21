#!/usr/bin/env bash

go work use -r ./components

go work use -r ./examples



read -t 2 -p "go work completed"