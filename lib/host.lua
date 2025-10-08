local host = {}

function id_match(what, id, id_like)
    if what == id then
        return true
    end

    if not id_like then
        return false
    end
    return id_like:match("(^" .. what .. "%s)|(%s" .. what .. "%s)|(%s" .. what .. "$)|(^" .. what .. "$)")
end

function host.path_join(...)
    local sep = package.config:sub(1, 1)  -- unter Windows: "\", sonst "/"
    local parts = { ... }
    local path = table.concat(parts, sep)

    path = path:gsub(sep .. "+", sep)

    if sep == "\\" then
        path = path:gsub("/", "\\")
    end

    return path
end

function host.read_file(path)
    local f, err = io.open(path, "r")
    if not f then
        return nil, string.format("Cannot open %q: %s", path, tostring(err))
    end

    local content = f:read("*a")
    f:close()
    return content
end

function host.exec_ext()
    if RUNTIME.osType:lower() == "windows" then
        return ".exe"
    end

    return ""
end

function host.with_exec_ext(fn)
    local ext = host.exec_ext()
    return fn .. ext
end

function host.exec(command)
    local ok, cause, code = os.execute(command)
    if type(ok) == "number" then
        code = ok
        if ok == 0 then
            ok = true
        else
            ok = false
        end
    end

    if ok == true then
        return
    end
    if cause then
        error(string.format("Execution of %q failed: [%s] %s", command, tostring(cause), tostring(code)))
    else
        error(string.format("Execution of %q failed: %s", command, tostring(code)))
    end
end

function host.mkdirs(name)
    if RUNTIME.osType:lower() == "windows" then
        host.exec(string.format([[powershell -NoProfile -Command ^New-Item -ItemType Directory -Force -Path '%s' ^| Out-Null^]], name))
        return name
    end

    host.exec(string.format("mkdir -p '%s'", name:gsub("'", "\\'")))
    return name
end

function host.mv(old, new)
    if RUNTIME.osType:lower() == "windows" then
        host.exec(string.format([[powershell -NoProfile -Command ^Move-Item -Path '%s' -Destination '%s' ^| Out-Null^]], old, new))
        return name
    end

    host.exec(string.format("mv '%s' '%s'", old:gsub("'", "\\'"), new:gsub("'", "\\'")))
    return name
end

function host.rm(name)
    if RUNTIME.osType:lower() == "windows" then
        host.exec(string.format([[powershell -NoProfile -Command ^Remove-Item -LiteralPath -Recurse -Force -ErrorAction SilentlyContinue -LiteralPath '%s' ^| Out-Null^]], name))
        return name
    end

    host.exec(string.format("rm -rf '%s'", name:gsub("'", "\\'")))
    return name
end

function mise_data_dir()
    local explicit = os.getenv("MISE_DATA_DIR")
    if explicit then
        return explicit
    end

    if RUNTIME.osType:lower() == "windows" then
        local lad = os.getenv("LOCALAPPDATA")
        if not lad then
            lad = host.path_join(os.getenv("USERPROFILE"), "AppData", "Local")
        end
        return host.path_join(lad, "mise")
    end

    local share = os.getenv("XDG_DATA_HOME")
    if not share then
        local hd = os.getenv("HOME")
        if not hd then
            hd = "~"
        end
        share = host.path_join(hd, ".local", "share")
    end

    return host.path_join(share, "mise")
end

function mise_download_dir()
    return host.path_join(mise_data_dir(), "downloads")
end

function vfox_cache_dir()
    local explicit = os.getenv("VFOX_CACHE")
    if explicit then
        return explicit
    end

    local home = os.getenv("VFOX_HOME")
    if not home then
        local user_home = os.getenv("HOME")
        if RUNTIME.osType:lower() == "windows" then
            user_home = os.getenv("USERPROFILE")
        end
        if not user_home then
            user_home = "~"
        end
        home = host.path_join(user_home, ".version-fox")
    end

    return host.path_join(home, "cache")
end

function host.cache_dir()
    local base
    local raOk = pcall(require,"archiver")
    if raOk then
        -- We're executed inside MISE... use this one as base. Because there does archiver does exist. In vfox not.
        base = mise_download_dir()
    else
        base = vfox_cache_dir()
    end

    return host.mkdirs(host.path_join(base, "echocat-vfox-mongod"))
end

function host.os()
    local plain = RUNTIME.osType:lower()
    if plain == "windows" then
        return "windows"
    end
    if plain == "darwin" then
        return "macos"
    end
    if plain == "linux" then
        return "linux"
    end

    error("Unsupported operating system: " .. plain)
end

function host.arch()
    local plain = RUNTIME.archType:lower()

    if plain == "x86_64" or plain == "amd64" or plain == "x64" then
        return "x86_64"
    end

    if plain == "arm64" or plain == "aarch64" then
        if RUNTIME.osType:lower() == "linux" then
            return "aarch64"
        end
        return "arm64"
    end

    error("Unsupported architecture: " .. plain)
end

function host.target()
    local osn = host.os()

    if osn == "windows" or osn == "macos" then
        return osn
    end

    if osn == "linux" then
        local osr = host.read_file("/etc/os-release")
        local id = osr:match('^ID="?(.-)"?\n') or osr:match('\nID="?(.-)"?\n')
        local id_like = osr:match('^ID_LIKE="?(.-)"?\n') or osr:match('\nID_LIKE="?(.-)"?\n')
        if not id then
            error("Illegal content of /etc/os-release: cannot find ID entry")
        end

        local version_id = osr:match('^VERSION_ID="?(.-)"?\n') or osr:match('\nVERSION_ID="?(.-)"?\n')
        if not version_id then
            error("Illegal content of /etc/os-release: cannot find VERSION_ID entry")
        end

        if id_match("ubuntu", id, id_like) then
            return "ubuntu" .. version_id:gsub("%.", "")
        end

        if id_match("debian", id, id_like) then
            return "debian" .. version_id:gsub("%.", "")
        end

        if id_match("rhel", id, id_like) then
            return "rhel" .. version_id:gsub("%.", "")
        end

        error("Unsupported linux distribution: " .. id .. "/" .. version_id)
    end

    error("Unsupported operating system: " .. osn)
end

return host
