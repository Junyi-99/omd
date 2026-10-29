# PROFILE=/dev/null: the Env snippet loads nvm, not ~/.zshrc edits.
# Expects nvm_load.sh to have run first (NVM_DIR set).
mkdir -p "$NVM_DIR" && curl -fsSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.8/install.sh | PROFILE=/dev/null bash && . "$NVM_DIR/nvm.sh" && nvm install 24
