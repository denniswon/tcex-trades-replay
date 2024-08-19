#!/bin/bash

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    *)          machine="UNKNOWN:${unameOut}"
esac

echo "Setting up on ${machine}"

if [ "$machine" == "Mac" ]; then
    # Install dependencies
    which -s brew
    if [[ $? != 0 ]] ; then
        # Install Homebrew
        echo "Installing Homebrew"
        ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
    else
        echo "Homebrew already installed. Running brew update"
        brew update
    fi

    # Install Postgres
    brew tap homebrew/core

    # Install Redis
    # Install Postgres
    brew list redis || brew install redis
    brew services restart redis

    # Optional: update redis configuration file to allow connections from anywhere.
    # This is not as secure as binding to localhost.
    # From `/opt/homebrew/etc/redis.conf`, uncomment line `bind 127.0.0.1 ::1`
    # Then, restart redis with `brew services restart redis`

    # Optional: configure a redis password
    # From `/opt/homebrew/etc/redis.conf`, uncomment line `# requirepass foobared

fi


