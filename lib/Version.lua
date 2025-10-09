local types = require("types")

local Version = {}
Version.__index = Version
Version.__type = "Version"

function Version:new(s)
    local result = setmetatable({}, self)

    if not s then
        return nil
    end

    if type(s) == "table" then
        if #s == 0 then
            return nil
        end
        for i = 1, #s do
            result[i] = s[i]
        end
        return result
    end

    if type(s) ~= "string" then
        error(("requires a string or table(array) to create a version from; but got %s"):format(type(s)))
        return nil
    end

    local i = 1

    for part in string.gmatch(s, "[^%.]+") do
        local num = tonumber(part)
        if not num then
            return nil
        end
        result[i] = num
        i = i + 1
    end

    if #result == 0 then
        return nil
    end

    return result
end

function Version:__tostring()
    return self:concat(".")
end

function Version:concat(sep)
    return table.concat(self, sep)
end

function Version.cmp(a, b)
    if type(a) == "string" then
        a = Version:new(a)
    end
    if type(b) == "string" then
        b = Version:new(b)
    end

    local aI, bI = types.instanceof(a, Version), types.instanceof(b, Version)

    if not aI and not bI then
        return 0
    end
    if not bI then
        return 1
    end
    if not aI then
        return -1
    end

    local max_len = math.max(#a, #b)

    for i = 1, max_len do
        local a_v = a[i] or 0
        local b_v = b[i] or 0
        if a_v > b_v then
            return 1
        elseif a_v < b_v then
            return -1
        end
    end

    return 0
end

return Version
