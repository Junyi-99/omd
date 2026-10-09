#!/bin/bash
export USER_BIN_HOME="${USER_CONFIG_HOME}/bin"
export USER_DEP_HOME="${USER_CONFIG_HOME}/dep"
export XDG_DATA_HOME="${USER_CONFIG_HOME}/data"
export XDG_CONFIG_HOME="${USER_CONFIG_HOME}/etc"
export XDG_CACHE_HOME="${USER_CONFIG_HOME}/cache"

# First-run only: create XDG dirs and ensure bin is executable.
# Gate on a sentinel file so we don't pay this cost every shell startup.
if [ ! -f "${XDG_CACHE_HOME}/.xdg-bootstrapped" ]; then
    mkdir -p "${USER_BIN_HOME}" "${USER_DEP_HOME}" "${XDG_CONFIG_HOME}" \
             "${XDG_DATA_HOME}" "${XDG_CACHE_HOME}" >/dev/null
    chmod -R +x "${USER_BIN_HOME}"
    touch "${XDG_CACHE_HOME}/.xdg-bootstrapped"
fi

export PATH="${USER_BIN_HOME}:${PATH}"
