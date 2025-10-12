function PLUGIN:PostInstall(ctx)
    local raOk, archiver = pcall(require,"vfox.archiver")
    if not raOk then
        archiver = require("archiver")
    end

    local host = require("host")
    local versions = require("versions")

    local sdkInfo = ctx.sdkInfo[PLUGIN.name]
    local version = versions.get(sdkInfo.version)
    local path = sdkInfo.path

    -- --------------------------------------------------------------------------------------------------
    -- BEGIN: Ugly workaround
    --
    -- ...because not all archives are extracted automatically by MISE. Especially not: .tgz :-(

    local origArchiveFn = host.path_join(path, version.url:match("([^/\\]+)$"))
    local archiveFn = origArchiveFn:gsub("%.tgz$", ".tar.gz")

    -- This will only work within MISE and if both file names are different.
    if host.is_mise() and origArchiveFn ~= archiveFn then
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

    -- END: Ugly workaround
    -- --------------------------------------------------------------------------------------------------

    local exe = host.path_join(path, "bin", host.with_exec_ext("mongod"))
    local success, err = pcall(host.exec, exe .. " -version")
    if not success then
        error(string.format("%s installation appears to be broken: got error while testing %s: %s", PLUGIN.name, exe, tostring(err)))
    end
end
