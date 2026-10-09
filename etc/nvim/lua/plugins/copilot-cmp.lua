return {
    {
        "zbirenbaum/copilot-cmp",
        config = function()
            local source = require("copilot_cmp.source")
            source.is_available = function(self)
                if self.client:is_stopped() or self.client.name ~= "copilot" then
                    return false
                end

                local clients = vim.lsp.get_clients({
                    bufnr = vim.api.nvim_get_current_buf(),
                    id = self.client.id,
                })
                return next(clients) ~= nil
            end

            require("copilot_cmp").setup()
        end,
    },
}
