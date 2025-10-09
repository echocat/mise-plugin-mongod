local http = require("http")
local json = require("json")
local Semver = require("Semver")
local host = require("host")
local Target = require("Target")
local Version = require("Version")

local versions = {}

local cache_ttl = 24 * 60 * 60 -- 12 hours

function versions.__fetch()
    local target = Target.host()
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

    local latest
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

        local candidate
        for _, download in ipairs(version.downloads) do
            if download.arch == arch then
                if download.edition == "base" or download.edition == "targeted" then
                    if download.archive and download.archive.url then
                        local _, download_target = pcall(Target.new, Target, download.target)
                        if target:equals(download_target) then
                            candidate = {
                                target_version = target.version,
                                download = download
                            }
                        elseif target:equals_base(download_target) and download_target.version and Version.cmp(target.version, download_target.version) > 0 then
                            if not candidate or Version.cmp(download_target.version, candidate.version) > 0 then
                                candidate = {
                                    target_version = download_target.version,
                                    download = download
                                }
                            end
                        end
                    end
                end
            end
        end

        if candidate then
            local sv = Semver:new(version.version)
            if version.production_release == true and sv and not latest or Semver.cmp(sv, latest) > 0 then
                latest = sv
            end

            result[version.version] = {
                note = note,
                release_notes = version.notes,
                edition = candidate.download.edition,
                url = candidate.download.archive.url,
                sha1 = candidate.download.archive.sha1,
                sha256 = candidate.download.archive.sha256,
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
    return host.path_join(host.cache_dir(), "versions-" .. Target.host_string() .. "-" .. host.arch() .. ".json")
end

function versions.__get_all()
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
        local vs, latest = versions.__fetch()
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

    local all, _ = versions.__get_all()

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
    local all, latest = versions.__get_all()
    local resolved_version = version

    if version == "latest" or version == "current" then
        if not latest then
            error("Currently there is no information about the latest version available. You need to explicitly point to a version.")
        end
        resolved_version = latest
    end

    local result = all[resolved_version]
    if not result then
        error(string.format("Version %s does not exist for %s/%s", resolved_version, Target.host_string(), host.arch()))
    end
    result["version"] = resolved_version
    return result
end

return versions
