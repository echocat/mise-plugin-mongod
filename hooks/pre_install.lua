function PLUGIN:PreInstall(ctx)
    local versions = require("versions")
    local host = require("host")
    local Target = require("Target")

    local version = versions.get(ctx.version)

    return {
        version = version.version,
        url = version.url,
        sha256 = version.sha256,
        note = string.format("Downloading %s/%s@%s ", Target.host_string(), host.arch(), version.version),
    }
end
