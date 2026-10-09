local M = {
    init = false,
}

local status = "Unknown"
local setup = function()
    local api = require("copilot.api")
    api.register_status_notification_handler(function(data)
        -- customize your message however you want
        if data.status == "Normal" then
            status = "Ready"
        elseif data.status == "InProgress" then
            status = "Pending"
        else
            if data.status ~= "" and data.status ~= nil then
                status = data.status -- might never actually be nil but just in case
            else
                status = "Offline"
            end
        end

        if status == "Ready" then
            status = ""
        elseif status == "Pending" then
            status = ""
        elseif status == "Offline" then
            status = ""
        else
            status = "Copilot: Unknown"
        end
    end)
end

M.get_status = function()
    if not M.init then
        setup()
        M.init = true
    end
    return status
end

return M
