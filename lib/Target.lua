local Version = require("Version")
local types = require("types")
local host = require("host")
local cos = require("os")

local Target = {}
Target.__index = Target
Target.__type = "Target"

function id_match(what, id, id_like)
    if what == id then
        return true
    end

    if not id_like then
        return false
    end
    return id_like:match("(^" .. what .. "%s)|(%s" .. what .. "%s)|(%s" .. what .. "$)|(^" .. what .. "$)")
end

function Target.__parse_version_year_month(s)
    if type(s) ~= "string" then
        return nil
    end
    if s:len() ~= 4 then
        return nil
    end

    local year = tonumber(s:sub(1, 2))
    local month = tonumber(s:sub(3, 4))

    if year == nil or month == nil then
        return nil
    end

    return Version:new({ year, month })
end

function Target.__format_version_year_month(s)
    if not s or #s == 0 then
        return ""
    end

    if #s ~= 2 or s[1] > 99 or s[1] < 0 or s[2] > 99 or s[2] < 0 then
        error(("Should format a version year month; but got: %s"):format(tostring(s)))
    end

    return ("%02d%02d"):format(s[1], s[2])
end

function Target.__parse_version_major_only(s)
    if type(s) ~= "string" then
        return nil
    end
    local major = tonumber(s)

    if major == nil then
        return nil
    end

    return Version:new({ major })
end

function Target.__format_version_major_only(s)
    if not s or #s == 0 then
        return ""
    end

    if #s ~= 1 or s[1] < 0 then
        error(("Should format a major only version; but got: %s"):format(tostring(s)))
    end

    return ("%d").format(s[1])
end

function Target.__parse_version_rhel(s)
    if type(s) ~= "string" then
        return nil
    end

    if s:len() < 1 or s:len() > 2 then
        return nil
    end

    local major = tonumber(s:sub(1, 1))
    if major == nil then
        return nil
    end

    local minor
    if s:len() > 1 then
        minor = tonumber(s:sub(2, 2))
        if minor == nil then
            return nil
        end
    end

    if minor ~= nil then
        return Version:new({ major, minor })
    end

    return Version:new({ major })
end

function Target.__format_version_rhel(s)
    if not s or #s == 0 then
        return ""
    end

    if #s < 1 or #s > 2 or (s[1] < 0 or s[1] > 9) or (#s > 1 and (s[2] < 0)) then
        error(("Should format a rhel version; but got: %s"):format(tostring(s)))
    end

    if #s > 1 and s[2] > 0 then
        return ("%d%d"):format(s[1], s[2])
    end

    return ("%d"):format(s[1])
end

Target.__distributions = {
    ubuntu = {
        parse_version = Target.__parse_version_year_month,
        format_version = Target.__format_version_year_month,
    },
    debian = {
        parse_version = Target.__parse_version_major_only,
        format_version = Target.__format_version_major_only,
    },
    amazon = {
        aliases = { "amzn" },
        parse_version = Target.__parse_version_major_only,
        format_version = Target.__format_version_major_only,
    },
    suse = {
        parse_version = Target.__parse_version_major_only,
        format_version = Target.__format_version_major_only,
    },
    rhel = {
        parse_version = Target.__parse_version_rhel,
        format_version = Target.__format_version_rhel,
    },
}

function Target:new(s)
    local result = setmetatable({}, self)

    if type(s) == "table" then
        if type(s.os) ~= "string" then
            error(("When creating Target with table field os - it needs to be a string; but got: %s - %s"):format(type(s.os), tostring(s.os)))
        end
        result.os = s.os

        if s.distribution then
            if type(s.distribution) ~= "string" then
                error(("When creating Target with table field distribution - it needs to be a distribution; but got: %s - %s"):format(type(s.distribution), tostring(s.distribution)))
            end
            result.distribution = s.distribution
        end

        if s.version then
            if not types.instanceof(s.version, Version) then
                error(("When creating Target with table field version - it needs to be a Version; but got: %s - %s"):format(type(s.version), tostring(s.version)))
            end
            result.version = s.version
        end

        return result
    end

    if type(s) ~= "string" then
        error(("When creating Target with an argument only strings or tables are allowed; but got: %s - %s"):format(type(s), tostring(s)))
    end

    if s == "windows" or s == "macos" then
        result.os = s
        return result
    end

    for prefix, settings in pairs(Target.__distributions) do
        if s:sub(1, #prefix) == prefix then
            result.os = "linux"
            result.distribution = prefix
            result.version = settings.parse_version(s:sub(#prefix + 1))

            if result.version == nil then
                error(("Version of target %s cannot be interpreted."):format(s))
            end

            return result
        end
    end

    error(("Unknown target %s."):format(s))
end

function Target:__tostring()
    if self.distribution then
        local d = Target.__distributions[self.distribution]
        if d then
            return self.distribution .. d.format_version(self.version)
        end

        -- Fallback in strange cases...
        return self.distribution .. tostring(self.version)
    end

    return self.os
end

function Target:equals_base(other)
    if not types.instanceof(other, Target) then
        return false
    end

    return self.os == other.os and self.distribution == other.distribution
end

function Target:equals(other)
    if not self:equals_base(other) then
        return false
    end

    return Version.cmp(self.version, other.version) == 0
end

function Target.host(os, os_release_fn)
    local overwriteEnv = cos.getenv("MONGOD_TARGET")
    if overwriteEnv then
        return Target:new(overwriteEnv)
    end

    if not os or os == "" then
        os = host.os()
    end

    if os == "windows" or os == "macos" then
        return Target:new({
            os = os,
        })
    end

    if os == "linux" then
        -- The following is only used in tests...
        local distributionOverwrite, versionOverwrite = RUNTIME.distributionType, RUNTIME.distributionVersion
        if distributionOverwrite and versionOverwrite then
            local version = Version:new(versionOverwrite)
            return Target:new({
                os = "linux",
                distribution = distributionOverwrite,
                version = version,
            })
        end

        if not os_release_fn or os_release_fn == "" then
            os_release_fn = "/etc/os-release"
        end

        local osr = host.read_file(os_release_fn)
        if osr then
            local id = osr:match('^ID="?(.-)"?\n') or osr:match('\nID="?(.-)"?\n')
            local id_like = osr:match('^ID_LIKE="?(.-)"?\n') or osr:match('\nID_LIKE="?(.-)"?\n')
            if not id then
                error("Illegal content of /etc/os-release: cannot find ID entry")
            end

            local version_id = osr:match('^VERSION_ID="?(.-)"?\n') or osr:match('\nVERSION_ID="?(.-)"?\n')
            if not version_id then
                error("Illegal content of /etc/os-release: cannot find VERSION_ID entry")
            end

            for distribution, settings in pairs(Target.__distributions) do
                local match = id_match(distribution, id, id_like)

                if not match and type(settings.aliases) == "table" then
                    for _, alias in ipairs(settings.aliases) do
                        if id_match(alias, id, id_like) then
                            match = true
                        end
                    end
                end

                if match then
                    local version = Version:new(version_id)

                    if version == nil then
                        error(("Don't know how to interpret VERSION_ID (%s) of current installation of distribution %s."):format(s, prefix))
                    end

                    return Target:new({
                        os = "linux",
                        distribution = distribution,
                        version = version,
                    })
                end
            end
            error("Unsupported linux distribution: " .. id .. "/" .. version_id)
        end
        error(("Unsupported linux distribution: %q does not exist"):format(os_release_fn))

    end
    error("Unsupported operating system: " .. os)
end

function Target.host_string()
    return tostring(Target.host())
end

return Target
