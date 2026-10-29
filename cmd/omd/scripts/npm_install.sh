# Expects nvm_load.sh to have run first. Input: OMD_NPM_ARGS (flags and package, already shell-formed).
command -v npm >/dev/null || { echo "npm not found: install Node.js (nvm) first" >&2; exit 1; }
eval "npm install -g $OMD_NPM_ARGS"
