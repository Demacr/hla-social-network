dofile("./names.lua")

-- Search by firstname and lastname
request = function()
    wrk.path = "/api/account/search?firstName=" .. firstNames[math.random(#firstNames)] .. "&lastName=" .. lastNames[math.random(#lastNames)]
    wrk.headers["Authorization"] = os.getenv("WRK_AUTH")
    return wrk.format("GET", wrk.path, wrk.headers)
end

-- Describe non-2XX responses
-- response = function(status, headers, body)
--     if status >= 300 then
--         io.write("Status: ".. status .."\n")
--         io.write("Body:\n")
--         io.write(body .. "\n")
--     end
-- end
