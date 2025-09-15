apt-get update && apt-get install -y locales && rm -rf /var/lib/apt/lists/*

localedef -i $LOCALE -f UTF-8 $LOCALE.UTF-8
update-locale LANG=$LOCALE.UTF-8
