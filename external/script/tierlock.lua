local tierlock = {}

function tierlock.letter(tier)
	if tier == nil then return nil end
	local group, suffix = tostring(tier):upper():match('^([UFDCBASXZ])([+%-]*)$')
	if group == nil or #suffix > 3 or ((group == 'U' or group == 'Z') and suffix ~= '') then return nil end
	return group
end

function tierlock.same(a, b, override)
	if override == true then return true end
	local left, right = tierlock.letter(a), tierlock.letter(b)
	return left ~= nil and left == right, left, right
end

return tierlock
