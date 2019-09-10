#!/bin/bash

go build -o service

export HUMAN_LOG=1

./service