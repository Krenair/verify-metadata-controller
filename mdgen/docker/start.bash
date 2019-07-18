#!/bin/bash
export LD_LIBRARY_PATH=/opt/cloudhsm/lib
export HSM_PARTITION=PARTITION_1
/opt/cloudhsm/bin/cloudhsm_client /opt/cloudhsm/etc/cloudhsm_client.cfg > cloudhsm.log 2>&1 &
bash
