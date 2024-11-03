#!/bin/bash

# Run your other command here
export PATH=/root/go/bin:/usr/lib/jvm/java-11-openjdk/bin:/work_dir/blackdaggerenv/bin:/root/go/bin:/sbin:/usr/bin:/root/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:/opt/bin:/usr/bin/core_perl:/root/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:/opt/bin:/usr/bin/core_perl:/usr/games/bin:
blackdagger server &
blackdagger pull devsecops
blackdagger pull cart
# Run gotty in the foreground
gotty -p 8090 -w --credential blackdagger:blackdagger /bin/bash
