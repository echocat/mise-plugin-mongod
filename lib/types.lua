local types = {}

function types.instanceof(obj, class)
    local mt = getmetatable(obj)
    if mt == class then
        return true
    end
    if mt and class and mt.__type == class.__type then
        return true
    end
    return false
end

return types