return {
    { "neovim/nvim-lspconfig" },
    { "williamboman/mason-lspconfig.nvim" },
    {
        "williamboman/mason.nvim",
        dependencies = { "neovim/nvim-lspconfig", "williamboman/mason-lspconfig.nvim" },
        config = function()
            require("mason").setup({
                ui = {
                    border = "rounded",
                    width = 0.7,
                    height = 0.8,
                    icons = {
                        -- package_installed = '󰺧',
                        -- package_pending = '',
                        -- package_uninstalled = '󰺭'
                        package_installed = "✓",
                        package_pending = "➜",
                        package_uninstalled = "✗",
                    },
                },
            })

            -- The order is important!
            require("mason-lspconfig").setup({
                ensure_installed = { "pyright", "lua_ls" },
                -- Enable installed servers through Neovim's native LSP API.
                automatic_enable = true,
            })
        end,
    },
}
