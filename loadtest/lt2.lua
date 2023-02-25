dofile("./names.lua")

-- Search by firstname and lastname
request = function()
    if math.random() >= 0.5 then
        wrk.path = "/api/account/search?firstName=" .. firstNames[math.random(#firstNames)] .. "&lastName=" .. lastNames[math.random(#lastNames)]
    else
        wrk.path = "/api/account/profile/" .. math.random(tonumber(os.getenv("WRK_MAX_ID")))
    end
    wrk.headers["Authorization"] = os.getenv("WRK_AUTH")
    return wrk.format("GET", wrk.path, wrk.headers)
end
