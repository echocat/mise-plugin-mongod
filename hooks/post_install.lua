function PLUGIN:PostInstall(ctx)
    local host = require("host")

    local sdkInfo = ctx.sdkInfo[PLUGIN.name]

    apply_ugly_workaround_if_required(sdkInfo)

    local exe = host.path_join(sdkInfo.path, "bin", host.with_exec_ext("mongod"))
    local success, err = pcall(host.exec, exe .. " -version")
    if not success then
        error(string.format("%s installation appears to be broken: got error while testing %s: %s", PLUGIN.name, exe, tostring(err)))
    end
end

-- --------------------------------------------------------------------------------------------------
-- Ugly workaround
--
-- ...because MISE prior 2025.10.8 cannot extract .tgz archives :-(
function apply_ugly_workaround_if_required(sdkInfo)
    local host = require("host")

    -- If this not not MISE (usually vfox directly), immediately return...
    if not host.is_mise() then
        return
    end

    local versions = require("versions")
    local archiver = require("archiver")

    local version = versions.get(sdkInfo.version)
    local path = sdkInfo.path

    local origArchiveFn = host.path_join(path, version.url:match("([^/\\]+)$"))
    local archiveFn = origArchiveFn:gsub("%.tgz$", ".tar.gz")

    -- This will only work if the source archive does exist and both filenames are different.
    if host.can_read(origArchiveFn) and origArchiveFn ~= archiveFn then
        host.mv(origArchiveFn, archiveFn)

        local adOk = pcall(archiver.decompress, archiveFn, path)
        if adOk then
            if RUNTIME.osType:lower() == "windows" then
                error("The workaround is currently not implemented in Windows.")
            end

            host.rm(archiveFn)
            -- Now move all contents from within up...
            host.exec(([[
find '%s' -mindepth 1 -maxdepth 1 -type d | while read -r dir; do
find "$dir" -mindepth 1 -maxdepth 1 -exec mv -f {} '%s' \;
done 2>/dev/null
]]):format(path:gsub("'", "\\'"), path:gsub("'", "\\'")))
            host.exec(([[
find '%s' -mindepth 1 -maxdepth 1 -type d -empty -delete 2>/dev/null
]]):format(path:gsub("'", "\\'")))
        end
    end
end
