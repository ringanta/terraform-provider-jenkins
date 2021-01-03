#!/bin/bash

set -euo pipefail

cd jenkins_home
ls -a | grep -v '\.$' | grep -v '^plugins' | xargs rm -fr

