function classify_by_message(tag, ts, record)
  if record["severity"] ~= nil then return 1, ts, record end
  local msg = record["message"] or record["log"] or record["content"]
  if msg ~= nil then
    local up = string.upper(msg)
    if string.find(up, "PAYMENT DECLINED", 1, true) then
      record["severity"] = "ERROR"
      record["loglevel"] = record["loglevel"] or "ERROR"
    elseif string.find(up, "OUT OF STOCK", 1, true) then
      record["severity"] = "WARN"
      record["loglevel"] = record["loglevel"] or "WARN"
    elseif string.find(up, "PROCESSED ORDER", 1, true) then
      record["severity"] = "INFO"
      record["loglevel"] = record["loglevel"] or "INFO"
    end
  end
  return 1, ts, record
end
