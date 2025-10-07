local http = require("http")
local json = require("json")
local host = require("host")

local versions = {}

local cache_ttl = 24 * 60 * 60 -- 12 hours

function fetch()
    local target = host.target()
    local arch = host.arch()
    local latest

    local resp, err = http.get({
        url = "https://downloads.mongodb.org/full.json",
    })

    if err ~= nil then
        error("Failed to fetch versions: " .. err)
    end
    if resp.status_code ~= 200 then
        error("GitHub API returned status " .. resp.status_code .. ": " .. resp.body)
    end

    local body = json.decode(resp.body)
    local result = {}

    for _, version in ipairs(body.versions) do
        local note
        if version.release_candidate == true then
            note = "pre-release"
        end
        if version.lts_release == true then
            note = "lts"
        end
        if version.current == true then
            note = "latest"
        end

        local downloadTarget
        for _, download in ipairs(version.downloads) do
            if download.target == target and download.arch == arch then
                if download.edition == "base" or download.edition == "targeted" then
                    if download.archive and download.archive.url then
                        downloadTarget = download
                    end
                end
            end
        end

        if downloadTarget then
            local sv = versions.parse_semver(version.version)
            if version.production_release == true and sv and not latest or versions.cmp_semver(sv, latest) > 0 then
                latest = sv
            end

            result[version.version] = {
                note = note,
                release_notes = version.notes,
                edition = downloadTarget.edition,
                url = downloadTarget.archive.url,
                sha1 = downloadTarget.archive.sha1,
                sha256 = downloadTarget.archive.sha256,
            }
        end
    end

    local latestStr
    if latest then
        latestStr = latest.string
    end
    return result, latestStr
end

function cache_file_name()
    local cache_dir = host.cache_dir()
    local target = host.target()
    local arch = host.arch()

    return host.path_join(cache_dir, "versions-" .. target .. "-" .. arch .. ".json")
end

function get_all()
    local now = os.time()
    local cache_fn = cache_file_name()

    local cache
    local cache_json, _ = host.read_file(cache_fn)
    if cache_json then
        local djOk, cached = pcall(json.decode, cache_json)
        if djOk and cached.created and (now - cached.created) < cache_ttl then
            cache = cached
        end
    end

    if not cache then
        local vs, latest = fetch()
        cache = {
            created = now,
            latest = latest,
            versions = vs,
        }

        local f, err = io.open(cache_fn, "w")
        if not f then
            error("Cannot open " .. cache_fn .. " for storing the cache inside: " .. tostring(err))
        end
        f:write(json.encode(cache))
        f:close()
    end

    return cache.versions, cache.latest
end

function versions.get_all()
    local result = {}

    local all, _ = get_all()

    for key, value in pairs(all) do
        value["version"] = key
        table.insert(result, value)
    end

    table.sort(result, function(a, b)
        return a.version > b.version
    end)

    return result
end

function versions.get(version)
    local all, latest = get_all()
    local target = version

    if version == "latest" or version == "current" then
        if not latest then
            error("Currently there is no information about the latest version available. You need to explicitly point to a version.")
        end
        target = latest
    end

    local result = all[target]
    if not result then
        error(string.format("Version %s does not exist for %s/%s", target, host.target(), host.arch()))
    end
    result["version"] = target
    return result
end

function versions.parse_semver(s)
    if type(s) ~= "string" then
        return nil
    end
    local M, m, p = s:match("^%s*[vV]?(%d+)%.(%d+)%.(%d+)%s*$")
    if not M then
        return nil
    end
    return {
        major = tonumber(M),
        minor = tonumber(m),
        patch = tonumber(p),
        string = s,
    }
end

function versions.cmp_semver(a, b)
    if type(a) == "string" then
        a = versions.parse_semver(a)
    end
    if type(b) == "string" then
        b = versions.parse_semver(b)
    end

    if not a and not b then
        return 0
    end
    if not b then
        return 0
    end
    if not a then
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

return versions
