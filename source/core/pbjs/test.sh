#!/bin/bash
echo "$OSTYPE"
if [[ "$OSTYPE" == "msys"* ]]; then
  ../../apps/worker/prebid_demand_builder/prebid_demand_builder.exe
else
  /home/assyrian/go/selfserve/source/apps/worker/prebid_demand_builder/prebid_demand_builder
fi