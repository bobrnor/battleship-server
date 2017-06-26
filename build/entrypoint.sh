#!/bin/sh

set -e

InitDataDir()
{
	echo "Initializing data dir"
	mkdir -p /external/logs \
	&& find /external -type d -exec chmod 755 {} + \
	&& find /external -type f -exec chmod 644 {} +
}

InitDataDir

exec "$@"
