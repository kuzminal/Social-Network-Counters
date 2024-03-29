box.cfg { listen = 3301 }

s = box.schema.space.create('message_counts', { if_not_exists = true })
s:format({
    { name = 'user_id', type = 'string' },
    { name = 'count', type = 'unsigned' }
})
s:create_index('primary', {
    if_not_exists = true,
    type = 'hash',
    parts = { 'user_id' }
})

s = box.schema.space.create('sessions', { if_not_exists = true })
s:format({
    { name = 'id', type = 'string' },
    { name = 'user_id', type = 'string' },
    { name = 'token', type = 'string' },
    { name = 'created_at', type = 'unsigned' }
})
s:create_index('primary', {
    if_not_exists = true,
    type = 'hash',
    parts = { 'id' }
})
s:create_index('ttl', {
    if_not_exists = true,
    parts = {'created_at'},
    unique = false
})
s:create_index('token', {
    if_not_exists = true,
    type = 'hash',
    parts = { 'token' }
})
s:create_index('user', {
    if_not_exists = true,
    type = 'hash',
    parts = { 'user_id' }
})

local fiber = require('fiber')

local TTL = 300000000 --5 минут на сессию
local YIELD_EVERY = 10


fiber.create(function()
    while true do
        fiber.sleep(YIELD_EVERY)
        box.space.sessions.index.ttl:pairs(fiber.time64() - TTL, {iterator = 'LE'})
           :each(function (x)
            box.space.sessions:delete(x.id)
        end)
    end
end)

function get_session_by_user_id(token)
    user_res = box.space.sessions.index.token:get(token)
    if user_res ~= nil then
        return user_res
    end
    return ""
end

function create_session(id, userId, token)
    tokenDb = box.space.sessions.index.user:get(userId)
    if tokenDb ~= nil then
        return tokenDb
    end
    res = box.space.sessions:insert { id, userId, token, fiber.time64() }
    if res ~= nil then
        return res
    end
    return ""
end

function incr_message(userId)
    counter = box.space.message_counts:get(userId)
    if counter ~= nil then
        counter = box.space.message_counts:update( userId, {{'+', 2, 1 }} )
        return counter[2]
    end
    counter = box.space.message_counts:insert { userId, 1 }
    return counter[2]
end

function decr_message(userId)
    counter = box.space.message_counts:get(userId)
    if counter ~= nil and counter[2] > 0 then
        counter = box.space.message_counts:update( userId, {{'-', 2, 1 }} )
        return counter[2]
    end
    return 0
end

function get_message_count(userId)
    counter = box.space.message_counts:get(userId)
    if counter ~= nil then
        return counter[2]
    end
    return 0
end

require 'console'.start()