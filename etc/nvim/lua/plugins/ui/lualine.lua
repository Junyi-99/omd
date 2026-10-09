return {
    "nvim-lualine/lualine.nvim",
    dependencies = { "nvim-tree/nvim-web-devicons", "SmiteshP/nvim-navic", "ofseed/copilot-status.nvim" },
    config = function()
        local status, lualine = pcall(require, "lualine")
        if not status then
            vim.notify("Something went wrong", "warn", lualine)
            return
        end

        local navic = require("nvim-navic")
        local curTime = function()
            return os.date("%H:%M:%S", os.time())
        end

        local hostname = function()
            local cache
            return function()
                if cache then
                    return cache
                end
                local data = vim.uv.os_gethostname()
                cache = string.gsub(data, "%.%w+$", "")
                return cache
            end
        end

        local copilot_status = require("copilot-status").get_status

        lualine.setup({
            options = {
                theme = "auto",
                component_separators = {
                    left = "|",
                    right = "|",
                },
                section_separators = {
                    left = " ",
                    right = "",
                },
            },
            extensions = { "nvim-tree", "toggleterm" },
            sections = {
                lualine_a = {
                    {
                        "mode",
                        separator = {
                            left = "",
                            right = "",
                        },
                        padding = {
                            left = 1,
                            right = 1,
                        },
                    },
                },
                lualine_b = { "branch", "diff", "diagnostics" },
                lualine_c = {
                    {
                        "filename",
                        path = 1,
                    },
                    {
                        function()
                            return navic.get_location()
                        end,
                        cond = function()
                            return navic.is_available()
                        end,
                    },
                },
                lualine_x = {
                    "encoding",
                    "filesize",
                    {
                        "fileformat",
                        symbols = {
                            unix = "", -- e712
                            dos = "", -- e70f
                            mac = "", -- e711
                        },
                    },
                    "filetype",
                },
                lualine_y = { { copilot_status }, { curTime } },
                -- lualine_z = { "location" },
                lualine_z = {
                    {
                        "location",
                        separator = {
                            left = "",
                        },
                        padding = {
                            left = 0,
                            right = 1,
                        },
                    },
                    {
                        "progress",
                        separator = {
                            right = "",
                        },
                        icon = {
                            "󰇽",
                            align = "left",
                        },
                        padding = {
                            left = 0,
                            right = 0,
                        },
                    },
                },
            },
            inactive_sections = {
                lualine_a = {
                    {
                        "mode",
                        separator = {
                            left = "",
                            right = "",
                        },
                        padding = {
                            left = 0,
                            right = 1,
                        },
                    },
                },
                lualine_b = {},
                lualine_c = { "filetype" },
                lualine_x = {},
                lualine_y = {},
                lualine_z = {},
            },
            tabline = {},
            winbar = {},
            inactive_winbar = {},
        })
    end,
}
