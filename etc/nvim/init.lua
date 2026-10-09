-- vim.g.mapleader = " " -- 设置 <leader> 为 <space>

-- Lazy plugin manager
local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.uv.fs_stat(lazypath) then
    vim.fn.system({
        "git",
        "clone",
        "--filter=blob:none",
        "https://github.com/folke/lazy.nvim.git",
        "--branch=stable", -- latest stable release
        lazypath,
    })
end

-- Disable netrw before plugins load to avoid startup race conditions.
vim.g.loaded_netrw = 1
vim.g.loaded_netrwPlugin = 1

vim.opt.rtp:prepend(lazypath)

local lazy_opts = {
    ui = {
        border = "rounded",
        title = "Plugin Manager",
        title_pos = "center",
    },
}

-- Lazy plugin manager
require("lazy").setup("plugins", lazy_opts)

vim.cmd([[
syntax on
set number
set nowrap
set tabstop=4
set cursorline

set shiftwidth=4
set expandtab

" Blink cursor on error instead of beeping (grr)
set visualbell]])

vim.opt.mousemoveevent = true
vim.opt.termguicolors = true

-- local function open_nvim_tree()
--     -- open the tree
--     require("nvim-tree.api").tree.open()
-- end

-- vim.api.nvim_create_autocmd({"VimEnter"}, {
--     callback = open_nvim_tree
-- })
