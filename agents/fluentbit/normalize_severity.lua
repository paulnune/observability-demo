function normalize_severity(tag, ts, record)
  local raw = record["severity"] or record["level"] or record["loglevel"]
  if raw ~= nil then
    local s = string.upper(tostring(raw))
    local map = {
      ["TRACE"]="TRACE",
      ["DEBUG"]="DEBUG",
      ["INFO"]="INFO", ["INFORMATION"]="INFO",
      ["WARN"]="WARN", ["WARNING"]="WARN",
      ["ERROR"]="ERROR", ["ERR"]="ERROR",
      ["FATAL"]="FATAL", ["CRITICAL"]="FATAL"
    }
    local norm = map[s] or "INFO"
    record["severity"] = norm
    if record["loglevel"] == nil then record["loglevel"] = norm end
  end
  return 1, ts, record
end
