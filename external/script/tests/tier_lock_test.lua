local function tierLetter(tier)
	local group, suffix = tostring(tier or 'U'):upper():match('^([UFDCBASXZ])([+%-]*)$')
	if group == nil or #suffix > 3 or ((group == 'U' or group == 'Z') and suffix ~= '') then return nil end
	return group
end
local function sameTierLock(a, b, override)
	return override == true or (tierLetter(a) ~= nil and tierLetter(a) == tierLetter(b))
end
assert(not sameTierLock('C', 'B'), 'C vs B must be rejected')
assert(sameTierLock('C-', 'C+'), 'C- vs C+ must be allowed')
assert(sameTierLock('C', 'B', true), 'override must allow C vs B')
print('tier_lock_test: PASS')
