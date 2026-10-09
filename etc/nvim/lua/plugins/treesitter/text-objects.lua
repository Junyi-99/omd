return {
    "nvim-treesitter/nvim-treesitter-textobjects",
    event = "VeryLazy",
    dependencies = { "nvim-treesitter/nvim-treesitter" },
    config = function()
        local textobjects = require("nvim-treesitter-textobjects")
        textobjects.setup({
            select = {
                lookahead = true,
            },
            move = {
                set_jumps = true,
            },
        })

        local select = require("nvim-treesitter-textobjects.select")
        local move = require("nvim-treesitter-textobjects.move")

        local function select_textobject(query)
            select.select_textobject(query, "textobjects")
        end

        for lhs, query in pairs({
            ["af"] = "@function.outer",
            ["if"] = "@function.inner",
            ["aj"] = "@conditional.outer",
            ["ij"] = "@conditional.inner",
            ["al"] = "@loop.outer",
            ["il"] = "@loop.inner",
            ["ac"] = "@class.outer",
            ["ic"] = "@class.inner",
        }) do
            vim.keymap.set({ "x", "o" }, lhs, function()
                select_textobject(query)
            end, { desc = "Select " .. query })
        end

        -- The plugin errors when the buffer has no treesitter parser, so skip
        -- it there (falling back to the builtin motion for ]] and [[).
        local function map_move(lhs, fn, query, desc)
            vim.keymap.set({ "n", "x", "o" }, lhs, function()
                if vim.treesitter.get_parser(0, nil, { error = false }) then
                    move[fn](query, "textobjects")
                elseif lhs == "]]" or lhs == "[[" then
                    vim.cmd("normal! " .. vim.v.count1 .. lhs)
                end
            end, { desc = desc })
        end

        map_move("]f", "goto_next_start", "@function.outer", "Next function start")
        map_move("]]", "goto_next_start", "@class.outer", "Next class start")
        map_move("[f", "goto_previous_start", "@function.outer", "Previous function start")
        map_move("[[", "goto_previous_start", "@class.outer", "Previous class start")
    end,
}
