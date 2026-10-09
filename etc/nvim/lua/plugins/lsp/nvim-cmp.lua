return {
    "hrsh7th/nvim-cmp",
    dependencies = {
        "neovim/nvim-lspconfig",
        "hrsh7th/cmp-nvim-lsp",
        "hrsh7th/cmp-buffer",
        "hrsh7th/cmp-path",
        "hrsh7th/cmp-cmdline", -- For luasnip users.
        "L3MON4D3/LuaSnip",
        "saadparwaiz1/cmp_luasnip",
        "hrsh7th/cmp-calc",
        "onsails/lspkind-nvim",
        "windwp/nvim-autopairs",
    },
    event = { "InsertEnter", "CmdlineEnter" },
    config = function()
        local luasnip = require("luasnip")
        local lspkind = require("lspkind")
        local autopairs = require("nvim-autopairs")

        local cmp_autopairs = require("nvim-autopairs.completion.cmp")
        local cmp = require("cmp")
        local has_words_before = function()
            unpack = unpack or table.unpack
            local line, col = unpack(vim.api.nvim_win_get_cursor(0))
            return col ~= 0 and vim.api.nvim_buf_get_lines(0, line - 1, line, true)[1]:sub(col, col):match("%s") == nil
        end

        autopairs.setup()
        cmp.event:on("confirm_done", cmp_autopairs.on_confirm_done())

        cmp.setup({
            snippet = {
                -- REQUIRED - 必须指定一个 snippet engine
                expand = function(args)
                    luasnip.lsp_expand(args.body)
                end,
            },

            -- 提示窗口带边框，视觉效果会更好一些
            window = {
                completion = cmp.config.window.bordered(),
                documentation = cmp.config.window.bordered(),
            },

            mapping = {
                ["<Tab>"] = cmp.mapping(function(fallback)
                    if cmp.visible() then
                        local entry = cmp.get_selected_entry()
                        if not entry then
                            cmp.select_next_item({
                                behavior = cmp.SelectBehavior.Select,
                            })
                        else
                            cmp.confirm()
                        end
                    elseif luasnip.expand_or_jumpable() then
                        luasnip.expand_or_jump()
                    elseif has_words_before() then
                        cmp.complete()
                    else
                        fallback()
                    end
                end, { "i", "s" }),
                ["<S-Tab>"] = cmp.mapping(function(fallback)
                    if cmp.visible() then
                        cmp.select_prev_item()
                    elseif luasnip.jumpable(-1) then
                        luasnip.jump(-1)
                    else
                        fallback()
                    end
                end, { "i", "s" }),
                ["<CR>"] = cmp.mapping({
                    i = function(fallback)
                        if cmp.visible() then
                            cmp.confirm({
                                select = true,
                            })
                        else
                            fallback()
                        end
                    end,
                    s = cmp.mapping.confirm({
                        select = true,
                    }),
                }),
                ["<Down>"] = cmp.mapping(function(fallback)
                    if cmp.visible() then
                        cmp.select_next_item()
                    else
                        fallback()
                    end
                end, { "i", "s" }),
                ["<Up>"] = cmp.mapping(function(fallback)
                    if cmp.visible() then
                        cmp.select_prev_item()
                    else
                        fallback()
                    end
                end, { "i", "s" }),
                ["<C-n>"] = cmp.mapping(function(fallback)
                    if cmp.visible() then
                        cmp.select_next_item()
                    else
                        fallback()
                    end
                end, { "i", "s" }),
                ["<C-p>"] = cmp.mapping(function(fallback)
                    if cmp.visible() then
                        cmp.select_prev_item()
                    else
                        fallback()
                    end
                end, { "i", "s" }),
                ["<C-e>"] = cmp.mapping(cmp.mapping.abort(), { "i", "s" }),
                ["<C-Space>"] = cmp.mapping(cmp.mapping.complete(), { "i", "s" }),
            },

            sources = cmp.config.sources({
                {
                    name = "nvim_lsp",
                },
                {
                    name = "luasnip",
                },
                {
                    name = "buffer",
                },
                {
                    name = "path",
                },
                {
                    name = "calc",
                },
                {
                    name = "copilot",
                },
            }),
            formatting = {
                format = lspkind.cmp_format({
                    mode = "symbol_text", -- show only symbol annotations
                    maxwidth = 50, -- prevent the popup from showing more than provided characters (e.g 50 will not show more than 50 characters)
                    ellipsis_char = "...", -- when popup menu exceed maxwidth, the truncated part would show ellipsis_char instead (must define maxwidth first)
                    symbol_map = { Copilot = "" },
                }),
            },
        })

        cmp.setup.filetype("gitcommit", {
            sources = cmp.config.sources({
                {
                    name = "git",
                },
                {
                    name = "buffer",
                },
            }),
        })

        -- 给 vim 的 : 命令提供补全
        cmp.setup.cmdline(":", {
            mapping = {
                -- 当按下 Tab 时，选择下一个补全项
                ["<Tab>"] = {
                    c = function()
                        if cmp.visible() then
                            cmp.select_next_item()
                        else
                            cmp.complete()
                        end
                    end,
                },
                -- 当按下 Shift + Tab 时，选择上一个补全项
                ["<S-Tab>"] = {
                    c = function()
                        if cmp.visible() then
                            cmp.select_prev_item()
                        else
                            cmp.complete()
                        end
                    end,
                },
                ["<Down>"] = {
                    c = function(fallback)
                        if cmp.visible() then
                            cmp.select_next_item()
                        else
                            fallback()
                        end
                    end,
                },
                ["<Up>"] = {
                    c = function(fallback)
                        if cmp.visible() then
                            cmp.select_prev_item()
                        else
                            fallback()
                        end
                    end,
                },
                -- 当按下 Ctrl + e 时，关闭补全窗口
                ["<C-e>"] = {
                    c = cmp.mapping.abort(),
                },
                -- 当按下回车时
                ["<CR>"] = {
                    c = cmp.mapping.confirm({
                        select = false,
                    }),
                },
            },
            sources = cmp.config.sources(
                { {
                    name = "path",
                } },
                { {
                    name = "cmdline",
                } }
            ),
        })

        -- Set up lspconfig.
        -- local capabilities = require('cmp_nvim_lsp').default_capabilities()
        -- vim.notify("cmp_nvim_lsp capabilities: " .. vim.inspect(capabilities))

        -- Replace <YOUR_LSP_SERVER> with each lsp server you've enabled.
        -- require('lspconfig')['<YOUR_LSP_SERVER>'].setup {
        --     capabilities = capabilities
        -- }
    end,
}
