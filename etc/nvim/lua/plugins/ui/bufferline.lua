return {
    "akinsho/bufferline.nvim",
    version = "*",
    dependencies = "nvim-tree/nvim-web-devicons",
    config = function()
        require("bufferline").setup({
            options = {
                offsets = {
                    {
                        filetype = "NvimTree",
                        text = "File Explorer",
                        text_align = "left",
                        separator = true,
                    },
                },
                diagnostics = "nvim_lsp",
                close_command = "bdelete! %d", -- 点击关闭按钮关闭
                right_mouse_command = "bdelete! %d", -- 右键点击关闭
                indicator = {
                    icon = "▎", -- 分割线
                    style = "underline",
                },
                hover = {
                    enabled = true,
                    delay = 200,
                    reveal = { "close" },
                },
            },
        })
    end,
}
