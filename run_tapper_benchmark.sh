#!/bin/bash

if [[ $(uname) == "Darwin" ]]
then
  export MIZU_HOME="/Users/nimrod/Projects/mizu"
  export MIZU_BENCHMARK_OUTPUT_DIR="/Users/nimrod/Temp/mizu-benchmark-results/33.0-dev25_no-redact_no-self-tap__$(date +%d-%m-%H-%M)"
else
  export MIZU_HOME="/home/nimrod/projects/mizu"
  export MIZU_BENCHMARK_OUTPUT_DIR="/home/nimrod/temp/mizu-benchmark-results/33.0-dev25_no-redact_no-self-tap__$(date +%d-%m-%H-%M)"
fi
export MIZU_BENCHMARK_CLIENT_PERIOD="5m"
export MIZU_BENCHMARK_RUN_COUNT="3"
export MIZU_BENCHMARK_QPS="500"
export MIZU_BENCHMARK_CLIENTS_COUNT="5"

if [[ $(uname) != "Darwin" ]]
then
  sudo setcap cap_sys_admin,cap_sys_resource,cap_dac_override,cap_net_raw,cap_net_admin=eip ./agent/build/mizuagent
fi

performance_analysis/run_tapper_benchmark.sh 
