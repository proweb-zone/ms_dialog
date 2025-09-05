#!lua name=dialog_lib

local function get_messages(keys, args)
   -- get_all_messages.lua
  local user_id = keys[1]
  local friend_id = keys[2]

  local dialog_key = user_id .. ':' .. friend_id .. ':dialog'

  -- Проверяем существование ключа
  if redis.call('EXISTS', dialog_key) == 0 then
    return cjson.encode({error = "Диалог не найден", code = 404})
  end

 -- Проверяем тип ключа
    local key_type = redis.call('TYPE', dialog_key).ok
    if key_type ~= 'zset' then
        return cjson.encode({
            error = "Ожидался zset, но получен: " .. key_type,
            code = 400
        })
    end

    -- Получаем все элементы из zset с scores
    -- ZRANGE key 0 -1 WITHSCORES
    local messages_with_scores = redis.call('ZRANGE', dialog_key, 0, -1, 'WITHSCORES')

    -- Обрабатываем результаты
    local result_messages = {}

    for i = 1, #messages_with_scores, 2 do
        local message_json = messages_with_scores[i]
        local score = messages_with_scores[i + 1]

        local success, message_data = pcall(cjson.decode, message_json)

        if success then
            table.insert(result_messages, {
                score = tonumber(score),
                timestamp = tonumber(score), -- обычно score = timestamp
                msg = tostring(message_data),
                raw_json = message_json
            })
        else
            table.insert(result_messages, {
                score = tonumber(score),
                raw_data = message_json,
                error = "Invalid JSON"
            })
        end
    end

    -- Сортируем по score (времени)
    table.sort(result_messages, function(a, b)
        return a.score < b.score  -- по возрастанию времени
    end)

    -- local result = {
    --     success = true,
    --     total_count = #result_messages,
    --     key_type = key_type,
    --     dialog_key = dialog_key,
    --     messages = result_messages,
    --     timestamp = redis.call('TIME')[1]
    -- }

    return cjson.encode(result_messages)
end

redis.register_function('get_messages', get_messages)
