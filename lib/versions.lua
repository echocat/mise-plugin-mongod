local http = require("http")
local json = require("json")
local host = require("host")
local file = require("file")

local versions = {}

local cache_ttl = 12 * 60 * 60 -- 12 hours

function fetch()
    local target = host.target()
    local arch = host.arch()

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
        local version_name = version.version

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
            result[version_name] = {
                note = note,
                release_notes = version.notes,
                edition = downloadTarget.edition,
                url = downloadTarget.archive.url,
                sha1 = downloadTarget.archive.sha1,
                sha256 = downloadTarget.archive.sha256,
            }
        end
    end

    return result
end

function cache_file_name()
    local cache_dir = host.cache_dir()
    local target = host.target()
    local arch = host.arch()

    return file.join_path(cache_dir, "versions-" .. target .. "-" .. arch .. ".json")
end

function get_all()
    local now = os.time()
    local cache_fn = cache_file_name()

    local cache
    local rfOk, cache_json = pcall(file.read, cache_fn)
    if rfOk then
        local djOk, cached = pcall(json.decode, cache_json)
        if djOk and cached.created and (now - cached.created) < cache_ttl then
            cache = cached
        end
    end

    if not cache then
        local v = fetch()
        cache = {
            created = now,
            versions = v,
        }

        local f, err = io.open(cache_fn, "w")
        if not f then
            error("Cannot open " .. cache_fn .. " for storing the cache inside: " .. tostring(err))
        end
        f:write(json.encode(cache))
        f:close()
    end

    return cache.versions
end

function versions.get_all()
    local result = {}

    for key, value in pairs(get_all()) do
        value["version"] = key
        table.insert(result, value)
    end

    table.sort(result, function(a, b)
        return a.version > b.version
    end)

    return result
end

function versions.get(version)
    local all = get_all()
    local result = all[version]
    if not result then
        error(string.format("Version %s does not exist for %s/%s", version, host.target(), host.arch()))
    end
    result["version"] = version
    return result
end

return versions
