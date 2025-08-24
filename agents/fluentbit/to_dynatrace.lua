function to_dynatrace(tag, ts, record)

  local msg = record["message"] or record["log"] or record["content"]

  if msg == nil then
    local sev = record["severity"] or record["loglevel"] or "INFO"
    local svc = record["service.name"] or record["service"] or "app"
    msg = string.format("[%s] %s event", sev, svc)
  end

  record["content"] = msg

  return 1, ts, record
end
