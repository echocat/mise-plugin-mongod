function PLUGIN:Available(ctx)
    local versions = require("versions")
    return versions.get_all()
end
