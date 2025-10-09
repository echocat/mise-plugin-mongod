local types = require("types")

local Semver = {}
Semver.__index = Semver
Semver.__type = "Semver"

function Semver:new(s)
    local result = setmetatable({}, self)

    if not s then
        return nil
    end

    if type(s) ~= "string" then
        error(("requires a string to create a semver from; but got %s"):format(type(s)))
        return nil
    end

    local M, m, p = s:match("^%s*[vV]?(%d+)%.(%d+)%.(%d+)%s*$")
    if not M then
        return nil
    end
    result.major = tonumber(M)
    result.minor = tonumber(m)
    result.patch = tonumber(p)

    return result
end

function Semver:__tostring()
    return ("%d.%d.%d"):format(self.major, self.minor, self.patch)
end

function Semver.cmp(a, b)
    if type(a) == "string" then
        a = Semver:new(a)
    end
    local ob = b
    if type(b) == "string" then
        b = Semver:new(b)
    end

    local aI, bI = types.instanceof(a, Semver), types.instanceof(b, Semver)

    if not aI and not bI then
        return 0
    end
    if not bI then
        return 1
    end
    if not aI then
        return -1
    end

    if not a.major and not b.major then
        return 0
    end
    if not b.major or a.major > b.major then
        return 1
    end
    if not a.major or a.major < b.major then
        return -1
    end

    if not a.minor and not b.minor then
        return 0
    end
    if not b.minor or a.minor > b.minor then
        return 1
    end
    if not a.minor or a.minor < b.minor then
        return -1
    end

    if not a.patch and not b.patch then
        return 0
    end
    if not b.patch or a.patch > b.patch then
        return 1
    end
    if not a.patch or a.patch < b.patch then
        return -1
    end

    return 0
end

return Semver
