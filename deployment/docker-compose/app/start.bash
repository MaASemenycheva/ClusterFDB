#!/bin/bash

set -eu;

/create_cluster_file.bash
FDB_CLUSTER_FILE="${FDB_CLUSTER_FILE:-/etc/foundationdb/fdb.cluster}"

# Attempt to connect. Configure the database if necessary.
if ! /usr/bin/fdbcli -C $FDB_CLUSTER_FILE --exec status --timeout 3 ; then
    echo "creating the database"
    if ! fdbcli -C $FDB_CLUSTER_FILE --exec "configure new single memory ; status" --timeout 10 ; then
        echo "Unable to configure new FDB cluster."
        exit 1
    fi
fi

/go/bin/foundationdb