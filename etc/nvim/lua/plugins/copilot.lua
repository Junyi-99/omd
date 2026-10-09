return {
    {
        "zbirenbaum/copilot.lua",
        cmd = "Copilot",
        event = "InsertEnter",
        build = ":Copilot auth",
        config = function()
            -- we put copilot into cmp sources.
            -- the copilot suggestion will show up in the completion menu
            require("copilot").setup({
                suggestion = {
                    enabled = false,
                }, -- Disable suggestions
                panel = {
                    enabled = false,
                }, -- Disable panel
            })
        end,
    },
}
