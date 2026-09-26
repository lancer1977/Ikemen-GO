local start = {}

--team side specific data storage
start.p = {{}, {}}
--cell data storage
start.c = {}
for i = 1, --[[gameOption('Config.Players')]]8 do
	table.insert(start.c, {selX = 0, selY = 0, cell = -1, randCnt = 0, randRef = nil})
end
--globally accessible temp data
start.challenger = 0
--overridable via -statsfile CLI flag (main.lua), default keeps normal shared stats file
start.statsFilePath = start.statsFilePath or 'save/stats.json'
--local variables
local restoreCursor = false
local selScreenEnd = false
local stageEnd = false
local stageRandom = false
local stageListNo = 0
local t_aiRamp = {}
local t_reservedChars = {{}, {}}
local t_portraitPriority = {1, 1}
local timerSelect = 0
local cursorActive = {}
local cursorDone = {}

--;===========================================================
--; COMMON FUNCTIONS
--;===========================================================
-- Returns the unified order table used by the active mode.
-- Priority: order<gameMode>, order<baseMode>, order
function start.f_getOrderChars(baseMode)
	local mode = gameMode() or ''
	local checked = {}
	for _, key in ipairs({mode, baseMode, 'default'}) do
		if key ~= nil and key ~= '' and not checked[key] then
			checked[key] = true
			if type(main.t_orderChars[key]) == 'table' then
				return main.t_orderChars[key]
			end
		end
	end
	return {}
end

--; ROSTER
--converts '.maxmatches' style table (key = order, value = max matches) to key = match number, value = subtable with char num and order data
function start.f_unifySettings(t, t_chars)
	local ret = {}
	for i = 1, #t do --for each order number
		if t_chars[i] ~= nil then --only if there are any characters available with this order
			local infinite = false
			local num = t[i]
			if num == -1 then --infinite matches
				num = #t_chars[i] --assign max amount of characters with this order
				infinite = true
			end
			for j = 1, num do --iterate up to max amount of matches versus characters with this order
				table.insert(ret, {['rmin'] = start.p[2].numChars, ['rmax'] = start.p[2].numChars, ['order'] = i})
			end
			if infinite then
				table.insert(ret, {['rmin'] = start.p[2].numChars, ['rmax'] = start.p[2].numChars, ['order'] = -1})
				break --no point in appending additional matches
			end
		end
	end
	return ret
end

-- start.t_makeRoster is a table storing functions returning table data used
-- by start.f_makeRoster function, depending on game mode. Can be appended via
-- external module, without conflicting with default scripts.
start.t_makeRoster = {}
function start.f_rosterMaxMatches(baseMode)
	local charData = start.f_getCharData(start.p[1].t_selected[1].ref)
	local mode = gameMode() or ''
	local t_chars = start.f_getOrderChars(baseMode)
	local fallback = start.p[2].teamMode == 0 and 'arcade' or 'team'
	local useMode = mode ~= '' and mode ~= 'arcade' and mode ~= 'team'
	local candidates = {}
	local function add(prefix, name)
		if name == nil or name == '' then
			return
		end
		if prefix ~= nil and prefix ~= '' then
			table.insert(candidates, prefix .. '_' .. name .. 'maxmatches')
		else
			table.insert(candidates, name .. 'maxmatches')
		end
	end
	if charData.maxmatches ~= nil then
		if useMode then
			add(charData.maxmatches, mode)
		end
		if baseMode ~= nil and baseMode ~= mode then
			add(charData.maxmatches, baseMode)
		end
		add(charData.maxmatches, fallback)
	end
	if useMode then
		add(nil, mode)
	end
	if baseMode ~= nil and baseMode ~= mode then
		add(nil, baseMode)
	end
	add(nil, fallback)
	for _, key in ipairs(candidates) do
		if main.t_selOptions[key] ~= nil then
			return start.f_unifySettings(main.t_selOptions[key], t_chars), t_chars
		end
	end
	return start.f_unifySettings(main.t_selOptions[fallback .. 'maxmatches'], t_chars), t_chars
end

start.t_makeRoster.arcade = function()
	return start.f_rosterMaxMatches('arcade')
end
start.t_makeRoster.teamcoop = start.t_makeRoster.arcade
start.t_makeRoster.netplayteamcoop = start.t_makeRoster.arcade
start.t_makeRoster.timeattack = start.t_makeRoster.arcade

start.t_makeRoster.survival = function()
	return start.f_rosterMaxMatches('survival')
end
start.t_makeRoster.survivalcoop = start.t_makeRoster.survival
start.t_makeRoster.netplaysurvivalcoop = start.t_makeRoster.survival

-- generates roster table
function start.f_makeRoster(t_ret)
	t_ret = t_ret or {}
	--prepare correct settings tables
	if start.t_makeRoster[gameMode()] == nil then
		panicError("\n" .. gameMode() .. " game mode unrecognized by start.f_makeRoster()\n")
	end
	local t, t_static = start.t_makeRoster[gameMode()]()
	--generate roster
	local t_removable = main.f_tableCopy(t_static) --copy into editable order table
	for i = 1, #t do --for each match number
		if t[i].order == -1 then --infinite matches for this order detected
			table.insert(t_ret, {-1}) --append infinite matches flag at the end
			break
		end
		if t_removable[t[i].order] ~= nil then
			if #t_removable[t[i].order] == 0 and main.forceRosterSize then
				t_removable = main.f_tableCopy(t_static) --allows character repetition, if needed to fill whole roster
			end
			if #t_removable[t[i].order] >= 1 then --there is at least 1 character with this order available
				local remaining = t[i].rmin - #t_removable[t[i].order]
				table.insert(t_ret, {}) --append roster table with new subtable
				local t_toinsert = {}
				local t_removableTemp = main.f_tableCopy(t_removable)
				for j = 1, math.random(math.min(t[i].rmin, #t_removableTemp[t[i].order]), math.min(t[i].rmax, #t_removableTemp[t[i].order])) do --for randomized characters count
					local rand = math.random(1, #t_removableTemp[t[i].order]) --randomize which character will be taken
					local ref = t_removableTemp[t[i].order][rand]
					if not main.charparam.single or not start.f_getCharData(ref).single then
						table.insert(t_toinsert, ref) --append character if 'single' param is not blocking larger team size
						table.remove(t_removableTemp[t[i].order], rand) --remove it from the t_removableTemp table
					else --otherwise only this character is added to roster
						t_toinsert = {ref}
						remaining = 0
						break
					end
				end
				for _, v in ipairs(t_toinsert) do
					table.insert(t_ret[#t_ret], v) --add such character into roster subtable
					main.f_tableRemove(t_removable[t[i].order], v) --and remove it from the available character pool
				end
				--fill the remaining slots randomly if there are not enough players available with this order
				while remaining > 0 do
					table.insert(t_ret[#t_ret], t_static[t[i].order][math.random(1, #t_static[t[i].order])])
					remaining = remaining - 1
				end
			end
		end
	end
	if gameOption('Debug.DumpLuaTables') then main.f_printTable(t_ret, 'debug/t_roster.txt') end
	return t_ret
end

--;===========================================================
--; AI RAMPING
-- start.t_aiRampData is a table storing functions returning variable data used
-- by start.f_aiRamp function, depending on game mode. Can be appended via
-- external module, without conflicting with default scripts.
start.t_aiRampData = {}
start.t_aiRampData.arcade = function()
	if start.p[2].teamMode == 0 then --Single
		return gameOption('Arcade.arcade.AIramp.start')[1], gameOption('Arcade.arcade.AIramp.start')[2], gameOption('Arcade.arcade.AIramp.end')[1], gameOption('Arcade.arcade.AIramp.end')[2]
	else --Simul / Turns / Tag
		return gameOption('Arcade.team.AIramp.start')[1], gameOption('Arcade.team.AIramp.start')[2], gameOption('Arcade.team.AIramp.end')[1], gameOption('Arcade.team.AIramp.end')[2]
	end
end
start.t_aiRampData.teamcoop = start.t_aiRampData.arcade
start.t_aiRampData.netplayteamcoop = start.t_aiRampData.arcade
start.t_aiRampData.timeattack = start.t_aiRampData.arcade
start.t_aiRampData.survival = function()
	return gameOption('Arcade.survival.AIramp.start')[1], gameOption('Arcade.survival.AIramp.start')[2], gameOption('Arcade.survival.AIramp.end')[1], gameOption('Arcade.survival.AIramp.end')[2]
end
start.t_aiRampData.survivalcoop = start.t_aiRampData.survival
start.t_aiRampData.netplaysurvivalcoop = start.t_aiRampData.survival

-- generates AI ramping table
function start.f_aiRamp(currentMatch)
	if start.t_aiRampData[gameMode()] == nil then
		panicError("\n" .. gameMode() .. " game mode unrecognized by start.f_aiRamp()\n")
	end
	local start_match, start_diff, end_match, end_diff = start.t_aiRampData[gameMode()]()
	local startAI = gameOption('Options.Difficulty') + start_diff
	if startAI > 8 then
		startAI = 8
	elseif startAI < 1 then
		startAI = 1
	end
	local endAI = gameOption('Options.Difficulty') + end_diff
	if endAI > 8 then
		endAI = 8
	elseif endAI < 1 then
		endAI = 1
	end
	if currentMatch == 1 then
		t_aiRamp = {}
	end
	for i = math.min(#t_aiRamp, currentMatch), math.max(#start.t_roster, currentMatch) do
		if i - 1 <= start_match then
			table.insert(t_aiRamp, startAI)
		elseif i - 1 <= end_match then
			local curMatch = i - (start_match + 1)
			table.insert(t_aiRamp, curMatch * (endAI - startAI) / (end_match - start_match) + startAI)
		else
			table.insert(t_aiRamp, endAI)
		end
	end
	if gameOption('Debug.DumpLuaTables') then main.f_printTable(t_aiRamp, 'debug/t_aiRamp.txt') end
end
--;===========================================================

--calculates AI level
function start.f_difficulty(player, offset)
	local t = {}
	if main.f_playerSide(player) == 1 then
		t = start.f_getCharData(start.p[1].t_selected[math.floor(player / 2 + 0.5)].ref)
	else
		t = start.f_getCharData(start.p[2].t_selected[math.floor(player / 2)].ref)
	end
	if t.ai ~= nil then
		return t.ai
	else
		local ref = main.f_playerSide(player) == 1
			and start.p[1].t_selected[math.floor(player / 2 + 0.5)].ref
			or start.p[2].t_selected[math.floor(player / 2)].ref
		if start.f_getCharFaction ~= nil and tostring(start.f_getCharFaction(ref) or ''):lower() == 'guardians' then
			return 8
		end
		return gameOption('Options.Difficulty') + offset
	end
end

--assigns AI level, remaps input
function start.f_remapAI(ai)
	--Offset
	local offset = 0
	if gameOption('Arcade.AI.Ramping') and main.aiRamp then
		if t_aiRamp[matchNo()] == nil then
			start.f_aiRamp(matchNo())
		end
		offset = t_aiRamp[matchNo()] - gameOption('Options.Difficulty')
	end
	local t_ex = {}
	for side = 1, 2 do
		if main.coop then
			for k, v in ipairs(start.p[side].t_selCmd) do
				if gameMode('versuscoop') then
					remapInput(v.player, v.cmd)
					setCom(v.player, 0)
					t_ex[v.player] = true
				else
					local pn = v.player * 2 - 1
					remapInput(pn, v.cmd)
					setCom(pn, 0)
					t_ex[pn] = true
				end
			end
		end
		if start.p[side].teamMode == 0 or start.p[side].teamMode == 2 then --Single or Turns
			if (not main.cpuSide[side] and not main.coop) or start.challenger > 0 or gameMode('training') then
				setCom(side, 0)
			else
				setCom(side, ai or start.f_difficulty(side, offset))
			end
		elseif start.p[side].teamMode == 1 then --Simul
			if not t_ex[side] then
				if (not main.cpuSide[side] and not main.coop) or start.challenger > 0 then
					setCom(side, 0)
				else
					setCom(side, ai or start.f_difficulty(side, offset))
				end
			end
			for i = side + 2, #start.p[side].t_selected * 2 do
				if not t_ex[i] and (i - 1) % 2 + 1 == side then
					remapInput(i, side) --P3/5/7 => P1 controls, P4/6/8 => P2 controls
					setCom(i, ai or start.f_difficulty(i, offset))
				end
			end
		else --Tag
			for i = side, #start.p[side].t_selected * 2 do
				if not t_ex[i] and (i - 1) % 2 + 1 == side then
					if (not main.cpuSide[side] and not main.coop) or start.challenger > 0 then
						remapInput(i, getRemapInput(side)) --P1/3/5/7 => P1 controls, P2/4/6/8 => P2 controls
						setCom(i, 0)
					else
						setCom(i, ai or start.f_difficulty(i, offset))
					end
				end
			end
		end
	end
end

--sets lifebar elements, round time, rounds to win
function start.f_setRounds(roundTime, t_rounds)
	setMotifElements(main.motif)
	setFightScreenElements(main.fightscreen)
	-- Round time
	local frames = fightScreenVar("time.framespercount")
	local p1FramesMul = 1
	local p2FramesMul = 1
	if start.p[1].teamMode == 3 then -- Tag
		p1FramesMul = start.p[1].numChars
	end
	if start.p[2].teamMode == 3 then -- Tag
		p2FramesMul = start.p[2].numChars
	end
	if (start.p[1].teamMode == 3 or start.p[2].teamMode == 3) and gameOption('Options.Tag.TimeScaling') > 0 then
		-- Calculate the maximum team size
		local maxTeamSize = math.max(p1FramesMul, p2FramesMul)
		-- Apply a base multiplier for team size
		local adjustedFrames = frames * (1 + (maxTeamSize - 1) * gameOption('Options.Tag.TimeScaling'))
		-- Enforce a minimum threshold to avoid overly short rounds
		frames = main.f_round(math.max(adjustedFrames, frames), 0)
	end
	setTimeFramesPerCount(frames)
	if roundTime ~= nil then
		setRoundTime(math.max(-1, roundTime * frames)) --round time predefined
	elseif main.charparam.time and start.f_getCharData(start.p[2].t_selected[1].ref).time ~= nil then --round time assigned as character param
		setRoundTime(math.max(-1, start.f_getCharData(start.p[2].t_selected[1].ref).time * frames))
	else --default round time
		setRoundTime(math.max(-1, main.roundTime * frames))
	end
	--Rounds to win. Determined by enemy team mode
	for side = 1, 2 do
		local enemy = 3 - side
		if start.p[enemy].teamMode == 2 then --Turns mode always uses team size
			setMatchWins(side, start.p[enemy].numChars)
		elseif t_rounds[side] ~= nil then --Use override if it exists
			setMatchWins(side, t_rounds[side])
		elseif enemy == 2 and main.charparam.rounds and start.f_getCharData(start.p[2].t_selected[1].ref).rounds ~= nil then --round num assigned as character param
			setMatchWins(side, start.f_getCharData(start.p[2].t_selected[1].ref).rounds)
		elseif start.p[enemy].teamMode == 1 then --default rounds num (Simul)
			setMatchWins(side, main.matchWins.simul[enemy])
		elseif start.p[enemy].teamMode == 3 then --default rounds num (Tag)
			setMatchWins(side, main.matchWins.tag[enemy])
		else --default rounds num (Single)
			setMatchWins(side, main.matchWins.single[enemy])
		end
		setMatchMaxDrawGames(side, main.matchWins.draw[side])
	end
	--timer / score counter
	local timer, t_score = start.f_prefightHUD()
	setFightScreenTimer(timer)
	setFightScreenScore(t_score[1], t_score[2])
end

local function f_listCharRefs(t)
	local ret = {}
	for i = 1, #t do
		table.insert(ret, start.f_getCharData(t[i].ref).char:lower())
	end
	return ret
end

--;===========================================================
-- Accumulators derived from game stats
--;===========================================================
-- Fold matches up to 'upto' (inclusive)
function start.f_accStats(upto)
	local ret = {
		win = {0,0}, lose = {0,0},
		time = { total = 0, matches = {} },
		score = { total = {0,0}, matches = {} },
		consecutive = {0,0},
	}
	local gameStats = getGameStats()
	local matches = (gameStats and gameStats.Matches) or {}
	local n = math.min(upto or #matches, #matches)
	local streak = {0,0}
	for i = 1, n do
		local m = matches[i] or {}
		ret.time.total = ret.time.total + (m.MatchTime or 0)
		if m.TotalScore then
			ret.score.total[1] = m.TotalScore[1] or ret.score.total[1]
			ret.score.total[2] = m.TotalScore[2] or ret.score.total[2]
		end
		local tr, sr = {}, {}
		if m.Rounds then
			for j, r in ipairs(m.Rounds) do
				tr[j] = r.Timer or 0
				sr[j] = { [1] = (r.Score and r.Score[1]) or 0, [2] = (r.Score and r.Score[2]) or 0 }
			end
		end
		table.insert(ret.time.matches, tr)
		table.insert(ret.score.matches, sr)
		if m.WinSide == 1 then
			ret.win[1] = ret.win[1] + 1; ret.lose[2] = ret.lose[2] + 1
			streak[1] = streak[1] + 1; streak[2] = 0
		elseif m.WinSide == 2 then
			ret.win[2] = ret.win[2] + 1; ret.lose[1] = ret.lose[1] + 1
			streak[2] = streak[2] + 1; streak[1] = 0
		end
		for s = 1, 2 do
			if streak[s] > ret.consecutive[s] then ret.consecutive[s] = streak[s] end
		end
	end
	return ret
end

-- Compute HUD timer/score to show at the beginning of the *next* match
function start.f_prefightHUD()
	-- "Next match index" is current matchNo(); we want totals of already-finished matches.
	local prev = math.max((matchNo() or 1) - 1, 0)
	local acc = start.f_accStats(prev)
	local t_score = {acc.score.total[1], acc.score.total[2]}
	local timer = acc.time.total
	if start.challenger > 0 and gameMode('versus') then
		return 0, {0, 0}
	end
	-- emulate resetScore-on-loss behavior for the next match HUD
	local gameStats = getGameStats()
	local last = (gameStats and gameStats.Matches and gameStats.Matches[prev]) or nil
	if last and main.resetScore and matchNo() ~= -1 then
		if last.WinSide == 2 then
			t_score[1] = acc.lose[1]
		end
	end
	return timer, t_score
end

--;===========================================================

--returns the next stage path from the given pool
function start.stageShuffleBag(id, pool)
	-- safety check: prevent nil or invalid pools
	if not pool or type(pool) ~= 'table' or #pool == 0 then
		return nil
	end

	-- safety check: prevent nil id
	id = id or 'defaultStageBag'
	start.shuffleStages = start.shuffleStages or {}
	start.shuffleStages[id] = start.shuffleStages[id] or {}

	if #start.shuffleStages[id] == 0 then
		local t = {}
		for i = 1, #pool do
			table.insert(t, i)
		end
		start.f_shuffleTable(t)
		-- prevent immediate repetition if the bag was just refilled
		if start.lastStageIdx and #pool > 1 and t[#t] == start.lastStageIdx then
			table.insert(t, 1, table.remove(t)) -- rotate
		end
		start.shuffleStages[id] = t
	end

	local idx = table.remove(start.shuffleStages[id])
	start.lastStageIdx = idx
	-- Pool entries are already resolved stage refs
	return pool[idx]
end

-- Return the next selectable stage in a complete rotation. The random tier
-- ladder must not inherit a fixed/random stage override from launchFight.
function start.stageCycle(id, pool)
	if not pool or type(pool) ~= "table" or #pool == 0 then
		return nil
	end
	id = id or "defaultStageCycle"
	start.stageCycles = start.stageCycles or {}
	local next = (start.stageCycles[id] or 0) + 1
	if next > #pool then
		next = 1
	end
	start.stageCycles[id] = next
	return pool[next]
end

--sets stage
function start.f_setStage(num, assigned)
	if gameMode("randomtierladder") and main.t_selectableStages then
		num = start.stageCycle("randomtierladder", main.t_selectableStages)
		stageListNo = 0
		stageRandom = false
		assigned = true
	elseif main.stageMenu then
		local pool = main.t_selectableStages
		if main.endlessRandomActive then
			num = start.stageShuffleBag('endlessRandomStage', pool)
			stageRandom = true
		elseif stageListNo == 0 then
			num = start.stageShuffleBag('stageMenu', pool)
			stageListNo = num -- comment out to randomize stage after each fight in survival mode, when random stage is chosen
			stageRandom = true
		else
			num = pool[stageListNo]
		end
		assigned = true
	end
	if not assigned then
		local sel = start.p[2] and start.p[2].t_selected and start.p[2].t_selected[1]
		local charData = sel and sel.ref and start.f_getCharData(sel.ref)
		if charData and charData.stage and #charData.stage > 0 and not (gameMode('training') and gameOption('Config.TrainingStage')) then --stage assigned as character param
			num = start.stageShuffleBag(charData.ref, charData.stage)
		elseif charData and main.stageOrder and main.t_orderStages[charData.order] then --stage assigned as stage order param
			num = start.stageShuffleBag(charData.order, main.t_orderStages[charData.order])
		elseif gameMode('training') and gameOption('Config.TrainingStage') ~= '' then --training stage
			num = start.f_getStageRef(gameOption('Config.TrainingStage'))
		else
			num = start.stageShuffleBag('includeStage', main.t_includeStage[1])
		end
	end
	selectStage(num)
	main.f_preloadBoostStage(num)
	return num
end

-- generate table with palette entries already used by this char ref
function start.f_setAssignedPal(ref, t_assignedPals)
	for side = 1, 2 do
		for k, v in pairs(start.p[side].t_selected) do
			if v.ref == ref then
				t_assignedPals[start.p[side].t_selected[k].pal] = true
			end
		end
	end
end

--remaps palette based on button press and character's keymap settings
function start.f_keyPalMap(ref, num)
	return start.f_getCharData(ref).pal_keymap[num] or num
end

-- returns palette number
function start.f_selectPal(ref, palno)
	-- generate table with palette entries already used by this char ref
	local t_assignedPals = {}
	start.f_setAssignedPal(ref, t_assignedPals)

	local charData = start.f_getCharData(ref)
	local availablePals = charData.pal

	-- selected palette by player input
	if palno ~= nil and palno > 0 then
		local mappedPal = start.f_keyPalMap(ref, palno)

		-- Check if the mapped palette is defined and not already used. (MUGEN doesn't do this)
		-- This leads to issues with certain characters who don't have the entire group 1's indices
		-- filled out, so it's been commented out for compatibility.

		-- local isDefined = false
		-- for _, p in ipairs(availablePals) do
		--     if p == mappedPal then
		--         isDefined = true
		--         break
		--     end
		-- end

		if not t_assignedPals[mappedPal] then
			return mappedPal
		end

		-- If the desired palette is not available, find the next available one.

		-- 1. Dynamically build the list of palettes to cycle through
		local cycleList = {1, 2, 3, 4, 5, 6}
		local customDefaults = false

		if charData.pal_defaults then
			local defaultsSet = {}
			for _, p_val in ipairs(charData.pal_defaults) do
				if p_val > 6 then
					-- To avoid duplicates in cycleList
					if not defaultsSet[p_val] then
						table.insert(cycleList, p_val)
						defaultsSet[p_val] = true
						customDefaults = true
					end
				end
			end
			if customDefaults then
				table.sort(cycleList) -- Ensure a consistent cycle order
			end
		end

		-- Exception: If a palette from 7 to 12 was chosen directly, cycle through all 12
		if mappedPal > 6 and not customDefaults then
			cycleList = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		end

		-- 2. Find the starting index for our search in the cycleList
		local startIndex = 1
		for i, p_val in ipairs(cycleList) do
			if p_val == mappedPal then
				startIndex = i
				break
			end
		end

		-- 3. Search for the next available palette in a circular manner
		for i = 1, #cycleList do
			-- Get the index for the next palette in the cycle
			local nextIndex = (startIndex - 1 + i) % #cycleList + 1
			local nextPal = cycleList[nextIndex]

			-- Check if this next palette is defined for the character
			local isNextDefined = false
			for _, p in ipairs(availablePals) do
				if p == nextPal then
					isNextDefined = true
					break
				end
			end

			-- If it's defined and not used, assign it.
			if isNextDefined and not t_assignedPals[nextPal] then
				return nextPal
			end
		end

		-- If all palettes in the cycle list are taken, return the originally mapped one as a fallback.
		return mappedPal

	-- default palette for AI or no-input selection
	elseif (not main.rotationChars and not gameOption('Arcade.AI.RandomColor')) or (main.rotationChars and not gameOption('Arcade.AI.SurvivalColor')) then
		for _, v in ipairs(charData.pal_defaults) do
			if not t_assignedPals[v] then
				return v
			end
		end
	end

	-- random palette
	local t = main.f_tableCopy(availablePals)
	if #t_assignedPals >= #t then -- not enough palettes for unique selection
		if #t > 0 then
			return t[math.random(1, #t)]
		else
			return 1
		end
	end
	main.f_tableShuffle(t)
	for _, v in ipairs(t) do
		if not t_assignedPals[v] then
			return v
		end
	end
	panicError("\n" .. charData.name .. " palette was not selected\n")
end

--returns player number
function start.f_getPlayerNo(side, member)
	if main.coop and not gameMode('versuscoop') then
		return side + member - 1
	end
	if side == 1 then
		return member * 2 - 1
	end
	return member * 2
end

--Convert number to name and get rid of the ""
function start.f_getName(ref, side)
	if ref == nil or start.f_getCharData(ref).hidden == 2 then
		return ''
	end
	if start.f_getCharData(ref).char == 'randomselect' or start.f_getCharData(ref).hidden == 3 then
		return motif.select_info['p' .. (side or 1)].name.random.text
	end
	return start.f_getCharData(ref).name
end

-- Resolve per-member motif tables (p3/p4/..), fallback to p1/p2 when undefined.
local function f_getMotifP(t, pn, side)
	if type(t) ~= "table" then return nil end
	local p = t["p" .. pn]
	if p ~= nil then return p end
	return t["p" .. side]
end

--reset temp data values
function start.f_resetTempData(t, subname)
	-- reset AnimData
	local seenTables = {}
	local seenAnim   = {}
	local function walk(t)
		if type(t) ~= "table" or seenTables[t] then return end
		seenTables[t] = true
		for k, v in pairs(t) do
			if k == "AnimData" and type(v) == "userdata" then
				if not seenAnim[v] then
					animReset(v)
					animUpdate(v)
					seenAnim[v] = true
				end
			elseif type(v) == "table" then
				walk(v)
			end
		end
	end
	walk(t)
	-- generate t_selTemp
	for side = 1, 2 do
		if #start.p[side].t_selTemp == 0 then
			for member, v in ipairs(start.p[side].t_selected) do
				table.insert(start.p[side].t_selTemp, {ref = v.ref})
			end
		end
		for member, v in ipairs(start.p[side].t_selTemp) do
			local pn = 2 * (member - 1) + side
			local pCfg = f_getMotifP(t, pn, side)
			if subname == '' then
				v.face_anim = pCfg.anim
				v.face_data = start.f_animGet(v.ref, side, member, pCfg, nil, true)
			else
				v.face_anim = pCfg[subname].anim
				v.face_data = start.f_animGet(v.ref, side, member, pCfg[subname], nil, true)
			end
			v.face2_anim = pCfg.face2.anim
			v.face2_data = start.f_animGet(v.ref, side, member, pCfg.face2, nil, true)
		end
		start.p[side].screenDelay = 0
	end
end

function start.f_animGet(ref, side, member, params, velParams, loop, srcAnim)
	if not ref then return nil end
	local velParams = velParams or params
	local pn = 2 * (member - 1) + side
	-- Animation/sprite priority order
	for _, v in ipairs({{params.anim, -1}, params.spr}) do
		local anim = v[1]
		if anim ~= nil and anim ~= -1 then
			-- Determine whether to apply palette
			local usePal = params.applypal or false
			-- Try to load the animation
			local a = animGetPreloadedCharData(ref, anim, v[2], loop)
			if a then
				local charData = start.f_getCharData(ref)
				local xscale = start.f_getCharData(ref).portraitscale * motif.info.localcoord[1] / start.f_getCharData(ref).localcoord
				local yscale = xscale
				if v[2] == -1 then
					xscale = xscale * (charData.cns_scale[1] or 1)
					yscale = yscale * (charData.cns_scale[2] or 1)
				end
				animSetLocalcoord(a, motif.info.localcoord[1], motif.info.localcoord[2])
				animSetLayerno(a, params.layerno)
				animSetVelocity(a, velParams.velocity[1], velParams.velocity[2])
				animSetMaxDist(a, velParams.maxdist[1], velParams.maxdist[2])
				animSetAccel(a, velParams.accel[1], velParams.accel[2])
				animSetFriction(a, velParams.friction[1], velParams.friction[2])
				animSetPos(a, 0, 0)
				animSetScale(a, params.scale[1] * xscale, params.scale[2] * yscale)
				animSetFacing(a, params.facing)
				animSetXShear(a, params.xshear)
				animSetAngle(a, params.angle)
				animSetXAngle(a, params.xangle)
				animSetYAngle(a, params.yangle)
				animSetProjection(a, params.projection)
				animSetFocalLength(a, params.focallength)
				animSetWindow(a, params.window[1], params.window[2], params.window[3], params.window[4])
				if srcAnim ~= nil then
					animApplyVel(a, srcAnim)
				end
				-- Apply palette if needed
				if usePal then
					local sel = start.p[side].t_selected[member]
					if sel and sel.ref then
						a = start.loadPalettes(a, ref, sel.pal)
					end
				end
				animUpdate(a)
				return a
			end
		end
	end
	return nil
end

--calculate portraits x pos
local function f_portraitsXCalc(side, member, paramsSide, params, skipOffset)
	local x = paramsSide.pos[1]
	if not skipOffset then
		x = x + params.offset[1]
	end
	if paramsSide.padding then
		return x + (2 * member - 1) * paramsSide.spacing[1] * paramsSide.num / (2 * math.min(paramsSide.num, math.max(start.p[side].numChars, #start.p[side].t_selected)))
	end
	return x + (member - 1) * paramsSide.spacing[1]
end

local function getParams(side, member, t, subname)
	local paramsSide = t['p' .. side]
	local pn = 2 * (member - 1) + side
	local params = f_getMotifP(t, pn, side)
	if subname and subname ~= '' then
		paramsSide = paramsSide[subname]
		params = params[subname]
	end
	return paramsSide, params
end

local function drawPortraitRandom(randomCfg)
	if not randomCfg then
		return false
	end
	local spr = randomCfg.spr
	if randomCfg.anim >= 0 or (spr and spr[1] >= 0 and spr[2] >= 0) then
		main.f_animPosDraw(randomCfg.AnimData)
		return true
	end
	return false
end

local function hasPortraitAnim(params)
	return params ~= nil and params.AnimData ~= nil and ((params.anim or -1) ~= -1 or (params.spr ~= nil and params.spr[1] ~= -1))
end

local function getPortraitDrawData(v, side, member, params, dataField)
	local data = v[dataField]
	local drawParams = params
	if not v.skipCurrent and data ~= nil then
		local state = v.ref ~= nil and getCharPreloadStatus(v.ref) or 'ready'
		if state ~= 'ready' and hasPortraitAnim(params.loading) then
			data = params.loading.AnimData
			drawParams = params.loading
		end
	elseif not v.skipCurrent and v.ref ~= nil then
		local state = getCharPreloadStatus(v.ref)
		if state == 'ready' then
			data = start.f_animGet(v.ref, side, member, params, nil, true)
			v[dataField] = data
		elseif hasPortraitAnim(params.loading) then
			data = params.loading.AnimData
			drawParams = params.loading
		end
	end
	return data, drawParams, drawParams == params.loading
end

local function drawPortraitLayer(t_portraits, side, t, subname, last, dataField)
	local lastIdx = #t_portraits
	-- "next player replaces previous one" case
	local idx = clamp(t_portraitPriority[side] or 1, 1, lastIdx)
	local paramsSide, params = getParams(side, idx, t, subname)
	if paramsSide.num == 1 and last then
		local v = t_portraits[idx]
		local data, drawParams, skipOffset = getPortraitDrawData(v, side, idx, params, dataField)
		if not v.skipCurrent and data ~= nil then
			main.f_animPosDraw(
				data,
				f_portraitsXCalc(side, 1, paramsSide, drawParams, skipOffset),
				paramsSide.pos[2] + (skipOffset and 0 or drawParams.offset[2])
			)
		end
		-- we're done for this layer in this mode
		return
	end
	-- stacked portraits up to num
	local order = {}
	for member = 1, lastIdx do
		local paramsSide, params = getParams(side, member, t, subname)
		order[#order + 1] = {m = member, o = params.draworder}
	end
	table.sort(order, function(a, b)
		if a.o == b.o then
			return a.m > b.m
		end
		return a.o < b.o
	end)
	for _, it in ipairs(order) do
		local member = it.m
		local paramsSide, params = getParams(side, member, t, subname)
		local v = t_portraits[member]
		local data, drawParams, skipOffset = getPortraitDrawData(v, side, member, params, dataField)
		if member <= paramsSide.num and not v.skipCurrent and data ~= nil then
			main.f_animPosDraw(
				data,
				f_portraitsXCalc(side, member, paramsSide, drawParams, skipOffset),
				paramsSide.pos[2] + (skipOffset and 0 or drawParams.offset[2]) + (member - 1) * paramsSide.spacing[2]
			)
		end
	end
end

-- draw portraits
function start.f_drawPortraits(t_portraits, side, t, subname, last, iconDone)
	if #t_portraits == 0 then
		return
	end
	-- reset skip flags
	for m = 1, #t_portraits do
		t_portraits[m].skipCurrent = false
	end
	-- draw random portraits (per member; required for co-op)
	for m = 1, #t_portraits do
		if t_portraits[m].inRandom then
			local pn = 2 * (m - 1) + side
			local pData = f_getMotifP(t, pn, side)
			-- face2 layer random portrait
			if pData.face2.random and drawPortraitRandom(pData.face2.random) then
				t_portraits[m].skipCurrent = true
			end
			-- primary face random portrait
			local baseFace = pData
			if subname and subname ~= '' then
				baseFace = baseFace[subname]
			end
			if baseFace.random and drawPortraitRandom(baseFace.random) then
				t_portraits[m].skipCurrent = true
			end
		end
	end
	-- face2 layer (if present)
	drawPortraitLayer(t_portraits, side, t, 'face2', last, 'face2_data')
	-- primary face layer
	drawPortraitLayer(t_portraits, side, t, subname, last, 'face_data')
	-- draw order icons (unchanged, still using main face params)
	if iconDone == nil then
		return
	end
	for member = 1, #t_portraits do
		local paramsSide, params = getParams(side, member, t, subname)
		if member > paramsSide.num then
			break
		end
		local animData = params.icon.AnimData
		if iconDone then
			animData = params.icon.done.AnimData
		end
		main.f_animPosDraw(
			animData,
			f_portraitsXCalc(side, member, paramsSide, params),
			paramsSide.pos[2] + params.offset[2] + (member - 1) * paramsSide.spacing[2]
		)
	end
end

--returns correct cell position after moving the cursor
function start.f_cellMovement(selX, selY, cmd, side, snd, dir)
	local tmpX = selX
	local tmpY = selY
	local found = false
	if getInput(cmd, motif.select_info.cell.up.key) or dir == 'U' then
		for i = 1, motif.select_info.rows do
			selY = selY - 1
			if selY < 0 then
				if motif.select_info.wrapping or dir ~= nil then
					selY = motif.select_info.rows - 1
				else
					selY = tmpY
				end
			end
			if dir ~= nil then
				found, selX = start.f_searchEmptyBoxes(selX, selY, side, -1)
			elseif (start.t_grid[selY + 1][selX + 1].char ~= nil or motif.select_info.moveoveremptyboxes) and start.t_grid[selY + 1][selX + 1].skip ~= 1 and (gameOption('Options.Team.Duplicates') or start.t_grid[selY + 1][selX + 1].char == 'randomselect' or not t_reservedChars[side][start.t_grid[selY + 1][selX + 1].char_ref]) and start.t_grid[selY + 1][selX + 1].hidden ~= 2 then
				break
			elseif motif.select_info.searchemptyboxesup then
				found, selX = start.f_searchEmptyBoxes(selX, selY, side, 1)
			end
			if found then
				break
			end
		end
	elseif getInput(cmd, motif.select_info.cell.down.key) or dir == 'D' then
		for i = 1, motif.select_info.rows do
			selY = selY + 1
			if selY >= motif.select_info.rows then
				if motif.select_info.wrapping or dir ~= nil then
					selY = 0
				else
					selY = tmpY
				end
			end
			if dir ~= nil then
				found, selX = start.f_searchEmptyBoxes(selX, selY, side, 1)
			elseif (start.t_grid[selY + 1][selX + 1].char ~= nil or motif.select_info.moveoveremptyboxes) and start.t_grid[selY + 1][selX + 1].skip ~= 1 and (gameOption('Options.Team.Duplicates') or start.t_grid[selY + 1][selX + 1].char == 'randomselect' or not t_reservedChars[side][start.t_grid[selY + 1][selX + 1].char_ref]) and start.t_grid[selY + 1][selX + 1].hidden ~= 2 then
				break
			elseif motif.select_info.searchemptyboxesdown then
				found, selX = start.f_searchEmptyBoxes(selX, selY, side, 1)
			end
			if found then
				break
			end
		end
	elseif getInput(cmd, motif.select_info.cell.left.key) or dir == 'B' then
		if dir ~= nil then
			found, selX = start.f_searchEmptyBoxes(selX, selY, side, -1)
		else
			for i = 1, motif.select_info.columns do
				selX = selX - 1
				if selX < 0 then
					if motif.select_info.wrapping then
						selX = motif.select_info.columns - 1
					else
						selX = tmpX
					end
				end
				if (start.t_grid[selY + 1][selX + 1].char ~= nil or motif.select_info.moveoveremptyboxes) and start.t_grid[selY + 1][selX + 1].skip ~= 1 and (gameOption('Options.Team.Duplicates') or start.t_grid[selY + 1][selX + 1].char == 'randomselect' or not t_reservedChars[side][start.t_grid[selY + 1][selX + 1].char_ref]) and start.t_grid[selY + 1][selX + 1].hidden ~= 2 then
					break
				end
			end
		end
	elseif getInput(cmd, motif.select_info.cell.right.key) or dir == 'F' then
		if dir ~= nil then
			found, selX = start.f_searchEmptyBoxes(selX, selY, side, 1)
		else
			for i = 1, motif.select_info.columns do
				selX = selX + 1
				if selX >= motif.select_info.columns then
					if motif.select_info.wrapping then
						selX = 0
					else
						selX = tmpX
					end
				end
				if (start.t_grid[selY + 1][selX + 1].char ~= nil or motif.select_info.moveoveremptyboxes) and start.t_grid[selY + 1][selX + 1].skip ~= 1 and (gameOption('Options.Team.Duplicates') or start.t_grid[selY + 1][selX + 1].char == 'randomselect' or not t_reservedChars[side][start.t_grid[selY + 1][selX + 1].char_ref]) and start.t_grid[selY + 1][selX + 1].hidden ~= 2 then
					break
				end
			end
		end
	end
	if (tmpX ~= selX or tmpY ~= selY) then
		if dir == nil then
			sndPlay(motif.Snd, snd[1], snd[2])
		end
	end
	return selX, selY
end

--used by above function to find valid cell in case of dummy character entries
function start.f_searchEmptyBoxes(x, y, side, direction)
	if direction > 0 then --right
		while true do
			x = x + 1
			if x >= motif.select_info.columns then
				return false, 0
			elseif start.t_grid[y + 1][x + 1].skip ~= 1 and start.t_grid[y + 1][x + 1].char ~= nil and (start.t_grid[y + 1][x + 1].char == 'randomselect' or not t_reservedChars[side][start.t_grid[y + 1][x + 1].char_ref]) and start.t_grid[y + 1][x + 1].hidden ~= 2 then
				return true, x
			end
		end
	elseif direction < 0 then --left
		while true do
			x = x - 1
			if x < 0 then
				return false, motif.select_info.columns - 1
			elseif start.t_grid[y + 1][x + 1].skip ~= 1 and start.t_grid[y + 1][x + 1].char ~= nil and (start.t_grid[y + 1][x + 1].char == 'randomselect' or not t_reservedChars[side][start.t_grid[y + 1][x + 1].char_ref]) and start.t_grid[y + 1][x + 1].hidden ~= 2 then
				return true, x
			end
		end
	end
end

-- Returns player cursor data
function start.f_getCursorData(pn)
	-- In coop/multi, p3+ may be undefined in motif; fallback to p1/p2 by parity.
	if main.coop then
		return motif.select_info['p' .. pn] or motif.select_info['p' .. ((pn - 1) % 2 + 1)]
	end
	return motif.select_info['p' .. ((pn - 1) % 2 + 1)]
end

-- Reset cursor animation for a specific slot only
local function resetCursorData(pn, store, param)
	local pData = start.f_getCursorData(pn)
	local cursorCfg = pData.cursor[param]
	local key = start.c[pn].selX .. '-' .. start.c[pn].selY
	local cursorParams = cursorCfg.default
	if cursorCfg[key] then
		cursorParams = cursorCfg[key]
	end
	local src = cursorParams.AnimData
	if not src then
		return
	end
	store[pn] = store[pn] or {}
	local cd = store[pn]
	cd.animCache = cd.animCache or {}
	local cache = cd.animCache[param]
	if cache == nil or cache.src ~= src then
		cache = {src = src, anim = animCopy(src)}
		cd.animCache[param] = cache
		if cache.anim then
			animReset(cache.anim)
			animUpdate(cache.anim)
		end
	end
	if cache.anim then
		animReset(cache.anim)
		animUpdate(cache.anim)
	end
end

-- Calculate cursor.tween
local function cursorTween(val, target, factor)
	if not factor or not target then
		return val
	end
	for i = 1, 2 do
		local t = target[i] or 0
		local f = math.min(math.abs(factor[i] or 0.5), 1)
		val[i] = val[i] + (t - val[i]) * f
	end
	return val
end

local function getCellOverride(col, row)
    local cells = motif.select_info.cell
    local exact = col .. '-' .. row
	local colWild = col .. '-*'
	local rowWild = '*-' .. row
    if cells[exact] then 
		return cells[exact] 
	end
    if cells[colWild] then 
		return cells[colWild] 
	end
    if cells[rowWild] then 
		return cells[rowWild] 
	end
    if cells['*-*'] then 
		return cells['*-*'] 
	end
    return nil
end

function getCellFacing(default, col, row)
	local override = getCellOverride(col, row)
	if override ~= nil and override.facing ~= 0 then
		return override.facing
	end
	return default
end

function getCellOffset(col, row)
	local override = getCellOverride(col, row)
	if override ~= nil and override.offset ~= nil then
		return override.offset
	end
	return {0, 0}
end

function getCellSpacing(col, row)
	local override = getCellOverride(col, row)
	if override ~= nil and override.spacing ~= nil then
		if override.spacing[1] ~= 0 then
			if override.spacing[2] == 0 then
				return {override.spacing[1], override.spacing[1]}
			end
			return override.spacing
		end
		if override.spacing[2] ~= 0 then
			return override.spacing
		end
	end
	return motif.select_info.cell.spacing
end

function getCellSkip(col, row)
	local override = getCellOverride(col, row)
	if override ~= nil then
		return override.skip
	end
	return false
end

function getCellTransform(col, row, paramName, default)
	local override = getCellOverride(col, row)
	if override ~= nil then
		local val = override[paramName]
		-- Table Validation
		if type(val) == "table" then
			if paramName == "scale" then
				if val[1] ~= 0 or val[2] ~= 0 then return val end
			else
				return val
			end
		elseif type(val) == "string" then
			if val ~= "" then 
				return val 
			end
		elseif type(val) == "number" then
			if val ~= 0 then 
				return val 
			end
		end
	end
	return default
end

--draw cursor
function start.f_drawCursor(pn, x, y, param, done)
	local pData = start.f_getCursorData(pn)
	-- select appropriate cursor table and initialize if needed
	local store = done and cursorDone or cursorActive
	store[pn] = store[pn] or {}
	local cd = store[pn]
	cd.currentPos  = cd.currentPos  or {0, 0}
	cd.targetPos   = cd.targetPos   or {0, 0}
	cd.startPos    = cd.startPos    or {0, 0}
	cd.slideOffset = cd.slideOffset or {0, 0}
	cd.init        = cd.init or false
	if not done then
		cd.snap = cd.snap or false -- only used by active cursors
	end
	-- calculate target cell coordinates using the pre-calculated grid
	local cellData = start.t_grid[y + 1] and start.t_grid[y + 1][x + 1]
	local baseX, baseY
	if cellData then
		-- cellData already includes all spacing and offsets
		baseX = motif.select_info.pos[1] + cellData.x
		baseY = motif.select_info.pos[2] + cellData.y
	end
	-- initialization or snap: set cursor directly
	if not cd.init or done or cd.snap then
		for i = 1, 2 do
			cd.currentPos[i] = (i == 1) and baseX or baseY
			cd.targetPos[i]  = cd.currentPos[i]
			cd.startPos[i]   = cd.currentPos[i]
			cd.slideOffset[i]= 0
		end
		cd.init, cd.snap = true, false
	-- new cell selected: recalc tween offsets
	elseif cd.targetPos[1] ~= baseX or cd.targetPos[2] ~= baseY then
		cd.startPos[1], cd.startPos[2] = cd.currentPos[1], cd.currentPos[2]
		cd.targetPos[1], cd.targetPos[2] = baseX, baseY
		cd.slideOffset[1] = cd.startPos[1] - baseX
		cd.slideOffset[2] = cd.startPos[2] - baseY
	end
	local t_factor = {
		pData.cursor.tween.factor[1],
		pData.cursor.tween.factor[2]
	}
	-- apply tween if enabled, otherwise snap to target
	if not done and t_factor[1] > 0 and t_factor[2] > 0 then
		cursorTween(cd.slideOffset, {0, 0}, t_factor)
	else
		cd.slideOffset[1], cd.slideOffset[2] = 0, 0
	end
	if pData.cursor.tween.wrap.snap then
		local dx = cd.targetPos[1] - cd.startPos[1]
		local dy = cd.targetPos[2] - cd.startPos[2]
		if math.abs(dx) > motif.select_info.cell.size[1] * (motif.select_info.columns - 1) or math.abs(dy) > motif.select_info.cell.size[2] * (motif.select_info.rows - 1) then
		cd.slideOffset[1], cd.slideOffset[2] = 0, 0	
		end
	end
	-- update final cursor position
	cd.currentPos[1] = cd.targetPos[1] + cd.slideOffset[1]
	cd.currentPos[2] = cd.targetPos[2] + cd.slideOffset[2]
	-- draw
	local params = pData.cursor[param].default
	local key = x .. '-' .. y
	if pData.cursor[param][key] ~= nil then
		params = pData.cursor[param][key]
	end
	local a = params.AnimData
	cd.animCache = cd.animCache or {}
	local cache = cd.animCache[param]
	if cache == nil or cache.src ~= a then
		cache = {src = a, anim = animCopy(a)}
		cd.animCache[param] = cache
		if cache.anim then
			animReset(cache.anim)
		end
	end
	a = cache.anim
	animSetFacing(a, getCellFacing(params.facing, x, y))
	local scale = getCellTransform(x, y, "scale", params.scale)
	animSetScale(a, scale[1], scale[2])
	animSetXShear(a, getCellTransform(x, y, "xshear", params.xshear))
	animSetAngle(a, getCellTransform(x, y, "angle", params.angle))
	animSetXAngle(a, getCellTransform(x, y, "xangle", params.xangle))
	animSetYAngle(a, getCellTransform(x, y, "yangle", params.yangle))
	animSetProjection(a, getCellTransform(x, y, "projection", params.projection))
	animSetFocalLength(a, getCellTransform(x, y, "focallength", params.focallength))
	animUpdate(a)
	main.f_animPosDraw(a, cd.currentPos[1], cd.currentPos[2], getCellFacing(params.facing, x, y))
end

-- snaps the cursor instantly to its target cell
local function f_snapCursor()
	for k, v in pairs(cursorActive) do
		v.snap = true
	end
end

--returns t_selChars table out of cell number
function start.f_selGrid(cell, slot)
	if main.t_selGrid[cell] == nil or #main.t_selGrid[cell].chars == 0 then
		local csCol = ((cell - 1) % motif.select_info.columns) + 1
		local csRow = math.floor((cell - 1) / motif.select_info.columns) + 1
		local cellCfg = motif.select_info.cell[(csCol - 1) .. '-' .. (csRow - 1)]
		if cellCfg ~= nil and cellCfg.skip then
			return {skip = 1}
		end
		return {}
	end
	return main.t_selChars[main.t_selGrid[cell].chars[(slot or main.t_selGrid[cell].slot)]]
end

--returns t_selChars table out of char ref
function start.f_getCharData(ref)
	if type(ref) ~= 'number' then
		return nil
	end
	return main.t_selChars[ref + 1]
end

local t_recordTierColors = {
	Z = {255, 215, 0},
	U = {128, 224, 255},
	X = {185, 0, 0},
	S = {255, 96, 96},
	A = {255, 160, 0},
	B = {224, 208, 0},
	C = {84, 188, 255},
	D = {180, 120, 255},
	F = {160, 160, 160},
}

start.t_recordTierOrder = {
	'U',
	'F---', 'F--', 'F-', 'F', 'F+', 'F++', 'F+++',
	'D---', 'D--', 'D-', 'D', 'D+', 'D++', 'D+++',
	'C---', 'C--', 'C-', 'C', 'C+', 'C++', 'C+++',
	'B---', 'B--', 'B-', 'B', 'B+', 'B++', 'B+++',
	'A---', 'A--', 'A-', 'A', 'A+', 'A++', 'A+++',
	'S---', 'S--', 'S-', 'S', 'S+', 'S++', 'S+++',
	'X---', 'X--', 'X-', 'X', 'X+', 'X++', 'X+++',
	'Z',
}

local t_recordTierIndex = {}
for i, tier in ipairs(start.t_recordTierOrder) do
	t_recordTierIndex[tier] = i
end

local function f_normalizeRecordTier(tier)
	tier = tostring(tier or 'U'):upper()
	if t_recordTierIndex[tier] ~= nil then
		return tier
	end
	local group = tier:sub(1, 1)
	if group == 'U' or group == 'F' or group == 'D' or group == 'C' or group == 'B' or group == 'A' or group == 'S' or group == 'X' or group == 'Z' then
		return group
	end
	return 'U'
end

local function f_recordTierDisplayText(tier)
	tier = f_normalizeRecordTier(tier)
	if tier == 'U' or tier == 'Z' then
		return tier
	end
	local group = tier:sub(1, 1)
	local suffix = tier:sub(2)
	if suffix == '---' then
		return 'LOW ' .. group .. '-'
	elseif suffix == '--' then
		return 'MID ' .. group .. '-'
	elseif suffix == '-' then
		return 'HIGH ' .. group .. '-'
	elseif suffix == '' then
		return 'MID ' .. group
	elseif suffix == '+' then
		return 'LOW ' .. group .. '+'
	elseif suffix == '++' then
		return 'MID ' .. group .. '+'
	elseif suffix == '+++' then
		return 'HIGH ' .. group .. '+'
	end
	return tier
end

local function f_recordWinRateTier(record)
	local matches = tonumber(record.matches) or 0
	local wins = tonumber(record.wins) or 0
	local losses = tonumber(record.losses) or 0
	if wins == 0 and losses == 0 then
		return 'U'
	elseif matches <= 0 then
		matches = wins + losses
	end
	if matches <= 0 then
		return 'U'
	elseif wins == 1 and losses == 0 then
		return 'D'
	elseif wins == 2 and losses == 0 then
		return 'D+++'
	elseif wins == 3 and losses == 0 then
		return 'C---'
	elseif wins == 4 and losses == 0 then
		return 'C+++'
	elseif wins == 5 and losses == 0 then
		return 'B---'
	elseif wins == 6 and losses == 0 then
		return 'B+++'
	elseif wins == 7 and losses == 0 then
		return 'A---'
	elseif wins == 8 and losses == 0 then
		return 'A+++'
	elseif wins == 9 and losses == 0 then
		return 'S---'
	elseif wins == 10 and losses == 0 then
		return 'S+++'
	elseif wins > 0 and losses == 0 and wins == matches then
		return 'Z'
	end
	local winRate = wins / matches
	local tierCutoffs = {
		{0.96, 'X+++'}, {0.94, 'X++'}, {0.92, 'X+'}, {0.90, 'X'}, {0.88, 'X-'}, {0.86, 'X--'}, {0.84, 'X---'},
		{0.82, 'S+++'}, {0.80, 'S++'}, {0.78, 'S+'}, {0.76, 'S'}, {0.74, 'S-'}, {0.72, 'S--'}, {0.70, 'S---'},
		{0.68, 'A+++'}, {0.66, 'A++'}, {0.64, 'A+'}, {0.62, 'A'}, {0.60, 'A-'}, {0.58, 'A--'}, {0.56, 'A---'},
		{0.54, 'B+++'}, {0.52, 'B++'}, {0.50, 'B+'}, {0.48, 'B'}, {0.46, 'B-'}, {0.44, 'B--'}, {0.42, 'B---'},
		{0.40, 'C+++'}, {0.38, 'C++'}, {0.36, 'C+'}, {0.34, 'C'}, {0.32, 'C-'}, {0.30, 'C--'}, {0.28, 'C---'},
		{0.26, 'D+++'}, {0.24, 'D++'}, {0.22, 'D+'}, {0.20, 'D'}, {0.18, 'D-'}, {0.16, 'D--'}, {0.14, 'D---'},
		{0.12, 'F+++'}, {0.10, 'F++'}, {0.08, 'F+'}, {0.06, 'F'}, {0.04, 'F-'}, {0.02, 'F--'},
	}
	for _, row in ipairs(tierCutoffs) do
		if winRate >= row[1] then
			return row[2]
		end
	end
	return 'F---'
end

local t_fightStatsCache = nil
local t_fightStatsCacheTime = nil
local t_fightStatsCachePath = nil
local function f_readFightStats()
	local now = type(gameTime) == 'function' and math.floor(gameTime() / 60) or os.time()
	if t_fightStatsCache ~= nil and t_fightStatsCacheTime == now and t_fightStatsCachePath == start.statsFilePath then
		return t_fightStatsCache
	end
	local ok, data = pcall(jsonDecode, start.statsFilePath)
	if not ok or type(data) ~= 'table' then
		data = {}
	end
	t_fightStatsCache = data
	t_fightStatsCacheTime = now
	t_fightStatsCachePath = start.statsFilePath
	return t_fightStatsCache
end

local function f_recordLookupKeys(data)
	local keys = {}
	local seen = {}
	local function addKey(key)
		if key == nil or key == '' or seen[key] then
			return
		end
		table.insert(keys, key)
		seen[key] = true
	end
	local function add(value)
		if value == nil then
			return
		end
		local key = tostring(value):gsub('\\', '/'):gsub('^%s+', ''):gsub('%s+$', ''):lower()
		if key ~= '' then
			addKey(key)
			addKey((key:gsub('[^a-z0-9]+', '')))
			addKey((key:gsub('^/', '')))
			addKey((key:gsub('%.def$', '')))
			local stem = key:gsub('^.*[/]', ''):gsub('%.def$', '')
			addKey(stem)
			addKey((stem:gsub('[^a-z0-9]+', '')))
			local folder = key:match('^([^/]+)/')
			if folder ~= nil then
				addKey(folder)
				addKey((folder:gsub('[^a-z0-9]+', '')))
			end
		end
	end
	add(data.recordKey)
	add(data.char)
	add(data.def)
	add(data.name)
	return keys
end

local function f_recordStatsVariants(key)
	local variants = {}
	local seen = {}
	local function add(value)
		if value == nil then
			return
		end
		local item = tostring(value):gsub('\\', '/'):gsub('^%s+', ''):gsub('%s+$', ''):lower()
		if item == '' or seen[item] then
			return
		end
		table.insert(variants, item)
		seen[item] = true
	end
	add(key)
	add((key or ''):gsub('^/', ''))
	add((key or ''):gsub('%.def$', ''))
	add((key or ''):gsub('^.*[/]', ''):gsub('%.def$', ''))
	local folder = tostring(key or ''):match('^([^/]+)/')
	add(folder)
	add((tostring(key or ''):gsub('[^a-z0-9]+', '')))
	return variants
end

local function f_recordCompactKey(value)
	return tostring(value or ''):gsub('\\', '/'):lower():gsub('[^a-z0-9]+', '')
end

local function f_bestRecordCandidate(current, rec, key)
	if type(rec) ~= 'table' then
		return current
	end
	local wins = tonumber(rec.wins or rec.win) or 0
	local losses = tonumber(rec.losses or rec.loss) or 0
	local matches = tonumber(rec.matches) or (wins + losses)
	local score = matches * 100000 + wins * 100 - losses
	if current == nil or score > current.score then
		return {
			score = score,
			record = rec,
			key = key,
			wins = wins,
			losses = losses,
			matches = matches,
		}
	end
	return current
end

-- A stats record that names its own fighter (record.char) must never be
-- credited to, or read for, a different fighter. Without this, a brand-new
-- fighter whose file/display name resembles an existing one (Mario_KF vs
-- mario) inherited that fighter's record and tier, so the ladder paired it in
-- the wrong tier and its results were written onto the other fighter.
local function f_recordOwnerNorm(value)
	local s = tostring(value or ''):lower():gsub('\\', '/'):gsub('^chars/', ''):gsub('%.def$', '')
	local a, b = s:match('^(.-)/(.-)$')
	if a ~= nil and a == b then
		s = a
	end
	return s
end
local function f_recordOwnedByOther(rec, ownerChar)
	if ownerChar == nil or type(rec) ~= 'table' or rec.char == nil or tostring(rec.char) == '' then
		return false
	end
	return f_recordOwnerNorm(rec.char) ~= f_recordOwnerNorm(ownerChar)
end
local t_statsCompactIndex = setmetatable({}, {__mode = 'k'})
local function f_findBestStatsRecord(source, keys, ownerChar)
	local best = nil
	for _, lookupKey in ipairs(keys) do
		for _, key in ipairs(f_recordStatsVariants(lookupKey)) do
			local rec = source[key] or source[key:upper()] or source[key:lower()]
			if not f_recordOwnedByOther(rec, ownerChar) then
				best = f_bestRecordCandidate(best, rec, key)
			end
		end
	end
	if best ~= nil then
		return best
	end
	-- Fuzzy fallback via a compact-key index built once per stats table.
	-- (Scanning every stats record per lookup made the ladder's tier pool
	-- take ~2 minutes between fights.)
	local index = t_statsCompactIndex[source]
	if index == nil then
		index = {}
		for statKey in pairs(source) do
			for _, key in ipairs(f_recordStatsVariants(statKey)) do
				local compact = f_recordCompactKey(key)
				local list = index[compact]
				if list == nil then
					index[compact] = {statKey}
				elseif list[#list] ~= statKey then
					table.insert(list, statKey)
				end
			end
		end
		t_statsCompactIndex[source] = index
	end
	local checked = {}
	for _, key in ipairs(keys) do
		local compact = f_recordCompactKey(key)
		if compact ~= '' and index[compact] ~= nil then
			for _, statKey in ipairs(index[compact]) do
				if not checked[statKey] then
					checked[statKey] = true
					local rec = source[statKey]
					if not f_recordOwnedByOther(rec, ownerChar) then
						best = f_bestRecordCandidate(best, rec, statKey)
					end
				end
			end
		end
	end
	return best
end

function start.f_getCharRecordKey(ref)
	local data = start.f_getCharData(ref)
	if data == nil then
		return nil
	end
	return tostring(data.recordKey or data.char or data.def or data.name or ''):lower()
end

function start.f_getCharRecord(ref)
	local data = start.f_getCharData(ref)
	local key = start.f_getCharRecordKey(ref)
	if data == nil then
		return {wins = 0, losses = 0, matches = 0, tier = 'U', key = key}
	end
	local stats = f_readFightStats()
	local source = type(stats.characters) == 'table' and stats.characters or stats
	local best = f_findBestStatsRecord(source, f_recordLookupKeys(data), data.char)
	if best ~= nil then
		return {
			wins = best.wins,
			losses = best.losses,
			matches = best.matches,
			tier = best.record.tier or best.record.rank,
			key = best.key,
		}
	end
	return {wins = 0, losses = 0, matches = 0, tier = 'U', key = key}
end

local t_factionIndex = nil
-- Keep the faction overlay reliable for legacy roster keys whose metadata
-- may be generated with a different token casing/path.
local t_factionDisplayOverrides = {
	['knuckles'] = 'Lancero',
	['riki'] = 'DBC',
}
local function f_addFactionIndexKey(index, key, faction)
	if key == nil or faction == nil or tostring(faction) == '' then
		return
	end
	local value = tostring(faction):gsub('%s+[Ff]action%s*$', ''):gsub('^%s+', ''):gsub('%s+$', '')
	if value == '' then
		return
	end
	local raw = tostring(key):gsub('\\', '/'):gsub('^%s+', ''):gsub('%s+$', ''):lower()
	if raw == '' then
		return
	end
	index[raw] = value
	index[raw:gsub('^/', '')] = value
	index[raw:gsub('%.def$', '')] = value
	index[raw:gsub('^.*[/]', ''):gsub('%.def$', '')] = value
	index[raw:gsub('[^a-z0-9]+', '')] = value
end

local function f_buildFactionIndex()
	local index = {}
	if main ~= nil and type(main.t_selChars) == 'table' then
		for _, data in ipairs(main.t_selChars) do
			if type(data) == 'table' and data.faction ~= nil then
				for _, key in pairs({data.recordKey, data.char, data.def, data.name, data.displayname}) do
					f_addFactionIndexKey(index, key, data.faction)
				end
			end
		end
	end
	local function loadRoster(path)
		local ok, rosterData = pcall(jsonDecode, path)
		if not ok or type(rosterData) ~= 'table' then
			return
		end
		local roster = rosterData.roster or rosterData.characters
		if type(roster) ~= 'table' then
			return
		end
		for key, info in pairs(roster) do
			if type(info) == 'table' and info.faction ~= nil then
				f_addFactionIndexKey(index, key, info.faction)
				for _, item in pairs({info.key, info.token, info.char, info.def, info.recordKey, info.name, info.displayName, info.originalName}) do
					f_addFactionIndexKey(index, item, info.faction)
				end
				if type(info.aliases) == 'table' then
					for _, alias in ipairs(info.aliases) do
						f_addFactionIndexKey(index, alias, info.faction)
					end
				end
			end
		end
	end
	loadRoster('save/roster-data.json')
	loadRoster('web/roster-data.json')
	loadRoster('/home/lancero7777/code/LanceroMugen.com/src/Web/data/roster-data.json')
	loadRoster('/home/lancero7777/code/LanceroMugen.com/LiveLancero/data/roster-data.json')
	return index
end

local t_factionRankIndex = nil
local function f_addFactionRankIndexKey(index, key, rank, total)
	if key == nil or rank == nil or total == nil then
		return
	end
	local raw = tostring(key):gsub('\\', '/'):gsub('^%s+', ''):gsub('%s+$', ''):lower()
	if raw == '' then
		return
	end
	local value = {rank = rank, total = total}
	index[raw] = value
	index[raw:gsub('^/', '')] = value
	index[raw:gsub('%.def$', '')] = value
	index[raw:gsub('^.*[/]', ''):gsub('%.def$', '')] = value
	index[raw:gsub('[^a-z0-9]+', '')] = value
end

local function f_buildFactionRankIndex()
	local index = {}
	local function loadRosterList(path)
		local ok, rosterData = pcall(jsonDecode, path)
		if not ok or type(rosterData) ~= 'table' then
			return nil
		end
		local roster = rosterData.roster or rosterData.characters
		if type(roster) ~= 'table' then
			return nil
		end
		return roster
	end
	local roster = nil
	for _, path in ipairs({
		'save/roster-data.json',
		'web/roster-data.json',
		'/home/lancero7777/code/LanceroMugen.com/src/Web/data/roster-data.json',
		'/home/lancero7777/code/LanceroMugen.com/LiveLancero/data/roster-data.json',
	}) do
		roster = loadRosterList(path)
		if roster ~= nil then
			break
		end
	end
	if roster == nil then
		return index
	end
	local factionTotals = {}
	for _, info in pairs(roster) do
		if type(info) == 'table' and info.faction ~= nil and info.factionRankNumber ~= nil then
			local faction = tostring(info.faction)
			local rank = tonumber(info.factionRankNumber) or 0
			factionTotals[faction] = math.max(factionTotals[faction] or 0, rank)
		end
	end
	for _, info in pairs(roster) do
		if type(info) == 'table' and info.faction ~= nil and info.factionRankNumber ~= nil then
			local faction = tostring(info.faction)
			local rank = tonumber(info.factionRankNumber)
			local total = factionTotals[faction]
			if rank ~= nil and total ~= nil then
				for _, item in pairs({info.key, info.token, info.char, info.def, info.recordKey, info.name, info.displayName, info.originalName}) do
					f_addFactionRankIndexKey(index, item, rank, total)
				end
				if type(info.aliases) == 'table' then
					for _, alias in ipairs(info.aliases) do
						f_addFactionRankIndexKey(index, alias, rank, total)
					end
				end
			end
		end
	end
	return index
end

function start.f_getCharFactionRank(ref)
	local data = start.f_getCharData(ref)
	if data == nil then
		return nil, nil
	end
	if t_factionRankIndex == nil then
		t_factionRankIndex = f_buildFactionRankIndex()
	end
	for _, key in ipairs(f_recordLookupKeys(data)) do
		local entry = t_factionRankIndex[key] or t_factionRankIndex[f_recordCompactKey(key)]
		if entry ~= nil then
			return entry.rank, entry.total
		end
	end
	return nil, nil
end

function start.f_getCharFaction(ref)
	local data = start.f_getCharData(ref)
	if data == nil then
		return nil
	end
	if data.faction ~= nil and tostring(data.faction) ~= '' then
		return tostring(data.faction):gsub('%s+[Ff]action%s*$', '')
	end
	if t_factionIndex == nil then
		t_factionIndex = f_buildFactionIndex()
	end
	for _, key in ipairs(f_recordLookupKeys(data)) do
		local faction = t_factionIndex[key] or t_factionIndex[f_recordCompactKey(key)]
		if faction ~= nil and faction ~= '' then
			return faction
		end
	end
	for _, key in ipairs(f_recordLookupKeys(data)) do
		local faction = t_factionDisplayOverrides[key] or t_factionDisplayOverrides[f_recordCompactKey(key)]
		if faction ~= nil and faction ~= '' then
			return faction
		end
	end
	local identity = tostring(data.name or data.char or data.def or ''):lower()
	if identity:find('knuckles', 1, true) then
		return 'Lancero'
	elseif identity:find('riki', 1, true) then
		return 'DBC'
	end
	return nil
end

function start.f_getRecordTier(record)
	if record == nil then
		return 'U'
	end
	if (tonumber(record.wins) or 0) > 0 or (tonumber(record.losses) or 0) > 0 or (tonumber(record.matches) or 0) > 0 then
		return f_recordWinRateTier(record)
	end
	local explicitTier = f_normalizeRecordTier(record.tier)
	return explicitTier
end

function start.f_getRecordTierDisplayText(recordOrTier)
	if type(recordOrTier) == 'table' then
		return f_recordTierDisplayText(start.f_getRecordTier(recordOrTier))
	end
	return f_recordTierDisplayText(recordOrTier)
end

function start.f_getRecordTierText(record)
	local display = start.f_getRecordTierDisplayText(record)
	local gap = (display == 'Z' or display == 'U') and '   ' or ' '
	return display .. gap .. 'Tier'
end

function start.f_getRecordTierColor(record)
	local tier = start.f_getRecordTier(record)
	local color = t_recordTierColors[tier] or t_recordTierColors[tier:sub(1, 1)] or t_recordTierColors.F
	return color[1], color[2], color[3]
end

start.t_forcedRecordGoodEvil = {
	['yuji itadori'] = 'PURE',
	['gojo'] = 'GOOD',
	['gill'] = 'EVIL',
	['gill-enhancedai'] = 'EVIL',
	['orochi gill'] = 'EVIL',
	['gill rr/gill.def'] = 'EVIL',
	['o_gill'] = 'EVIL',
	['frazard'] = 'EVIL',
	['hashirama'] = 'GOOD',
	['nintendo_punchout_littlemac'] = 'NEUTRAL',
	['little mac'] = 'NEUTRAL',
	['dempsy ippo'] = 'GOOD',
	['ippo makunouchi'] = 'GOOD',
	['deadly'] = 'EVIL',
	['trixie lulamoon'] = 'EVIL',
	['the great and powerful trixie'] = 'EVIL',
	['ryu'] = 'IDPT',
	['terry99m'] = 'IDPT',
	['terry99m/terry99m.def'] = 'IDPT',
	['sf3_gill_k'] = 'EVIL',
	['wario'] = 'IDPT',
	['scorpionjuan'] = 'IDPT',
	['link'] = 'IDPT',
	['son goku rn'] = 'PURE',
	['yga'] = 'IDPT',
	['robert-keyser'] = 'IDPT',
	['frieza_bt'] = 'EVIL',
	['frieza'] = 'EVIL',
	['freezaz2'] = 'EVIL',
	['goldenfreezaz2'] = 'EVIL',
	['emperor_frieza'] = 'EVIL',
	['cooler final'] = 'CHAOTIC',
	['ui_gokuz2'] = 'PURE',
	['whatsappgoku'] = 'PURE',
	['kakarotto ssj3'] = 'PURE',
	['goku'] = 'PURE',
	['goku/goku.def'] = 'PURE',
	['gokuz2'] = 'PURE',
	['songoku'] = 'PURE',
	['ue vegeta z2i'] = 'ANTI-HERO',
	['armor_vegetaz2/armor_vegetaz2'] = 'ANTI-HERO',
	['discordvegeta'] = 'ANTI-HERO',
	['majin vegeta (dbz)'] = 'ANTI-HERO',
	['vegeta'] = 'ANTI-HERO',
	['sonic the hedgehog'] = 'GOOD',
	['ai-sonicmhii/ai-sonicmhii'] = 'GOOD',
	['sonic_tp njpt'] = 'GOOD',
	['broly'] = 'CHAOTIC',
	['broly lssj3'] = 'CHAOTIC',
	['broly-mugen10'] = 'CHAOTIC',
	['broly z2'] = 'CHAOTIC',
	['brolynew/brolynew'] = 'CHAOTIC',
	['brolynew/brolynew.def'] = 'CHAOTIC',
	['brolyz2/brolyz2'] = 'CHAOTIC',
	['brolyz2/brolyz2.def'] = 'CHAOTIC',
	['new nightmare broly'] = 'CHAOTIC',
	['kenshiro_c2'] = 'PURE',
	['kenshiro-kofm'] = 'GOOD',
	['kenshironew'] = 'PURE',
	['kenshirou_ai-patch'] = 'GOOD',
	['kenoh'] = 'EVIL',
	['conqueror raoh_s3'] = 'EVIL',
	['raoh'] = 'EVIL',
	['alexex'] = 'ANTI-HERO',
	['banzoku-alex'] = 'ANTI-HERO',
	['sf3_alex'] = 'ANTI-HERO',
	['rei_ai-patch'] = 'GOOD',
	['bug_rei/bug_rei.def'] = 'GOOD',
	['bug_rei'] = 'GOOD',
	['bug rei'] = 'GOOD',
	['definitive rei ai alt'] = 'GOOD',
	['rei-hnk-kofa'] = 'GOOD',
	['rei'] = 'GOOD',
	['kenshiro'] = 'PURE',
	['kenshiro-hell'] = 'GOOD',
	['toki-kofa'] = 'PURE',
	['toki_ai-patch'] = 'PURE',
	['toki'] = 'PURE',
	['shin_ai-patch'] = 'ANTI-HERO',
	['crazy-shin/crazy-shin.def'] = 'ANTI-HERO',
	['crazy-shin'] = 'ANTI-HERO',
	['shin'] = 'ANTI-HERO',
	['jagi'] = 'ANTI-HERO',
	['jagi_r'] = 'ANTI-HERO',

	['alteramiba'] = 'NEUTRAL',
	['homerjsimpson/homerjsimpson.def'] = 'NEUTRAL',
	['homerjsimpson'] = 'NEUTRAL',
	['homer j simpson'] = 'NEUTRAL',
	['patrick'] = 'NEUTRAL',
	['sash_lilac'] = 'NEUTRAL',
	['zerohaohmaru'] = 'NEUTRAL',
	['haohmaru'] = 'NEUTRAL',
	['haohmaru/haohmaru.def'] = 'NEUTRAL',
	['haohmaru kofm'] = 'NEUTRAL',
	['piccoloz2'] = 'NEUTRAL',
	['hakufu sonsaku'] = 'NEUTRAL',
	['blackprincess_tohka'] = 'NEUTRAL',
	['mamiya'] = 'NEUTRAL',
	['yuri kyaku'] = 'NEUTRAL',
	['a-shi'] = 'NEUTRAL',
	['morrigan aensland'] = 'NEUTRAL',
	['kanae'] = 'NEUTRAL',
	['ultimate c.falcon'] = 'NEUTRAL',
	['ryuko2nd'] = 'NEUTRAL',
	['morrigan_pots'] = 'NEUTRAL',
	['awk_morrigan'] = 'NEUTRAL',
	['felicia redhot'] = 'NEUTRAL',
	['felicia_aipatch'] = 'NEUTRAL',
	['felicia'] = 'NEUTRAL',
	['botan'] = 'NEUTRAL',
	['coco/coco/coco.def'] = 'NEUTRAL',
	['pekora'] = 'NEUTRAL',
	['suisei'] = 'NEUTRAL',
	['juggernault'] = 'NEUTRAL',
	['helder'] = 'ANTI-HERO',
	['zeroreika'] = 'NEUTRAL',
	['zeroreika/zeroreika.def'] = 'NEUTRAL',
	['yuki'] = 'NEUTRAL',
	['bloodñ+killer'] = 'NEUTRAL',
	['perfect weapon mb-02'] = 'NEUTRAL',
	['majin starfish v3'] = 'NEUTRAL',
	['d-donald'] = 'NEUTRAL',
	['gogetasuper4'] = 'NEUTRAL',
	['cvtwstriderhien'] = 'NEUTRAL',
	['striderhiryu'] = 'NEUTRAL',
	['ace'] = 'NEUTRAL',
	['cvsiori'] = 'NEUTRAL',
	['crazy_vega'] = 'NEUTRAL',
	['kenshin himura'] = 'NEUTRAL',
	['fei-long'] = 'NEUTRAL',
	['jin(the evil awakens 2)'] = 'EVIL',
	['heihachi'] = 'EVIL',
	['iori'] = 'EVIL',
	['clone blood igniz'] = 'NEUTRAL',
	['mizuchi'] = 'NEUTRAL',
	['beterryu/ryu.def'] = 'NEUTRAL',
	['omegath'] = 'NEUTRAL',
	['theking'] = 'GOOD',
	['morshu'] = 'ANTI-HERO',
	['morshu/morshu.def'] = 'ANTI-HERO',

	['donald'] = 'CHAOTIC',
	['donald/donald.def'] = 'CHAOTIC',
	['awakened-clark'] = 'CHAOTIC',
	['kof_orochi_shermie'] = 'CHAOTIC',
	['frozen-yashiro'] = 'CHAOTIC',
	['genericbrad'] = 'CHAOTIC',
	['colonel sanders'] = 'CHAOTIC',
	['hitto'] = 'CHAOTIC',
	['bulla/bulla'] = 'CHAOTIC',
	['evilsagat/definition.def'] = 'EVIL',
	['k-feilong_wls'] = 'EVIL',
	['evilcarlos'] = 'EVIL',
	['evildanx'] = 'EVIL',
	['super sonic(the evil awakens 2)/super sonic(the evil awakens 2).def'] = 'EVIL',
	['daredevil'] = 'EVIL',
	['evil donald/evil donald.def'] = 'EVIL',
	['mariops'] = 'EVIL',
	['devil_mario_prime'] = 'EVIL',
	['119way-e-ryu'] = 'EVIL',
	['cyberryuevil'] = 'EVIL',
	["evilryusf'2/evilryusf'2.def"] = 'EVIL',
	['lord evil ken'] = 'EVIL',
	['lord evil ryu'] = 'EVIL',
	['evil ryu alpha (f.f)'] = 'EVIL',
	['donald_return'] = 'EVIL',
	['fire combat donazen/fire combat donazen.def'] = 'CHAOTIC',
	['gamer of 2021/gamer of 2021.def'] = 'CHAOTIC',
	['gamer of 2021'] = 'CHAOTIC',
	['agent mac'] = 'CHAOTIC',
	['fake phantom donald/fake phantom donald.def'] = 'CHAOTIC',
	['fake phantom donald'] = 'CHAOTIC',
	['donald_solo_a5/donald_solo_a5.def'] = 'CHAOTIC',
	['donald_solo_a5'] = 'CHAOTIC',
	['donald solo a5'] = 'CHAOTIC',
	['souldonaldrainbow/souldonaldrainbow.def'] = 'CHAOTIC',
	['souldonaldrainbow'] = 'CHAOTIC',
	['rainbow souldonald v1.1'] = 'CHAOTIC',
	['super evildonald/super evildonald.def'] = 'EVIL',
	['super evildonald'] = 'EVIL',
	['super evil donald'] = 'EVIL',
	['the melancholy/the melancholy.def'] = 'CHAOTIC',
	['the melancholy'] = 'CHAOTIC',
	['winter'] = 'CHAOTIC',
	['donald_solo_a5_aero/donald_solo_a5_aero.def'] = 'CHAOTIC',
	['ice donald/ice donald.def'] = 'CHAOTIC',
	['spicy donald/spicy donald.def'] = 'CHAOTIC',
	['donald_miku/miku.def'] = 'CHAOTIC',
	['donald_solo_1st/donald_solo_1st.def'] = 'CHAOTIC',
	['donald_solo_2nd_alpha6/donald_solo_2nd_alpha6.def'] = 'CHAOTIC',
	['d-donald/d-donald.def'] = 'EVIL',}

start.t_forcedRecordAlignment = {
	-- Final Fantasy pack (2026-09-24)
	['bulk_final_fantasy_sephiroth_e6a4b544/sephiroth.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_seph2016_35f91241/seph2016.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_seph2016_2c3745d9/seph2016.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_seph2016_08f79d2c/seph2016.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kam_lanaut_fc823b46/kam\'lanaut.def'] = 'Faction Final Fantasy',
	-- Vocaloid Vengeance pack (2026-09-24)
	['bulk_vocaloid_d4miku_97d4e56f/d4miku.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_dark_rin_03d5145b/dark_rin.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_lily_e6fc22c4/lily.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_kaito_a9288fa4/kaito.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_kizna_0bef895b/kizna.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_luka_6d2fa0e0/luka.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_takoruka_c6c53d48/takoruka.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_tda_maid_miku_0a5a1aa5/tda_maid_miku.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_jam band_f9a469be/jam band.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_mikummd_t_e37ffa56/mikummd_t.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_miku-m_0429b2fa/miku-m-1.0.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_miku-m_7e788de2/miku-m-win.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_ju_miku_12e7534d/ju_miku.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_miku_light_f1f0beaf/miku_light.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_miku1_5fa7b40c/miku1.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_mikusi_717c52cd/mikusi.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_mikummd_l_e2fb589f/mikummd_l.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_miku_3f25ab0e/miku.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_hatsune_miku_784b5e69/hatsune_miku.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_miku_kfm_e3b554dd/miku_kfm.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_miku5_08_13_2010_e668964d/miku5.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_miku _miku_n pop__157e5490/miku (miku\'n pop).def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_element_miku_f752ce32/element_miku.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_miku_92483e52/miku.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_mikudayo_fd96ae63/mikudayo.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_rin _miku_n pop__78193aeb/rin (miku\'n pop).def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_backup_fb5018b2/rin (miku\'n pop).def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_yowane haku_9d1a3c0a/yowane haku(ai).def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_yowane haku_f527ebeb/yowane haku.def'] = 'Faction Vocaloid Vengeance',
	['bulk_vocaloid_a-zundamon_e37478d8/a-zundamon.def'] = 'Faction Vocaloid Vengeance',
	-- every Kung Fu Man edit joins KFM (2026-09-24)
	['kung fu man-kofm'] = 'Faction KFM',
	['misinterpreted kfm'] = 'Faction KFM',
	['dgkung fu man/dgkung fu man.def'] = 'Faction KFM',
	['test-kfm-z-v2/test-kfm-z-v2.def'] = 'Faction KFM',
	['wr-kfm/wr-kfm.def'] = 'Faction KFM',
	['kfm_cap/kfm_cap.def'] = 'Faction KFM',
	['kungfuman/kungfuman.def'] = 'Faction KFM',
	['sf4kfm/sf4kfm.def'] = 'Faction KFM',
	['kfm_pack_!yfm_21_01_04/!yfm.def'] = 'Faction KFM',
	['kfm_pack_!kfm_mhx/!kfm_mhx.def'] = 'Faction KFM',
	['kfm_pack_0kfm(tommy-gun)/0kfm.def'] = 'Faction KFM',
	['kfm_pack_250_kfmf/250_kfmf.def'] = 'Faction KFM',
	['kfm_pack_8-chance_man/8-chance man.def'] = 'Faction KFM',
	['kfm_pack_808080/808080.def'] = 'Faction KFM',
	['kfm_pack_ac_kfm/kfm.def'] = 'Faction KFM',
	['kfm_pack_boss_kung_fu_man/boss kung fu man.def'] = 'Faction KFM',
	['kfm_pack_butter_for_the_win/butterguy.def'] = 'Faction KFM',
	['kfm_pack_ccs_kfm/ccs kfm.def'] = 'Faction KFM',
	['kfm_pack_ccs_kfm_the_final_update!!!!/ccs kfm.def'] = 'Faction KFM',
	['kfm_pack_chain-combo_man/cc-man.def'] = 'Faction KFM',
	['kfm_pack_cheaperdankfm/cheaperdankfm.def'] = 'Faction KFM',
	['kfm_pack_cheapestdankfm/cheapestdankfm.def'] = 'Faction KFM',
	['kfm_pack_christmaskungfuman/christmaskungfuman.def'] = 'Faction KFM',
	['kfm_pack_clone/clone.def'] = 'Faction KFM',
	['kfm_pack_colonel_mcpain/colmcpain.def'] = 'Faction KFM',
	['kfm_pack_crap_hj/hj.def'] = 'Faction KFM',
	['kfm_pack_cvs_fiunn_v1.0/cvsfiunn.def'] = 'Faction KFM',
	['kfm_pack_cvs_fiunn_v1.0_(1)/cvsfiunn.def'] = 'Faction KFM',
	['kfm_pack_dhmalamltotalt2_charactertemplate3_kfmsuperdemogenstylededit/charactertemplate3_kfmsuperdemogenstylededit.def'] = 'Faction KFM',
	['kfm_pack_dar_kfm/dar_kfm.def'] = 'Faction KFM',
	['kfm_pack_dankfm/dankfm.def'] = 'Faction KFM',
	['kfm_pack_dank_fu_man/dank_fu_man.def'] = 'Faction KFM',
	['kfm_pack_darker_than_black_kfm/darker than black kung fu man.def'] = 'Faction KFM',
	['kfm_pack_devil_kung_fu_man/devil-kfm.def'] = 'Faction KFM',
	['kfm_pack_dirty_kung_fu_man/dkfm.def'] = 'Faction KFM',
	['kfm_pack_divekickkfm/divekickkfm.def'] = 'Faction KFM',
	['kfm_pack_divinekfg/divinekfg.def'] = 'Faction KFM',
	['kfm_pack_divinekfm/divinekfm.def'] = 'Faction KFM',
	['kfm_pack_dkfm0/dkfm0.def'] = 'Faction KFM',
	['kfm_pack_esp-k/esp-k.def'] = 'Faction KFM',
	['kfm_pack_elecbyte2001kfm/kfm.def'] = 'Faction KFM',
	['kfm_pack_f1_kung_fu_man/f1kfm.def'] = 'Faction KFM',
	['kfm_pack_gfm/kfm.def'] = 'Faction KFM',
	['kfm_pack_gt_battle_mix_template_by_gtfoxn6y/kfm.def'] = 'Faction KFM',
	['kfm_pack_gatotyuman/gatotyuman.def'] = 'Faction KFM',
	['kfm_pack_ghost_kung_fu_man/ghost kung fu man.def'] = 'Faction KFM',
	['kfm_pack_giudizio_finale/giudizio_finale.def'] = 'Faction KFM',
	['kfm_pack_guardian_kfm/guardian kfm.def'] = 'Faction KFM',
	['kfm_pack_hyper_omega_cheapest_lolol_sector8_kung_fu_man/hyper omega cheapest lolol sector8 kung fu man.def'] = 'Faction KFM',
	['kfm_pack_holy_kfm/holy kfm!.def'] = 'Faction KFM',
	['kfm_pack_infinite_jump_kung_fu_man/infinite jump kung fu man.def'] = 'Faction KFM',
	['kfm_pack_invisible_kung_fu_man/invisible kung fu man.def'] = 'Faction KFM',
	['kfm_pack_jjbatemplate_edit/jjbatemplate_edit.def'] = 'Faction KFM',
	['kfm_pack_jeremiah_kung-fu-man/jeremiah kung-fu-man.def'] = 'Faction KFM',
	['kfm_pack_kfgirl-kofm/kfgirl-kofm.def'] = 'Faction KFM',
	['kfm_pack_kfkfkfkfkfm_by/kfkfkfkfkfm.def'] = 'Faction KFM',
	['kfm_pack_kfm-master/kfm-master.def'] = 'Faction KFM',
	['kfm_pack_kfm-master_(1)/kfm-master.def'] = 'Faction KFM',
	['kfm_pack_kfm-ultra_v1.6_director\'s_cut/kfm-ultra.def'] = 'Faction KFM',
	['kfm_pack_kfm-xi/kfm-xi.def'] = 'Faction KFM',
	['kfm_pack_kfm-type-c/kfm-type-c.def'] = 'Faction KFM',
	['kfm_pack_kfmsvkedit/kfmsvkedit.def'] = 'Faction KFM',
	['kfm_pack_kfmsvkedit2/kfmsvkedit2.def'] = 'Faction KFM',
	['kfm_pack_kfm_(gtorangemugen_style_template)/kfm.def'] = 'Faction KFM',
	['kfm_pack_kfm_(orangestrikers_mugen_style_template)/kfm.def'] = 'Faction KFM',
	['kfm_pack_kfm_(snk_arrange)_by_mage_[playable]-ramon_garcia/kfm_k.def'] = 'Faction KFM',
	['kfm_pack_kfm_09/kfm.def'] = 'Faction KFM',
	['kfm_pack_kfm_irc_(kyouakufightman)/kfm_irc.def'] = 'Faction KFM',
	['kfm_pack_kfm_iwbtg/kfm_iwbtg.def'] = 'Faction KFM',
	['kfm_pack_kfm_mvc_template_edit/kfm.def'] = 'Faction KFM',
	['kfm_pack_kfm_father/kfm father.def'] = 'Faction KFM',
	['kfm_pack_kfma1/kfma4a.def'] = 'Faction KFM',
	['kfm_pack_kfma3/kfma4a3.def'] = 'Faction KFM',
	['kfm_pack_kk_deuces_wild_snapshot_3mt8xbl/kk deuces wild.def'] = 'Faction KFM',
	['kfm_pack_kofkfm_ver1.5/kofkfm.def'] = 'Faction KFM',
	['kfm_pack_kasufire-man/kasufire-man.def'] = 'Faction KFM',
	['kfm_pack_kfm21/kfm21.def'] = 'Faction KFM',
	['kfm_pack_kian_type-kfm/kian type-kfm.def'] = 'Faction KFM',
	['kfm_pack_killer_kung_fu_man/killerkfm.def'] = 'Faction KFM',
	['kfm_pack_kingfancyman/kingfancyman.def'] = 'Faction KFM',
	['kfm_pack_kingfancyman2/kingfancyman2.def'] = 'Faction KFM',
	['kfm_pack_kingfancyman2_(1)/kingfancyman2.def'] = 'Faction KFM',
	['kfm_pack_kung-fu_men/kung-fu men.def'] = 'Faction KFM',
	['kfm_pack_kungfumanreturn/kungfuman_return.def'] = 'Faction KFM',
	['kfm_pack_kungfumanreturn_(1)/kungfuman_return.def'] = 'Faction KFM',
	['kfm_pack_kungfumanspecial/kfmsp.def'] = 'Faction KFM',
	['kfm_pack_kungfuman_return/kungfuman_return.def'] = 'Faction KFM',
	['kfm_pack_kungfustick/kungfustick.def'] = 'Faction KFM',
	['kfm_pack_kung_fu_girl(jg_edition)/kung fu girl.def'] = 'Faction KFM',
	['kfm_pack_kung_fu_head/kfh.def'] = 'Faction KFM',
	['kfm_pack_kung_fu_man_snk_style_by_mage_[with_ai_by_an_unknown_author.]/kfm_k.def'] = 'Faction KFM',
	['kfm_pack_kung_fu_paragas/kung fu paragas.def'] = 'Faction KFM',
	['kfm_pack_lm/lm.def'] = 'Faction KFM',
	['kfm_pack_lakitu_kfm/lakitu kfm.def'] = 'Faction KFM',
	['kfm_pack_legend_kfm/legend kfm.def'] = 'Faction KFM',
	['kfm_pack_limit_kfm/limit kfm.def'] = 'Faction KFM',
	['kfm_pack_lkfmoo/lkfm oo.def'] = 'Faction KFM',
	['kfm_pack_mekakfm/mekakfm.def'] = 'Faction KFM',
	['kfm_pack_mpkfm/mpkfm.def'] = 'Faction KFM',
	['kfm_pack_master_of_arts/master of arts.def'] = 'Faction KFM',
	['kfm_pack_mathu_oka_man/mathu oka man.def'] = 'Faction KFM',
	['kfm_pack_metal_kfm/metal kfm.def'] = 'Faction KFM',
	['kfm_pack_misinterpreted_kfm/misinterpreted kfm.def'] = 'Faction KFM',
	['kfm_pack_monst_style_kfm_by_y77/kfm-monsterstrike.def'] = 'Faction KFM',
	['kfm_pack_mre/mre.def'] = 'Faction KFM',
	['kfm_pack_nuga/rkfm.def'] = 'Faction KFM',
	['kfm_pack_nkfmen/nkfmen.def'] = 'Faction KFM',
	['kfm_pack_orochi-kfm/orochi-kfm.def'] = 'Faction KFM',
	['kfm_pack_quick_man/quick man.def'] = 'Faction KFM',
	['kfm_pack_r-man/r-man.def'] = 'Faction KFM',
	['kfm_pack_recover_man/recover_man.def'] = 'Faction KFM',
	['kfm_pack_seboll100_kfm/shin evil burning orochi lolol lv 60 kfm.def'] = 'Faction KFM',
	['kfm_pack_shin_crazy_burning_arien_lolol_chaos_type_kfm/shin crazy burning arien lolol chaos type kfm.def'] = 'Faction KFM',
	['kfm_pack_ss-kfm_2nd/ss-kfm_2nd.def'] = 'Faction KFM',
	['kfm_pack_szm/szm.def'] = 'Faction KFM',
	['kfm_pack_super_fast_kfm/super fast kfm.def'] = 'Faction KFM',
	['kfm_pack_super_kung_fu_man_version_1.0/skfm.def'] = 'Faction KFM',
	['kfm_pack_symbiote_kung_fu_man/symbiotekfm.def'] = 'Faction KFM',
	['kfm_pack_tcomg/tcomg.def'] = 'Faction KFM',
	['kfm_pack_teleport_man/teleport_man.def'] = 'Faction KFM',
	['kfm_pack_the_kung_fu_man_[update]/the kung fu man.def'] = 'Faction KFM',
	['kfm_pack_ultimate_humans/ultimate_humans.def'] = 'Faction KFM',
	['kfm_pack_urusaikfm/urusaikfm.def'] = 'Faction KFM',
	['kfm_pack_useless_man/useless_man.def'] = 'Faction KFM',
	['kfm_pack_valentine_man/valentine man.def'] = 'Faction KFM',
	['kfm_pack_wakamoto_man_m.u.g.e.n/wmm.def'] = 'Faction KFM',
	['kfm_pack_weird_fu_man_remake/kfm.def'] = 'Faction KFM',
	['kfm_pack_whiteness_of_g_kfm/whiteness of g kfm.def'] = 'Faction KFM',
	['kfm_pack_why/why.def'] = 'Faction KFM',
	['kfm_pack_youtube_poop(ytp)_kung_fu_man/ytpkfm.def'] = 'Faction KFM',
	['kfm_pack_kfmc/kfmc.def'] = 'Faction KFM',
	['kfm_pack_another_kfm/another_kfm.def'] = 'Faction KFM',
	['kfm_pack_another_kfm_(1)/another_kfm.def'] = 'Faction KFM',
	['kfm_pack_arcadekfm/arcadekfm.def'] = 'Faction KFM',
	['kfm_pack_bender/bender.def'] = 'Faction KFM',
	['kfm_pack_blackness_of_g_kfm/blackness of g kfm.def'] = 'Faction KFM',
	['kfm_pack_brokken/brokken.def'] = 'Faction KFM',
	['kfm_pack_chuukfm/chuukfm.def'] = 'Faction KFM',
	['kfm_pack_colmcpain/colmcpain.def'] = 'Faction KFM',
	['kfm_pack_deadkfm/deadkfm.def'] = 'Faction KFM',
	['kfm_pack_dkfm/dkfm.def'] = 'Faction KFM',
	['kfm_pack_ekfm/ekfm.def'] = 'Faction KFM',
	['kfm_pack_ekfm2009/ekfm2009.def'] = 'Faction KFM',
	['kfm_pack_evil/evil.def'] = 'Faction KFM',
	['kfm_pack_flgss/flgm.def'] = 'Faction KFM',
	['kfm_pack_giantkfm/giantkfm.def'] = 'Faction KFM',
	['kfm_pack_gpm/gpm.def'] = 'Faction KFM',
	['kfm_pack_gpm_(1)/gpm.def'] = 'Faction KFM',
	['kfm_pack_infinitekfm/infinitekfm.def'] = 'Faction KFM',
	['kfm_pack_jfmjan2022/jfm.def'] = 'Faction KFM',
	['kfm_pack_jjkfm/jjkfm.def'] = 'Faction KFM',
	['kfm_pack_karate/karate.def'] = 'Faction KFM',
	['kfm_pack_kf-hippie/kf-hippie.def'] = 'Faction KFM',
	['kfm_pack_kfblanka/kfblanka.def'] = 'Faction KFM',
	['kfm_pack_kfd/kfd.def'] = 'Faction KFM',
	['kfm_pack_kfdedede/kfdedede.def'] = 'Faction KFM',
	['kfm_pack_kff/kff.def'] = 'Faction KFM',
	['kfm_pack_kfg/kfg.def'] = 'Faction KFM',
	['kfm_pack_kfh/kfh.def'] = 'Faction KFM',
	['kfm_pack_kfm/kfm.def'] = 'Faction KFM',
	['kfm_pack_kfm\'09/kfm\'09.def'] = 'Faction KFM',
	['kfm_pack_kfm(the_brutal)/kfm.def'] = 'Faction KFM',
	['kfm_pack_kfm-hyperrion/kfm-hyperrion.def'] = 'Faction KFM',
	['kfm_pack_kfm-ultra_2/kfm-ultra.def'] = 'Faction KFM',
	['kfm_pack_kfm1/kfm1.def'] = 'Faction KFM',
	['kfm_pack_kfm11/kfm11.def'] = 'Faction KFM',
	['kfm_pack_kfm12/kfm12.def'] = 'Faction KFM',
	['kfm_pack_kfm14/kfm14.def'] = 'Faction KFM',
	['kfm_pack_kfm16/kfm16.def'] = 'Faction KFM',
	['kfm_pack_kfm17/kfm17.def'] = 'Faction KFM',
	['kfm_pack_kfm22/kfm22.def'] = 'Faction KFM',
	['kfm_pack_kfm22_(1)/kfm22.def'] = 'Faction KFM',
	['kfm_pack_kfm23/kfm23.def'] = 'Faction KFM',
	['kfm_pack_kfm3/kfm3.def'] = 'Faction KFM',
	['kfm_pack_kfm4/kfm4.def'] = 'Faction KFM',
	['kfm_pack_kfm63/kfm63.def'] = 'Faction KFM',
	['kfm_pack_kfm7/kfm7.def'] = 'Faction KFM',
	['kfm_pack_kfm9/kfm9.def'] = 'Faction KFM',
	['kfm_pack_kfmtaw/kfmtaw.def'] = 'Faction KFM',
	['kfm_pack_kfm_18/kfm_18.def'] = 'Faction KFM',
	['kfm_pack_kfm_3d3aa/kfm.def'] = 'Faction KFM',
	['kfm_pack_kfm_9eab0/kfm.def'] = 'Faction KFM',
	['kfm_pack_kfm_a_tgm/kfm_a_tgm.def'] = 'Faction KFM',
	['kfm_pack_kfm_altz/kfm_altz.def'] = 'Faction KFM',
	['kfm_pack_kfm_arc/kfm_arc.def'] = 'Faction KFM',
	['kfm_pack_kfm_g/kfm_g.def'] = 'Faction KFM',
	['kfm_pack_kfm_kof/kfm_kof.def'] = 'Faction KFM',
	['kfm_pack_kfm_neo/kfm_neo.def'] = 'Faction KFM',
	['kfm_pack_kfm_type-s/kfm type-s.def'] = 'Faction KFM',
	['kfm_pack_kfm_at/kfm_at.def'] = 'Faction KFM',
	['kfm_pack_kfm_bt/kfm_bt.def'] = 'Faction KFM',
	['kfm_pack_kfm_cap/kfm_cap.def'] = 'Faction KFM',
	['kfm_pack_kfm_ht/kfm_ht.def'] = 'Faction KFM',
	['kfm_pack_kfm_old/kfm_old.def'] = 'Faction KFM',
	['kfm_pack_kfm_sadclaps-20230426t174112z-001/kfm_sadclaps.def'] = 'Faction KFM',
	['kfm_pack_kfm_version_2/kfm1.def'] = 'Faction KFM',
	['kfm_pack_kfm_wx_v1.01/kfm_wx.def'] = 'Faction KFM',
	['kfm_pack_kfmario/kfmario.def'] = 'Faction KFM',
	['kfm_pack_kfmaster/kfmaster.def'] = 'Faction KFM',
	['kfm_pack_kfmbison/kfmbison.def'] = 'Faction KFM',
	['kfm_pack_kfmex/kfmex.def'] = 'Faction KFM',
	['kfm_pack_kfmgb/kfmgb.def'] = 'Faction KFM',
	['kfm_pack_kfmm/kfmm.def'] = 'Faction KFM',
	['kfm_pack_kfmsp/kfmsp.def'] = 'Faction KFM',
	['kfm_pack_kfmt1/kfmt1.def'] = 'Faction KFM',
	['kfm_pack_kfmtaw/kfmtaw.def'] = 'Faction KFM',
	['kfm_pack_kfmw21_v3/kung fu man wlanmaniax 2021.def'] = 'Faction KFM',
	['kfm_pack_kfmx/kfmx.def'] = 'Faction KFM',
	['kfm_pack_kfr/kfr.def'] = 'Faction KFM',
	['kfm_pack_kfr_787ef/kfr.def'] = 'Faction KFM',
	['kfm_pack_koryufm/koryufm.def'] = 'Faction KFM',
	['kfm_pack_kurofuku/kurofuku.def'] = 'Faction KFM',
	['kfm_pack_mako/mako.def'] = 'Faction KFM',
	['kfm_pack_mkfm/mkfm.def'] = 'Faction KFM',
	['kfm_pack_mm_fiunn/mm_fiunn.def'] = 'Faction KFM',
	['kfm_pack_moon/moon.def'] = 'Faction KFM',
	['kfm_pack_nam_uf_gnuk_is_facing_the_right_way_youre_the_one_thats_backwards/rkfm.def'] = 'Faction KFM',
	['kfm_pack_oakfm/oakfm.def'] = 'Faction KFM',
	['kfm_pack_orbkfm/orbkfm.def'] = 'Faction KFM',
	['kfm_pack_ped/ped.def'] = 'Faction KFM',
	['kfm_pack_pkfm/pkfm.def'] = 'Faction KFM',
	['kfm_pack_robokangfu/robokangfu.def'] = 'Faction KFM',
	['kfm_pack_sf2kfm/sf2kfm.def'] = 'Faction KFM',
	['kfm_pack_sillykfm2009/sillykfm.def'] = 'Faction KFM',
	['kfm_pack_the-kung-fu-man/the-kung fu man.def'] = 'Faction KFM',
	['kfm_pack_tushou/tushou.def'] = 'Faction KFM',
	['kfm_pack_wfa/wfa.def'] = 'Faction KFM',
	['kfm_pack_wmm/wmm.def'] = 'Faction KFM',
	['kfm_pack_xkfm_v1.2/xkfm.def'] = 'Faction KFM',
	['kfm_pack_yanagi-no_v4.0/yanagi-no.def'] = 'Faction KFM',
	['kfm_pack_zoom/zoom.def'] = 'Faction KFM',
	-- TMNT pack (2026-09-24)
	['bulk_tmnt_alopex_8483335f/alopex.def'] = 'Faction TMNT',
	['bulk_tmnt_dsj_april_5acbb37d/dsj_april.def'] = 'Faction TMNT',
	['bulk_tmnt_tmntsraprildx_937180f7/tmntsraprildx.def'] = 'Faction TMNT',
	['bulk_tmnt_dsjs_mitsu_85e16b86/dsjs_mitsu.def'] = 'Faction TMNT',
	['bulk_tmnt_bebop_11378b10/bebop.def'] = 'Faction TMNT',
	['bulk_tmnt_donatellotmnt_eb37b6b0/donatellotmnt.def'] = 'Faction TMNT',
	['bulk_tmnt_foot_f3f8408b/foot.def'] = 'Faction TMNT',
	['bulk_tmnt_leo_c68802c2/leo.def'] = 'Faction TMNT',
	['bulk_tmnt_leo_v1_1_ by matti_44830a78/leo.def'] = 'Faction TMNT',
	['bulk_tmnt_leonardo_6d0a3633/leonardo.def'] = 'Faction TMNT',
	['bulk_tmnt_leonardosi_de4cf92e/leonardosi.def'] = 'Faction TMNT',
	['bulk_tmnt_leo_by_barany1027_501f3b9c/leo_by_barany1027.def'] = 'Faction TMNT',
	['bulk_tmnt_leonardo2_65b60b30/leonardo2.def'] = 'Faction TMNT',
	['bulk_tmnt_mvc_leonardo_985b8acc/mvc_leonardo.def'] = 'Faction TMNT',
	['bulk_tmnt_leo_56d313d1/leo.def'] = 'Faction TMNT',
	['bulk_tmnt_leonardo_91a11231/leonardo.def'] = 'Faction TMNT',
	['bulk_tmnt_michelangelo_7c13d38b/tatoruzu.def'] = 'Faction TMNT',
	['bulk_tmnt_mike_d095ced9/mike.def'] = 'Faction TMNT',
	['bulk_tmnt_tfgaf-tmnt2007-mrbigrobot_41c2d14a/tfgaf-tmnt2007-mrbigrobot.def'] = 'Faction TMNT',
	['bulk_tmnt_tmnt_nightfall_e3d1b1f0/tmnt_nightfall.def'] = 'Faction TMNT',
	['bulk_tmnt_rocksteady_v1_3_1_ by matti_e69d18c7/rocksteady.def'] = 'Faction TMNT',
	['bulk_tmnt_shredder_fccc737e/dcat_shredder.def'] = 'Faction TMNT',
	['bulk_tmnt_shredder_v2_5_ by matti_2a128f96/shred.def'] = 'Faction TMNT',
	['bulk_tmnt_slash_a495f5ff/slash.def'] = 'Faction TMNT',
	['bulk_tmnt_supershredder_86f7ec07/supershredder.def'] = 'Faction TMNT',
	['bulk_tmnt_wwyatsraph_sebastian_1c0da5d7/wwyatsraph_sebastian.def'] = 'Faction TMNT',
	-- Nintendo pack (2026-09-24)
	['bulk_super_mario_bros_9voltsi_6294f7e6/9voltsi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_aidrmario_55b1f4f5/aidrmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_andersonkenya1 v2_12b02933/andersonkenya1 v2.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros__babymario_.def_3f440fdb/(babymario).def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_baby luigi 2.def_0c678bbf/baby luigi 2.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_baby mario.def_8597989e/baby mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_beanie_e0a4032e/beanie.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_b-mario_1fa42071/b-mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_birdo_c7ca1b3f/birdo.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_birdo_fa8b9afc/birdo.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_black yoshi_15a64c84/black yoshi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_black yoshi_fc236b25/metalyoshi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_deviling_a6852a69/deviling.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_boomboxer_5498d311/boomboxer.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_boshi_7933a6e9/boshi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_bowser_59a5b06f/bowser.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_giga bowser _rivals of aether__5f6325c1/giga bowser (rivals of aether).def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_newbowserjr_10f45b44/newbowserjr.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_bowsette_5eb553ad/bowsette.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_countcannoli_68256d05/countcannoli.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_daisy_c2a2eda6/daisy.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_paperdaisy_v2_079d949d/paperdaisy_v2.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_brd_dantek_696aec44/brd_dantek.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_dark ashley_dadd92ba/dark ashley.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_dark daisy_19c5c382/dark daisy.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_dark_mario_bfe94840/dark_mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_drmario_5837837e/drmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_doctormario_ecfca444/doctormario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_golf master ella_d453ce02/golf master ella.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_eskimo koopa_8d4d90a7/eskimo koopa.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_evil supermario64_9878d179/evil supermario64.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_fennel_b5a9b756/fennel.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_fire_mario_ed5a5f1c/fatmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_fire_mario_03d3eabd/fire_mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_fire_mario_b4e9863f/metal mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_firemarioex_dba24588/firemarioex.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_fortran_da3aa822/fortran.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_godmodderlavenderjuice_af63e487/godmodderlavenderjuice.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_goom_5a20d579/goom.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_goomba_f4b54500/goomba.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_grand dad_cbdd0040/grand dad.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_luigi_0251ae75/luigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_greaterwario_271cbfb3/greaterwario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_hallo_weerio__553c45c0/hallo\'weerio!.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_hammerbro_90933ff1/hammerbro.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_hammermario_ca1fd0a8/hammermario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_hnk_90467d7b/smallkamek.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_hnk_f92106bb/hnk.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_hnk_a433261b/kamek.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_hnk_1ea18e45/rival.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_hdimentio_e51d3d0e/hdimentio.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_ice_luigi_f3a35218/ice_luigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_jrtroopa_33d81b0f/jrtroopa.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_killer_wario_00a58198/killer_wario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kamek_34b091e1/smallkamek.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kamek_b06b98ac/hnk.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kamek_02018b9e/kamek.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kamek_bca545f8/rival.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_king_boo_84915b9b/king_boo.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_pesky hammer bro_aa766fb6/pesky hammer bro.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_pesky koopa_a3db066a/pesky koopa.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_koopa_djhannibal_e21ff8af/koopa_djhannibal.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_ssbbkoopa_50fed133/ssbbkoopa.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_tfg-pm-koopatrol_965e1081/tfg-pm-koopatrol.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kfmario_f1b1ce7e/kfmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_sloppio_f3421dd4/sloppio.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_bootlegrussianluigi_a6e9dd4c/bootlegrussianluigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_dctemplate_bdd95012/dctemplate.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_luigismb_999ae72f/luigismb.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toadswap luigi_fb09f31c/playerundertaletemplate.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_luigi_fe6d555a/luigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_superluigi3_4bd32d46/superluigi3.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_luuiiiigiiiiiiiiiiii__a4ad7ad8/luuiiiigiiiiiiiiiiii!!.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_rated_m_for_mario_092a304f/rated_m_for_mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_m4ri0000_8dfcf885/m4ri0000.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mx_8d297bb9/eoh.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_yaridovich_4eedd9ca/yaridovich.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toon_weegee_51add305/toon_weegee.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mallow_da683128/mallow.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_supermamaluigi_10bba9db/supermamaluigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_0smw_cc7f2e6d/0smw.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_backup_d01c6098/supermariowar.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_b_f_mario_a8c008ed/b_f_mario(game).def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_b_f_mario_67e723c9/b_f_mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_cdimario_59a511f5/cdimario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_captainlou_947a36bf/captainlou.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_gb_mario_f84520b2/gb_mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_luigi_78eae733/luigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mariossf_dd1627f5/mariossf.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mtm_malleo_8b8ca099/mtm_malleo.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_bowser_ddab71d6/bowser.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario_a4d32d42/mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario_03991809/mariomariomario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario_d2473dfd/mario_kf.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario_d1423ca3/mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_sma4mario_fc69c4f1/sma4mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario70v2_dd9dd1c2/mario70v2.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mariomvdk_1b26b3e2/mariomvdk.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mariosmb_5c670457/mariosmb.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mariossb_4b72bc54/mariossb.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario_combat_deluxe_mario_c651f1c8/mario_combat_deluxe_mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario_smwtv_add47d34/mario_smwtv.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario_smwtv_1e8bf897/mario_smwtv10.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario_sebastian_edef9c9a/mario_sebastian.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario_tb_e91e5fd2/mario_tb.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_nsmbw_57fe64ec/nsmbw.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_nteam_4536f7b3/nteam.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_paper mario_30cd6d0f/paper mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_pesky plumber super mario_c3c7d395/pesky plumber super mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_powerstarmario_beta__b872e8e8/powerstarmario(beta).def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_sdsmmm_0c32b881/sdsmmm.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_sm3_ec4d008c/sm3.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_stickmario_camrensp_173e8244/stickmario_camrensp.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_supermario64_56491ab2/supermario64.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_supermariokart_c31f0ba6/supermariokart.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_undertoad mario_972e9c5e/undertoad mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_vga_mario_0d13c375/vga_mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_vga_mario_54630f73/vga_mario_tag.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_vgf_mario_revamp_ff42daca/vgf!mario_revamp.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_wh2_mario_1cf73656/wh2_mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_crew_fc6d3441/crew.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_dkmario_a3e35b43/dkmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kofmariotwins_8a2697df/kofmariotwins.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_luigi_paperjam_248b3b76/luigi_paperjam.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_marioax_8f8cd110/marioax.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mariobros_aca2c1a9/mariobros.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_smb2mario_3d46b81d/smb2mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_smwmario3_24e23cc9/smwmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_ssbbmario - copy_71187284/ssbbmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_ssbbmario_9191bdfb/ssbbmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_charamario_super_mario_kun_2de2a494/charamario_super_mario_kun.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_cvtwmario_722625db/cvtwmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_marioblicky_01c08c71/marioblicky.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_marioii_71536be9/marioii.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario.exe_23c1f78a/mario.exe.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_sml mario punching bag_36454f01/sml mario punching bag.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_05.06.2021_3.21.48_am__1abccdae/svcryu.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_svcryu_d10aa84d/svcryu.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_newmasao_16602e57/newmasao.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_nesmario_3cc6350f/nesmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mechamario_df5e02b8/mechamario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_meka_goomba_d0e84b7a/meka_goomba.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_metal mario_fb928feb/metal mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_metal_mario_kart_0fd7af31/metal_mario_kart.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_metalyoshi_fef2f0f7/metalyoshi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_midbus_dcdf9c41/midbus.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_muncher_9211165e/muncher.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kf-x_mmario_d30da643/kf-x_mmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kf-x_mmario_228a27ad/kf-x_mmario_gkai.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_myrio_60aa1851/myrio.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_neilella_sebastian_10981dfc/neilella_sebastian.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_peach_c58d6ec3/peach.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_pow-mario_66f14a89/pow-mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mvc_peach_14f2842e/mvc_peach.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_p_shroob_0e19d978/p_shroob.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_r-luigi_20d2a2ae/r-luigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_rescuesquadtoad_e0859a94/rescuesquadtoad.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_retarded-luigi_bc11bc39/retarded-luigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_retarded goomba_42f091d5/retarded goomba.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_retarded kof mario_54f4c184/retarded kof mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_retarded mario_5ce7edd4/retarded mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_backup_a0fb756f/eq-princesspeach.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_retarded princess peach_87231412/retarded princess peach.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_retardedsupermario_a2b70af9/retardedsupermario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_retarded waluigi_a96097d3/retarded waluigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_retarded wario_1773674b/retarded wario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_retarded wario_148455dd/wario2.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_pic1_4d6d7e80/pic1.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_riuky_51464c9f/riuky.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_riuky_f5bf944d/shy guy.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_rozadv_7cc39a53/rozadv.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_smg4 853_78c388c9/smg4 853.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_beyond ssjb princess peach_417c7a1b/beyond ssjb princess peach.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_shadoo_9fe74292/noob_saibot.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_shy_guy_98425ce0/shyguy.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_somari_e57f6420/somari.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_speedrunner mario_20b63864/yahoo.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_spondagemario_1124a6f7/spondagemario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_spondagemario_179dbe2b/spondagemario10.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_sr pelo_a8999504/sr pelo.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_superlitmario_sebastian_94593e8a/superlitmario_sebastian.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_findluigi_47fb1ce1/findluigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_smw - bala gigante_4203fb86/smw - bala gigante.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_donaldmario1_5462614a/donaldmario1.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_superbettersoulchild5_f9747048/superbettersoulchild5.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_symbmario_a49f0997/metal mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_symbmario_41abb21c/symbmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_symbiotemariomaker_876f71ea/symbiotemariomaker.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_symbiote minimario_f3a4a146/symbiote minimario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_symbiote minimario_6067a402/metmario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_symboite hotel mario_7ff27ba6/symboite hotel mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_talkingwonderflower_17d52e73/talkingwonderflower.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_the_master_e0765b9d/the_master.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_universe-waluigi_bfcda76b/universe-waluigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_tippi_292689bd/tippi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toad_33d01920/toad.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toad chibi_2e644017/toad chibi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toad_smw_b5081f5f/toad_smw.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toad_smw_839509e8/toad_smw10.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toad_f116557f/toad.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_tfv_8b4c6446/tfv.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toadette_smw_5a9cb95c/toadette_smw.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toadette_smw_006762ca/toadette_smw10.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kftoadette_f890eecd/kftoadette.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toy_mario_7e237379/toy luigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_toy_mario_c33628c8/toy mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_c-mario_1bbfd76b/c-mario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_nes waluigi_77672883/nes waluigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_waluigi_3a92f33e/waluigi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_mario _mario _ wario__a48085f2/mario (mario & wario).def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_rudytheclown_c7904275/rudytheclown.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_wario_5d217b0e/wario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_warioedit_171a7e1f/warioedit.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_kfwario_1adbc5a4/kfwario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_wario_181b5519/wario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_the wario appartion_e09b99df/the wario appartion.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_wart_06ceb028/wart.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_weegee_c9df7822/weegee.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_yosh__25d3ebba/yosh!.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_super_wario_2394a517/super_wario.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_yoshi_tb_fa4f967a/yoshi_tb.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_smw_yoshi_7f5b7c96/smw_yoshi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_smw_yoshi_adfe98d8/smw_yoshinoshoes.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_yoshi_975ad65f/yoshi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_yoshi__hokuto_hyakuretsu_ken_2b178616/yoshi.def'] = 'Faction Nintendo',
	['bulk_super_mario_bros_a boss_50b8dae9/a boss.def'] = 'Faction Nintendo',
	-- Sailor Moon pack (2026-09-24)
	['bulk_sailor_moon_castornpollux_6d87119a/castornpollux.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_g64sailormoon_612164cb/g64sailormoon.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_emerald_2ff42446/emerald.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_esmoon_70c60de1/esmoon.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_beryl_12d252d3/beryl.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_beryl_7893057b/beryl.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_rubeus_8b8a8357/rubeus.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_supersailorchibimoon_36fd5203/supersailorchibimoon.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_jupiter_9f5d3afe/jupiter.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailor_jupiter_6e4e91e1/sailor_jupiter.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailormars_84b0fe08/sailormars.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_mars_871fa625/mars.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailor_mercury_41fa9ff5/sailor_mercury.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailor_mercury_49fc8b68/sailor_mercury.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_mercury_45564872/mercury.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailor_moon_ea_hu_386d6bda/sailor_moon.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailor moon_8046f78b/sailor moon.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailor-moon_6a729a68/sailor-moon.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailormoon_886095dc/sailormoon.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailorneptune_aeabc33c/sailorneptune.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_hyper sailor neptune _hsm and mugenverse__2d93e007/hyper sailor neptune (hsm and mugenverse).def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailorpluto_871d760d/sailorpluto.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailorpluto_ef8edbb8/sailorpluto.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailor_saturn_acda13b6/sailor_saturn.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailorsaturn_e4cfc62d/sailorsaturn.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_sailor_venus_acbfb925/sailor_venus.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_venus_2475f0cf/venus.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_supersailormoon_462466c6/supersailormoon.def'] = 'Faction Sailor Moon',
	['bulk_sailor_moon_zoicite_9514e5da/zoi.def'] = 'Faction Sailor Moon',
	-- Final Fantasy pack (2026-09-24)
	['bulk_final_fantasy_ace_cf13df7e/ace.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_aerith_26548cc4/aeris0.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_agrias_1.0_5c676427/agrias_1.0.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_agrias_1.0_98a1c891/agrias_1.1.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy__ffxi_ mithra muskeeter akaneko_d65b1b1b/(ffxi) mithra muskeeter akaneko.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_akatsuki_cloud_2a9f798d/akatsuki_cloud.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy__ffxi_ hume prince aldo_39015015/(ffxi) hume prince aldo.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_amber_7024bbc1/amber.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ami_36d749bb/ami.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_areuhat_50380405/areuhat.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ff12 ashe_1090f9af/ff12 ashe.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_atomos_6e7152ee/atomos.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_auron_9d624dd4/auron.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ayame_4a87e36c/ayame.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_bahamut_67e1e687/bahamut.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_hd_bahamutsin_85839207/hd_bahamutsin.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_barbarccia_92a2d3d8/barbarccia.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_bartz_6c8439a3/bartz.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_benjamin_6dad67db/benjamin.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_g_blackmage_cdab44a1/g_blackmage.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_bom_8514648c/bom.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_caitsith_17dd0436/zackcc.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_celes_3ed86a46/celes.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_charlotte _holy knight__588c2c39/charlotte (holy knight).def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_chloe_46dcc51b/chloe.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_chocobo_e8de3b46/chocobo.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_christine_83764121/christine.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_cid_c9f1163e/cid.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_citra _ folka_79c59429/citra & folka.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_cloud_s_f84f3f19/cloud.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_cloud of darkness _dissidia__6066eba2/cloud of darkness (dissidia).def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_colibri_8eb7d009/colibri.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_colibri_c61c5c84/colibri_ai.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_cyan_14526801/cyan.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_dark mage_aaf43e82/dark mage.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_dietlinde_13cc369b/dietlinde.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_edgar_17c066d4/edgar.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_edge_22bceb08/edge.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_elfreeda_19149295/elfreeda.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_enuo_e40fe753/enuo.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_gogo_92852647/gogo.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_monk_16dca3ef/monk.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_firion_b3dd2569/firion.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_gadalar_946c874c/gadalar.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_galka_353f6e74/galka.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kefka god_29d4b07e/kefka god.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_gulool_53a402bf/gulool.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_hellhouse_81877d6b/hellhouse.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_high_redcat_f8cb15ec/high_redcat.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_hbc_24eb7776/hbc.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ishil_788e8fee/ishil.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_joachim_ef39afb0/joachim.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_cloud_a70ba472/cloud.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_sephiroth1_57b88761/sephiroth1.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kain_9f185e5e/kain.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kam_lanaut_fc823b46/kam\'lanaut.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_karakuri_8be77603/karakuri.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kefka2_07333576/kefka2.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kingbehinnmoth_ba4f225b/kingbehinnmoth.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kitanaininjya_82f49eb9/kitanaininjya.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kitanaininjya-t_0cd2564a/kitanaininjya-t.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kitanaitaru_ebb2e11c/kitanaitaru.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_kurokishineko_be90e625/kurokishineko.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_leon_cfff92f9/leon.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_leonoyne_479ee36a/leonoyne.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_leopold_78bea2fd/leopold.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ffs summon leviathan_db23e9c9/ffs summon leviathan.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_lightning_a0caaa75/lightning.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_locke_33d58377/locke.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy__ffxi_ tarutaru mage lutette_e38e866a/(ffxi) tarutaru mage lutette.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_matt_610694e8/maat.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_majic-pot_94abf88a/majic-pot.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_mandoraf-f_1a89662e/mandoraf-f.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_melusine_c889b64b/melusine.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_mog_6b39a81b/mog.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_montblanc_dcbc9ec5/montblanc.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_moogle_7d3a69c1/moogle.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_mumor_b5407655/mumor.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_naitou_39c83aec/naitou.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy__ffxi_ tarutaru white knight_00add39a/(ffxi) tarutaru white knight.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_naji_2e80b986/naji.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_necorsair_d86cb846/necorsair.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_apolyonnecrongiygas_40443ac5/apolyonnecrongiygas.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_nekomo_33966dde/kungfumithra.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_nekomo_5d46486f/nekomo.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_nekonin_056148d3/nekonin.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_noctis_1.0_ef6dfdf9/noctis_1.0.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ff5_omega_9fbf60bb/ff5_omega.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_onekonin_0fb6e6f1/onekonin.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ff12 penelo_0f167b9b/ff12 penelo.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_plague_f0ece515/plague.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_prishe_f82c01fc/prishe.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ffs summon ramuh_d06343e4/ffs summon ramuh.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_renna_2ac1ba97/renna.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_rikku_4565f81f/rikku.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ff8 rinoa heartilly_5df0cd44/ff8 rinoa heartilly.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_robber_crab_cc9cef81/robber_crab.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_rughadjeen_71093200/rughadjeen.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_rydia_05b9b109/rydia.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ryusan_dc0375ea/ryusan.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy__ffxi_ tarutaru lancer ryuttan_124b67c2/(ffxi) tarutaru lancer ryuttan.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_mash_7e2a4c2c/mash.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_sarutabarutamusou_e51a498d/sarutabarutamusou.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_mg-sephiroth_e451de76/mg-sephiroth.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_shadow ff_e22249b1/shadow ff.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_shantotto_ed7cebc3/shantotto.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ff5_shinryu_107f3a7e/ff5_shinryu.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_shiroma_05f739fa/shiroma.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_sieghard_47da4343/sieghard.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_sirohime_2af98c71/sirohime.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_sister_d861f9f5/sister.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_skaha_abf7bfec/skaha.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_squall_da9aa22d/squall.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_tidus_dcf54e51/tidus.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_zerotifa_bdbd59c5/zerotifa.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy__ffxi_ summon giant cocoachan_1097d67d/(ffxi) summon giant cocoachan.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_z4tonberry_b6ebb065/z4tonberry.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_elena_d13c37fa/elena.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ff7 turk reno_6e87b6a3/ff7 turk reno.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy__final-fantasy_ turk tseng_02ee568b/(final-fantasy) turk tseng.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ultimate-high-legs_5d7b2c8d/ultimate-high-legs.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_ultros_68b0823c/ultros.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_vaan_35ccae95/normalvaan.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_vaan_26311808/turbovaan.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_vv_2e9a1f40/vv.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_vivi_080f5026/vivi.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_warrior of light _dissidia__e5ba63ab/warrior of light (dissidia).def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_worm_2384b110/worm.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_yashtola_89f3c5fb/yashtola.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy__final-fantasy_ yazoo_fd85ab72/(final-fantasy) yazoo.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_yuna_432a2880/yuna.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_yunalesca_eec26ca7/yunalesca.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_zack_ff8e89da/zack.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_zeid_d9f2d138/zeid.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_zidane_983f434c/zidane.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_zile_e0f54608/zile.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_zombone_50d86885/zombone_1.0.def'] = 'Faction Final Fantasy',
	['bulk_final_fantasy_zombone_21b2920d/zombone_1.1.def'] = 'Faction Final Fantasy',
	-- bulk download packs get their own factions (2026-09-23)
	['bulk_bleach_aaroniero arruruerie_a30d5835/aaroniero arruruerie.def'] = 'Faction Bleach',
	['bulk_bleach_ashido kano_f2c84565/ashido kano.def'] = 'Faction Bleach',
	['bulk_bleach_bankai ichigo_c7f4c753/bankai ichigo.def'] = 'Faction Bleach',
	['bulk_bleach_byakuya kid_1b48131e/byakuya kid.def'] = 'Faction Bleach',
	['bulk_bleach_byakuya_e234f33c/byakuya.def'] = 'Faction Bleach',
	['bulk_bleach_captain shigekuni_a1c47041/captain shigekuni.def'] = 'Faction Bleach',
	['bulk_bleach_ganju_f43ec31c/ganju.def'] = 'Faction Bleach',
	['bulk_bleach_granfisher_18828d6b/granfisher.def'] = 'Faction Bleach',
	['bulk_bleach_grimmjow_fee2fda5/grimmjow.def'] = 'Faction Bleach',
	['bulk_bleach_halibel_1a4af7df/halibel.def'] = 'Faction Bleach',
	['bulk_bleach_hirako_b4535d5a/hirako.def'] = 'Faction Bleach',
	['bulk_bleach_hisagi_8bf50856/hisagi.def'] = 'Faction Bleach',
	['bulk_bleach_ichigo new mask_fdea41b0/ichigo new mask.def'] = 'Faction Bleach',
	['bulk_bleach_ichigo vastolorde_2ba6579f/ichigo vastolorde.def'] = 'Faction Bleach',
	['bulk_bleach_ichigo vizard_5fabb7c8/ichigo vizard.def'] = 'Faction Bleach',
	['bulk_bleach_ichigo_4fdedc33/ichigo.def'] = 'Faction Bleach',
	['bulk_bleach_kenpachi_c5449ce7/kenpachi.def'] = 'Faction Bleach',
	['bulk_bleach_kira izuru_36a5d34f/kira.def'] = 'Faction Bleach',
	['bulk_bleach_kyouraku_334d5ba7/kyouraku.def'] = 'Faction Bleach',
	['bulk_bleach_lilynette_225f3db1/lilynette.def'] = 'Faction Bleach',
	['bulk_bleach_mikkaku_aafd96a8/mikkaku.def'] = 'Faction Bleach',
	['bulk_bleach_nell_84358ca5/nell.def'] = 'Faction Bleach',
	['bulk_bleach_nemu_a2f1b6fb/nemu.def'] = 'Faction Bleach',
	['bulk_bleach_nnoitora_9d3123a2/nnoitora.def'] = 'Faction Bleach',
	['bulk_bleach_orihime_c1344beb/orihime.def'] = 'Faction Bleach',
	['bulk_bleach_rangiku_5592e6e0/rangiku.def'] = 'Faction Bleach',
	['bulk_bleach_sosuke_aizen_fe81fa8b/sosuke_aizen.def'] = 'Faction Bleach',
	['bulk_bleach_stark_e134f220/stark.def'] = 'Faction Bleach',
	['bulk_bleach_stormex-ogihci_99b3288d/stormex-ogihci.def'] = 'Faction Bleach',
	['bulk_bleach_super saiyan ichigo_3037f129/super saiyan ichigo.def'] = 'Faction Bleach',
	['bulk_bleach_super saiyan toshiro hitsugaya_4c34478f/super saiyan toshiro hitsugaya.def'] = 'Faction Bleach',
	['bulk_bleach_szayelaporro grantz_4a365cef/szayelaporro grantz.def'] = 'Faction Bleach',
	['bulk_bleach_tensa zangetsu_4eedece6/tensa zangetsu.def'] = 'Faction Bleach',
	['bulk_bleach_toshiro_hitsugaya_26f39530/toshiro_hitsugaya.def'] = 'Faction Bleach',
	['bulk_bleach_ulquiorra_c3f97297/ulquiorra.def'] = 'Faction Bleach',
	['bulk_bleach_unohana_5e1747a2/unohana.def'] = 'Faction Bleach',
	['bulk_bleach_urahara_52ef8043/urahara.def'] = 'Faction Bleach',
	['bulk_bleach_uryu_1af2ae5d/uryu.def'] = 'Faction Bleach',
	['bulk_bleach_yachiru_253b7dc7/yachiru.def'] = 'Faction Bleach',
	['bulk_bleach_yasutora_sado_6ec8c214/yasutora_sado.def'] = 'Faction Bleach',
	['bulk_bleach_yoruichi_61beaab0/yoruichi.def'] = 'Faction Bleach',
	['bulk_bleach_zommari_15deb87f/zommari.def'] = 'Faction Bleach',
	['bulk_contra__char_black_viper_27d190f7/black viper.def'] = 'Faction Contra',
	['bulk_contra__char_contra4_tank_0b8f7377/contra4_tank.def'] = 'Faction Contra',
	['bulk_contra__char_gomeramos_king_e6580bab/gomeramosking.def'] = 'Faction Contra',
	['bulk_contra__char_jappa_691da997/jappa.def'] = 'Faction Contra',
	['bulk_contra__char_rocket_tank_0225d580/rockettank.def'] = 'Faction Contra',
	['bulk_contra_babalu_destructoid_9b6de2bc/babalu.def'] = 'Faction Contra',
	['bulk_contra_billrizer rabbitgentleman_billrizer_60125a59/billrizer.def'] = 'Faction Contra',
	['bulk_contra_char_shadow_beast_kimkoh_contra-shadowbeastkimkou_54b26c7c/contra-shadowbeastkimkou.def'] = 'Faction Contra',
	['bulk_contra_chc-fang_chc-fang_657a7370/chc-fang.def'] = 'Faction Contra',
	['bulk_contra_contra3_cyberskeleton_287dc6e8/contra3_cyberskeleton.def'] = 'Faction Contra',
	['bulk_contra_contra_4e07f719/contra.def'] = 'Faction Contra',
	['bulk_contra_lance_bfc1740c/lance.def'] = 'Faction Contra',
	['bulk_fatal_fury_alfred_988ae852/alfred.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_andy_314534c9/andy.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_andy_9c7379d8/andy2.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_andy_ff3rb_9e073e4b/andy_ff3rb.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_axelhawk_28e66b1b/axelhawk.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_billy_c55d06a4/billy.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_bob wilson_ab4386c9/bob_pro.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_bob wilson_acce9bbe/bob wilson.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_bruno_dcd84d0a/bruno.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_bud2_0c7682ff/bud2.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_cheng_0d31fd22/cheng.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_chonrei_6172f341/chonrei.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_chonshu_509e6a31/chonshu.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_dsj_andybogard_1e873098/dsj_andybogard.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_dsj_cheng_7a35d308/dsj_cheng.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_duck-king_bee59977/duck-king.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_duck_king_83e806d7/duck_king.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_duck_king_dm_54a2d040/duck_king_dm.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_fatalfury2_billy_7ddd9473/fatalfury2_billy.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_andy_822b2637/ff1_andy.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_billy_0b7b681a/ff1_billy.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_duck_db2d531f/ff1_duck.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_geese_b6742cdd/ff1_geese.def'] = 'Faction GSD',
	['bulk_fatal_fury_ff1_hwa_638b3e4f/ff1_hwa.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_joe_e9951468/ff1_joe.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_michael_24b54d81/ff1_michael.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_raiden_f4b89a8a/ff1_raiden.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_richard_0f3ec9f2/ff1_richard.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_terry_6129ee63/ff1_terry.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ff1_tung_a1b1de82/ff1_tung.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_franco_685ca23b/franco.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_geese_howard_ffs_1b433d49/ssf2tgeese.def'] = 'Faction GSD',
	['bulk_fatal_fury_gspkrauser_37c55654/gspkrauser.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_gspryo_5b5b2895/gspryo.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_gsptung_c7ce1a0d/gsptung.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_hon-fu_0be1c6fa/hon-fu2.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_hon-fu_9afc06d6/hon-fu.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_hwa jai kof_401c7934/hwa_jai.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ironryuji_2e681cbf/ironryuji.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_joe_7070b74a/joe.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_kim_9054f23e/kim.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_kimffs_b7364de9/kimffs.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_kof_richard_myer_6b4688ea/kof_richard_myer.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_laurence_d4cba9c1/laurence.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_mai shiranui ffs_e2420576/maiffs.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_mai_25b13995/mai.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_mai_9bbd87d4/mai2.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_mai_ff3rb_7e241d6a/mai_ff3rb.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_mary_ff3rb_85046f4f/mary_ff3rb.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_michael max kof_76f37f90/michael_max.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_mini_alfred_ac536177/mini_alfred.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_minmei-bonus_6c54775a/bonus.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_minmei_franco_fc52e999/minmei_franco.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_missimar x heihachi_ecd8c963/mx9900.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_mr_bear_5c45e206/mr_bear.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rb2_geese_0ccc783e/rb2_geese.def'] = 'Faction GSD',
	['bulk_fatal_fury_rb2_geese_be8cb902/rb2_geese.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rb2bash_bfca9b5c/rb2bash.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rb2bob_4b785569/rb2bob.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rb2bob_bc3ccdc1/rb2bob2.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rb2kim_4032bb7a/rb2kim.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rb2kim_74532950/rb2kim2.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rb2laurence_ce6b8b60/rb2laurence.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rb_kim_10ee3b4e/rb_kim.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rbmary_06e73767/rbmary.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rbmary_4449888f/rbmary2.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_rbterry_1b88769b/rbterry.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_robotnik_f0d79bef/robotnik.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ryo another_fdda1110/ryo another.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_shade_f4131245/shade.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_sokaku_3b27befa/sokaku.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_sokaku_f243a125/sokaku.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_ssf2x_terry_1e24d667/ssf2x_terry.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_terry bogard first contact_33079177/pocket-terry.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_terry_ff3rb_6983fe19/terry_ff3rb.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_tora_116d4efb/tora.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_tung_rb_93870f7a/tung_rb.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_white_4e05af16/white.def'] = 'Faction Fatal Fury',
	['bulk_fatal_fury_wood_bonus_game_051da9f6/wood_bonus_game.def'] = 'Faction Fatal Fury',
	['bulk_fire_emblem_amelia_4cf31edc/amelia.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_cecilia_7df4e174/cecilia.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_chrom_c0881a2f/chrom.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_corrin_sebastian_8959581b/corrin_sebastian.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_dorcas_03bb0e98/dorcas.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_erk_a55946bc/erk.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_fe_joshua_7164eb53/fe_joshua.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_fn_db3b2f5f/fn.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_general1_ebdd45c0/general1.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_harken_4c7949ff/harken.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_hawkeye_80cca499/hawkeye.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_hector_22baa3e5/hector.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_ike_572aed32/ike.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_jaffar_379ba61b/jaffar.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_karel_94fc5c0d/karel.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_legault_fab17838/legault.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_legendary hero noire_e4ea910d/legendary hero noire.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_lucina svk infinite_820d9819/lucina svk infinite.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_lyndis_2f20e582/lyndis.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_lyndis_8cdcc045/seizi_ai.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_male corrin _fire emblem__1aa047b5/male corrin (fire emblem).def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_marth_1940104a/marth.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_raven_beta__b177fa23/raven.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_robin _fire emblem__9d479599/robin (fire emblem).def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_roy_abd3d71f/roy-jp.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_roy_fe6d6f44/roy-en.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_sain_4e43cd57/sain.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_tharjaremake_ffe4f719/tharja.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_thief_903cd60a/thief.def'] = 'Faction Fire Emblem',
	['bulk_fire_emblem_will_1470ebea/will.def'] = 'Faction Fire Emblem',
	['bulk_metroid_ff_samus_a33ca787/ff_samus.def'] = 'Faction Metroid',
	['bulk_metroid_ridley_3d0211ce/ridley.def'] = 'Faction Metroid',
	['bulk_metroid_samus2_b305c45d/samus2.def'] = 'Faction Metroid',
	['bulk_ninja_gaiden_cvs_ryuhayabusa_fe6cb49e/cvs_ryuhayabusa.def'] = 'Faction Ninja Gaiden',
	['bulk_ninja_gaiden_ff_ryu_6cd85419/ff_ryu.def'] = 'Faction Ninja Gaiden',
	['bulk_ninja_gaiden_kenhayabusanes_44734999/kenhayabusanes.def'] = 'Faction Ninja Gaiden',
	['bulk_ninja_gaiden_ninjagaiden_bossrush_d3096d20/ninjagaiden_bossrush.def'] = 'Faction Ninja Gaiden',
	['bulk_ninja_gaiden_ryuhayabushanes_8cac3b65/ryuhayabushanes.def'] = 'Faction Ninja Gaiden',
	['bulk_ranma_12__mathias_69623b2d/!!mathias.def'] = 'Faction Ranma',
	['bulk_ranma_12_akane_ff426a68/akane.def'] = 'Faction Ranma',
	['bulk_ranma_12_cologne_28d605f4/cologne.def'] = 'Faction Ranma',
	['bulk_ranma_12_genma_human_4a9e64e3/genma_human.def'] = 'Faction Ranma',
	['bulk_ranma_12_gosunkugi_b5448e34/gosunkugi.def'] = 'Faction Ranma',
	['bulk_ranma_12_herb_1f6a6388/herb.def'] = 'Faction Ranma',
	['bulk_ranma_12_hinakoranma_5ff4e465/hinakoranma.def'] = 'Faction Ranma',
	['bulk_ranma_12_king_0226a504/king.def'] = 'Faction Ranma',
	['bulk_ranma_12_kodachi_b6cad7b7/kodachi.def'] = 'Faction Ranma',
	['bulk_ranma_12_kouchou_486aada0/kouchou.def'] = 'Faction Ranma',
	['bulk_ranma_12_kuno_403bff4d/kuno.def'] = 'Faction Ranma',
	['bulk_ranma_12_mariko_b0c52393/mariko.def'] = 'Faction Ranma',
	['bulk_ranma_12_ranma2_7d5cf6b8/ranma2.def'] = 'Faction Ranma',
	['bulk_ranma_12_ranma_nibbunoichi_c39d3b62/ranma_nibbunoichi.def'] = 'Faction Ranma',
	['bulk_ranma_12_rouge_9891699e/rouge.def'] = 'Faction Ranma',
	['bulk_ranma_12_ryouga_194b7c15/ryouga.def'] = 'Faction Ranma',
	['bulk_ranma_12_sfchappo_a0e495db/sfchappo.def'] = 'Faction Ranma',
	['bulk_ranma_12_shampoo_161642fb/shampoo.def'] = 'Faction Ranma',
	['bulk_ranma_12_tarou_90b42c46/tarou.def'] = 'Faction Ranma',
	['bulk_ranma_12_ukyo by unknown_68dc6dc4/ukyo by unknown.def'] = 'Faction Ranma',
	['bulk_scp_foundation_14-02-2016 10-7-17 pm_b37745a7/gonzalo torchia the loser crybaby.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_682beta_49d4ef64/682beta.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_gatorade_a2857cf4/gatorade.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-028-jp_17dbaa5d/scp-028-jp.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-040-jp-kn_5ad2bcdd/scp-040-jp-kn.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-049-2863 plague doctor-gashadokuro_168d3798/scp-049-2863 plague doctor-gashadokuro.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-049_4e9883ba/scp-049.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-064_fa3fb08c/scp-064.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-066_24bb15fd/scp-066.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-073_69a4fd75/scp-073.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-106_9660f257/scp-106.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-173_89e64932/scp-173.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-1768_b6929603/scp-1768.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-2317_d9c2b770/scp-2317.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-2521_5306e1f3/scp-2521.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-3008_39c8097a/scp-3008.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-354_833d451e/scp-354.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-3812_6aea4bd4/scp-3812.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-3999_8633be38/scp-3999.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-4666_e1ddeebf/scp-4666.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-575_1176a5f2/scp-575.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp-910-jp_cbec1b06/scp-910-jp.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp096beta_bitch__9e5aa897/096.def'] = 'Faction SCP Foundation',
	['bulk_scp_foundation_scp999_14d11275/scp999.def'] = 'Faction SCP Foundation',
	['bulk_shovel_knight_blackk_34e72e16/blackk.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_dark_reize_c78478ac/dark_reize.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_donovan_b15d9ab1/donovan.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_goldarmor_737f9584/goldarmor.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_goldarmor_85a9e45c/dragon_goldarmor.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_hoverhaft_0fb2082d/hoverhaft.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_king knight_66a67c08/king knight.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_king knight_68e29f9c/king knight ultimate supreme.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_liquid_samurai_1c80292c/liquid_samurai.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_mole knight_e7782db2/mole knight.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_mona_ee4b5234/mona.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_mr_hat_4a21a25a/mr_hat.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_phantom striker_0ca5e32c/phantom striker.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_plague knight_7cbe7915/plague knight.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_plague knight_ef1cd9c7/plague knight boomtech.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_polar knight_af47d10f/polar knight.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_propeller knight female_88133169/propeller knight female.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_propeller knight_bd28dd54/propeller knight.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_reizem1.0_528775b0/reize.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_shield knight_aed3c68c/shield  knight.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_shovelknight_65cee04e/shovelknight.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_specterknight_e37112bd/specterknight.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_tinker knight_3b339294/tinker knight.def'] = 'Faction Shovel Knight',
	['bulk_shovel_knight_treasure knight_8ff9cf46/treasure knight.def'] = 'Faction Shovel Knight',
	['bulk_spongebob_200000_c58b47ef/200000.def'] = 'Faction SpongeBob',
	['bulk_spongebob_abrasivespongebob_81af4894/abrasivespongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_barakabob_5279fbb7/barakabob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_barnacle_boy_caue_5eec2147/barnacleboy.def'] = 'Faction SpongeBob',
	['bulk_spongebob_classicspongebob_sebastian_d2423eff/classicspongebob_sebastian.def'] = 'Faction SpongeBob',
	['bulk_spongebob_custom_cheap_spongebob_c53bed78/custom_cheap_spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_darkpatrickremastered_deeea7cc/darkpatrickremastered.def'] = 'Faction SpongeBob',
	['bulk_spongebob_dennis_caue_892b4dc4/dennis.def'] = 'Faction SpongeBob',
	['bulk_spongebob_doodle_94653c98/doodle.def'] = 'Faction SpongeBob',
	['bulk_spongebob_doodlebob_3313eead/doodlebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_doodlebob_929ffda6/doodlebob1.1.def'] = 'Faction SpongeBob',
	['bulk_spongebob_doodlebob_b51a5f0e/doodlebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_doodlebob_mvc_7d595fd0/doodlebob_mvc.def'] = 'Faction SpongeBob',
	['bulk_spongebob_fish_19560fce/fish.def'] = 'Faction SpongeBob',
	['bulk_spongebob_gary_a67a4735/gary.def'] = 'Faction SpongeBob',
	['bulk_spongebob_gary_borderoflife_527c12c9/garythekingsnailboss.def'] = 'Faction SpongeBob',
	['bulk_spongebob_gary_borderoflife_71627b06/garythekingsnail.def'] = 'Faction SpongeBob',
	['bulk_spongebob_garythesnail_aaae9c85/garythesnail.def'] = 'Faction SpongeBob',
	['bulk_spongebob_heartman_f3e8207a/heartman.def'] = 'Faction SpongeBob',
	['bulk_spongebob_holandes herrante_334c0118/holandes herrante.def'] = 'Faction SpongeBob',
	['bulk_spongebob_karate_spongebob_d4b50311/karate_spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_kfsponge_9bbba190/kfsponge.def'] = 'Faction SpongeBob',
	['bulk_spongebob_ki.spongebob_74b5a849/ki.spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_krabs_9ca11b2e/krabs.def'] = 'Faction SpongeBob',
	['bulk_spongebob_larry_0d895be5/larry.def'] = 'Faction SpongeBob',
	['bulk_spongebob_larryrr_2ec6dfd5/larryrr.def'] = 'Faction SpongeBob',
	['bulk_spongebob_man_ray_caue_cab70cf4/manray.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mermaid man v2_5cce437f/eoh.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mermaid_man_caue_9d5c5fa1/mermaidman.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mermaidman_3db85d30/mermaidman.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mgnfnsquidward_191f9832/squidward.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mgnfnsquidward_3412dc6c/mgnfnsquidward.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mixrtaitei2_0410c075/mixrtaitei2.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mixrtaitei2_d140300e/mixrtaitei2.def'] = 'Faction SpongeBob',
	['bulk_spongebob_moarkrabs_4a264233/moarkrabs.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mr.krabs_82a38211/mr.krabs.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mr.krabs_f172b5b0/mr.krabs.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mrkrabs_e74572e4/mrkrabs.def'] = 'Faction SpongeBob',
	['bulk_spongebob_mrkrabs_edc48f21/mrkrabs.def'] = 'Faction SpongeBob',
	['bulk_spongebob_nat peterson_5f07859a/nat peterson.def'] = 'Faction SpongeBob',
	['bulk_spongebob_new spongebob_3eb3d2de/new spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_newpatrick_9d7566d7/newpatrick.def'] = 'Faction SpongeBob',
	['bulk_spongebob_newspongebob_4176b83d/newspongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_newsquidward_f230fb58/newsquidward.def'] = 'Faction SpongeBob',
	['bulk_spongebob_omega spongebob_2ae1112b/omega spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_omega spongebob_ab769dea/omega spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_oni-dark_spongebob_738e7054/oni-dark_spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_orochi spongebob_2f5e1b13/orochi spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_patricio beta 0.81_415db37b/patricio beta 0.81.def'] = 'Faction SpongeBob',
	['bulk_spongebob_patrick1105_551b260b/patrick1105.def'] = 'Faction SpongeBob',
	['bulk_spongebob_patrick_1662c7e9/pat.def'] = 'Faction SpongeBob',
	['bulk_spongebob_patrick_ecc96348/patrick.def'] = 'Faction SpongeBob',
	['bulk_spongebob_pearl without ai_b82c22ed/pearl without ai.def'] = 'Faction SpongeBob',
	['bulk_spongebob_pearl_52ae3c23/pearl.def'] = 'Faction SpongeBob',
	['bulk_spongebob_planit eater patrick_0bfa5a4b/planit eater patrick.def'] = 'Faction SpongeBob',
	['bulk_spongebob_plankton_2fcc1267/plankton (2).def'] = 'Faction SpongeBob',
	['bulk_spongebob_plankton_541887e7/plankton.def'] = 'Faction SpongeBob',
	['bulk_spongebob_plankton_d0f1b86a/plankton.def'] = 'Faction SpongeBob',
	['bulk_spongebob_plankton_eb5f49a9/plankton.def'] = 'Faction SpongeBob',
	['bulk_spongebob_plankton_f2e53757/plankton.def'] = 'Faction SpongeBob',
	['bulk_spongebob_plankton_ffbe9b18/plankton.def'] = 'Faction SpongeBob',
	['bulk_spongebob_retarded ki spongebob_8dd27297/retarded ki spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_retarded spongebob_504b24ec/retarded spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_robot_chicken_spongebob_ex_c98d5c9b/rc_sb-ex.def'] = 'Faction SpongeBob',
	['bulk_spongebob_sandy bochechas_1640e128/sandy bochechas..def'] = 'Faction SpongeBob',
	['bulk_spongebob_sandy bochechas_2f8130f5/sandy bochechas.def'] = 'Faction SpongeBob',
	['bulk_spongebob_sandy cheeks_5b306576/sandy cheeks.def'] = 'Faction SpongeBob',
	['bulk_spongebob_sandy_dc1ad352/sandy.def'] = 'Faction SpongeBob',
	['bulk_spongebob_sb2_spongebob_2ae09a15/sb2_spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_sbob_f7c5e155/sbob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_sbobthesnge - copy _2__81cdce70/sbobthesnge - copy (2).def'] = 'Faction SpongeBob',
	['bulk_spongebob_sbobthesnge - copy _2__fc822044/sbobthesnge - copy (2).def'] = 'Faction SpongeBob',
	['bulk_spongebob_spongebob1105_f1204acf/spongebob1105.def'] = 'Faction SpongeBob',
	['bulk_spongebob_spongebob2_bf3a4ec3/spongebob2.def'] = 'Faction SpongeBob',
	['bulk_spongebob_spongebob_4347a844/spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_spongebob_69b44ea5/spongebob.def'] = 'Faction SpongeBob',
	['bulk_spongebob_spongebobsanchez_4ec38eb5/spongebobsanchez.def'] = 'Faction SpongeBob',
	['bulk_spongebob_spongebobsquarepantsff_e8564c15/spongebobsquarepantsff.def'] = 'Faction SpongeBob',
	['bulk_spongebob_squidward_58fb5c0b/squidward.def'] = 'Faction SpongeBob',
	['bulk_spongebob_sub-patrick_71d85357/sub-patrick.def'] = 'Faction SpongeBob',
	['bulk_spongebob_superspongebobkart_0a35f2e6/superspongebobkart.def'] = 'Faction SpongeBob',
	['bulk_spongebob_swat_fish_8a32dcd2/swatfish.def'] = 'Faction SpongeBob',
	['bulk_spongebob_tanic sandy_2a9d0ced/tanic sandy.def'] = 'Faction SpongeBob',
	['bulk_street_fighter_definitive_adon.iv-2_53aa8ed6/adon.iv-2.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_adonsf1_6ed7644c/adonsf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_birdiesf1_887daece/birdiesf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_cammy-vi_75cd23b6/cammy-vi.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_chun-li-vi_dd9ae34f/chun-li-vi.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_chun-li_ex_sfa_blackjack_68073f9a/chunliex.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_dee_jay_ssf2t_blackjack_69174404/deejay.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_dudley_karmacharmeleon_ca256691/dudley.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_eaglesf1_1e271390/eaglesf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_gekisf1_7e2f903c/gekisf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_gensf1_5a48acb8/gensf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_guile-vi_3a01da1c/guile-vi.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_guy_3c9756ed/guy.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_joesf1_2798d6e6/joesf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_kanekosugi_9b89c5f6/kanekosugi.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_ken-vi_9381d9cc/ken-vi.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_kensf1_93321177/kensf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_laura_3df29ac2/laura.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_leesf1_21ed9335/leesf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_mikesf1_4d9adb4d/mikesf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_retsusf1_8b7e6e7e/retsusf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_ryu-vi_fdd44a15/ryu-vi.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_ryusf1_b668c713/ryusf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_sagatsf1_ad24f4b7/sagatsf1.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_sf2_chun-li_masa_067f8c12/ssf2xchunli.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_sf2_chun-li_masa_9026b3a2/ssf2xchunli2.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_sf3xshinakuma_a423d8af/sf3xshinakuma.def'] = 'Faction SF Definitive',
	['bulk_street_fighter_definitive_violent ken sfv_d4a4e153/violent ken sfv.def'] = 'Faction SF Definitive',
	-- Z-andy joins Lancero; Guardian Heroes pack joins DBC (2026-09-23)
	['bulk_fatal_fury_z-andy_5f02c965/z-andy.def'] = 'Faction Lancero',
	['bulk_guardian_heroes_silver golden_c3fb0067/g-silver.def'] = 'Faction DBC',
	['bulk_guardian_heroes_samuel han_b0fdb6d0/han.def'] = 'Faction DBC',
	['bulk_guardian_heroes_randy_53f270c0/randy.def'] = 'Faction DBC',
	['bulk_guardian_heroes_nicole_f9ef6596/nicole.def'] = 'Faction DBC',
	['bulk_guardian_heroes_p-silver_be710a85/p-silver.def'] = 'Faction DBC',
	['bulk_guardian_heroes_valgar_3b7ab5b1/valgar.def'] = 'Faction DBC',
	['bulk_guardian_heroes_valgar_0e335f8b/protovalgar.def'] = 'Faction DBC',
	['bulk_guardian_heroes_valgar_6b748041/skyvalgar.def'] = 'Faction DBC',
	['bulk_guardian_heroes_valgar_610202ca/valgarxbla.def'] = 'Faction DBC',
	['ryu'] = 'Faction Leader Lancero',
	['thor_avx'] = 'Faction Knight Lancero',
	['066_doge'] = 'Faction Leader Lancero',
	['mr.krabs-ai_patched (hotfix)'] = 'Faction Juggernaut Lancero',
	['sf3_gill_k'] = 'Faction Lancero',
	['rei'] = 'Faction Guardian Lancero',
	['piccolo'] = 'Faction Titan Lancero',
	['songoku'] = 'Faction Lancero',
	['pxzzero'] = 'Faction Lancero',
	['the shotofusion ryuken'] = 'Faction Juggernaut Lancero',
	['electric'] = 'Faction Lancero',
	['wario'] = 'Faction Lancero',
	['orochi leona-kofm'] = 'Faction Lancero',
	['leona_ys'] = 'Faction Lancero',
	['helder'] = 'Faction Creator Lancero',
	['meld_boss/meldboss.def'] = 'Faction Lancero',
	['coco/coco/coco.def'] = 'Faction Titan Lancero',
	['miku'] = 'Faction Captain Lancero',
	['miku/miku.def'] = 'Faction Captain Lancero',
	['knuckles'] = 'Faction Lancero',
	['vega-l'] = 'Faction Lancero',
	['vega-cv2'] = 'Faction Knight Lancero',
	['terry-kof98'] = 'Faction Lancero',
	['ssj3_gokuz2i'] = 'Faction Lancero',
	['kakarotto ssj3'] = 'Faction Lancero',
	['roy/roy.def'] = 'Faction Lancero',
	['roy'] = 'Faction Lancero',
	['pewdiepie/pewdiepie.def'] = 'Faction Lancero',
	['raditz_kn.edit'] = 'Faction Lancero',
	['wr-danielsagat'] = 'Faction Lancero',
	['daniel sagat'] = 'Faction Lancero',
	['k-feilong_wls'] = 'Faction Lancero',
	['0d_hanzou'] = 'Faction Lancero',
	['pip7eop'] = 'Faction Competitive',
	['pip7eop/pip7eop.def'] = 'Faction Competitive',
	["piplup's 7th evolution"] = 'Faction Competitive',
	['super mario 64'] = 'Faction Competitive',
	['super mario 64/super mario 64.def'] = 'Faction Competitive',
	['hyper hayato'] = 'Faction Competitive',
	['hyper hayato/hyper hayato.def'] = 'Faction Competitive',
	['burai1.1'] = 'Faction Competitive',
	['burai_yamamoto'] = 'Faction Competitive',
	['bardock_bt'] = 'Faction Competitive',
	['bardock_bt/bardock_bt'] = 'Faction Competitive',
	['dark storm k'] = 'Faction Competitive',
	['fei-long'] = 'Faction Competitive',
	['snk fei-long'] = 'Faction Competitive',
	['sasuke-kun'] = 'Faction Competitive',
	['sasuke kun'] = 'Faction Competitive',
	['lauren'] = 'Faction Competitive',
	['lauren revolt'] = 'Faction Competitive',
	['mizuchi'] = 'Faction Competitive',
	['lord ryu'] = 'Faction Competitive',
	["tta'goku mui"] = 'Faction Competitive',
	['goku mui'] = 'Faction Competitive',
	['metaknight_zzzz'] = 'Faction Competitive',
	['metaknight_zzzz/metaknight_zzzz.def'] = 'Faction Competitive',
	['metaknight'] = 'Faction Competitive',
	['meta knight'] = 'Faction Competitive',
	['brawl meta knight'] = 'Faction Competitive',
	['kunio'] = 'Faction DBC',
	['misako'] = 'Faction DBC',
	['riki'] = 'Faction DBC',
	['kyoko'] = 'Faction DBC',
	['serena'] = 'Faction DBC',
	['serena/serena.def'] = 'Faction DBC',
	['serena corsair'] = 'Faction DBC',
	['serena by ryon'] = 'Faction DBC',
	['shinryu 2.0'] = 'Faction King DBC',
	['shin ryu'] = 'Faction DBC',
	['scorpionjuan'] = 'Faction General DBC',
	['akuma'] = 'Faction DBC',
	['ultimate_ryu'] = 'Faction DBC',
	['ultimate_ryu/ultimate_ryu.def'] = 'Faction DBC',
	['kenai'] = 'Faction DBC',
	['kenai/kenai.def'] = 'Faction DBC',
	['han_samuel'] = 'Faction Juggernaut DBC',
	['avgn'] = 'Faction Titan DBC',
	['bass'] = 'Faction Captain DBC',
	['donkey_kong/donkey kong.def'] = 'Faction DBC',
	['donkey_kong'] = 'Faction DBC',
	['donkeykongff/donkeykongff.def'] = 'Faction DBC',
	['donkeykongff'] = 'Faction DBC',
	['blanka/blanka.def'] = 'Faction DBC',
	['blanka'] = 'Faction DBC',
	['earl/earl.def'] = 'Faction DBC',
	['earl'] = 'Faction DBC',
	['big rappin\' earl'] = 'Faction DBC',
	['toejam/toejam.def'] = 'Faction DBC',
	['toejam'] = 'Faction DBC',
	['masta dj tj'] = 'Faction DBC',
	['kain/definition.def'] = 'Faction DBC',
	['kain'] = 'Faction DBC',
	['xtr_gillius/xtr_gillius.def'] = 'Faction DBC',
	['xtr_gillius'] = 'Faction DBC',
	['gillius'] = 'Faction DBC',
	['axel'] = 'Faction Last Resort DBC',
	['axel/axel.def'] = 'Faction Last Resort DBC',
	['axel_sor4'] = 'Faction Captain DBC',
	['blaze-sor4'] = 'Faction Knight DBC',
	['so_blaze'] = 'Faction Knight DBC',
	['shiva'] = 'Faction King DBC',
	['ganondorf'] = 'Faction Leader WTG',
	['hektan'] = 'Faction General WTG',
	['morshu'] = 'Faction Juggernaut WTG',
	['theking'] = 'Faction King WTG',
	['link'] = 'Faction Guardian WTG',
	['meta-knight'] = 'Faction Knight WTG',
	['meta_knight'] = 'Faction Knight WTG',
	['zelda mythos'] = 'Faction WTG',
	['sonic the hedgehog'] = 'Faction Leader Cory',
	['miles \'tails\' prower'] = 'Faction Guardian Cory',
	['shadow'] = 'Faction Juggernaut Cory',
	['ai-amyrose'] = 'Faction Captain Cory',
	['jin kazama'] = 'Faction Titan JR',
	['waluigi'] = 'Faction Last Resort Cory',
	['super better mario'] = 'Faction King Cory',
	['super better luigi'] = 'Faction Knight Cory',
	['wannakof'] = 'Faction General Cory',
	['amytea2'] = 'Faction Cory',
	['amytea2/amytea2.def'] = 'Faction Cory',
	['ayumi'] = 'Faction Cory',
	['ayumi/ayumi.def'] = 'Faction Cory',
	['blazetea2'] = 'Faction Cory',
	['blazetea2/blazetea2.def'] = 'Faction Cory',
	['darksonic'] = 'Faction Cory',
	['darksonic/darksonic.def'] = 'Faction Cory',
	['enerjak (tea2)'] = 'Faction Cory',
	['enerjak (tea2)/enerjak (tea2).def'] = 'Faction Cory',
	['exeller'] = 'Faction Cory',
	['exeller/exeller.def'] = 'Faction Cory',
	['exetior'] = 'Faction Cory',
	['exetior/exetior.def'] = 'Faction Cory',
	['exslayer'] = 'Faction Cory',
	['exslayer/exslayer.def'] = 'Faction Cory',
	['knucklestea2'] = 'Faction Cory',
	['knucklestea2/knucklestea2.def'] = 'Faction Cory',
	['metalsonictea2'] = 'Faction Cory',
	['metalsonictea2/metalsonictea2.def'] = 'Faction Cory',
	['rougetea2'] = 'Faction Cory',
	['rougetea2/rougetea2.def'] = 'Faction Cory',
	['sally'] = 'Faction Cory',
	['sally/sally.def'] = 'Faction Cory',
	['scourge(tea2)'] = 'Faction Cory',
	['scourge(tea2)/scourge(tea2).def'] = 'Faction Cory',
	['seelkadoomtea2'] = 'Faction Cory',
	['seelkadoomtea2/seelkadoomtea2.def'] = 'Faction Cory',
	['sonicmodern'] = 'Faction Cory',
	['sonicmodern/sonicmodern.def'] = 'Faction Cory',
	['sonictea'] = 'Faction Cory',
	['sonictea/sonictea.def'] = 'Faction Cory',
	['super sonic(the evil awakens 2)'] = 'Faction Cory',
	['super sonic(the evil awakens 2)/super sonic(the evil awakens 2).def'] = 'Faction Cory',
	['supersonic'] = 'Faction Cory',
	['supersonic/supersonic.def'] = 'Faction Cory',
	['surgetea2'] = 'Faction Cory',
	['surgetea2/surgetea2.def'] = 'Faction Cory',
	['tea2sonic'] = 'Faction Cory',
	['tea2sonic/tea2sonic.def'] = 'Faction Cory',
	['tailstea2'] = 'Faction Cory',
	['tailstea2/tailstea2.def'] = 'Faction Cory',
	['doda'] = 'Faction King Do555',
	['donald_return'] = 'Faction Do555',
	['badmarioff'] = 'Faction Do555',
	['your computer'] = 'Faction Do555',
	['your computer/your computer.def'] = 'Faction Do555',
	['acchi/acchi.def'] = 'Faction Do555',
	['acchi'] = 'Faction Do555',
	['alpacabeta36/alpacabeta36.def'] = 'Faction Do555',
	['alpacabeta36'] = 'Faction Do555',
	['alpaca'] = 'Faction Do555',
	['alpaca (beta)'] = 'Faction Do555',
	['fatbandit/fatbandit.def'] = 'Faction Do555',
	['fatbandit'] = 'Faction Do555',
	['fat bandit'] = 'Faction Do555',
	['fb'] = 'Faction Do555',
	['rogue/rogue.def'] = 'Faction Do555',
	['rogue'] = 'Faction Do555',
	['aos/aos.def'] = 'Faction Do555',
	['aos'] = 'Faction Do555',
	['blackfungus/blackfungus.def'] = 'Faction Do555',
	['blackfungus'] = 'Faction Do555',
	['bubblesk/bubblesk.def'] = 'Faction Do555',
	['bubblesk'] = 'Faction Do555',
	['hammergearman/hammergearman.def'] = 'Faction Do555',
	['hammergearman'] = 'Faction Do555',
	['howlingwolf/howlingwolf.def'] = 'Faction Do555',
	['howlingwolf'] = 'Faction Do555',
	['javelinknight/javelinknight.def'] = 'Faction Do555',
	['javelinknight'] = 'Faction Do555',
	['kabegami36/kabegami36.def'] = 'Faction Do555',
	['kabegami36'] = 'Faction Do555',
	['manikaai/manikaai.def'] = 'Faction Do555',
	['manikaai'] = 'Faction Do555',
	['oguriai/oguriai.def'] = 'Faction Do555',
	['oguriai'] = 'Faction Do555',
	['rainbouyoubi/rainbouyoubi.def'] = 'Faction Do555',
	['rainbouyoubi'] = 'Faction Do555',
	['twister/twister.def'] = 'Faction Do555',
	['twister'] = 'Faction Do555',
	['karla'] = 'Faction Do555',
	['karla/karla.def'] = 'Faction Do555',
	['ds_b.b.hood ex'] = 'Faction Do555',
	['dorothy'] = 'Faction Do555',
	['cfj_android-21'] = 'Faction Juggernaut Do555',
	['milus lilith'] = 'Faction Do555',
	['daviddde'] = 'Faction Do555',
	['frostfreezer'] = 'Faction Do555',
	['linspire'] = 'Faction Do555',
	['tomodachi'] = 'Faction Knight Do555',
	['saburever1'] = 'Faction Titan Do555',
	['katie blue'] = 'Faction Do555',
	['rajaasupergirl'] = 'Faction Do555',
	['cameraman'] = 'Faction Captain Do555',
	['darkness-athena'] = 'Faction Do555',
	['athena2ndact'] = 'Faction Do555',
	['athena2ndact/athena2ndact.def'] = 'Faction Do555',
	['angel element'] = 'Faction Guardian Do555',
	['1000sk-2'] = 'Faction Do555',
	['ghostface'] = 'Faction Do555',
	['wyverncom'] = 'Faction General Do555',
	['medoi'] = 'Faction Last Resort Do555',
	['banditheartless'] = 'Faction Do555',
	['bouncywild'] = 'Faction Do555',
	['xena'] = 'Faction Do555',
	['jshermione'] = 'Faction Do555',
	['leia'] = 'Faction Do555',
	['vegeta primal ego'] = 'Faction Do555',
	['dark=morrigan'] = 'Faction Do555',
	['dark=morrigan_second'] = 'Faction Do555',
	['dark=morrigan_second/dark=morrigan_second.def'] = 'Faction Do555',
	['dark_morrigan_second'] = 'Faction Do555',
	['dark_morrigan_second/dark_morrigan_second.def'] = 'Faction Do555',
	['armeddragon'] = 'Faction Do555',
	['armeddragon/armeddragon.def'] = 'Faction Do555',
	['armed dragon'] = 'Faction Do555',
	['defender'] = 'Faction Do555',
	['defender/defender.def'] = 'Faction Do555',
	['gearlady'] = 'Faction Do555',
	['gearlady/gearlady.def'] = 'Faction Do555',
	['gear lady'] = 'Faction Do555',
	['gorgeousdancer'] = 'Faction Do555',
	['gorgeousdancer/gorgeousdancer.def'] = 'Faction Do555',
	['gorgeous dancer'] = 'Faction Do555',
	['kirarabeauty'] = 'Faction Do555',
	['kirarabeauty/kirarabeauty.def'] = 'Faction Do555',
	['kirara beauty'] = 'Faction Do555',
	['triplesnakes'] = 'Faction Do555',
	['triplesnakes/triplesnakes.def'] = 'Faction Do555',
	['triple snakes'] = 'Faction Do555',
	['fire combat donazen/fire combat donazen.def'] = 'Faction Donald',
	['donald'] = 'Faction Donald',
	['donald/donald.def'] = 'Faction Donald',
	['fire combat donazen'] = 'Faction Donald',
	['gamer of 2021/gamer of 2021.def'] = 'Faction Donald',
	['gamer of 2021'] = 'Faction Donald',
	['agent mac'] = 'Faction Donald',
	['fake phantom donald/fake phantom donald.def'] = 'Faction Donald',
	['fake phantom donald'] = 'Faction Donald',
	['donald_solo_a5/donald_solo_a5.def'] = 'Faction Donald',
	['donald_solo_a5'] = 'Faction Donald',
	['donald solo a5'] = 'Faction Donald',
	['souldonaldrainbow/souldonaldrainbow.def'] = 'Faction Donald',
	['souldonaldrainbow'] = 'Faction Donald',
	['rainbow souldonald v1.1'] = 'Faction Donald',
	['super evildonald/super evildonald.def'] = 'Faction Donald',
	['super evildonald'] = 'Faction Donald',
	['super evil donald'] = 'Faction Donald',
	['the melancholy/the melancholy.def'] = 'Faction Donald',
	['the melancholy'] = 'Faction Donald',
	['winter'] = 'Faction Donald',
	['winter/winter.def'] = 'Faction Donald',
	['donald_solo_a5_aero/donald_solo_a5_aero.def'] = 'Faction Donald',
	['donald_solo_a5_aero'] = 'Faction Donald',
	['evil donald/evil donald.def'] = 'Faction Donald',
	['evil donald'] = 'Faction Donald',
	['ice donald/ice donald.def'] = 'Faction Donald',
	['ice donald'] = 'Faction Donald',
	['spicy donald/spicy donald.def'] = 'Faction Donald',
	['spicy donald'] = 'Faction Donald',
	['donald_miku/miku.def'] = 'Faction Donald',
	['ronald mc miku'] = 'Faction Donald',
	['donald_solo_1st/donald_solo_1st.def'] = 'Faction Donald',
	['donald_solo_1st'] = 'Faction Donald',
	['donald_solo_2nd_alpha6/donald_solo_2nd_alpha6.def'] = 'Faction Donald',
	['donald_solo_2nd_alpha6'] = 'Faction Donald',
	['madness donald'] = 'Faction Donald',
	['d-donald/d-donald.def'] = 'Faction Donald',
	['d-donald'] = 'Faction Donald',
	['angry ryu'] = 'Faction General Momori',
	['bbl_tricky'] = 'Faction Leader Momori',
	['dark side ryu'] = 'Faction Momori',
	['darkvega'] = 'Faction Momori',
	['darkflare'] = 'Faction Momori',
	['demon ken (f.f. style)'] = 'Faction Momori',
	["evilryusf'2/evilryusf'2.def"] = 'Faction Momori',
	['broly-mugen10'] = 'Faction Juggernaut Momori',
	['daimon-kof98'] = 'Faction Momori',
	['divineskullo'] = 'Faction Momori',
	['evil ryu alpha (f.f)'] = 'Faction Momori',
	['omnimanff(edit)'] = 'Faction Momori',
	['thanos'] = 'Faction Momori',
	['god thanos'] = 'Faction Momori',
	['god thanos/god thanos.def'] = 'Faction Momori',
	['hornet'] = 'Faction Momori',
	['hornet/hornet.def'] = 'Faction Momori',
	['the mighty thor'] = 'Faction Momori',
	['immortal thor'] = 'Faction Momori',
	['immortal thor/immortal thor.def'] = 'Faction Momori',
	['infernogami'] = 'Faction Momori',
	['infernogami/infernogami.def'] = 'Faction Momori',
	['thunder-akatsuki'] = 'Faction Momori',
	['thunder-akatsuki 1.1'] = 'Faction Momori',
	['thunder-akatsuki 1.1/thunder-akatsuki 1.1.def'] = 'Faction Momori',
	['diana'] = 'Faction Momori',
	['wonder_womangoldenage'] = 'Faction Momori',
	['wonder_womangoldenage/wonder_womangoldenage.def'] = 'Faction Momori',
	['a_kaigen'] = 'Faction General Momori',
	['astonishing_cyclopsai'] = 'Faction Titan Momori',
	['c.kyo.blood-ckofm'] = 'Faction Momori',
	['cfalconai'] = 'Faction Momori',
	['cooler final'] = 'Faction Momori',
	['coolsroyai'] = 'Faction Knight Momori',
	['deadpoolai'] = 'Faction Last Resort Momori',
	['frieza_bt'] = 'Faction Momori',
	["goukenssf'2+"] = 'Faction Momori',
	['kaoru_hanayama'] = 'Faction Captain Momori',
	['kratosai'] = 'Faction Guardian Momori',
	['lisa a. moon'] = 'Faction King Momori',
	['ryougi'] = 'Faction Momori',
	['ryuga7b'] = 'Faction Momori',
	['sc_sain'] = 'Faction Momori',
	['scorpionai'] = 'Faction Momori',
	['subzeroai'] = 'Faction Momori',
	['ue vegeta z2i'] = 'Faction Momori',
	['weissai'] = 'Faction Momori',
	['yga'] = 'Faction Momori',
	['divergent sila divine general mahoraga'] = 'Faction Momori',
	['blake v3-1.1'] = 'Faction Momori',
	['yagami-instinct'] = 'Faction Momori',
	['yagami-instinct(yuuhi)'] = 'Faction Momori',
	["abyss'mega's/abyss'mega's.def"] = 'Faction Momori',
	['deathtenshi/deathtenshi.def'] = 'Faction Momori',
	['fl_night/fl_night.def'] = 'Faction Momori',
	['raiun/raiun.def'] = 'Faction Momori',
	['raika/raika.def'] = 'Faction Momori',
	['vladrose/vladrose.def'] = 'Faction Momori',
	['ogre ex'] = 'Faction Leader Momori',
	['119way-e-ryu'] = 'Faction Momori',
	['ascended ryu lvl 3'] = 'Faction Momori',
	['chaos-ckofm'] = 'Faction Momori',
	['crazy sakura'] = 'Faction Momori',
	['darth vader'] = 'Faction Momori',
	['extreme iron fist'] = 'Faction Momori',
	['gk_darkchunli'] = 'Faction Momori',
	['giano/giano.def'] = 'Faction Momori',
	['god orochi ryu'] = 'Faction Momori',
	['god urien'] = 'Faction Momori',
	['godrugal'] = 'Faction Momori',
	['godzilla-p'] = 'Faction Momori',
	['gogeta ssj'] = 'Faction Momori',
	['guile-l'] = 'Faction Momori',
	['jotaro over heaven'] = 'Faction Momori',
	['kff'] = 'Faction Momori',
	['khriz blood-kofm'] = 'Faction Momori',
	['lady bison-c'] = 'Faction Momori',
	['legendaryryu v2'] = 'Faction Momori',
	['legulus-ckofm'] = 'Faction Lancero',
	['mbs'] = 'Faction Momori',
	['mg_kratos'] = 'Faction Momori',
	['mvc2_strider_hiryu'] = 'Faction Momori',
	['mx(mario 85)'] = 'Faction Momori',
	['mace windu'] = 'Faction Momori',
	['master gouken'] = 'Faction Momori',
	['mastersakura'] = 'Faction Momori',
	['metal ryu'] = 'Faction Momori',
	['neokei-ckofm'] = 'Faction Momori',
	['oni by lessard'] = 'Faction Momori',
	['onslaughter'] = 'Faction Momori',
	['potspomni'] = 'Faction Boxing',
	['papyrus(exe)'] = 'Faction Momori',
	['professor zoom'] = 'Faction Momori',
	['shin gouki'] = 'Faction Momori',
	['shinlvl2akuma'] = 'Faction Momori',
	['speedshadow'] = 'Faction Momori',
	['uv_dan'] = 'Faction Momori',
	['zeus'] = 'Faction Momori',
	['atomicthouther'] = 'Faction Momori',
	['cyberdan2077'] = 'Faction Momori',
	['magneticbison'] = 'Faction Momori',
	['sf3_ken_c'] = 'Faction Momori',
	['sf4_ryu_c'] = 'Faction Momori',
	['so_sonic'] = 'Faction Momori',
	['momori-dog'] = 'Faction Momori',
	['momori-dog/momori-dog.def'] = 'Faction Momori',
	['batman'] = 'Faction Leader BMB',
	['superman'] = 'Faction General BMB',
	['daredevil'] = 'Faction Juggernaut BMB',
	['skeletor'] = 'Faction Last Resort BMB',
	['megatron'] = 'Faction BMB',
	['billy mmpr'] = 'Faction Titan BMB',
	['spider_morales_by_gartanham'] = 'Faction King BMB',
	['red_ranger'] = 'Faction Guardian BMB',
	['red_ranger/red_ranger.def'] = 'Faction Guardian BMB',
	['pac-man'] = 'Faction Captain BMB',
	['tommy'] = 'Faction Knight BMB',
	['twilight sparkle'] = 'Faction Leader MGM',
	['rainbow dash'] = 'Faction Titan MGM',
	['pinkie pie'] = 'Faction General MGM',
	['applejack'] = 'Faction MGM',
	['trixie lulamoon'] = 'Faction Juggernaut MGM',
	['the great and powerful trixie'] = 'Faction Juggernaut MGM',
	['sonic_tp njpt'] = 'Faction Guardian MGM',
	['shadic njpt'] = 'Faction Captain MGM',
	['son goku rn'] = 'Faction Last Resort MGM',
	['dgblossom'] = 'Faction MGM',
	['dgbubbles'] = 'Faction MGM',
	['dgbuttercup'] = 'Faction King MGM',
	['blossommvc'] = 'Faction Knight MGM',
	['bubblesmvc'] = 'Faction MGM',
	['spiderman-kofm'] = 'Faction MGM',
	['hugo'] = 'Faction Leader GSD',
	['bomber_hugo/bomber_hugo.def'] = 'Faction GSD',
	['hugo_gm2011/hugo.def'] = 'Faction GSD',
	['hugo_andore_gm/hugo.def'] = 'Faction GSD',
	['bomber_hugo'] = 'Faction GSD',
	['geese-rotd/geese-rotd.def'] = 'Faction GSD',
	['nightmare-geese/nightmare-geese.def'] = 'Faction GSD',
	['geesehoward.nw[jjjong1917]/geesehoward.nw[jjjong1917].def'] = 'Faction GSD',
	['n.geese-max'] = 'Faction GSD',
	['hugo andore'] = 'Faction GSD',
	['homero/homero.def'] = 'Faction GSD',
	['homer j simpson'] = 'Faction GSD',
	['kazuya'] = 'Faction Juggernaut JR',
	['zhou lee'] = 'Faction Guardian GSD',
	['geese_ys'] = 'Faction Last Resort GSD',
	['geese'] = 'Faction Knight GSD',
	['geese-kof98'] = 'Faction General GSD',
	['tearyu'] = 'Faction Titan GSD',
	['i-geese'] = 'Faction Captain GSD',
	['akame_ga_kill_1.10'] = 'Faction Leader B-Man',
	['keyser-aunthmer'] = 'Faction Captain B-Man',
	['pain'] = 'Faction Juggernaut B-Man',
	['rock-sp'] = 'Faction Guardian B-Man',
	['rock-sp/rock-sp (hard).def'] = 'Faction Last Resort B-Man',
	['shana'] = 'Faction B-Man',
	['asgore ver.s/asgore ver.s.def'] = 'Faction B-Man',
	['pheromosa_ex'] = 'Faction B-Man',
	['colonel_fire'] = 'Faction Titan B-Man',
	['constantine'] = 'Faction B-Man',
	['g-barbatoslupusrex'] = 'Faction B-Man',
	['ms-06s zaku ii'] = 'Faction B-Man',
	['nevermore-rcmugen_ver0.8/nevermore-rin.def'] = 'Faction Knight B-Man',
	['nightmare-gamera_ex'] = 'Faction B-Man',
	['omega_rugal'] = 'Faction General B-Man',
	['order_ky'] = 'Faction B-Man',
	['robert-keyser'] = 'Faction King B-Man',
	['agnika_soul'] = 'Faction B-Man',
	['yuji itadori'] = 'Faction Leader Dee',
	['gojo'] = 'Faction Last Resort Dee',
	['sukuna shinjuku'] = 'Faction Juggernaut Dee',
	['aizen'] = 'Faction King Dee',
	['madara'] = 'Faction Knight Dee',
	['hashirama'] = 'Faction Captain Dee',
	['boxer'] = 'Faction Juggernaut Boxing',
	['boxing guy'] = 'Faction Knight Boxing',
	['makunoushi ippo'] = 'Faction Titan Boxing',
	['stevefox'] = 'Faction General Boxing',
	['tj comboai'] = 'Faction Boxing',
	['heavyd_ai'] = 'Faction Boxing',
	['m.bison_mx'] = 'Faction Boxing',
	['michael-max-kofa'] = 'Faction Boxing',
	['rickrb'] = 'Faction Leader Boxing',
	['ippo-enhancedai'] = 'Faction Captain Boxing',
	['dempsy ippo'] = 'Faction King Boxing',
	['dudley'] = 'Faction Boxing',
	['ultimate_balrog'] = 'Faction Juggernaut Boxing',
	['ultimate_balrog/ultimate_balrog.def'] = 'Faction Juggernaut Boxing',
	['nintendo_punchout_littlemac'] = 'Faction Guardian Boxing',
	['ippo'] = 'Faction Boxing',
	['rainbowdudley'] = 'Faction Boxing',
	['deadly'] = 'Faction Last Resort Boxing',
	['kenshiro_c2'] = 'Faction Juggernaut Hokuto',
	['rei_ai-patch'] = 'Faction Hokuto',
	['bug_rei/bug_rei.def'] = 'Faction Hokuto',
	['bug_rei'] = 'Faction Hokuto',
	['bug rei'] = 'Faction Hokuto',
	['definitive rei ai alt'] = 'Faction General Hokuto',
	['kenoh'] = 'Faction Hokuto',
	['conqueror raoh_s3'] = 'Faction Last Resort Hokuto',
	['kenshiro-kofm'] = 'Faction Knight Hokuto',
	['rei-hnk-kofa'] = 'Faction King Hokuto',
	['toki-kofa'] = 'Faction Guardian Hokuto',
	['kenshironew'] = 'Faction Leader Hokuto',
	['kenshirou_ai-patch'] = 'Faction Hokuto',
	['shin_ai-patch'] = 'Faction Titan Hokuto',
	['crazy-shin/crazy-shin.def'] = 'Faction Hokuto',
	['crazy-shin'] = 'Faction Hokuto',
	['jagi'] = 'Faction Dee',
	['jagi_r'] = 'Faction Dee',
	['toki_ai-patch'] = 'Faction Captain Hokuto',
	['goldship_hrt'] = 'Faction Leader Umamusume',
	['mac'] = 'Faction King Umamusume',
	['opera'] = 'Faction Juggernaut Umamusume',
	['takion'] = 'Faction Last Resort Umamusume',
	['tokai teio'] = 'Faction General Umamusume',
	['kfmperointernet'] = 'Faction Leader KFM',
	['kfmperointernet/kfmperointernet.def'] = 'Faction Leader KFM',
	['kung fu man'] = 'Faction KFM',
	['kung fu internet.'] = 'Faction Leader KFM',
	['super kfm'] = 'Faction King KFM',
	['super kfm/super kfm.def'] = 'Faction King KFM',
	['super kung fu man'] = 'Faction King KFM',
	['brokenbosskfm'] = 'Faction Juggernaut KFM',
	['brokenbosskfm/brokenbosskfm.def'] = 'Faction Juggernaut KFM',
	['b.t'] = 'Faction Juggernaut KFM',
	['druggy kung fu man'] = 'Faction Juggernaut KFM',
	['violent kung fu man'] = 'Faction Last Resort KFM',
	['violent kung fu man/violent kung fu man.def'] = 'Faction Last Resort KFM',
	['ai kfm'] = 'Faction Last Resort KFM',
	['kfmt1'] = 'Faction General KFM',
	['kfmt1/kfmt1.def'] = 'Faction General KFM',
	['kfm t1'] = 'Faction General KFM',
	['teletransporte_kfm'] = 'Faction Captain KFM',
	['teletransporte_kfm/teletransporte_kfm.def'] = 'Faction Captain KFM',
	['yafm'] = 'Faction KFM',
	['yafm/yafm.def'] = 'Faction KFM',
	['yang fu man'] = 'Faction KFM',
	['tearyu'] = 'Faction Titan GSD',
	['zeroryu'] = 'Faction SF',
	['the shotofusion ryuken'] = 'Faction Juggernaut SF',
	['ryu'] = 'Faction Leader SF',
	['shinryu 2.0'] = 'Faction King SF',
	['shin ryu'] = 'Faction SF',
	['chun_li_jj'] = 'Faction SF',
	['chun-li'] = 'Faction SF',
	['super ryu'] = 'Faction SF',
	['mugen ryup'] = 'Faction Knight SF',
	['hcryu'] = 'Faction SF',
	['hcryu/hcryu.def'] = 'Faction SF',
	['vega-sm'] = 'Faction Knight SF',
	['ultimate_ryu'] = 'Faction Leader Retro SF',
	['kenai'] = 'Faction Captain Retro SF',
	['ultimate_ken'] = 'Faction General Retro SF',
	['ultimate_ken/ultimate_ken.def'] = 'Faction General Retro SF',
	['ultimate_chun_li'] = 'Faction Guardian Retro SF',
	['ultimate_chun_li/ultimate_chun_li.def'] = 'Faction Guardian Retro SF',
	['ultimate_deejay'] = 'Faction Retro SF',
	['ultimate_deejay/ultimate_deejay.def'] = 'Faction Retro SF',
	['ultimate_guile'] = 'Faction Retro SF',
	['ultimate_guile/ultimate_guile.def'] = 'Faction Retro SF',
	['ultimate_sagat'] = 'Faction Retro SF',
	['ultimate_sagat/ultimate_sagat.def'] = 'Faction Retro SF',
	['ultimate_zangief'] = 'Faction Titan Retro SF',
	['ultimate_balrog'] = 'Faction Juggernaut Retro SF',
	['ultimate_blanka'] = 'Faction Last Resort Retro SF',
	['ultimate_blanka/ultimate_blanka.def'] = 'Faction Last Resort Retro SF',
	['ultimate_e_honda'] = 'Faction Retro SF',
	['ultimate_e_honda/ultimate_e_honda.def'] = 'Faction Retro SF',
	['ultimate_dhalsim'] = 'Faction Retro SF',
	['ultimate_dhalsim/ultimate_dhalsim.def'] = 'Faction Retro SF',
	['mvc2_cammy delta'] = 'Faction SF',
	['vega'] = 'Faction Knight SF',
	['original_vega_k/original_vega_k.def'] = 'Faction Knight SF',
	['vega ii'] = 'Faction Knight SF',
	['mvc2_blanka'] = 'Faction SF',
	['akumagg'] = 'Faction SF',
	['guy'] = 'Faction Titan SF',
	['holy ryu3'] = 'Faction SF',
	['god akuma'] = 'Faction SF',
	['akuma'] = 'Faction SF',
	['godakuma'] = 'Faction SF',

	['turbomodegouki'] = 'Faction SF',
	['cyberryuevil'] = 'Faction SF',
	['xryu'] = 'Faction SF',
	['dragon ryu'] = 'Faction Guardian SF',
	['beterryu'] = 'Faction SF',
	['dark god ryu'] = 'Faction SF',
	['ryuuken'] = 'Faction General SF',
	['alexex'] = 'Faction SF',
	['dudley'] = 'Faction SF',
	['rainbowdudley'] = 'Faction SF',
	['burnken'] = 'Faction SF',
	['hken'] = 'Faction Last Resort SF',
	['one ken'] = 'Faction SF',
	['ken_sf3_marvel'] = 'Faction SF',
	['gouken'] = 'Faction SF',
	['fury ken master'] = 'Faction SF',
	['fury ken master/fury ken master.def'] = 'Faction SF',
	['angry ryu'] = 'Faction SF',
	['dark side ryu'] = 'Faction SF',
	['darkvega'] = 'Faction Knight SF',
	["evilryusf'2/evilryusf'2.def"] = 'Faction SF',
	['lord evil ryu'] = 'Faction Captain SF',
	["goukenssf'2+"] = 'Faction SF',
	['demon ken (f.f. style)'] = 'Faction SF',
	['lord evil ken'] = 'Faction SF',
	['violent_ken'] = 'Faction SF',
	['ultimate_ryu_sf1/ultimate_ryu_sf1.def'] = 'Faction Knight Retro SF',
	['vegageegus'] = 'Faction Knight SF',
	['banzoku-alex'] = 'Faction SF',
	['sf3_alex'] = 'Faction SF',
	['another god akuma/another god akuma'] = 'Faction SF',
	['balrog/balrog'] = 'Faction SF',
	['2024sf3gouki'] = 'Faction SF',
	['kaizou shingouki'] = 'Faction SF',
	['chunli/chunli'] = 'Faction SF',
	['cvs_cammy/cvs_cammy'] = 'Faction SF',
	['cyberryu/cyberryu'] = 'Faction SF',
	['evil ryu alpha (f.f)'] = 'Faction SF',
	['dead_vega'] = 'Faction Knight SF',
	['ai-bison/ai-bison'] = 'Faction SF',
	['ssf2x_gouki'] = 'Faction King Retro SF',
	['cvszangief_ex/cvszangief_ex'] = 'Faction SF',
	['daigoken/daigoken'] = 'Faction SF',
	['hakuma'] = 'Faction SF',
	['ogre ex'] = 'Faction SF',
	['real guy'] = 'Faction SF',
	['hugo'] = 'Faction Leader GSD',
	['super ken'] = 'Faction SF',
	['dragon-ken'] = 'Faction SF',
	['m.bison_mx'] = 'Faction SF',

	['gill'] = 'Faction Leader Domination',
	['armor_vegetaz2/armor_vegetaz2'] = 'Faction Domination',
	['brolynew/brolynew'] = 'Faction Domination',
	['d-donald'] = 'Faction Donald',
	['d-donald/d-donald'] = 'Faction Donald',
	['d-donald/d-donald.def'] = 'Faction Donald',
	['terry99m'] = 'Faction Leader KOF',
	['terry99m/terry99m.def'] = 'Faction Leader KOF',
	['legendary bogard 1.0'] = 'Faction KOF',
	['terry bogard'] = 'Faction KOF',
	['zerohaohmaru'] = 'Faction KOF',
	['haohmaru'] = 'Faction KOF',
	['haohmaru/haohmaru.def'] = 'Faction KOF',
	['haohmaru kofm'] = 'Faction KOF',
	['haohmaru-kof98'] = 'Faction Last Resort KOF',
	['cvshaohmaru_ex/cvshaohmaru_ex'] = 'Faction KOF',
	['cvshaohmaru_ex'] = 'Faction KOF',
	['kim_bx'] = 'Faction Captain KOF',
	['itf_kim'] = 'Faction KOF',
	['kim kaphwan'] = 'Faction KOF',
	['kim dong hwan'] = 'Faction KOF',
	['kim jae hoon'] = 'Faction KOF',
	['kim maree'] = 'Faction Leader Do555',
	['helder'] = 'Faction Creator Lancero',
	['jin(the evil awakens 2)'] = 'Faction JR',
	['heihachi'] = 'Faction JR',
	['brolyz2/brolyz2'] = 'Faction Domination',
	['awakened-clark'] = 'Faction Domination',
	['kof_orochi_shermie'] = 'Faction Domination',
	['frozen-yashiro'] = 'Faction Domination',
	['genericbrad'] = 'Faction Domination',
	['hitto'] = 'Faction Domination',
	['ai-cellta/ai-cellta.def'] = 'Faction Domination',
	['ai-cellta'] = 'Faction Domination',
	['gill-enhancedai'] = 'Faction Domination',
	['orochi gill'] = 'Faction Domination',
	['gill rr/gill.def'] = 'Faction Domination',
	['omegath1'] = 'Faction Domination',
	['bonzibuddy'] = 'Faction Domination',
	['thegodofchucknorris/thegodofchucknorris.def'] = "Faction King Gods",
	['thegodofchucknorris'] = "Faction King Gods",
	['norris'] = "Faction King Gods",
	['weegee'] = "Faction Leader Gods",
	['weegee/weegee.def'] = "Faction Leader Gods",
	['black_apocalypse/black_apocalypse'] = "Faction Knight Gods",
	['black_apocalypse'] = "Faction Knight Gods",
	['holyhiryu'] = "Faction Titan Gods",
	['adel-boss/adel-boss'] = "Faction Guardian Gods",
	['adel-boss'] = "Faction Guardian Gods",
	['o_gill'] = "Faction Juggernaut Gods",
	['omegath'] = "Faction General Gods",
	['cfj_law/cfj_law.def'] = 'Faction Leader JR',
	['cfj_law'] = 'Faction Leader JR',
	['cfj_yoshimitsu/cfj_yoshimitsu.def'] = 'Faction General JR',
	['cfj_yoshimitsu'] = 'Faction General JR',
	['cvs_hwoarangtag/cvs_hwoarangtag.def'] = 'Faction Juggernaut JR',
	['cvs_hwoarangtag'] = 'Faction Juggernaut JR',
	['cvs_kumatag/cvs_kumatag.def'] = 'Faction Captain JR',
	['cvs_kumatag'] = 'Faction Captain JR',
	['hwoarangrmh/hwoarangrmh.def'] = 'Faction Knight JR',
	['hwoarangrmh'] = 'Faction Knight JR',
	['paul/definition.def'] = 'Faction Guardian JR',
	['paul'] = 'Faction Guardian JR',
	['mk1_cage/mk1_cage.def'] = 'Faction JR',
	['mk1_cage'] = 'Faction JR',
	['johnny cage mk1'] = 'Faction JR',
	['johnny cage'] = 'Faction JR',
	['mk1_liu-kang/mk1_liu-kang.def'] = 'Faction JR',
	['mk1_liu-kang'] = 'Faction JR',
	['liu-kang mk1'] = 'Faction JR',
	['liu-kang'] = 'Faction JR',
	['mk1_scorpion/mk1_scorpion.def'] = 'Faction JR',
	['mk1_scorpion'] = 'Faction JR',
	['scorpion mk1'] = 'Faction JR',
	['scorpion'] = 'Faction JR',
	['jin'] = 'Faction JR',
	['jin_kazama'] = 'Faction Titan JR',
	-- SSB Melee faction (generated)
	['superluigi_zm'] = 'Faction Leader SSB Melee',
	['superluigi_zm/superluigi_zm.def'] = 'Faction Leader SSB Melee',
	['superluigi'] = 'Faction Leader SSB Melee',
	['luigi'] = 'Faction Leader SSB Melee',
	['strong superluigi/weak superluigi.def'] = 'Faction SSB Melee',
	['strong superluigi'] = 'Faction SSB Melee',
	['w-super luigi'] = 'Faction SSB Melee',
	['kirby'] = 'Faction Captain WTG',
	['kirby/kirby.def'] = 'Faction Captain WTG',
	['119way-mario'] = 'Faction SSB Melee',
	['119way-mario/119way-mario'] = 'Faction SSB Melee',
	['bowser/bowser'] = 'Faction Guardian SSB Melee',
	['morshu/morshu.def'] = 'Faction Juggernaut WTG',
	['119way-mario/119way-mario.def'] = 'Faction SSB Melee',
	['super.mario'] = 'Faction SSB Melee',
	['mario'] = 'Faction SSB Melee',
	['super better mario'] = 'Faction King Cory',
	['super better mario v3'] = 'Faction King Cory',
	['mariops'] = 'Faction SSB Melee',
	['devil mario'] = 'Faction Juggernaut SSB Melee',
	['super bad mario'] = 'Faction SSB Melee',
	['bad mario'] = 'Faction SSB Melee',
	['powerstarmario'] = 'Faction Titan SSB Melee',
	['powerstar mario'] = 'Faction Titan SSB Melee',
	['devil_mario_prime'] = 'Faction Juggernaut SSB Melee',
	['new madness mario'] = 'Faction Last Resort SSB Melee',
	['thouthermario'] = 'Faction SSB Melee',
	['thouther mario'] = 'Faction SSB Melee',
	['bowser'] = 'Faction Guardian SSB Melee',
	['vga_bowser'] = 'Faction Guardian SSB Melee',
	['vga_bowser/vga_bowser.def'] = 'Faction Guardian SSB Melee',
	['link'] = 'Faction Guardian WTG',
	['ganondorf'] = 'Faction General SSB Melee',
	['morshu'] = 'Faction Juggernaut WTG',
	['galacta'] = 'Faction Titan SSB Melee',
	['galacta/galacta.def'] = 'Faction Titan SSB Melee',
	['genocide route kirby'] = 'Faction Titan SSB Melee',
	['galacta knight'] = 'Faction Titan SSB Melee',

	-- Joke Central faction (generated)
	['reggie-skatore'] = 'Faction Leader Joke Central',
	['reggie skatore'] = 'Faction Leader Joke Central',
	['reggie-skatore/reggie-skatore.def'] = 'Faction Leader Joke Central',
	['ghettowar'] = 'Faction Juggernaut Joke Central',
	['ghettowar/ghettowar.def'] = 'Faction Juggernaut Joke Central',
	['ghettowar/ghettowar10.def'] = 'Faction Juggernaut Joke Central',
	['ghetto warmachine'] = 'Faction Juggernaut Joke Central',
	['ghetto war machine'] = 'Faction Juggernaut Joke Central',
	['wario'] = 'Faction King Joke Central',
	['wario/wario.def'] = 'Faction King Joke Central',
	['wario cvs2'] = 'Faction King Joke Central',
	['alteramiba'] = 'Faction Last Resort Joke Central',
	['alteramiba/alteramiba.def'] = 'Faction Last Resort Joke Central',
	['alteramiba/alteramiba_10.def'] = 'Faction Last Resort Joke Central',
	['alter amiba'] = 'Faction Last Resort Joke Central',
	['066_doge'] = 'Faction Captain Joke Central',
	['066_doge/066_doge.def'] = 'Faction Captain Joke Central',
	['doge'] = 'Faction Captain Joke Central',
	['flying shibe such scare'] = 'Faction Captain Joke Central',
	['cblanka/cblanka.def'] = 'Faction Joke Central',
	['cblanka'] = 'Faction Joke Central',
	['coffee blanka'] = 'Faction Joke Central',
	['homerddr2.2'] = 'Faction Joke Central',
	['homer x'] = 'Faction Joke Central',
	['lyndis/seizi_ai.def'] = 'Faction Joke Central',
	['lyndis'] = 'Faction Joke Central',
	['homerjsimpson/homerjsimpson.def'] = 'Faction Joke Central',
	['homerjsimpson'] = 'Faction Joke Central',
	['homer j simpson'] = 'Faction Joke Central',
	['ban green xi'] = 'Faction Joke Central',
	['green'] = 'Faction Joke Central',
	['gogan-if'] = 'Faction Joke Central',
	['gogan'] = 'Faction Joke Central',
	['ultimate_aizen/ultimate_aizen.def'] = 'Faction Bleach',
	['ultimate_aizen'] = 'Faction Bleach',
	['ultimate aizen'] = 'Faction Bleach',

	-- SF3 faction generated membership overrides.
	['sf3_ibuki/ibuki.def'] = 'Faction SF3',
	['sf3_ibuki'] = 'Faction SF3',
	['sf3_remy/sf3_remy.def'] = 'Faction SF3',
	['sf3_remy'] = 'Faction SF3',
	['sf3_yang/yang.def'] = 'Faction SF3',
	['sf3_yang'] = 'Faction SF3',
	['sf3_yun/yun.def'] = 'Faction SF3',
	['sf3_yun'] = 'Faction SF3',
	['119way-e-ryu'] = 'Faction SF3',
	['ascended ryu lvl 3'] = 'Faction SF3',
	['gk_darkchunli'] = 'Faction SF3',
	['god orochi ryu'] = 'Faction SF3',
	['legendaryryu v2'] = 'Faction SF3',
	['master gouken'] = 'Faction SF3',
	['metal ryu'] = 'Faction SF3',
	['sf3_ken_c'] = 'Faction SF3',
	['sf4_ryu_c'] = 'Faction SF3',
	['tearyu'] = 'Faction SF3',
	['kenshirou_ai-patch'] = 'Faction SF3',
	['kenshiro_c2'] = 'Faction SF3',
	['zeroryu'] = 'Faction SF3',
	['shinryu 2.0'] = 'Faction SF3',
	['shin ryu'] = 'Faction SF3',
	['chun_li_jj'] = 'Faction SF3',
	['chun-li'] = 'Faction SF3',
	['mugen ryup'] = 'Faction SF3',
	['hcryu/hcryu.def'] = 'Faction SF3',
	['hcryu'] = 'Faction SF3',
	['holy ryu3'] = 'Faction SF3',
	['cyberryuevil'] = 'Faction SF3',
	['xryu'] = 'Faction SF3',
	['dragon ryu'] = 'Faction SF3',
	['dark god ryu'] = 'Faction SF3',
	['burnken'] = 'Faction SF3',
	['hken'] = 'Faction SF3',
	['one ken'] = 'Faction SF3',
	['ken_sf3_marvel'] = 'Faction SF3',
	['gouken'] = 'Faction SF3',
	['fury ken master/fury ken master.def'] = 'Faction SF3',
	['fury ken master'] = 'Faction SF3',
	['angry ryu'] = 'Faction SF3',
	["goukenssf'2+"] = 'Faction SF3',
	['dark side ryu'] = 'Faction SF3',
	['demon ken (f.f. style)'] = 'Faction SF3',
	["evilryusf'2/evilryusf'2.def"] = 'Faction SF3',
	["evilryusf'2"] = 'Faction SF3',
	['lord evil ken'] = 'Faction SF3',
	['lord evil ryu'] = 'Faction SF3',
	['violent_ken'] = 'Faction SF3',
	['chunli/chunli'] = 'Faction SF3',
	['chunli'] = 'Faction SF3',
	['cyberryu/cyberryu'] = 'Faction SF3',
	['cyberryu'] = 'Faction SF3',
	['daigoken/daigoken'] = 'Faction SF3',
	['daigoken'] = 'Faction SF3',
	['evil ryu alpha (f.f)'] = 'Faction SF3',
	['super ken'] = 'Faction SF3',
	['dragon-ken'] = 'Faction SF3',
	['kenai/kenai.def'] = 'Faction SF3',
	['kenai'] = 'Faction SF3',
	['ultimate_ken/ultimate_ken.def'] = 'Faction SF3',
	['ultimate_ken'] = 'Faction SF3',
	['beterryu/ryu.def'] = 'Faction SF3',
	['beterryu'] = 'Faction SF3',
	['ibuki'] = 'Faction SF3',
	['sf3 remy'] = 'Faction SF3',
	['remy'] = 'Faction SF3',
	['yang'] = 'Faction SF3',
	['yun'] = 'Faction SF3',

	-- DBZ roster/faction imports
	['armor_vegetaz2'] = 'Faction DBZ',
	['normal vegeta z2'] = 'Faction DBZ',
	['vegeta'] = 'Faction DBZ',
	['dbs broly movie goku'] = 'Faction DBZ',
	['blizzard style goku'] = 'Faction DBZ',
	['discordvegeta'] = 'Faction DBZ',
	['discord vegeta'] = 'Faction DBZ',
	['god_goku'] = 'Faction DBZ',
	['god son goku'] = 'Faction DBZ',
	['goku (saiyan god)'] = 'Faction DBZ',
	['gogeta ssj'] = 'Faction DBZ',
	['gogeta'] = 'Faction DBZ',
	['gogeta ssj4'] = 'Faction DBZ',
	['gogeta ssj4 evolution by franciynaldo'] = 'Faction DBZ',
	['gogetasuper4'] = 'Faction DBZ',
	['super gogeta4'] = 'Faction DBZ',
	['supergogeta4'] = 'Faction DBZ',
	['goku/goku.def'] = 'Faction DBZ',
	['goku'] = 'Faction DBZ',
	['goku(the evil awakens 2)'] = 'Faction DBZ',
	['goku black ssr'] = 'Faction DBZ',
	['goku black super saiyan rose ver1.1'] = 'Faction DBZ',
	['goku legend ssj'] = 'Faction DBZ',
	['goku ssj'] = 'Faction DBZ',
	['goku ultra instinto dominado'] = 'Faction DBZ',
	['gokussj4_eb'] = 'Faction DBZ',
	['goku_super4'] = 'Faction DBZ',
	['goku ssj4'] = 'Faction DBZ',
	['gokuub22v3/gokuub22v3.def'] = 'Faction DBZ',
	['gokuub22v3'] = 'Faction DBZ',
	['hyper goku'] = 'Faction DBZ',
	['gokuz2'] = 'Faction DBZ',
	['gokuz2_1.0-ai'] = 'Faction DBZ',
	['goku z2'] = 'Faction DBZ',
	['kidgokuz2i'] = 'Faction DBZ',
	['kid goku z2i'] = 'Faction DBZ',
	['kid goku v2'] = 'Faction DBZ',
	['majin vegeta (dbz)'] = 'Faction DBZ',
	['majin vegeta'] = 'Faction DBZ',
	['piccolo'] = 'Faction DBZ',
	['piccoloz2'] = 'Faction DBZ',
	['piccolo z2'] = 'Faction DBZ',
	['r-ssj4goku'] = 'Faction DBZ',
	['son goku rn'] = 'Faction DBZ',
	['son goku'] = 'Faction DBZ',
	['songoku'] = 'Faction DBZ',
	['songoku_us'] = 'Faction DBZ',
	['songoku_us ver3.4'] = 'Faction DBZ',
	['ssg_gokuz2'] = 'Faction DBZ',
	['god goku z2'] = 'Faction DBZ',
	['god goku'] = 'Faction DBZ',
	['ssgss_gokuz2'] = 'Faction DBZ',
	['ssjb god goku z2'] = 'Faction DBZ',
	['ssjb goku'] = 'Faction DBZ',
	['ssj 2 gokuz2'] = 'Faction DBZ',
	['ssj 2 goku z2'] = 'Faction DBZ',
	['ssj3_gokuz2i'] = 'Faction DBZ',
	['ssj3 goku z2i'] = 'Faction DBZ',
	['ssj3 goku'] = 'Faction DBZ',
	['ssj_goku_legendaryz2'] = 'Faction DBZ',
	['ssj goku legendary z2'] = 'Faction DBZ',
	['ssj goku legendary'] = 'Faction DBZ',
	['ssjblue_vegettoz2'] = 'Faction DBZ',
	['ssjb vegetto z2'] = 'Faction DBZ',
	['ssjb vegetto'] = 'Faction DBZ',
	['super goku boss ultimate'] = 'Faction DBZ',
	['super vegetto4'] = 'Faction DBZ',
	['supergogeta ver4.0'] = 'Faction DBZ',
	['super gogeta'] = 'Faction DBZ',
	['supervegettoai'] = 'Faction DBZ',
	['super vegetto'] = 'Faction DBZ',
	[' vegito ai'] = 'Faction DBZ',
	['tta\'goku mui'] = 'Faction DBZ',
	['goku mui'] = 'Faction DBZ',
	['ue vegeta z2i'] = 'Faction DBZ',
	['ultra ego vegetaz2i'] = 'Faction DBZ',
	['u.e vegeta'] = 'Faction DBZ',
	['ui_gokuz2'] = 'Faction DBZ',
	['ultra instinct goku z2'] = 'Faction DBZ',
	['ui_goku'] = 'Faction DBZ',
	['ultra instinct goku'] = 'Faction DBZ',
	['universal tournament goku'] = 'Faction DBZ',
	['the g.o.a.t of mugen'] = 'Faction DBZ',
	['vegeta by alunfla'] = 'Faction DBZ',
	['vegeta all forms'] = 'Faction DBZ',
	['vegeta primal ego'] = 'Faction DBZ',
	['vegeta super transform 1.6'] = 'Faction DBZ',
	['vegeta-super4(patch)'] = 'Faction DBZ',
	['super vegeta4'] = 'Faction DBZ',
	['super vegeta4(patch)'] = 'Faction DBZ',
	['vegetto'] = 'Faction DBZ',
	['vegettoblue'] = 'Faction DBZ',
	['vegetto blue'] = 'Faction DBZ',
	['vegettoblue_kn'] = 'Faction DBZ',
	['vegettoblue_kn.edit'] = 'Faction DBZ',
	['vegettoz2'] = 'Faction DBZ',
	['vegetto z2'] = 'Faction DBZ',
	['whatsappgoku'] = 'Faction DBZ',
	['whatsapp goku'] = 'Faction DBZ',
	['xeno gogetassj4'] = 'Faction DBZ',
	['xeno gogeta'] = 'Faction DBZ',
	['lancer'] = 'Faction DBZ',
	['xeno_goku op'] = 'Faction DBZ',
	['xeno goku'] = 'Faction DBZ',
	['xeno goku op'] = 'Faction DBZ',
	['nappa'] = 'Faction DBZ',
	['nappa/nappa.def'] = 'Faction DBZ',
	['raditz'] = 'Faction DBZ',
	['raditz/raditz.def'] = 'Faction DBZ',
	['raditz_kn.edit'] = 'Faction DBZ',
}

local t_factionColors = {
	competitive = {255, 235, 90},
	lancero = {80, 220, 255},
	['joke central'] = {255, 130, 220},
	momori = {255, 160, 190},
	dbc = {255, 115, 70},
	cory = {100, 255, 170},
	bmb = {175, 135, 255},
	do555 = {255, 105, 105},
	donald = {255, 70, 70},
	mgm = {235, 205, 255},
	['b-man'] = {95, 175, 255},
	gsd = {255, 185, 80},
	wtg = {120, 255, 215},
	['ssb melee'] = {170, 220, 255},
	dee = {230, 230, 95},
	boxing = {245, 135, 85},
	kof = {255, 205, 70},
	hokuto = {200, 120, 255},
	umamusume = {255, 145, 105},
	kfm = {80, 255, 120},
	sf = {85, 170, 255},
	['retro sf'] = {115, 210, 255},
	sf3 = {80, 135, 255},
	domination = {255, 90, 90},
	gods = {255, 230, 80},
	dbz = {255, 165, 45},
	cheapies = {255, 85, 210},
}

local function f_normalizeFactionName(faction)
	local value = tostring(faction or ''):gsub('%s+[Ff]action%s*$', ''):gsub('^%s+', ''):gsub('%s+$', '')
	return value
end

function start.f_getFactionColor(faction)
	local key = f_normalizeFactionName(faction):lower()
	local color = t_factionColors[key]
	if color == nil then
		color = t_factionColors[key:gsub('[^a-z0-9]+', ' '):gsub('^%s+', ''):gsub('%s+$', '')]
	end
	if color == nil then
		color = {215, 235, 255}
	end
	return color[1], color[2], color[3]
end

function start.f_getActiveCharRef(side)
	if start.p == nil or start.p[side] == nil or start.p[side].t_selected == nil then
		return nil
	end
	if type(currentCharRef) == 'function' then
		local ref = currentCharRef(side)
		if ref and ref >= 0 and start.f_getCharData(ref) then
			return ref
		end
	end
	local selected = start.p[side].t_selected
	if #selected == 0 then
		return nil
	end
	if start.p[side].teamMode == 2 and type(turnsActiveMember) == 'function' then
		local member = turnsActiveMember(side) or 1
		member = math.max(1, math.min(#selected, member))
		if selected[member] ~= nil and selected[member].ref ~= nil then
			return selected[member].ref
		end
	end
	return selected[1].ref
end

function start.f_getActiveCharFaction(side)
	local ref = start.f_getActiveCharRef(side)
	local faction = ref ~= nil and start.f_getCharFaction(ref) or nil
	if faction ~= nil and tostring(faction) ~= '' then
		return f_normalizeFactionName(faction)
	end
	if start.p == nil or start.p[side] == nil or start.p[side].t_selected == nil then
		return nil
	end
	for _, selected in ipairs(start.p[side].t_selected or {}) do
		if selected.ref ~= nil then
			faction = start.f_getCharFaction(selected.ref)
			if faction ~= nil and tostring(faction) ~= '' then
				return f_normalizeFactionName(faction)
			end
		end
	end
	return nil
end

local function f_discordMatchNo()
	if type(matchno) == 'function' then
		return matchno()
	elseif type(matchNo) == 'function' then
		return matchNo()
	end
	return 0
end

local function f_discordSideRefs(side)
	local refs = {}
	local seen = {}
	if start.p == nil or start.p[side] == nil or start.p[side].t_selected == nil then
		return refs
	end
	for _, selected in ipairs(start.p[side].t_selected or {}) do
		local ref = selected.ref
		if ref ~= nil and not seen[ref] then
			table.insert(refs, ref)
			seen[ref] = true
		end
	end
	return refs
end

local function f_discordRefsSignature(refs)
	local parts = {}
	for _, ref in ipairs(refs or {}) do
		table.insert(parts, tostring(ref))
	end
	return table.concat(parts, ',')
end

local function f_discordEventSide(side, refs)
	local ret = {}
	for _, ref in ipairs(refs or f_discordSideRefs(side)) do
		local data = start.f_getCharData(ref)
		if data ~= nil then
			local record = start.f_getCharRecord(ref)
			table.insert(ret, {
				ref = ref,
				char = data.char,
				name = data.name or data.displayname or data.char,
				wins = record.wins or 0,
				losses = record.losses or 0,
				matches = record.matches or 0,
				tier = start.f_getRecordTier(record),
			})
		end
	end
	return ret
end

local function f_jsonEscape(value)
	return tostring(value)
		:gsub('\\', '\\\\')
		:gsub('"', '\\"')
		:gsub('\b', '\\b')
		:gsub('\f', '\\f')
		:gsub('\n', '\\n')
		:gsub('\r', '\\r')
		:gsub('\t', '\\t')
end

local function f_jsonEncode(value)
	local valueType = type(value)
	if valueType == 'nil' then
		return 'null'
	elseif valueType == 'boolean' then
		return value and 'true' or 'false'
	elseif valueType == 'number' then
		return tostring(value)
	elseif valueType == 'string' then
		return '"' .. f_jsonEscape(value) .. '"'
	elseif valueType ~= 'table' then
		return 'null'
	end
	local maxIndex = 0
	local count = 0
	for key, _ in pairs(value) do
		count = count + 1
		if type(key) == 'number' and key > 0 and math.floor(key) == key then
			if key > maxIndex then
				maxIndex = key
			end
		else
			maxIndex = nil
			break
		end
	end
	local parts = {}
	if maxIndex ~= nil and maxIndex == count then
		for i = 1, maxIndex do
			table.insert(parts, f_jsonEncode(value[i]))
		end
		return '[' .. table.concat(parts, ',') .. ']'
	end
	for key, item in pairs(value) do
		table.insert(parts, f_jsonEncode(tostring(key)) .. ':' .. f_jsonEncode(item))
	end
	return '{' .. table.concat(parts, ',') .. '}'
end

local function f_updateFightRecordForRef(stats, ref, result)
	local data = start.f_getCharData(ref)
	if data == nil then
		return nil
	end
	if type(stats.characters) ~= 'table' then
		stats.characters = {}
	end
	local best = f_findBestStatsRecord(stats.characters, f_recordLookupKeys(data), data.char)
	local key = (best and best.key) or start.f_getCharRecordKey(ref)
	if key == nil or key == '' then
		return nil
	end
	local record = stats.characters[key]
	if type(record) ~= 'table' then
		record = {}
		stats.characters[key] = record
	end
	record.wins = tonumber(record.wins or record.win) or 0
	record.losses = tonumber(record.losses or record.loss) or 0
	record.matches = tonumber(record.matches) or (record.wins + record.losses)
	local oldTier = start.f_getRecordTier(record)
	local oldWins = record.wins
	local oldLosses = record.losses
	if result == 'win' then
		record.wins = record.wins + 1
	elseif result == 'loss' then
		record.losses = record.losses + 1
	else
		return nil
	end
	record.matches = record.wins + record.losses
	record.tier = f_recordWinRateTier(record)
	record.recordKey = key
	record.char = data.char
	record.name = data.name or data.displayname or data.char
	return {
		ref = ref,
		key = key,
		name = record.name,
		result = result,
		oldTier = oldTier,
		newTier = record.tier,
		oldWins = oldWins,
		oldLosses = oldLosses,
		wins = record.wins,
		losses = record.losses,
		matches = record.matches,
	}
end

local function f_updateFightRecords(winnerSide)
	winnerSide = tonumber(winnerSide)
	if winnerSide ~= 1 and winnerSide ~= 2 then
		return
	end
	local currentMatch = f_discordMatchNo()
	local currentRound = type(roundno) == 'function' and roundno() or 0
	local loserSide = winnerSide == 1 and 2 or 1
	local winnerRefs = f_discordSideRefs(winnerSide)
	local loserRefs = f_discordSideRefs(loserSide)
	local sig = tostring(currentMatch) .. ':' .. tostring(currentRound) .. ':' .. tostring(winnerSide) .. ':' .. f_discordRefsSignature(winnerRefs) .. '>' .. f_discordRefsSignature(loserRefs)
	if start.recordLastUpdateSig == sig then
		return
	end
	start.recordLastUpdateSig = sig
	local ok, stats = pcall(jsonDecode, start.statsFilePath)
	if not ok or type(stats) ~= 'table' then
		stats = {}
	end
	local changed = false
	local changes = {}
	for _, ref in ipairs(winnerRefs) do
		local change = f_updateFightRecordForRef(stats, ref, 'win')
		if change ~= nil then
			table.insert(changes, change)
			changed = true
		end
	end
	for _, ref in ipairs(loserRefs) do
		local change = f_updateFightRecordForRef(stats, ref, 'loss')
		if change ~= nil then
			table.insert(changes, change)
			changed = true
		end
	end
	if changed then
		stats.updatedAt = os.date('!%Y-%m-%dT%H:%M:%SZ')
		main.f_fileWrite(start.statsFilePath, f_jsonEncode(stats) .. '\n', 'w+')
		t_fightStatsCache = nil
		t_fightStatsCacheTime = nil
		start.lastTierChangeSummary = changes
		start.tierChangeDisplayUntil = (type(gameTime) == 'function' and gameTime() or 0) + 600
	end
end

local function f_appendDiscordMatchRecordEvent(winnerSide)
	winnerSide = tonumber(winnerSide)
	if winnerSide ~= 1 and winnerSide ~= 2 then
		return
	end
	local currentMatch = f_discordMatchNo()
	local currentRound = type(roundno) == 'function' and roundno() or 0
	local loserSide = winnerSide == 1 and 2 or 1
	local winnerRefs = f_discordSideRefs(winnerSide)
	local loserRefs = f_discordSideRefs(loserSide)
	local sig = tostring(currentMatch) .. ':' .. tostring(currentRound) .. ':' .. tostring(winnerSide) .. ':' .. f_discordRefsSignature(winnerRefs) .. '>' .. f_discordRefsSignature(loserRefs)
	if start.discordLastMatchRecordSig == sig then
		return
	end
	start.discordLastMatchRecordSig = sig
		local currentMode = type(gameMode) == 'function' and gameMode() or tostring(gamemode or '')
		local event = {
			schema = 1,
			source = 'IKEMEN_DEV',
			createdAt = os.date('!%Y-%m-%dT%H:%M:%SZ'),
			mode = currentMode,
		match = currentMatch,
		type = 'match',
		winnerSide = winnerSide,
		winners = f_discordEventSide(winnerSide, winnerRefs),
		losers = f_discordEventSide(loserSide, loserRefs),
			p1 = f_discordEventSide(1),
			p2 = f_discordEventSide(2),
		}
	main.f_fileWrite('save/match_events.jsonl', f_jsonEncode(event) .. '\n', 'a+')
end

function start.f_recordCharMatchResult(winnerSide)
	f_updateFightRecords(winnerSide)
	f_appendDiscordMatchRecordEvent(winnerSide)
end

local function f_tierChangeMakeText(x, y, align, r, g, b, scale)
	return text:create({
		font = 3,
		bank = 0,
		align = align,
		text = '',
		x = x,
		y = y,
		scaleX = scale,
		scaleY = scale,
		r = r,
		g = g,
		b = b,
		a = 255,
		height = -1,
		xshear = 0,
		angle = 0,
		window = nil,
		defsc = false,
	})
end

local function f_tierChangeText(change)
	local arrow = change.oldTier ~= change.newTier and ' -> ' or ' = '
	return string.format('%s: %s%s%s  %d-%d', change.name or '?', start.f_getRecordTierDisplayText(change.oldTier), arrow, start.f_getRecordTierDisplayText(change.newTier), change.wins or 0, change.losses or 0)
end

function start.f_drawTierChangeOverlay()
	if start.lastTierChangeSummary == nil or #start.lastTierChangeSummary == 0 then
		return
	end
	if start.tierChangeDisplayUntil ~= nil and type(gameTime) == 'function' and gameTime() > start.tierChangeDisplayUntil then
		return
	end
	if start.txt_tierChangeOverlay == nil then
		start.txt_tierChangeOverlay = {
			titleShadow = f_tierChangeMakeText(642, 102, 0, 0, 0, 0, 1.55),
			titleValue = f_tierChangeMakeText(640, 100, 0, 255, 235, 80, 1.55),
			linesShadow = {},
			linesValue = {},
		}
		for i = 1, 4 do
			start.txt_tierChangeOverlay.linesShadow[i] = f_tierChangeMakeText(642, 126 + i * 18, 0, 0, 0, 0, 1.15)
			start.txt_tierChangeOverlay.linesValue[i] = f_tierChangeMakeText(640, 124 + i * 18, 0, 255, 255, 255, 1.15)
		end
	end
	local title = 'RECORD / TIER UPDATE'
	start.txt_tierChangeOverlay.titleShadow:update({text = title, x = 642, y = 102, align = 0, r = 0, g = 0, b = 0})
	start.txt_tierChangeOverlay.titleShadow:draw()
	start.txt_tierChangeOverlay.titleValue:update({text = title, x = 640, y = 100, align = 0, r = 255, g = 235, b = 80})
	start.txt_tierChangeOverlay.titleValue:draw()
	for i = 1, 4 do
		local change = start.lastTierChangeSummary[i]
		local line = change ~= nil and f_tierChangeText(change) or ''
		local r, g, b = 255, 255, 255
		if change ~= nil and change.oldTier ~= change.newTier then
			local color = t_recordTierColors[change.newTier] or t_recordTierColors[tostring(change.newTier or ''):sub(1, 1)] or t_recordTierColors.F
			r, g, b = color[1], color[2], color[3]
		end
		start.txt_tierChangeOverlay.linesShadow[i]:update({text = line, x = 642, y = 126 + i * 18, align = 0, r = 0, g = 0, b = 0})
		start.txt_tierChangeOverlay.linesShadow[i]:draw()
		start.txt_tierChangeOverlay.linesValue[i]:update({text = line, x = 640, y = 124 + i * 18, align = 0, r = r, g = g, b = b})
		start.txt_tierChangeOverlay.linesValue[i]:draw()
	end
end

--returns stage ref out of def filename
function start.f_getStageRef(def)
	if def == '' then
		return getStageNo()
	end
	if main.t_stageDef[def:lower()] == nil then
		 main.f_addStage(def)
	end
	return main.t_stageDef[def:lower()]
end

--returns char ref out of def filename
function start.f_getCharRef(def)
	if main.t_charDef[def:lower()] == nil then
		if not main.f_addChar(def .. ', order = 0, ordersurvival = 0, exclude = 1', true, false) then
			panicError("\nUnable to add character. No such file or directory: " .. def .. "\n")
		end
	end
	return main.t_charDef[def:lower()]
end

--returns teammode int from string
function start.f_stringToTeamMode(tm)
	if tm == 'single' then
		return 0
	elseif tm == 'simul' then
		return 1
	elseif tm == 'turns' then
		return 2
	elseif tm == 'tag' then
		return 3
	end
	return nil
end

--returns formatted clear time string
function start.f_clearTimeText(text, totalSec)
	local h = tostring(math.floor(totalSec / 3600))
	local m = tostring(math.floor((totalSec / 3600 - h) * 60))
	local s = tostring(math.floor(((totalSec / 3600 - h) * 60 - m) * 60))
	local x = tostring(math.floor((((totalSec / 3600 - h) * 60 - m) * 60 - s) * 100))
	if string.len(m) < 2 then
		m = '0' .. m
	end
	if string.len(s) < 2 then
		s = '0' .. s
	end
	if string.len(x) < 2 then
		x = '0' .. x
	end
	return text:gsub('%%h', h):gsub('%%m', m):gsub('%%s', s):gsub('%%x', x)
end

--returns formatted record text table
function start.f_getRecordText()
	local text = motif.select_info.record.text[gameMode()]
	if text == nil then
		return ""
	end
	local stats = jsonDecode(start.statsFilePath)
	if stats.modes == nil or stats.modes[gameMode()] == nil or stats.modes[gameMode()].ranking == nil or stats.modes[gameMode()].ranking[1] == nil then
		return ""
	end
	local t = stats.modes[gameMode()].ranking[1]
	--time
	text = start.f_clearTimeText(text, t.time)
	--score
	text = text:gsub('%%p', tostring(t.score))
	--win
	text = text:gsub('%%r', tostring(t.win))
	--char name
	local name = '?' --in case character being removed from roster
	if main.t_charDef[t.chars[1]] ~= nil then
		name = start.f_getCharData(main.t_charDef[t.chars[1]]).name
	end
	text = text:gsub('%%c', name)
	--player name
	text = text:gsub('%%n', t.name)
	return text
end

local selectStatsOverlay = nil

local function f_initSelectStatsOverlay()
	if selectStatsOverlay ~= nil then
		return selectStatsOverlay
	end
	local font = fontNew('font/default-3x5.def', -1)
	local function makeText(x, y, align, r, g, b)
		local t = textImgNew()
		textImgSetFont(t, font)
		textImgSetBank(t, 0)
		textImgSetAlign(t, align)
		textImgSetColor(t, r, g, b, 255)
		textImgSetScale(t, 1, 1)
		textImgSetPos(t, x, y)
		return t
	end
	local function makeRect()
		local r = rectNew()
		rectSetColor(r, 0, 0, 0)
		rectSetAlpha(r, 220, 0)
		rectSetLayerno(r, 2)
		return r
	end
	selectStatsOverlay = {
		box = {makeRect(), makeRect()},
		label = {makeText(8, 171, 1, 220, 220, 220), makeText(312, 171, -1, 220, 220, 220)},
		value = {makeText(8, 181, 1, 255, 225, 120), makeText(312, 181, -1, 255, 225, 120)},
	}
	return selectStatsOverlay
end

local function f_selectStatsNormalizeKey(value)
	if value == nil then
		return nil
	end
	value = tostring(value):gsub('\\', '/'):gsub('^%s+', ''):gsub('%s+$', '')
	if value == '' then
		return nil
	end
	return value:lower()
end

local function f_selectStatsKeys(ref)
	local keys = {}
	local seen = {}
	local function addKey(key)
		if key == nil or key == '' or seen[key] then
			return
		end
		table.insert(keys, key)
		seen[key] = true
	end
	local function add(value)
		local key = f_selectStatsNormalizeKey(value)
		if key ~= nil then
			addKey(key)
			addKey((key:gsub('[^a-z0-9]+', '')))
			addKey((key:gsub('%.def$', '')))
			local stem = key:gsub('^.*[/]', ''):gsub('%.def$', '')
			addKey(stem)
			addKey((stem:gsub('[^a-z0-9]+', '')))
			local folder = key:match('^([^/]+)/')
			if folder ~= nil then
				addKey(folder)
				addKey((folder:gsub('[^a-z0-9]+', '')))
			end
		end
	end
	if ref ~= nil then
		local ok, data = pcall(start.f_getCharData, ref)
		if ok and data ~= nil then
			add(data.char)
			add(data.def)
			add(data.name)
		end
	end
	return keys
end

local function f_selectStatsRead()
	local ok, stats = pcall(jsonDecode, start.statsFilePath)
	if not ok or type(stats) ~= 'table' then
		return {}
	end
	return stats
end

local function f_selectStatsRecord(stats, ref)
	local sources = {}
	if type(stats.characters) == 'table' then
		table.insert(sources, stats.characters)
	end
	table.insert(sources, stats)
	local keys = f_selectStatsKeys(ref)
	local best = nil
	for _, source in ipairs(sources) do
		local ownerData = start.f_getCharData(ref)
		local candidate = f_findBestStatsRecord(source, keys, ownerData and ownerData.char or nil)
		if candidate ~= nil then
			best = f_bestRecordCandidate(best, candidate.record, candidate.key)
		end
	end
	if best ~= nil then
		return best.record
	end
	return {}
end

local function f_selectStatsValue(rec, ...)
	for _, key in ipairs({...}) do
		if rec[key] ~= nil then
			return rec[key]
		end
	end
	return nil
end

local function f_selectStatsText(stats, ref)
	local rec = f_selectStatsRecord(stats, ref)
	local wins = tonumber(f_selectStatsValue(rec, 'wins', 'win', 'Wins', 'Win')) or 0
	local losses = tonumber(f_selectStatsValue(rec, 'losses', 'loss', 'Losses', 'Loss')) or 0
	local tier = f_selectStatsValue(rec, 'tier', 'Tier', 'rank', 'Rank') or 'U'
	return tostring(wins) .. ' - ' .. tostring(losses) .. ' - ' .. tostring(tier)
end

function start.f_drawSelectStatsOverlay(counter)
	local overlay = f_initSelectStatsOverlay()
	local stats = f_selectStatsRead()
	for side = 1, 2 do
		local ref = nil
		if #start.p[side].t_selTemp > 0 then
			local idx = t_portraitPriority[side] or #start.p[side].t_selTemp
			ref = start.p[side].t_selTemp[idx] and start.p[side].t_selTemp[idx].ref
		end
		local x1 = side == 1 and 0 or 160
		local x2 = side == 1 and 160 or 320
		rectSetWindow(overlay.box[side], x1, 168, x2, 192)
		rectDraw(overlay.box[side], 2)
		textImgReset(overlay.label[side], {'text'})
		textImgSetText(overlay.label[side], 'stats')
		textImgDraw(overlay.label[side], 2)
		textImgReset(overlay.value[side], {'text'})
		textImgSetText(overlay.value[side], f_selectStatsText(stats, ref))
		textImgDraw(overlay.value[side], 2)
	end
end

--cursor sound data, play cursor sound
function start.f_playWave(ref, name, g, n, loops)
	if g < 0 or n < 0 then return 0 end
	if name == 'stage' then
		local a = main.t_selStages[ref].attachedChar
		if a == nil or a.sound == nil then
			return 0
		end
		if main.t_selStages[ref][name .. '_wave_data'] == nil then
			main.t_selStages[ref][name .. '_wave_data'] = waveNew(a.dir .. a.sound, g, n, loops or -1)
		end
		wavePlay(main.t_selStages[ref][name .. '_wave_data'], g, n)
	else
		local sound = start.f_getCharData(ref).sound
		if sound == nil or sound == '' then
			return 0
		end
		local key = name .. '_wave_data_' .. g .. '_' .. n
		if start.f_getCharData(ref)[key] == nil then
			start.f_getCharData(ref)[key] =
				waveNew(start.f_getCharData(ref).dir .. start.f_getCharData(ref).sound, g, n, loops or -1)
		end
		wavePlay(start.f_getCharData(ref)[key], g, n)
	end
end

--removes char with particular ref from table
function start.f_excludeChar(t, ref)
	for _, list in pairs(t) do
		if type(list) == 'table' then
			for i = #list, 1, -1 do
				if list[i] == ref then
					table.remove(list, i)
				end
			end
		end
	end
	return t
end

--shuffles a table in-place (using synced RNG)
function start.f_shuffleTable(t, last)
	for i = #t, 2, -1 do
		local j = math.random(i)
		t[i], t[j] = t[j], t[i]
	end
	-- prevent first element from repeating the last of previous cycle
	if last and #t > 1 and t[#t] == last then
		-- swap the first element with a random other position
		local swap = math.random(#t - 1)
		t[#t], t[swap] = t[swap], t[#t]
	end
end

--returns random char ref
function start.f_randomChar(pn)
	if #main.t_randomChars == 0 then
		return nil
	end
	start.shuffleChars = start.shuffleChars or {}

	if not start.shuffleChars[pn] or #start.shuffleChars[pn] == 0 then
		local last = start.lastRandomChar and start.lastRandomChar[pn]
		local t = {}
		for _, v in ipairs(main.t_randomChars) do
			if gameOption('Options.Team.Duplicates') or not t_reservedChars[pn][v] then
				table.insert(t, v)
			end
		end
		start.f_shuffleTable(t, last)
		start.shuffleChars[pn] = t
	end
	-- draw one char from the bag
	local result = table.remove(start.shuffleChars[pn])
	-- store the last drawn value
	start.lastRandomChar = start.lastRandomChar or {}
	start.lastRandomChar[pn] = result
	return result
end

--return true if slot is selected, update start.t_grid
function start.f_slotSelected(cell, side, cmd, player, x, y)
	if cmd == nil then
		return false, false
	end
	local original_slot = main.t_selGrid[cell].slot
	if #main.t_selGrid[cell].chars > 0 then
		-- select.def 'slot' parameter special keys detection
		for _, cmdType in ipairs({'select', 'next', 'previous'}) do
			if main.t_selGrid[cell][cmdType] ~= nil then
				for k, v in pairs(main.t_selGrid[cell][cmdType]) do
					if commandGetState(cmd, k) then
						if cmdType == 'next' then
							local ok = false
							for i = main.t_selGrid[cell].slot + 1, #v do
								if start.f_getCharData(start.f_selGrid(cell, v[i]).char_ref).hidden < 2 then
									main.t_selGrid[cell].slot = v[i]
									ok = true
									break
								end
							end
							if not ok then
								for i = 1, main.t_selGrid[cell].slot - 1 do
									if start.f_getCharData(start.f_selGrid(cell, v[i]).char_ref).hidden < 2 then
										main.t_selGrid[cell].slot = v[i]
										ok = true
										break
									end
								end
							end
							if ok then
								sndPlay(motif.Snd, motif.select_info['p' .. side].swap.snd[1], motif.select_info['p' .. side].swap.snd[2])
							end
						elseif cmdType == 'previous' then
							local ok = false
							for i = main.t_selGrid[cell].slot -1, 1, -1 do
								if start.f_getCharData(start.f_selGrid(cell, v[i]).char_ref).hidden < 2 then
									main.t_selGrid[cell].slot = v[i]
									ok = true
									break
								end
							end
							if not ok then
								for i = #v, main.t_selGrid[cell].slot + 1, -1 do
									if start.f_getCharData(start.f_selGrid(cell, v[i]).char_ref).hidden < 2 then
										main.t_selGrid[cell].slot = v[i]
										ok = true
										break
									end
								end
							end
							if ok then
								sndPlay(motif.Snd, motif.select_info['p' .. side].swap.snd[1], motif.select_info['p' .. side].swap.snd[2])
							end
						else --select
							main.t_selGrid[cell].slot = v[math.random(#v)]
							start.c[player].selRef = start.f_selGrid(cell).char_ref
						end
						start.t_grid[y + 1][x + 1].char = start.f_selGrid(cell).char
						start.t_grid[y + 1][x + 1].char_ref = start.f_selGrid(cell).char_ref
						start.t_grid[y + 1][x + 1].hidden = start.f_selGrid(cell).hidden
						start.t_grid[y + 1][x + 1].skip = start.f_selGrid(cell).skip
						return cmdType == 'select', (main.t_selGrid[cell].slot ~= original_slot)
					end
				end
			end
		end
	end
	-- returns true on pressed key if current slot is not blocked by TeamDuplicates feature
	return main.f_btnPalNo(cmd) > 0 and (not t_reservedChars[side][start.t_grid[y + 1][x + 1].char_ref] or start.t_grid[start.c[player].selY + 1][start.c[player].selX + 1].char == 'randomselect'),false
end

--generate start.t_grid table, assign row and cell to main.t_selChars
local cnt = motif.select_info.columns + 1
local row = 1
local col = 0
start.t_grid = {[row] = {}}
for i = 1, motif.select_info.rows * motif.select_info.columns do
	if i == cnt then
		row = row + 1
		cnt = cnt + motif.select_info.columns
		start.t_grid[row] = {}
	end
	col = #start.t_grid[row] + 1
	local cell_spacing = getCellSpacing(col - 1, row - 1)
	local cell_offset = getCellOffset(col - 1, row - 1)
	start.t_grid[row][col] = {
		x = (col - 1) * (motif.select_info.cell.size[1] + cell_spacing[1]) + cell_offset[1],
		y = (row - 1) * (motif.select_info.cell.size[2] + cell_spacing[2]) + cell_offset[2]
	}
	if start.f_selGrid(i).char ~= nil then
		start.t_grid[row][col].char = start.f_selGrid(i).char
		start.t_grid[row][col].char_ref = start.f_selGrid(i).char_ref
		start.t_grid[row][col].hidden = start.f_selGrid(i).hidden
		for j = 1, #main.t_selGrid[i].chars do
			start.f_selGrid(i, j).row = row
			start.f_selGrid(i, j).col = col
		end
	end
	local overrideSkip = getCellSkip(col - 1, row - 1)
	if start.f_selGrid(i).skip == 1 or overrideSkip then
		start.t_grid[row][col].skip = 1
	end
end
if gameOption('Debug.DumpLuaTables') then main.f_printTable(start.t_grid, 'debug/t_grid.txt') end

local function updateCommon(common, add)
	for k, values in pairs(common) do
		if values ~= nil and #values > 0 then
			local optionName = 'Common.' .. k:gsub('^%l', string.upper)
			local t = gameOption(optionName)
			if add then
				local existing = {}
				local changed = false
				for _, v in ipairs(t) do
					existing[v] = true
				end
				for _, v in ipairs(values) do
					if v ~= '' and not existing[v] then
						table.insert(t, v)
						existing[v] = true
						changed = true
					end
				end
				if changed then
					modifyGameOption(optionName, t)
				end
			else
				local toRemove = {}
				local changed = false
				for _, v in ipairs(values) do
					if v ~= '' then
						toRemove[v] = true
					end
				end
				for i = #t, 1, -1 do
					if toRemove[t[i]] then
						table.remove(t, i)
						changed = true
					end
				end
				if changed then
					modifyGameOption(optionName, t)
				end
			end
		end
	end
end

-- return amount of life to recover
local function f_lifeRecovery(lifeMax, fighter)
	local ret = hook.runFirst("start.f_lifeRecovery", lifeMax, fighter)
	if ret ~= nil then return ret end
	local bonus = lifeMax * gameOption('Options.Turns.Recovery.Bonus') / 100
	local base = lifeMax * gameOption('Options.Turns.Recovery.Base') / 100
	return base + main.f_round(timeRemaining() / (timeRemaining() + timeElapsed()) * bonus)
end

-- match persistence
function start.f_matchPersistence()
	-- checked only after at least 1 match
	if matchNo() >= 2 then
		local gameStats = getGameStats()
		local matches = (getGameStats().Matches) or {}
		local idx = #matches
		if idx <= (start.matchPersistenceStatsIdx or 0) then
			return start.p[1].numChars
		end
		start.matchPersistenceStatsIdx = idx
		local roundStats = matches[idx].Rounds
		-- set 'existed' flag (decides if var/fvar should be persistent between matches)
		if roundStats then
			for _, round in ipairs(roundStats) do
				for side = 1, 2 do
					local fighters = (round.Fighters and round.Fighters[side]) or {}
					for _, f in ipairs(fighters) do
						local memberIdx = (f.MemberNo or 0) + 1
						if start.p[side].t_selected[memberIdx] ~= nil then
							start.p[side].t_selected[memberIdx].existed = true
						end
					end
				end
			end
		end
		-- if defeated members should be skipped from next match, or if life should be maintained
		if main.dropDefeated or main.persistLife then
			local turnsOffset = start.p[1].turnsOffset or 0
			-- Turns
			if start.p[1].teamMode == 2 then
				--for each round in the last match
				if roundStats then
					for _, round in ipairs(roundStats) do
						-- P1 active fighter snapshot for the round
						local f1 = (round.Fighters and round.Fighters[1] and round.Fighters[1][1]) or nil
						if f1 then
							local memberIdx = (f1.MemberNo or 0) + 1
							-- if defeated
							if f1.KO and (f1.Life or 0) <= 0 then
								-- Keep full roster for lifebar faces + hiscore, but advance turnsOffset
								-- so defeated members are skipped in the next match.
								if main.dropDefeated then
									turnsOffset = math.max(turnsOffset, memberIdx)
									if start.p[1].t_selected[memberIdx] ~= nil then
										start.p[1].t_selected[memberIdx].life = 0
									end
								-- or resurrect and recover character's life
								elseif main.persistLife then
									start.p[1].t_selected[memberIdx].life = math.max(1, f_lifeRecovery(f1.LifeMax or 0, f1))
								end
							-- otherwise maintain character's life
							elseif main.persistLife then
								start.p[1].t_selected[memberIdx].life = f1.Life or start.p[1].t_selected[memberIdx].life
							end
						end
					end
				end
				if main.dropDefeated then
					start.p[1].turnsOffset = turnsOffset
					local total = #start.p[1].t_selected
					start.p[1].numChars = math.max(0, total - turnsOffset)
				end
			-- Single / Simul / Tag
			else
				-- for each player data in the last round (new format)
				if roundStats and #roundStats > 0 then
					local lastRound = roundStats[#roundStats]
					for side = 1, 2 do
						local fighters = (lastRound.Fighters and lastRound.Fighters[side]) or {}
						-- only check player-controlled side
						if not main.cpuSide[side] then
							for _, f in ipairs(fighters) do
								local memberIdx = (f.MemberNo or 0) + 1
								-- if defeated
								if f.KO and (f.Life or 0) <= 0 then
									if main.persistLife then
										start.p[1].t_selected[memberIdx].life = math.max(1, f_lifeRecovery(f.LifeMax or 0, f))
									end
								-- otherwise maintain character's life
								elseif main.persistLife then
									start.p[1].t_selected[memberIdx].life = f.Life or start.p[1].t_selected[memberIdx].life
								end
							end
						end
					end
				end
			end
		end
	end
	return start.p[1].numChars
end

--start game
function start.f_game(common)
	clearColor(0, 0, 0)
	if gameOption('Debug.DumpLuaTables') and start ~= nil then
		main.f_printTable(start.p, 'debug/t_p.txt')
	end
	if gameMode('training') then
		menu.f_trainingReset()
	end
	start.saltyBetRoundResetDone = false
	start.saltyBetHoldStartGameTime = nil
	local winner = -1
	winner, start.challenger = game()
	if start.challenger == 0 then
		local winnerSide = tonumber(winner)
		if winnerSide ~= 1 and winnerSide ~= 2 then
			winnerSide = type(getWinnerTeam) == 'function' and getWinnerTeam() or winnerSide
		end
		f_updateFightRecords(winnerSide)
		f_appendDiscordMatchRecordEvent(winnerSide)
	end
	if gameOption('Debug.DumpLuaTables') then
		main.f_printTable(getGameStats(), 'debug/t_gameStats.txt')
	end
	main.f_restoreInput()
	updateCommon(common, false)
	if shutdown() then
		clearColor(0, 0, 0)
		os.exit()
	end
	return winner
end

--;===========================================================
--; MODES LOOP
--;===========================================================
function start.f_selectMode()
	start.f_selectReset(true)
	while true do
		--select screen
		if gameOption('Config.BootLoadingMode') == 1 then
			main.f_waitForPreloads()
		end
		if not start.f_selectScreen() then
			bgReset(motif[main.background].BGDef)
			fadeInInit(motif[main.group].fadein.FadeData)
			playBgm({source = "motif.title", interrupt = true})
			return
		end
		-- lua file with custom arcade path detection
		local path = main.luaPath
		if main.charparam.arcadepath then
			if start.f_getCharData(start.p[1].t_selected[1].ref).arcadepath ~= '' then
				path = start.f_getCharData(start.p[1].t_selected[1].ref).arcadepath
			end
			path = hook.runFirst("start.f_selectMode.luaPath", path) or path
			if path ~= '' and path ~= main.luaPath then
				if not main.f_fileExists(path) then
					panicError("\n" .. start.f_getCharData(start.p[1].t_selected[1].ref).name .. " arcadepath doesn't exist: " .. path .. "\n")
				end
			end
		end
		local customArcadePath = main.charparam.arcadepath and path ~= main.luaPath
		--first match
		if start.reset then
			-- Save current remap state. main.f_restoreInput() should restore to this.
			main.f_saveBaseRemapInput()
			if customArcadePath then
				main.t_availableChars = main.f_tableCopy(main.t_orderChars.default)
			else
				main.t_availableChars = main.f_tableCopy(start.f_getOrderChars())
			end
			--generate default roster
			if main.makeRoster and not customArcadePath then
				start.t_roster = start.f_makeRoster()
			else
				start.t_roster = {}
			end
			--generate AI ramping table
			if main.aiRamp then
				start.f_aiRamp(1)
			end
			start.reset = false
		end
		--external script execution
		local oldCustomArcadePath = start.customArcadePath
		start.customArcadePath = customArcadePath
		assert(loadfile(path))()
		start.customArcadePath = oldCustomArcadePath
		--infinite matches flag detected
		if main.makeRoster and start.t_roster[matchNo()] ~= nil and start.t_roster[matchNo()][1] == -1 then
			table.remove(start.t_roster, matchNo())
			start.t_roster = start.f_makeRoster(start.t_roster)
			if main.aiRamp then
				start.f_aiRamp(matchNo())
			end
		--otherwise
		else
			if matchNo() == -1 then --no more matches left
				-- hiscore & stats handled in Go; returns (cleared, place)
				local cleared, place = computeRanking(gameMode())
				if main.motif.hiscore and place > 0 then
					main.f_hiscore(gameMode(), place)
				end
				--credits
				if cleared and main.storyboard.credits and motif.end_credits.enabled and main.f_fileExists(motif.end_credits.storyboard) then
					main.f_storyboard(motif.end_credits.storyboard)
				end
				--game over
				if main.storyboard.gameover and motif.game_over_screen.enabled and main.f_fileExists(motif.game_over_screen.storyboard) then
					if cleared or not main.motif.continuescreen or (not continued() and motif.continue_screen.gameover.enabled) then
						main.f_storyboard(motif.game_over_screen.storyboard)
					end
				end
				--exit to main menu
				if main.exitSelect then
					if motif.files.intro.storyboard ~= '' and not motif.attract_mode.enabled then
						main.f_storyboard(motif.files.intro.storyboard)
					end
				end
				start.exit = start.exit or main.exitSelect or not main.selectMenu[1]
			end
			if start.exit then
				bgReset(motif[main.background].BGDef)
				fadeInInit(motif[main.group].fadein.FadeData)
				playBgm({source = "motif.title", interrupt = true})
				start.exit = false
				return
			end
			if not continued() or esc() then
				start.f_selectReset(false)
			else
				t_reservedChars = {{}, {}}
			end
		end
	end
end

--resets various data
function start.f_selectReset(hardReset, preserveProgress)
	esc(false)
	if not preserveProgress then
		resetGameStats()
		start.matchPersistenceStatsIdx = 0
		setMatchNo(1)
		if main.elimination then
			setWinCount(1, 0)
			setWinCount(2, 0)
		end
		setConsecutiveWins(1, 0)
		setConsecutiveWins(2, 0)
	end
	local col = 1
	local row = 1
	for i = 1, #main.t_selGrid do
		if i > motif.select_info.columns * row then
			row = row + 1
			col = 1
		end
		if main.t_selGrid[i].slot ~= 1 then
			main.t_selGrid[i].slot = 1
			start.t_grid[row][col].char = start.f_selGrid(i).char
			start.t_grid[row][col].char_ref = start.f_selGrid(i).char_ref
			start.t_grid[row][col].hidden = start.f_selGrid(i).hidden
			start.t_grid[row][col].skip = start.f_selGrid(i).skip
		end
		col = col + 1
	end
	if hardReset then
		if motif.select_info.stage.randomselect == 0 or motif.select_info.stage.randomselect == 2 then
			stageListNo = 1
		else
			stageListNo = 0
		end
		restoreCursor = false
		--cursor start cell
		for i = 1, gameOption('Config.Players') do
			if start.f_getCursorData(i).cursor.startcell[1] < motif.select_info.rows then
				start.c[i].selY = start.f_getCursorData(i).cursor.startcell[1]
			else
				start.c[i].selY = 0
			end
			if start.f_getCursorData(i).cursor.startcell[2] < motif.select_info.columns then
				start.c[i].selX = start.f_getCursorData(i).cursor.startcell[2]
			else
				start.c[i].selX = 0
			end
			start.c[i].cell = -1
			start.c[i].randCnt = 0
			start.c[i].randRef = nil
		end
	end
	if stageRandom then
		stageListNo = 0
		stageRandom = false
	end
	for side = 1, 2 do
		if hardReset then
			start.p[side].numSimul = math.max(2, gameOption('Options.Simul.Min'))
			start.p[side].numTag = math.max(2, gameOption('Options.Tag.Min'))
			start.p[side].numTurns = math.max(2, gameOption('Options.Turns.Min'))
			start.p[side].teamMenu = 1
			start.p[side].t_cursor = {}
			start.p[side].teamMode = 0
		end
		start.p[side].numSimul = math.min(start.p[side].numSimul, gameOption('Options.Simul.Max'))
		start.p[side].numTag = math.min(start.p[side].numTag, gameOption('Options.Tag.Max'))
		start.p[side].numTurns = math.min(start.p[side].numTurns, gameOption('Options.Turns.Max'))
		start.p[side].numChars = 1
		start.p[side].teamEnd = main.cpuSide[side] and (side == 2 or not main.cpuSide[1]) and main.forceChar[side] == nil
		start.p[side].selEnd = not main.selectMenu[side]
		start.p[side].t_selected = {}
		start.p[side].t_selTemp = {}
		start.p[side].t_selCmd = {}
		-- Tracks how many leading Turns members are already defeated across matches.
		-- Used by Survival Turns to skip them without removing from roster.
		start.p[side].turnsOffset = 0
		hook.run("start.f_selectReset.side", side, hardReset, preserveProgress, start.p[side])
	end
	for _, v in ipairs(start.c) do
		v.cell = -1
	end
	selScreenEnd = false
	stageEnd = false
	t_reservedChars = {{}, {}}
	cursorActive = {}
	cursorDone = {}
	if main.preload ~= nil and main.preload.charHighlight ~= nil then
		for _, ref in pairs(main.preload.charHighlight) do
			queueCharPreload(ref, 1)
		end
		main.preload.charHighlight = {}
	end
	t_portraitPriority = {1, 1}
	if start.challenger == 0 and not preserveProgress then
		start.t_roster = {}
		start.reset = true
	end
	menu.movelistChar = 1
	hook.run("start.f_selectReset")
end

-- Which command list should drive menu navigation for a given side.
function start.f_menuCmd(side)
	if not main.coop and side == 2 and main.cpuSide[2] then
		return 1
	end
	return side
end

local function makeChallengerResumeSnapshot(pendingFightData, stageNo)
	local resume = {
		p = main.f_tableCopy(start.p),
		c = main.f_tableCopy(start.c),
		roster = main.f_tableCopy(start.t_roster),
		availableChars = main.f_tableCopy(main.t_availableChars),
		pendingFight = main.f_tableCopy(pendingFightData or {}),
		matchNo = matchNo(),
		matchPersistenceStatsIdx = start.matchPersistenceStatsIdx or 0,
		p1ConsecutiveWins = getConsecutiveWins(1),
		p2ConsecutiveWins = getConsecutiveWins(2),
		gameStatsJson = getGameStatsJson(),
		teamarcade = main.teamarcade,
	}
	if stageNo ~= nil then
		resume.pendingFight.stageNo = stageNo
	end
	return resume
end

local function applyWinnerSelectMemory(resume, winnerSide, challengerState)
	if resume == nil or challengerState == nil or winnerSide < 1 or winnerSide > 2 then
		return
	end
	resume.p = resume.p or {{}, {}}
	resume.c = resume.c or {}
	local srcSide = challengerState.p and challengerState.p[winnerSide] or nil
	if srcSide ~= nil then
		resume.p[1] = resume.p[1] or {}
		-- Preserve the winner's last confirmed selection cells / team menu state.
		resume.p[1].t_cursor = main.f_tableCopy(srcSide.t_cursor or {})
		if srcSide.teamMenu ~= nil then resume.p[1].teamMenu = srcSide.teamMenu end
		if srcSide.teamMode ~= nil then resume.p[1].teamMode = srcSide.teamMode end
		if srcSide.numSimul ~= nil then resume.p[1].numSimul = srcSide.numSimul end
		if srcSide.numTag ~= nil then resume.p[1].numTag = srcSide.numTag end
		if srcSide.numTurns ~= nil then resume.p[1].numTurns = srcSide.numTurns end
	end
	if challengerState.c ~= nil then
		-- Promote the winner-side cursor slots onto the new arcade P1 slots:
		-- side 1: 1/3/5/7 -> 1/3/5/7
		-- side 2: 2/4/6/8 -> 1/3/5/7
		for srcPn = winnerSide, gameOption('Config.Players'), 2 do
			if challengerState.c[srcPn] ~= nil then
				local dstPn = srcPn
				if winnerSide == 2 then
					dstPn = srcPn - 1
				end
				resume.c[dstPn] = main.f_tableCopy(challengerState.c[srcPn])
			end
		end
	end
	hook.run("start.f_selectChallenger.resume", resume, winnerSide, challengerState)
end

function start.f_selectChallenger(resume)
	esc(false)
	resume = resume or makeChallengerResumeSnapshot()
	local challengerCmd = start.challenger
	local arcadeP1Controller = getRemapInput(1)
	local challengerController = getRemapInput(challengerCmd)

	-- Start challenger match from a clean state.
	main.f_default()

	-- Build a clean, non-chained input mapping:
	--   cmd list 1 -> arcade P1 physical controller
	--   cmd list 2 -> challenger physical controller
	-- Keep the rest as a swap-based permutation (no dual-feeding / chaining).
	local function swapCmd(a, b)
		if a == b then return end
		local ra = getRemapInput(a)
		local rb = getRemapInput(b)
		remapInput(a, rb)
		remapInput(b, ra)
	end
	-- Put arcade P1 controller on cmd 1.
	swapCmd(1, arcadeP1Controller)
	-- Find which cmd currently owns challengerController and swap it onto cmd 2.
	local holder = 2
	for i = 1, gameOption('Config.Players') do
		if getRemapInput(i) == challengerController then
			holder = i
			break
		end
	end
	swapCmd(2, holder)
	-- Keep challenger flag (>0) so versus() enters challenger mode
	start.challenger = challengerCmd

	main.t_itemname.versus()

	start.f_selectReset(false)
	if not start.f_selectScreen() then
		start.exit = true
		start.challenger = 0
		main.f_restoreInput()
		return false
	end
	local ok = launchFight{challenger = true}
	local challengerSelectState = {
		p = main.f_tableCopy(start.p),
		c = main.f_tableCopy(start.c),
	}
	local challengerWinner = getWinnerTeam()

	start.challenger = 0
	if not ok or challengerWinner < 1 or challengerWinner > 2 or start.exit or esc() then
		main.f_restoreInput()
		return false
	end
	-- Winner of the challenger match becomes the new arcade P1 owner.
	local resumeController = arcadeP1Controller
	if challengerWinner == 2 then
		resumeController = challengerController
	end

	-- Rebuild arcade mode with the winner's controller mapped to P1.
	main.teamarcade = resume.teamarcade
	main.f_default()
	setLastInputController(resumeController)
	main.t_itemname.arcade()
	main.f_saveBaseRemapInput()

	-- Restore interrupted arcade progress before going back to select screen.
	applyWinnerSelectMemory(resume, challengerWinner, challengerSelectState)
	start.p = main.f_tableCopy(resume.p)
	start.c = main.f_tableCopy(resume.c)
	start.t_roster = main.f_tableCopy(resume.roster)
	main.t_availableChars = main.f_tableCopy(resume.availableChars)
	setGameStatsJson(resume.gameStatsJson)
	setMatchNo(resume.matchNo)
	start.matchPersistenceStatsIdx = resume.matchPersistenceStatsIdx or 0
	setConsecutiveWins(1, resume.p1ConsecutiveWins)
	setConsecutiveWins(2, resume.p2ConsecutiveWins)
	start.reset = false
	start.exit = false

	-- Preserve the interrupted opponent. P1 is cleared for re-selection,
	-- but P2 stays selected so resuming does not roll a new random opponent.
	local p2 = main.f_tableCopy(start.p[2])
	start.f_selectReset(false, true)
	start.p[2] = p2
	start.reset = false

	if not start.f_selectScreen() then
		start.exit = true
		return false
	end

	-- Resume the exact interrupted arcade fight using the original caller args.
	return launchFight(resume.pendingFight)
end

local function buildMusicParams(data)
	local out = {}
	for k, v in pairs(data) do
		if type(k) == "string" and k:match("music$") then
			if type(v) == "string" then
				out[#out + 1] = k .. "=" .. v
			elseif type(v) == "table" and #v > 0 then
				local first = v[1]
				-- If table looks like { "path.mp3", 100, 123, 456 } => positional args
				local positional = (type(first) == "string") and (#v == 1 or type(v[2]) ~= "string")
				if positional then
					local pieces = {}
					for i = 1, #v do
						pieces[i] = tostring(v[i])
					end
					-- space-separated to avoid commas inside the value
					out[#out + 1] = k .. "=" .. table.concat(pieces, " ")
				else
					-- Treat as multiple candidate tracks: {"a.mp3","b.mp3",...}
					for i = 1, #v do
						out[#out + 1] = k .. "=" .. tostring(v[i])
					end
				end
			end
		end
	end
	return table.concat(out, ", ")
end

function launchFight(data)
	local data = data or {}
	local t = {}
	if continued() then -- on rematch all arguments are ignored and values are restored from last match
		t = main.f_tableCopy(start.launchFightSav)
		start.p[2].t_selTemp = {} -- in case it's not cleaned already (preserved p2 side during select screen)
	else -- otherwise take all arguments and settings into account
		t.p1numchars = start.p[1].numChars
		t.p1teammode = start.p[1].teamMode
		t.p2numchars = start.p[2].numChars
		t.p2teammode = start.p[2].teamMode
		t.challenger = main.f_arg(data.challenger, false)
		t.continue = main.f_arg(data.continue, main.motif.continuescreen)
		t.quickcontinue = (not main.selectMenu[1] and not main.selectMenu[2]) or main.f_arg(data.quickcontinue, main.quickContinue or gameOption('Options.QuickContinue'))
		t.order = data.order or 1
		t.orderselect = {main.f_arg(data.p1orderselect, main.orderSelect[1]), main.f_arg(data.p2orderselect, main.orderSelect[2])}
		t.p1char = data.p1char or {}
		t.p1pal = data.p1pal
		t.p1rounds = data.p1rounds or nil
		t.p2char = data.p2char or {}
		t.p2pal = data.p2pal
		t.p2rounds = data.p2rounds or nil
		t.exclude = data.exclude or {}
		t.musicParams = buildMusicParams(data)
		t.stage = data.stage or ''
		t.ai = data.ai or nil
		t.vsscreen = main.f_arg(data.vsscreen, main.motif.vsscreen)
		t.victoryscreen = main.f_arg(data.victoryscreen, main.motif.victoryscreen)
		t.winscreen = main.f_arg(data.winscreen, main.motif.winscreen)
		--t.frames = data.frames or fightScreenVar("time.framespercount")
		t.roundtime = data.time or nil
		t.lua = data.lua or ''
		t.stageNo = start.f_getStageRef(t.stage)
		t.stageAssigned = data.stageNo ~= nil or t.stage ~= ''
		t.stageNo = data.stageNo or start.f_getStageRef(t.stage)
		-- Absolute isolation guard for randomtierladder. The selected refs are
		-- authoritative; no fight may proceed when their tier letters differ.
		if gameMode() == 'randomtierladder' and data.p1ref ~= nil and data.p2ref ~= nil then
			local function tierLetter(ref)
				local record = start.f_getCharRecord(ref)
				local tier = tostring(start.f_getRecordTier(record) or record.tier or 'U'):upper()
				local group, suffix = tier:match('^([UFDCBASXZ])([+%-]*)$')
				if group == nil or #suffix > 3 or ((group == 'U' or group == 'Z') and suffix ~= '') then
					return nil
				end
				return group
			end
			local p1tier = tierLetter(data.p1ref)
			local p2tier = tierLetter(data.p2ref)
			if not main.f_sameTierLock(p1tier, p2tier, getCommandLineValue('-tierlockoverride') ~= nil) then
				printConsole('launchFight: BLOCKED randomtierladder cross-tier pair [' .. p1tier .. '] vs [' .. p2tier .. ']')
				return false
			end
		end
		start.p[1].numChars = data.p1numchars or math.max(start.p[1].numChars, #t.p1char)
		start.p[1].teamMode = start.f_stringToTeamMode(data.p1teammode) or start.p[1].teamMode
		start.p[2].numChars = data.p2numchars or math.max(start.p[2].numChars, #t.p2char)
		start.p[2].teamMode = start.f_stringToTeamMode(data.p2teammode) or start.p[2].teamMode
		t.p1numchars = start.f_matchPersistence()
		-- add P1 chars forced via function arguments (ignore char param restrictions)
		local reset = false
		local cnt = 0
		for _, v in main.f_sortKeys(t.p1char) do
			if not reset then
				start.p[1].t_selected = {}
				start.p[1].t_selTemp = {}
				reset = true
			end
			cnt = cnt + 1
			local ref = (cnt == 1 and data.p1ref ~= nil) and data.p1ref or start.f_getCharRef(v)
			table.insert(start.p[1].t_selected, {
				ref = ref,
				pal = t.p1pal or start.f_selectPal(ref),
				pn = start.f_getPlayerNo(1, #start.p[1].t_selected + 1),
				--cursor = {},
			})
			hook.run("start.launchFight.selected", 1, #start.p[1].t_selected, start.p[1].t_selected[#start.p[1].t_selected], start.p[1], t)
			main.t_availableChars = start.f_excludeChar(main.t_availableChars, ref)
		end
		if #start.p[1].t_selected == 0 then
			panicError("\n" .. "launchFight(): no valid P1 characters\n")
			start.exit = true
			return false -- return to main menu
		end
		-- add P2 chars forced via function arguments (ignore char param restrictions)
		local onlyme = false
		cnt = 0
		for _, v in main.f_sortKeys(t.p2char) do
			cnt = cnt + 1
			local ref = (cnt == 1 and data.p2ref ~= nil) and data.p2ref or start.f_getCharRef(v)
			table.insert(start.p[2].t_selected, {
				ref = ref,
				pal = t.p2pal or start.f_selectPal(ref),
				pn = start.f_getPlayerNo(2, #start.p[2].t_selected + 1),
				--cursor = {},
			})
			hook.run("start.launchFight.selected", 2, #start.p[2].t_selected, start.p[2].t_selected[#start.p[2].t_selected], start.p[2], t)
			main.t_availableChars = start.f_excludeChar(main.t_availableChars, ref)
			if not onlyme then onlyme = start.f_getCharData(ref).single end
		end
		-- add remaining P2 chars of particular order if there are still free slots in the selected team mode
		if main.cpuSide[2] and #start.p[2].t_selected < start.p[2].numChars and not onlyme then
			-- get list of available chars
			local t_chars = main.f_tableCopy(main.t_availableChars)
			-- remove chars temporary excluded from this match
			for _, v in ipairs(t.exclude) do
				t_chars = start.f_excludeChar(t_chars, start.f_getCharRef(v))
			end
			-- remove chars with 'single' param if some characters are forced into team
			if #start.p[2].t_selected > 0 then
				for _, v in ipairs(t_chars[t.order]) do
					if start.f_getCharData(v).single then
						t_chars = start.f_excludeChar(t_chars, v)
					end
				end
			end
			-- fill free slots
			local t_remaining = main.f_tableCopy(t_chars)
			local t_tmp = {}
			for i = #start.p[2].t_selected, start.p[2].numChars - 1 do
				if t_chars[t.order] ~= nil and #t_chars[t.order] > 0 then
					local rand = math.random(1, #t_chars[t.order])
					local ref = t_chars[t.order][rand]
					if not start.f_getCharData(ref).single then
						table.remove(t_chars[t.order], rand)
						table.insert(t_tmp, ref)
					else --one entry if 'single' param is detected on any opponent
						t_tmp = {ref}
						onlyme = true
						break
					end
				end
			end
			-- not enough unique characters of particular order, take into account only if skiporder parameter = false
			while not t.skiporder and #t_tmp + #start.p[2].t_selected < start.p[2].numChars and not onlyme and t_remaining[t.order] ~= nil and #t_remaining[t.order] > 0 do
				table.insert(t_tmp, t_remaining[t.order][math.random(1, #t_remaining[t.order])])
			end
			-- append remaining characters
			for _, v in ipairs(t_tmp) do
				table.insert(start.p[2].t_selected, {
					ref = v,
					pal = start.f_selectPal(v),
					pn = start.f_getPlayerNo(2, #start.p[2].t_selected + 1),
					--cursor = {},
				})
				hook.run("start.launchFight.selected", 2, #start.p[2].t_selected, start.p[2].t_selected[#start.p[2].t_selected], start.p[2], t)
				main.t_availableChars = start.f_excludeChar(main.t_availableChars, v)
			end
			-- team conversion if 'single' param is set on randomly added chars
			if onlyme and #start.p[2].t_selected > 1 then
				panicError("Unexpected launchFight state.\nPlease write down everything that lead to this error and report it to K4thos.\n")
				--[[for i = 1, #start.p[2].t_selected do
					if not start.f_getCharData(start.p[2].t_selected[i].ref).single then
						table.insert(main.t_availableChars[t.order], start.p[2].t_selected[i].ref)
						table.remove(start.p[2].t_selected, k)
					end
				end]]
			end
		end
		if onlyme then
			start.p[2].numChars = #start.p[2].t_selected
		end
		-- skip match if needed
		if #start.p[2].t_selected < start.p[2].numChars then
			start.p[2].t_selected = {}
			start.p[2].t_selTemp = {}
			printConsole("launchFight(): not enough P2 characters, skipping execution")
			setMatchNo(matchNo() + 1)
			return true --continue lua code execution
		end
	end
	clearSelected()
	local ok = false
	local loopCount = 0
	while true do
		-- fight initialization
		setTeamMode(1, start.p[1].teamMode, start.p[1].numChars)
		setTeamMode(2, start.p[2].teamMode, start.p[2].numChars)
		start.f_remapAI(t.ai)
		start.f_setRounds(t.roundtime, {t.p1rounds, t.p2rounds})
		t.stageNo = start.f_setStage(t.stageNo, t.stageAssigned or continued() or loopCount > 0)
		local challengerResume = nil
		if not t.challenger then
			-- Snapshot before game() runs. If a challenger interrupts the fight, this is the last clean arcade state.
			challengerResume = makeChallengerResumeSnapshot(data, t.stageNo)
		end
		local common = {lua = {}}
		if t.lua ~= '' then
			table.insert(common.lua, t.lua)
		end
		-- Hooks may mutate "common" in place before loading starts.
		hook.run("launchFight", common, t, data)
		updateCommon(common, true)
		-- Resolve match-scoped params before VS can start background loading.
		local winscreen = main.f_arg(t.winscreen, main.motif.winscreen)
		if winscreen and data.winscreen == nil then
			if start.customArcadePath or (main.makeRoster and start.t_roster[matchNo() + 1] ~= nil) then
				winscreen = false
			end
		end
		local loadStartParams = main.f_tableCopy(t)
		loadStartParams.winscreen = winscreen

		if not start.f_selectVersus(t.vsscreen, t.orderselect, loadStartParams) then break end
		-- If VS started background loading, do not restart the loader here.
		if gameOption('Config.VsScreenLoading') and start.bgLoadStarted then
			clearAllSound()
		elseif not start.f_selectLoading(loadStartParams) then
			break
		end
		start.f_game(common)
		clearColor(motif.selectbgdef.bgclearcolor[1], motif.selectbgdef.bgclearcolor[2], motif.selectbgdef.bgclearcolor[3])
		if start.exit or start.characterchange then
			start.characterchange = false
			break
		-- here comes a new challenger
		elseif start.challenger > 0 then
			if t.challenger then -- end function called by f_arcadeChallenger() regardless of outcome
				ok = not start.exit and not esc()
				break
			else
				return start.f_selectChallenger(challengerResume)
			end
		-- player exit the game via ESC
		elseif getWinnerTeam() == -1 then
			if not main.selectMenu[1] and not main.selectMenu[2] then
				setMatchNo(-1)
			end
			break
		-- player lost in modes that ends after 1 lose
		elseif getWinnerTeam() ~= 1 and main.elimination then
			setMatchNo(-1)
			break
		-- player won or continuing is disabled
		elseif getWinnerTeam() == 1 or not t.continue then
			start.p[2].t_selected = {}
			start.p[2].t_selTemp = {}
			setMatchNo(matchNo() + 1)
			ok = true -- continue lua code execution
			break
		-- continue = no
		elseif not continued() then
			setMatchNo(-1)
			break
		-- continue = yes
		elseif not t.quickcontinue then -- if 'Quick Continue' is disabled
			for i = 1, 2 do
				for _, v in ipairs(start.p[i].t_selCmd) do
					v.selectState = 0
				end
			end
			start.p[1].t_selected = {}
			start.p[1].t_selTemp = {}
			start.p[1].selEnd = false
			start.launchFightSav = main.f_tableCopy(t)
			--start.p[2].t_selTemp = {} -- uncomment to disable enemy team showing up in select screen
			selScreenEnd = false
			return
		end
		start.challenger = 0
		loopCount = loopCount + 1
	end
	-- restore original values
	start.p[1].numChars = t.p1numchars
	start.p[1].teamMode = t.p1teammode
	start.p[2].numChars = t.p2numchars
	start.p[2].teamMode = t.p2teammode
	return ok
end

function launchStoryboard(path)
	if path == nil or path == '' then
		return false
	end
	main.f_storyboard(path)
	return true
end

function codeInput(name)
	return commandGetState(getLastInputController(), name)
end

--;===========================================================
--; SELECT SCREEN
--;===========================================================
local function refreshActiveFacePortraits()
	for side = 1, 2 do
		for k, v in ipairs(start.p[side].t_selCmd) do
			local member = main.f_tableLength(start.p[side].t_selected) + k
			if main.coop and (side == 1 or gameMode('versuscoop')) then
				member = k
			end
			local st = start.p[side].t_selTemp[member]
			local player = v.player
			local selRef = start.c[player].selRef
			if v.selectState == 0 and st ~= nil and selRef ~= nil and st.ref == selRef then
				local state = getCharPreloadStatus(selRef)
				if state == 'ready' then
					local pn = 2 * (member - 1) + side
					local pCfg = f_getMotifP(motif.select_info, pn, side)
					local updated = false
					if st.face_data == nil then
						st.face_anim = pCfg.face.anim
						st.face_data = start.f_animGet(selRef, side, member, pCfg.face, nil, true, st.face_data)
						updated = updated or st.face_data ~= nil
					end
					if st.face2_data == nil then
						st.face2_anim = pCfg.face2.anim
						st.face2_data = start.f_animGet(selRef, side, member, pCfg.face2, nil, true, st.face2_data)
						updated = updated or st.face2_data ~= nil
					end
					if updated then
						start.needUpdateDrawList = true
					end
				end
			end
		end
	end
end

function start.updateDrawList()
	local drawList = {}

	for row = 1, motif.select_info.rows do
		for col = 1, motif.select_info.columns do
			local cellIndex = (row - 1) * motif.select_info.columns + col
			local t = start.t_grid[row][col]
			local c = col - 1
			local r = row - 1

			if t.skip ~= 1 then
				local charData = start.f_selGrid(cellIndex)
				local function getTransforms(base)
					return {
						facing      = getCellFacing(base.facing, c, r),
						scale       = getCellTransform(c, r, "scale", base.scale),
						xshear      = getCellTransform(c, r, "xshear", base.xshear),
						angle       = getCellTransform(c, r, "angle", base.angle),
						xangle      = getCellTransform(c, r, "xangle", base.xangle),
						yangle      = getCellTransform(c, r, "yangle", base.yangle),
						projection  = getCellTransform(c, r, "projection", base.projection),
						focallength = getCellTransform(c, r, "focallength", base.focallength)
					}
				end

				if (charData and charData.char ~= nil and (charData.hidden == 0 or charData.hidden == 3)) or motif.select_info.showemptyboxes then
					local item = getTransforms(motif.select_info.cell.bg)
					item.anim = motif.select_info.cell.bg.AnimData
					item.x = motif.select_info.pos[1] + t.x
					item.y = motif.select_info.pos[2] + t.y
					table.insert(drawList, item)
				end

				if charData and (charData.char == 'randomselect' or charData.hidden == 3) then
					local item = getTransforms(motif.select_info.cell.random)
					item.anim = motif.select_info.cell.random.AnimData
					item.x = motif.select_info.pos[1] + t.x + motif.select_info.portrait.offset[1]
					item.y = motif.select_info.pos[2] + t.y + motif.select_info.portrait.offset[2]
					table.insert(drawList, item)
				end

				if charData and charData.char_ref ~= nil and charData.hidden == 0 and charData.char ~= 'randomselect' then
					local portrait = motif.select_info.portrait
					local loadingPortrait = false
					if getCharPreloadStatus(charData.char_ref) ~= 'ready' and hasPortraitAnim(portrait.loading) then
						portrait = portrait.loading
						loadingPortrait = true
					end
					local item = getTransforms(portrait)
					item.anim = loadingPortrait and portrait.AnimData or charData.cell_data
					item.x = motif.select_info.pos[1] + t.x
					item.y = motif.select_info.pos[2] + t.y
					if not loadingPortrait then
						item.x = item.x + portrait.offset[1]
						item.y = item.y + portrait.offset[2]
					end
					-- apply cell scale override while preserving portrait resolution factor
					-- loading portrait comes from system.sff, so don't apply character localcoord scaling to it
					if item.scale ~= nil and not loadingPortrait then
						local charInfo = main.t_selChars[charData.char_ref + 1]
						if charInfo then
							local portraitScale = charInfo.portraitscale or 1
							local charLocalcoord = charInfo.localcoord or motif.info.localcoord[1]
							-- recompute resolution compensation factor
							local resFix = portraitScale * motif.info.localcoord[1] / charLocalcoord
							item.scale = {
								item.scale[1] * resFix,
								item.scale[2] * resFix
							}
						end
					end
					table.insert(drawList, item)
				end
			end
		end
	end

	return drawList
end

local function tickScreenDelay(side)
	if start.p[side].screenDelay <= 0 then
		return false
	end
	local selectComplete = start.p[1].selEnd and start.p[2].selEnd and start.p[1].teamEnd and start.p[2].teamEnd
	local canSkip = not start.p[side].inPalMenu and (
		not start.p[side].selEnd
		or (selectComplete and (not main.stageMenu or stageEnd))
	)
	if canSkip then
		local owners = {}
		local seen = {}
		if main.coop and (side == 1 or gameMode('versuscoop')) then
			for _, v in ipairs(start.p[side].t_selCmd) do
				if v.cmd ~= nil and not seen[v.cmd] then
					table.insert(owners, v.cmd)
					seen[v.cmd] = true
				end
			end
		else
			local cmd = start.f_menuCmd(side)
			if cmd ~= nil then
				table.insert(owners, cmd)
			end
		end
		if #owners > 0 and getInput(owners, motif.select_info.done.key) then
			start.p[side].screenDelay = 0
			return true
		end
	end
	start.p[side].screenDelay = start.p[side].screenDelay - 1
	return false
end

start.needUpdateDrawList = false
function start.f_selectScreen()
	if (not main.selectMenu[1] and not main.selectMenu[2]) or selScreenEnd then
		return true
	end
	bgReset(motif.selectbgdef.BGDef)
	fadeInInit(motif.select_info.fadein.FadeData)
	local fadeOutStarted = false
	playBgm({source = "motif.select", interrupt = true})
	start.f_resetTempData(motif.select_info, 'face')
	f_snapCursor()
	local stageActiveCount = 0
	local stageActiveState = false
	timerSelect = 0
	start.escFlag = false
	local t_teamMenu = {{}, {}}
	local blinkCount = 0
	local counter = 0 - motif.select_info.fadein.time
	local timerReset = false
	local stageTextData = motif.select_info.stage.active.TextSpriteData
	-- generate team mode items table
	for side = 1, 2 do
		-- read display names for the current gameMode (or default)
		local params = motif.select_info.teammenu.itemname.default
		if motif.select_info.teammenu.itemname[gameMode()] ~= nil then
			params = motif.select_info.teammenu.itemname[gameMode()]
		end
		-- read itemname_order for the current gameMode (or default)
		local itemname_order = motif.select_info.teammenu.itemname_order.default
		if motif.select_info.teammenu.itemname_order[gameMode()] ~= nil then
			itemname_order = motif.select_info.teammenu.itemname_order[gameMode()]
		end
		-- map itemname -> mode (kept from old defaults)
		local modeByName = {
			single = 0,
			simul  = 1,
			turns  = 2,
			tag    = 3,
		}
		-- itemname_order lists exactly what to render, in the correct order
		for _, name in ipairs(itemname_order) do
			local itemname = name
			local mode = modeByName[itemname]
			if mode ~= nil and main.teamMenu[side][itemname] then
				table.insert(t_teamMenu[side], {
					itemname    = itemname,
					displayname = params[itemname],
					mode        = mode,
				})
			end
		end
		hook.run("start.selectScreen.teamMenu", side, t_teamMenu[side], params, itemname_order)
	end

	textImgReset(motif.select_info.record.TextSpriteData)
	textImgSetText(motif.select_info.record.TextSpriteData, start.f_getRecordText())

	local staticDrawList = start.updateDrawList()
	start.needUpdateDrawList = false

	while not selScreenEnd do
		main.f_preloadTick(4)
		refreshActiveFacePortraits()
		counter = counter + 1
		--draw clearcolor
		clearColor(motif.selectbgdef.bgclearcolor[1], motif.selectbgdef.bgclearcolor[2], motif.selectbgdef.bgclearcolor[3])
		--draw layerno = 0 backgrounds
		bgDraw(motif.selectbgdef.BGDef, 0)
		--draw title
		textImgDraw(motif.select_info.title.TextSpriteData)
		--draw portraits
		for side = 1, 2 do
			if #start.p[side].t_selTemp > 0 then
				start.f_drawPortraits(start.p[side].t_selTemp, side, motif.select_info, 'face', true)
			end
		end
		--draw cell art
		if start.needUpdateDrawList then
			staticDrawList = start.updateDrawList()
			start.needUpdateDrawList = false 
		end
		batchDraw(staticDrawList)
		--draw done cursors
		for side = 1, 2 do
			local persist = motif.select_info['p' .. side].cursor.persist 
			local totalSelected = #start.p[side].t_selected
			local drawnCells = {} -- Track drawn cells to avoid duplicates
			for k, v in pairs(start.p[side].t_selected) do
				if v.cursor ~= nil then
					--get cell coordinates
					local x = v.cursor[1]
					local y = v.cursor[2]
					local t = start.t_grid[y + 1][x + 1]
					--retrieve proper cell coordinates in case of random selection
					--TODO: doesn't work with slot feature
					--if (t.char == 'randomselect' or t.hidden == 3) --[[and not gameOption('Options.Team.Duplicates')]] then
					--	x = start.f_getCharData(v.ref).col - 1
					--	y = start.f_getCharData(v.ref).row - 1
					--	t = start.t_grid[y + 1][x + 1]
					--end
					--render only if cell is not hidden
					if t.hidden ~= 1 and t.hidden ~= 2 then
						local shouldDraw = false
						if main.coop or persist then
							shouldDraw = true
						end
						if start.p[side].selEnd and k == totalSelected then
							shouldDraw = true
						end
						if shouldDraw then
							local cellKey = x .. ',' .. y
							if not drawnCells[cellKey] then
								start.f_drawCursor(v.pn, x, y, 'done', true)
								drawnCells[cellKey] = true
							end
						end
					end
				end
			end
		end
		--team and select menu
		if blinkCount < motif.select_info.p2.cursor.switchtime then
			blinkCount = blinkCount + 1
		else
			blinkCount = 0
		end
		local screenDelayInterrupted = false
		for side = 1, 2 do
			if not start.p[side].teamEnd then
				start.f_teamMenu(side, t_teamMenu[side])
			elseif not start.p[side].selEnd then
				--for each player with active controls
				for k, v in ipairs(start.p[side].t_selCmd) do
					local member = main.f_tableLength(start.p[side].t_selected) + k
					if main.coop and (side == 1 or gameMode('versuscoop')) then
						member = k
					end
					--member selection
					local drawUpdateFlag
					v.selectState, drawUpdateFlag = start.f_selectMenu(side, v.cmd, v.player, member, v.selectState)
					if drawUpdateFlag then
						start.needUpdateDrawList = true
					end
					--draw active cursor
					if side == 2 and motif.select_info.p2.cursor.blink then
						local sameCell = false
						for _, v2 in ipairs(start.p[1].t_selCmd) do							
							if start.c[v.player].cell == start.c[v2.player].cell and v.selectState == 0 and v2.selectState == 0 then
								if blinkCount == 0 then
									start.c[v.player].blink = not start.c[v.player].blink
								end
								sameCell = true
								break
							end
						end
						if not sameCell then
							start.c[v.player].blink = false
						end
					end
					local cursorState = 'active'
					if v.selectState > 0 and motif.select_info.paletteselect > 0 then
						local cursorData = motif.select_info['p' .. side].cursor
						if cursorData.preview and cursorData.preview.default and 
						(cursorData.preview.default.anim ~= -1 or cursorData.preview.default.spr[1] ~= -1) then
						--cursorState when palmenu is active
							cursorState = 'preview'
						else
							cursorState = 'done'
						end
					end
					if v.selectState < 4 and start.f_selGrid(start.c[v.player].cell + 1).hidden ~= 1 and not start.c[v.player].blink then
						start.f_drawCursor(v.player, start.c[v.player].selX, start.c[v.player].selY, cursorState, false)
					end
				end
			end
			--delayed screen transition for the duration of face_done_anim or selection sound
			if tickScreenDelay(side) then
				screenDelayInterrupted = true
			end
			--exit select screen
			for _, v in ipairs(start.p[side].t_selCmd) do
				if not start.escFlag and (esc() or (getInput(v.cmd, motif.select_info.cancel.key) and not start.p[side].inPalMenu)) then
					sndPlay(motif.Snd, motif.select_info.cancel.snd[1], motif.select_info.cancel.snd[2])
					fadeOutInit(motif.select_info.fadeout.FadeData)
					fadeOutStarted = true
					start.escFlag = true
				end
			end
			if start.p[side].inPalMenu then
				local palActive = false
				if motif.select_info.paletteselect > 0 then
					for _, sv in ipairs(start.p[side].t_selCmd) do
						if sv.selectState == 1 then
							palActive = true
							break
						end
					end
				end
				if not palActive then
					start.p[side].inPalMenu = false
				end
			end
		end
		--draw names
		for side = 1, 2 do
			if #start.p[side].t_selTemp > 0 then
				for i = 1, #start.p[side].t_selTemp do
					if i <= motif.select_info['p' .. side].name.num then
						local name = ''
						if motif.select_info['p' .. side].name.num == 1 then
							name = start.f_getName(start.p[side].t_selTemp[#start.p[side].t_selTemp].ref, side)
						else
							name = start.f_getName(start.p[side].t_selTemp[i].ref, side)
						end
						textImgReset(motif.select_info['p' .. side].name.TextSpriteData)
						textImgAddPos(
							motif.select_info['p' .. side].name.TextSpriteData,
							(i - 1) * motif.select_info['p' .. side].name.spacing[1],
							(i - 1) * motif.select_info['p' .. side].name.spacing[2]
						)
						textImgSetText(motif.select_info['p' .. side].name.TextSpriteData, name)
						textImgDraw(motif.select_info['p' .. side].name.TextSpriteData)
					end
				end
			end
		end
		--team and character selection complete
		if start.p[1].selEnd and start.p[2].selEnd and start.p[1].teamEnd and start.p[2].teamEnd then
			restoreCursor = true
			if main.stageMenu and not stageEnd then --Stage select
				start.p[1].screenDelay, start.p[2].screenDelay = 0, 0
				start.f_stageMenu()
				if not timerReset then
					timerSelect = motif.select_info.timer.displaytime
					timerReset = true
				end
			elseif start.p[1].screenDelay <= 0 and start.p[2].screenDelay <= 0 and not fadeOutStarted then
				fadeOutInit(motif.select_info.fadeout.FadeData)
				fadeOutStarted = true
			end
			--draw stage portrait
			if main.stageMenu then
				--draw stage portrait background
				main.f_animPosDraw(motif.select_info.stage.portrait.bg.AnimData)
				--draw stage portrait (random)
				if stageListNo == 0 then
					main.f_animPosDraw(motif.select_info.stage.portrait.random.AnimData)
				--draw stage portrait loaded from stage SFF
				else
					local stageRef = main.t_selectableStages[stageListNo]
					local portrait = motif.select_info.stage.portrait
					local anim = main.t_selStages[stageRef].anim_data
					local loadingPortrait = false
					if getStagePreloadStatus(stageRef) ~= 'ready' and hasPortraitAnim(portrait.loading) then
						portrait = portrait.loading
						anim = portrait.AnimData
						loadingPortrait = true
					end
					local x = motif.select_info.stage.pos[1]
					local y = motif.select_info.stage.pos[2]
					if not loadingPortrait then
						x = x + portrait.offset[1]
						y = y + portrait.offset[2]
					end
					main.f_animPosDraw(anim, x, y)
				end
				if not stageEnd then
					local canConfirmStage = (getInput(-1, motif.select_info.done.key) and not screenDelayInterrupted) or timerSelect == -1
					if canConfirmStage then
						local preloadReady = true
						if stageListNo > 0 then
							local stageRef = main.t_selectableStages[stageListNo]
							local state = getStagePreloadStatus(stageRef)
							preloadReady = state == 'ready'
							if not preloadReady then
								main.f_preloadBoostStage(stageRef)
							else
								if main.f_materializeStagePortrait(stageRef) then
									start.needUpdateDrawList = true
								end
							end
						end
						if not preloadReady and getInput(-1, motif.select_info.done.key) and not screenDelayInterrupted then
							sndPlay(motif.Snd, motif.select_info.cancel.snd[1], motif.select_info.cancel.snd[2])
						end
						if preloadReady then
							sndPlay(motif.Snd, motif.select_info.stage.done.snd[1], motif.select_info.stage.done.snd[2])
							stageTextData = motif.select_info.stage.done.TextSpriteData
							stageEnd = true
						end
					elseif stageActiveCount < motif.select_info.stage.active.switchtime then --delay change
						stageActiveCount = stageActiveCount + 1
					else
						if stageActiveState then
							stageActiveState = false
							stageTextData = motif.select_info.stage.active2.TextSpriteData
						else
							stageActiveState = true
							stageTextData = motif.select_info.stage.active.TextSpriteData
						end
						stageActiveCount = 0
					end
				end
				--draw stage name
				local stage_text = motif.select_info.stage.random.text
				if stageListNo ~= 0 then
					stage_text = main.f_formatBySpec(motif.select_info.stage.text, {i = stageListNo, s = main.t_selStages[main.t_selectableStages[stageListNo]].name})
				end
				textImgReset(stageTextData)
				textImgSetText(stageTextData, stage_text)
				textImgDraw(stageTextData)
			end
		else
			--draw record text
			textImgDraw(motif.select_info.record.TextSpriteData)
		end
		--draw timer
		if motif.select_info.timer.count ~= -1 and (not start.p[1].teamEnd or not start.p[2].teamEnd or not start.p[1].selEnd or not start.p[2].selEnd or (main.stageMenu and not stageEnd)) and counter >= 0 then
			timerSelect = main.f_drawTimer(timerSelect, motif.select_info.timer)
		end
		-- hook
		hook.run("start.f_selectScreen")
		--draw layerno = 1 backgrounds
		bgDraw(motif.selectbgdef.BGDef, 1)
		start.f_drawSelectStatsOverlay(counter)
		--frame transition
		if not fadeActive() and (fadeOutStarted or start.escFlag) then
			selScreenEnd = true
			break --skip last frame rendering
		end
		refresh()
	end
	return not start.escFlag
end

--;===========================================================
--; TEAM MENU
--;===========================================================
local t_teamActiveCount = {0, 0}
local t_teamActiveState = {false, false}

function start.f_teamMenu(side, t)
	if #t == 0 then
		start.p[side].teamEnd = true
		-- Team menu has no renderable entries (e.g. itemname_order hides them).
		-- Still allow character selection for this side if enabled.
		if not start.p[side].selEnd and #start.p[side].t_selCmd == 0 then
			table.insert(start.p[side].t_selCmd, {cmd = start.f_menuCmd(side), player = side, selectState = 0})
		end
		return
	end
	--skip selection if only 1 team mode is available and team size is fixed
	if #t == 1 and (t[1].itemname == 'single' or (t[1].itemname == 'simul' and main.numSimul[1] == main.numSimul[2]) or (t[1].itemname == 'turns' and main.numTurns[1] == main.numTurns[2]) or (t[1].itemname == 'tag' and main.numTag[1] == main.numTag[2])) then
		if t[1].itemname == 'single' then
			start.p[side].numChars = 1
		elseif t[1].itemname == 'simul' then
			start.p[side].numChars = start.p[side].numSimul
		elseif t[1].itemname == 'turns' then
			start.p[side].numChars = start.p[side].numTurns
		elseif t[1].itemname == 'tag' then
			start.p[side].numChars = start.p[side].numTag
		end
		start.p[side].teamMode = t[1].mode
		start.p[side].teamEnd = true
	--otherwise display team mode selection
	else
		--Commands
		local t_cmd = {}
		if main.coop then
			for i = 1, gameOption('Config.Players') do
				if not gameMode('versuscoop') or (i - 1) % 2 + 1 == side then
					table.insert(t_cmd, i)
				end
			end
		else
			t_cmd = {start.f_menuCmd(side)}
		end
		--Calculate team cursor position
		if start.p[side].teamMenu > #t then
			start.p[side].teamMenu = 1
		end
		if #t > 1 and getInput(t_cmd, motif.select_info['p' .. side].teammenu.previous.key) then
			if start.p[side].teamMenu > 1 then
				sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.move.snd[1], motif.select_info['p' .. side].teammenu.move.snd[2])
				start.p[side].teamMenu = start.p[side].teamMenu - 1
			elseif motif.select_info.teammenu.move.wrapping then
				sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.move.snd[1], motif.select_info['p' .. side].teammenu.move.snd[2])
				start.p[side].teamMenu = #t
			end
		elseif #t > 1 and getInput(t_cmd, motif.select_info['p' .. side].teammenu.next.key) then
			if start.p[side].teamMenu < #t then
				sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.move.snd[1], motif.select_info['p' .. side].teammenu.move.snd[2])
				start.p[side].teamMenu = start.p[side].teamMenu + 1
			elseif motif.select_info.teammenu.move.wrapping then
				sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.move.snd[1], motif.select_info['p' .. side].teammenu.move.snd[2])
				start.p[side].teamMenu = 1
			end
		else
			local handled = hook.runFirst("start.f_teamMenu.input", side, t, t_cmd)
			if handled then
				-- handled by external module
			elseif t[start.p[side].teamMenu].itemname == 'simul' then
				if getInput(t_cmd, motif.select_info['p' .. side].teammenu.subtract.key) then
					if start.p[side].numSimul > main.numSimul[1] then
						sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.value.snd[1], motif.select_info['p' .. side].teammenu.value.snd[2])
						start.p[side].numSimul = start.p[side].numSimul - 1
					end
				elseif getInput(t_cmd, motif.select_info['p' .. side].teammenu.add.key) then
					if start.p[side].numSimul < main.numSimul[2] then
						sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.value.snd[1], motif.select_info['p' .. side].teammenu.value.snd[2])
						start.p[side].numSimul = start.p[side].numSimul + 1
					end
				end
			elseif t[start.p[side].teamMenu].itemname == 'turns' then
				if getInput(t_cmd, motif.select_info['p' .. side].teammenu.subtract.key) then
					if start.p[side].numTurns > main.numTurns[1] then
						sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.value.snd[1], motif.select_info['p' .. side].teammenu.value.snd[2])
						start.p[side].numTurns = start.p[side].numTurns - 1
					end
				elseif getInput(t_cmd, motif.select_info['p' .. side].teammenu.add.key) then
					if start.p[side].numTurns < main.numTurns[2] then
						sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.value.snd[1], motif.select_info['p' .. side].teammenu.value.snd[2])
						start.p[side].numTurns = start.p[side].numTurns + 1
					end
				end
			elseif t[start.p[side].teamMenu].itemname == 'tag' then
				if getInput(t_cmd, motif.select_info['p' .. side].teammenu.subtract.key) then
					if start.p[side].numTag > main.numTag[1] then
						sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.value.snd[1], motif.select_info['p' .. side].teammenu.value.snd[2])
						start.p[side].numTag = start.p[side].numTag - 1
					end
				elseif getInput(t_cmd, motif.select_info['p' .. side].teammenu.add.key) then
					if start.p[side].numTag < main.numTag[2] then
						sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.value.snd[1], motif.select_info['p' .. side].teammenu.value.snd[2])
						start.p[side].numTag = start.p[side].numTag + 1
					end
				end
			end
		end
		--Exit during team menu
		if not start.escFlag and (esc() or getInput(-1, motif.select_info.cancel.key)) then
			esc(false)
			sndPlay(motif.Snd, motif.select_info.cancel.snd[1], motif.select_info.cancel.snd[2])
			fadeOutInit(motif.select_info.fadeout.FadeData)
			fadeOutStarted = true
			start.escFlag = true
		end
		--Draw team background
		main.f_animPosDraw(motif.select_info['p' .. side].teammenu.bg.default.AnimData)
		--Draw team title
		if side == 2 and main.cpuSide[2] then
			main.f_animPosDraw(motif.select_info['p' .. side].teammenu.enemytitle.AnimData)
			textImgDraw(motif.select_info['p' .. side].teammenu.enemytitle.TextSpriteData)
		else
			main.f_animPosDraw(motif.select_info['p' .. side].teammenu.selftitle.AnimData)
			textImgDraw(motif.select_info['p' .. side].teammenu.selftitle.TextSpriteData)
		end
		--Draw team cursor
		main.f_animPosDraw(
			motif.select_info['p' .. side].teammenu.item.cursor.AnimData,
			(start.p[side].teamMenu - 1) * motif.select_info['p' .. side].teammenu.item.spacing[1],
			(start.p[side].teamMenu - 1) * motif.select_info['p' .. side].teammenu.item.spacing[2]
		)
		local teammenu = motif.select_info['p' .. side].teammenu
		local itemCfg = teammenu.item
		local valueCfg = teammenu.value
		local spacingX, spacingY = itemCfg.spacing[1], itemCfg.spacing[2]
		local uppercase = itemCfg.uppercase
		local gm = gameMode()
		for i = 1, #t do
			local x = (i - 1) * spacingX
			local y = (i - 1) * spacingY
			local itemname = t[i].itemname
			local itemText = main.f_itemnameUpper(t[i].displayname, uppercase)
			local textData = itemCfg.TextSpriteData
			--Draw team items
			if i == start.p[side].teamMenu then
				if t_teamActiveCount[side] < itemCfg.active.switchtime then --delay change
					t_teamActiveCount[side] = t_teamActiveCount[side] + 1
				else
					if t_teamActiveState[side] then
						t_teamActiveState[side] = false
					else
						t_teamActiveState[side] = true
					end
					t_teamActiveCount[side] = 0
				end
				--Draw team active item background
				if teammenu.active.bg[gm .. '-' .. itemname] ~= nil then
					main.f_animPosDraw(teammenu.active.bg[gm .. '-' .. itemname].AnimData)
				elseif teammenu.active.bg[itemname] ~= nil then
					main.f_animPosDraw(teammenu.active.bg[itemname].AnimData)
				end
				--Draw team active item font
				if t_teamActiveState[side] then
					textData = itemCfg.active2.TextSpriteData
				else
					textData = itemCfg.active.TextSpriteData
				end
			else
				--Draw team not active item background
				if teammenu.bg[gm .. '-' .. itemname] ~= nil then
					main.f_animPosDraw(teammenu.bg[gm .. '-' .. itemname].AnimData)
				elseif teammenu.bg[itemname] ~= nil then
					main.f_animPosDraw(teammenu.bg[itemname].AnimData)
				end
			end
			--Draw item font
			textImgReset(textData)
			textImgAddPos(textData, x, y)
			textImgSetText(textData, itemText)
			textImgDraw(textData)
			--Draw team icons
			if itemname == 'simul' then
				for j = 1, main.numSimul[2] do
					local vx = x + (j - 1) * valueCfg.spacing[1]
					local vy = y + (j - 1) * valueCfg.spacing[2]
					if j <= start.p[side].numSimul then
						main.f_animPosDraw(valueCfg.icon.AnimData, vx, vy)
					else
						main.f_animPosDraw(valueCfg.empty.icon.AnimData, vx, vy)
					end
				end
			elseif itemname == 'turns' then
				for j = 1, main.numTurns[2] do
					local vx = x + (j - 1) * valueCfg.spacing[1]
					local vy = y + (j - 1) * valueCfg.spacing[2]
					if j <= start.p[side].numTurns then
						main.f_animPosDraw(valueCfg.icon.AnimData, vx, vy)
					else
						main.f_animPosDraw(valueCfg.empty.icon.AnimData, vx, vy)
					end
				end
			elseif itemname == 'tag' then
				for j = 1, main.numTag[2] do
					local vx = x + (j - 1) * valueCfg.spacing[1]
					local vy = y + (j - 1) * valueCfg.spacing[2]
					if j <= start.p[side].numTag then
						main.f_animPosDraw(valueCfg.icon.AnimData, vx, vy)
					else
						main.f_animPosDraw(valueCfg.empty.icon.AnimData, vx, vy)
					end
				end
			else
				hook.runFirst("start.f_teamMenu.drawItemValue", side, t, i, x, y, valueCfg)
			end
		end
		--Confirmed team selection
		if getInput(t_cmd, motif.select_info['p' .. side].teammenu.done.key) or timerSelect == -1 then
			timerSelect = motif.select_info.timer.displaytime
			start.p[1].screenDelay, start.p[2].screenDelay = 0, 0
			sndPlay(motif.Snd, motif.select_info['p' .. side].teammenu.done.snd[1], motif.select_info['p' .. side].teammenu.done.snd[2])
			local confirmData = hook.runFirst("start.f_teamMenu.confirm", side, t, t_cmd)
			if confirmData ~= nil then
				local defaultMode = t[start.p[side].teamMenu].mode
				local teamMode = confirmData.teamMode or defaultMode
				local numChars = confirmData.numChars or start.p[side].numChars
				for k, v in pairs(confirmData) do
					start.p[side][k] = v
				end
				start.p[side].teamMode = teamMode
				start.p[side].numChars = numChars
			elseif t[start.p[side].teamMenu].itemname == 'single' then
				start.p[side].teamMode = t[start.p[side].teamMenu].mode
				start.p[side].numChars = 1
			elseif t[start.p[side].teamMenu].itemname == 'simul' then
				start.p[side].teamMode = t[start.p[side].teamMenu].mode
				start.p[side].numChars = start.p[side].numSimul
			elseif t[start.p[side].teamMenu].itemname == 'turns' then
				start.p[side].teamMode = t[start.p[side].teamMenu].mode
				start.p[side].numChars = start.p[side].numTurns
			elseif t[start.p[side].teamMenu].itemname == 'tag' then
				start.p[side].teamMode = t[start.p[side].teamMenu].mode
				start.p[side].numChars = start.p[side].numTag
			end
			start.p[side].teamEnd = true
		end
	end
	--t_selCmd table appending once team mode selection is finished
	if start.p[side].teamEnd then
		if main.coop and (side == 1 or gameMode('versuscoop')) then
			for i = 1, start.p[side].numChars do
				if gameMode('versuscoop') then
					if side == 1 then
						table.insert(start.p[side].t_selCmd, {cmd = getRemapInput(i * 2 - 1), player = start.f_getPlayerNo(side, #start.p[side].t_selCmd + 1), selectState = 0})
					else
						table.insert(start.p[side].t_selCmd, {cmd = getRemapInput(i * 2), player = start.f_getPlayerNo(side, #start.p[side].t_selCmd + 1), selectState = 0})
					end
				else
					table.insert(start.p[1].t_selCmd, {cmd = getRemapInput(i), player = start.f_getPlayerNo(side, #start.p[1].t_selCmd + 1), selectState = 0})
				end
			end
		else
			table.insert(start.p[side].t_selCmd, {cmd = start.f_menuCmd(side), player = start.f_getPlayerNo(side, #start.p[side].t_selCmd + 1), selectState = 0})
		end
	end
end

--===========================================================
--; PALETTE SELECT
--===========================================================
LoadedPals = {}
-- Tracks which characters have already had their palettes loaded to avoid redundant loading
local function ifCharPalsLoaded(ref)
	for _, v in ipairs(LoadedPals) do
		if v == ref then
			return true
		end
	end
	table.insert(LoadedPals, ref)
	return false
end
-- Loads palettes for a character if needed, prepares the animation, and applies the palette
function start.loadPalettes(a, ref, pal)
	if not ifCharPalsLoaded(ref) then
		animLoadPalettes(a, ref)
	end
	local srcAnim = a
	a = animPrepare(a, ref)
	animApplyVel(a, srcAnim)
	a = animSetColorPalette(a, pal)
	return a
end

--===========================================================
-- Draw Palette Menu
--===========================================================
function start.f_palMenuDraw(side, member, curIdx, validIdx ,maxIdx)
	local charData = start.f_getCharData(start.p[side].t_selTemp[member].ref)
	if not charData or not charData.pal then return end
	local palIndex = start.p[side].t_selTemp[member].pal
	local totalPals = #charData.pal
	local pn = 2 * (member - 1) + side
	local pCfg = f_getMotifP(motif.select_info, pn, side)
	local displayText = (curIdx == maxIdx) and pCfg.palmenu.random.text or tostring(validIdx)
	-- bg
	main.f_animPosDraw(pCfg.palmenu.bg.AnimData)
	-- draw number
	textImgReset(pCfg.palmenu.number.TextSpriteData)
	textImgSetText(pCfg.palmenu.number.TextSpriteData, displayText)
	textImgDraw(pCfg.palmenu.number.TextSpriteData)
	-- draw text
	textImgReset(pCfg.palmenu.text.TextSpriteData)
	textImgDraw(pCfg.palmenu.text.TextSpriteData)
end

--returns a random palette (using synced RNG)
function start.f_randomPal(charRef, validPals)
	start.shufflePals = start.shufflePals or {}
	start.shufflePals[charRef] = start.shufflePals[charRef] or {}

	if #start.shufflePals[charRef] == 0 then
		local last = start.lastRandomPal and start.lastRandomPal[charRef]
		local t = {}
		for _, v in ipairs(validPals) do
			table.insert(t, v)
		end
		start.f_shuffleTable(t, last)
		start.shufflePals[charRef] = t
	end
	-- draw one palette from the bag
	local result = table.remove(start.shufflePals[charRef])
	-- store last drawn palette
	start.lastRandomPal = start.lastRandomPal or {}
	start.lastRandomPal[charRef] = result
	return result
end

local function resolvePalConflict(side, charRef, pal)
	local charData = start.f_getCharData(charRef)
	if not charData or not charData.pal then
		return pal
	end
	local usedPals = {}
	for s = 1, 2 do
		for _, sel in ipairs(start.p[s].t_selected) do
			if sel.ref == charRef and sel.pal then
				usedPals[sel.pal] = true
			end
		end
	end
	-- if the chosen palette is not used, keep it
	if not usedPals[pal] then
		return validatePal(pal, charRef)
	end
	-- if it's in use, try to find the next free one
	local maxPal = gameOption('Config.PaletteMax')
	for i = pal + 1, maxPal do
	if not usedPals[i] then
			return validatePal(i, charRef)
		end
	end
	for i = 1, pal - 1 do
		if not usedPals[i] then
			return validatePal(i, charRef)
		end
	end

	return validatePal(pal, charRef)
end

local function applyPalette(sel, charData, palIndex)
	if sel.face_data then
		local srcAnim = sel.face_data
		sel.face_data = animSetColorPalette(sel.face_data, palIndex)
	end
	if sel.face2_data then
		local srcAnim = sel.face2_data
		sel.face2_data = animSetColorPalette(sel.face2_data, palIndex)
	end
end

-- palette select menu
function start.f_palMenu(side, cmd, player, member, selectState)
	local st = start.p[side].t_selTemp[member]
	local charRef = st.ref
	local charData = start.f_getCharData(charRef)
	local pn = 2 * (member - 1) + side
	local pCfg = f_getMotifP(motif.select_info, pn, side)
	-- initialize palette list and index if character changed or not yet set
	if st.validPalsCharRef ~= charRef or not st.validPals then
		local valid, seen, cur = {}, {}, validatePal(1, charRef)
		valid[1], seen[cur] = cur, true
		for i = 1, #charData.pal do
			local nextp = validatePal(cur + 1, charRef)
			if seen[nextp] then break end
			table.insert(valid, nextp)
			seen[nextp], cur = true, nextp
		end
		st.validPals, st.validPalsCharRef = valid, charRef
		-- set current index to match current palette (or default to first)
		local curPal = st.pal or valid[1]
		st.currentIdx = 1
		for i, p in ipairs(valid) do
			if p == curPal then st.currentIdx = i; break end
		end
	end

	local validPals = st.validPals
	local curIdx = st.currentIdx or 1
	local pal = st.pal or validPals[curIdx]
	local maxIdx = #validPals + 1
	start.p[side].inPalMenu = true

	-- accept selection
	local autoConfirm = #validPals <= 1
	if autoConfirm or getInput(cmd, motif.select_info['p' .. side].palmenu.done.key) or timerSelect == -1 then
		-- TODO: There's an issue here where when the palette is selected there will be 1 frame without any cursor
		-- Since the "done" cursor only appears in the next frame
		pal = (curIdx == maxIdx) and (start.c[player].randPalPreview or start.f_randomPal(charRef, validPals)) or validPals[curIdx]
		st.pal, st.currentIdx = pal, curIdx

		-- done anim after pal confirmation - primary face
		local done_anim = pCfg.face.done.anim
		local done_spr = pCfg.face.done.spr
		local preview_anim = pCfg.palmenu.preview.anim
		if done_anim ~= preview_anim or done_spr[1] ~= -1 then
			if (st.face_anim ~= done_anim or done_spr[1] ~= -1) and (main.coop or motif.select_info['p' .. side].face.num > 1 or main.f_tableLength(start.p[side].t_selected) + 1 == start.p[side].numChars) then
				local a = start.f_animGet(start.c[player].selRef, side, member, pCfg.face.done, pCfg.face, false, st.face_data)
				if a then
					st.face_data = start.loadPalettes(a, charRef, pal)
					animUpdate(st.face_data)
					start.p[side].screenDelay = math.min(120, math.max(start.p[side].screenDelay, animGetLength(st.face_data)))
				end
			end
		end

		-- face2 "done" anim after pal confirmation
		local done_anim2 = pCfg.face2.done.anim
		if st.face2_anim ~= done_anim2 and (main.coop or motif.select_info['p' .. side].face2.num > 1 or main.f_tableLength(start.p[side].t_selected) + 1 == start.p[side].numChars) then
			local a = start.f_animGet(start.c[player].selRef, side, member, pCfg.face2.done, pCfg.face2, false, st.face2_data)
			if a then
				st.face2_data = start.loadPalettes(a, charRef, pal)
				animUpdate(st.face2_data)
				start.p[side].screenDelay = math.min(120, math.max(start.p[side].screenDelay, animGetLength(st.face2_data)))
			end
		end
		selectState = 3
		start.f_playWave(start.c[player].selRef, 'cursor', motif.select_info['p' .. side].select.snd[1], motif.select_info['p' .. side].select.snd[2])
		-- Skip cursor sound for auto-confirmation, since we already played the select character confirmation sound
		if not autoConfirm then
			sndPlay(motif.Snd, motif.select_info['p' .. side].palmenu.done.snd[1], motif.select_info['p' .. side].palmenu.done.snd[2])
		end
	 -- next palette
	elseif getInput(cmd, motif.select_info['p' .. side].palmenu.next.key) then
		curIdx = (curIdx == maxIdx) and 1 or curIdx + 1
		st.currentIdx = curIdx
		if curIdx < maxIdx then
			applyPalette(st, charData, validPals[curIdx])
		end
		sndPlay(motif.Snd, motif.select_info['p' .. side].palmenu.value.snd[1], motif.select_info['p' .. side].palmenu.value.snd[2])
	-- previous palette
	elseif getInput(cmd, motif.select_info['p' .. side].palmenu.previous.key) then
		curIdx = (curIdx == 1) and maxIdx or curIdx - 1
		st.currentIdx = curIdx
		if curIdx < maxIdx then
			applyPalette(st, charData, validPals[curIdx])
		end
		sndPlay(motif.Snd, motif.select_info['p' .. side].palmenu.value.snd[1], motif.select_info['p' .. side].palmenu.value.snd[2])
	-- cancel
	elseif getInput(cmd, motif.select_info['p' .. side].palmenu.cancel.key) then
		st.face_data = start.f_animGet(start.c[player].selRef, side, member, motif.select_info['p' .. pn].face, nil, true, st.face_data)
		st.face2_data = start.f_animGet(start.c[player].selRef, side, member, motif.select_info['p' .. pn].face2, nil, true, st.face2_data)
		selectState = 0
		st.currentIdx = nil
		st.validPals = nil
		sndPlay(motif.Snd, motif.select_info['p' .. side].palmenu.cancel.snd[1], motif.select_info['p' .. side].palmenu.cancel.snd[2])
	end
	-- random hotkey
	if getInput(cmd, motif.select_info['p' .. side].palmenu.random.key) then
		curIdx, st.currentIdx = maxIdx, maxIdx
	end
	-- random preview update
	if st.currentIdx == maxIdx then
		if not start.c[player].randPalCnt or start.c[player].randPalCnt <= 0 then
			start.c[player].randPalCnt = motif.select_info.palmenu.random.switchtime
			start.c[player].randPalPreview = start.f_randomPal(charRef, validPals)
			if motif.select_info.palmenu.random.applypal then
				applyPalette(st, charData, start.c[player].randPalPreview)
			else
				applyPalette(st, charData, 1)
			end
			sndPlay(motif.Snd, motif.select_info['p' .. side].palmenu.value.snd[1], motif.select_info['p' .. side].palmenu.value.snd[2])
		else
			start.c[player].randPalCnt = start.c[player].randPalCnt - 1
		end
	end
	start.f_palMenuDraw(side, member, curIdx, validPals[curIdx], maxIdx)
	return selectState
end

--;===========================================================
--; SELECT MENU
--;===========================================================
function start.f_selectMenu(side, cmd, player, member, selectState)
	local needUpdateDrawList = false
	--predefined selection
	if main.forceChar[side] ~= nil then
		local t = {}
		for _, v in ipairs(main.forceChar[side]) do
			if t[v] == nil then
				t[v] = ''
			end
			table.insert(start.p[side].t_selected, {
				ref = v,
				pal = start.f_selectPal(v),
				--pn = start.f_getPlayerNo(side, #start.p[side].t_selected + 1),
				--cursor = = {},
			})
		end
		start.p[side].selEnd = true
		return 0, false
	--manual selection
	elseif not start.p[side].selEnd then
		local pn = 2 * (member - 1) + side
		local pCfg = f_getMotifP(motif.select_info, pn, side)
		--cell not selected yet
		if selectState == 0 then
			--restore cursor coordinates
			if restoreCursor then
				-- remove entries if stored cursors exceeds team size
				if #start.p[side].t_cursor > start.p[side].numChars then
					for i = #start.p[side].t_cursor, start.p[side].numChars + 1, -1 do
						start.p[side].t_cursor[i] = nil
					end
				end
				-- restore saved position
				if start.p[side].t_cursor[member] ~= nil then
					local selX = start.p[side].t_cursor[member].x
					local selY = start.p[side].t_cursor[member].y
					if gameOption('Options.Team.Duplicates') or t_reservedChars[side][start.t_grid[selY + 1][selX + 1].char_ref] == nil then
						start.c[player].selX = selX
						start.c[player].selY = selY
					end
					start.p[side].t_cursor[member] = nil
				end
			end
			--calculate current position
			start.c[player].selX, start.c[player].selY = start.f_cellMovement(start.c[player].selX, start.c[player].selY, cmd, side, start.f_getCursorData(player).cursor.move.snd)
			start.c[player].cell = start.c[player].selX + motif.select_info.columns * start.c[player].selY
			start.c[player].selRef = start.f_selGrid(start.c[player].cell + 1).char_ref
			main.f_preloadSetCharHighlight(player, start.c[player].selRef)
			-- temp data not existing yet
			if start.p[side].t_selTemp[member] == nil then
				t_portraitPriority[side] = member
				table.insert(start.p[side].t_selTemp, {
					ref = start.c[player].selRef,
					cell = start.c[player].cell,
					inRandom = false,
					face_anim = pCfg.face.anim,
					face_data = start.f_animGet(start.c[player].selRef, side, member, pCfg.face, nil, true),
					face2_anim = pCfg.face2.anim,
					face2_data = start.f_animGet(start.c[player].selRef, side, member, pCfg.face2, nil, true),
				})
			else
				local updateAnim = false
				local slotSelected,slotChanged = start.f_slotSelected(start.c[player].cell + 1, side, cmd, player, start.c[player].selX, start.c[player].selY)
				local timerExpired = motif.select_info.timer.count ~= -1 and timerSelect == -1
				needUpdateDrawList = slotChanged
				local velCopy = false
				if slotChanged then
					start.c[player].selRef = start.f_selGrid(start.c[player].cell + 1).char_ref
				end
				if timerExpired then
					if start.c[player].selRef == nil or main.t_selChars[start.c[player].selRef + 1] == nil then
						start.c[player].selRef = start.f_randomChar(side)
					end
				end
				-- cursor changed position or character change within current slot
				if start.p[side].t_selTemp[member].cell ~= start.c[player].cell or start.p[side].t_selTemp[member].ref ~= start.c[player].selRef then
					if start.p[side].t_selTemp[member].cell ~= start.c[player].cell then
						t_portraitPriority[side] = member
					end
					--start.p[side].t_selTemp[member].pal = 1
					start.p[side].t_selTemp[member].ref = start.c[player].selRef
					start.p[side].t_selTemp[member].cell = start.c[player].cell
					start.p[side].t_selTemp[member].face_anim = pCfg.face.anim
					start.p[side].t_selTemp[member].face2_anim = pCfg.face2.anim
					if start.f_getCursorData(player).cursor.reset then
						resetCursorData(player, cursorActive, 'active')
					end
					updateAnim = true
				end
				-- cursor at randomselect cell
				if start.f_selGrid(start.c[player].cell + 1).char == 'randomselect' or start.f_selGrid(start.c[player].cell + 1).hidden == 3 then
					local wasRandom = start.p[side].t_selTemp[member].inRandom
					start.p[side].t_selTemp[member].inRandom = true
					velCopy = wasRandom
					-- re-entering random slot: reset random overlay anims so sliding restarts
					if not wasRandom then
						local pData = pCfg
						if pData.face2.random then
							animReset(pData.face2.random.AnimData)
							animUpdate(pData.face2.random.AnimData)
						end
						if pData.face.random then
							animReset(pData.face.random.AnimData)
							animUpdate(pData.face.random.AnimData)
						end
					end
					if start.c[player].randCnt > 0 then
						start.c[player].randCnt = start.c[player].randCnt - 1
						start.c[player].selRef = start.c[player].randRef
					else
						if motif.select_info.random.move.snd.cancel then
							sndStop(motif.Snd, start.f_getCursorData(player).random.move.snd[1], start.f_getCursorData(player).random.move.snd[2])
						end
						sndPlay(motif.Snd, start.f_getCursorData(player).random.move.snd[1], start.f_getCursorData(player).random.move.snd[2])
						start.c[player].randCnt = motif.select_info.cell.random.switchtime
						start.c[player].selRef = start.f_randomChar(side)
						if start.c[player].randRef ~= start.c[player].selRef or start.p[side].t_selTemp[member].face_data == nil then
							updateAnim = true
							start.c[player].randRef = start.c[player].selRef
						end
					end
				else
					start.p[side].t_selTemp[member].inRandom = false
				end
				main.f_preloadSetCharHighlight(player, start.c[player].selRef)
				-- update anim data
				if updateAnim then
					local face_data = velCopy and start.p[side].t_selTemp[member].face_data or nil
					local face2_data = velCopy and start.p[side].t_selTemp[member].face2_data or nil
					start.p[side].t_selTemp[member].face_data = start.f_animGet(start.c[player].selRef, side, member, pCfg.face, nil, true, face_data)
					start.p[side].t_selTemp[member].face2_data = start.f_animGet(start.c[player].selRef, side, member, pCfg.face2, nil, true, face2_data)
				end
				-- cell selected or select screen timer reached 0
				local canConfirm = (slotSelected and start.f_selGrid(start.c[player].cell + 1).char ~= nil and start.f_selGrid(start.c[player].cell + 1).hidden ~= 2) or timerExpired
				if canConfirm then
					local preloadReady = true
					if start.c[player].selRef ~= nil then
						local state = getCharPreloadStatus(start.c[player].selRef)
						preloadReady = state == 'ready'
						if not preloadReady then
							main.f_preloadBoostChar(start.c[player].selRef)
						else
							main.f_materializeCharByRef(start.c[player].selRef)
						end
					end
					if not preloadReady then
						if slotSelected then
							sndPlay(motif.Snd, motif.select_info.cancel.snd[1], motif.select_info.cancel.snd[2])
						end
						canConfirm = false
					end
				end
				if canConfirm then
					if motif.select_info.paletteselect ~= 0 then
						timerSelect = motif.select_info.timer.displaytime
					end
					sndPlay(motif.Snd, start.f_getCursorData(player).cursor.done.default.snd[1], start.f_getCursorData(player).cursor.done.default.snd[2])
					if motif.select_info.paletteselect == 0 then
						start.f_playWave(start.c[player].selRef, 'cursor', motif.select_info['p' .. side].select.snd[1], motif.select_info['p' .. side].select.snd[2])
					end
					if motif.select_info.paletteselect > 0 then
						resetCursorData(player, cursorActive, 'done')
					end
					start.p[side].t_selTemp[member].pal = main.f_btnPalNo(cmd)
					start.p[side].t_selTemp[member].inRandom = false
					if start.p[side].t_selTemp[member].pal == nil or start.p[side].t_selTemp[member].pal == 0 then
						start.p[side].t_selTemp[member].pal = 1
					end
					-- done anim helper
					local function setDoneAnim(ref, side, member, params, velParams, targetField)
						local a = start.f_animGet(ref, side, member, params, velParams, false, start.p[side].t_selTemp[member][targetField])
						if a then
							start.p[side].t_selTemp[member][targetField] = a
							start.p[side].screenDelay = math.min(120, math.max(start.p[side].screenDelay, animGetLength(a)))
						end
					end
					-- if select anim differs from done anim and coop or pX.face.num allows to display more than 1 portrait or it's the last team member
					local done_anim = pCfg.face.done.anim
					local done_anim2 = pCfg.face2.done.anim
					local done_spr = pCfg.face.done.spr
					local palmenu_preview_anim = pCfg.palmenu.preview.anim
					local face_anim = start.p[side].t_selTemp[member].face_anim
					local face2_anim = start.p[side].t_selTemp[member].face2_anim
					local canShow = main.coop or motif.select_info['p' .. side].face.num > 1 or main.f_tableLength(start.p[side].t_selected) + 1 == start.p[side].numChars
					local canShow2 = main.coop or motif.select_info['p' .. side].face2.num > 1 or main.f_tableLength(start.p[side].t_selected) + 1 == start.p[side].numChars
					-- primary face "done" / preview
					if (face_anim ~= done_anim or done_spr[1] ~= -1) and canShow then
						if motif.select_info.paletteselect == 0 and (done_anim ~= -1 or done_spr[1] ~= -1) then
							setDoneAnim(start.c[player].selRef, side, member, pCfg.face.done, pCfg.face, 'face_data')
						elseif palmenu_preview_anim ~= -1 and motif.select_info.paletteselect ~= 0 then
							start.f_playWave(start.c[player].selRef, 'cursor', motif.select_info['p' .. side].palmenu.preview.snd[1], motif.select_info['p' .. side].palmenu.preview.snd[2])
							setDoneAnim(start.c[player].selRef, side, member, pCfg.palmenu.preview, pCfg.face, 'face_data')
						end
					end
					-- face2 "done" anim
					if face2_anim ~= done_anim and canShow2 and done_anim2 ~= -1 then
						setDoneAnim(start.c[player].selRef, side, member, pCfg.face2.done, pCfg.face2, 'face2_data')
					end

					start.p[side].t_selTemp[member].ref = start.c[player].selRef
					local charRef = start.p[side].t_selTemp[member].ref
					local charData = start.f_getCharData(charRef)
					local pal = start.p[side].t_selTemp[member].pal
					local finalPal

					if motif.select_info.paletteselect == 1 then
						finalPal = 1
					elseif motif.select_info.paletteselect == 2 then
						finalPal = pal
					elseif motif.select_info.paletteselect == 3 then
						finalPal = start.f_keyPalMap(charRef, pal)
					else
						finalPal = start.f_keyPalMap(charRef, pal)
					end

					-- resolve visual palette conflict
					finalPal = resolvePalConflict(side, charRef, finalPal)

					if motif.select_info.paletteselect > 0 then
						start.p[side].t_selTemp[member].pal = finalPal
					end

					if start.p[side].t_selTemp[member].face_data ~= nil then
						local applyFlag = pCfg.face.applypal
						if applyFlag then
							start.p[side].t_selTemp[member].face_data = start.loadPalettes(start.p[side].t_selTemp[member].face_data, charRef, finalPal)
							animUpdate(start.p[side].t_selTemp[member].face_data)
						end
					end
					if start.p[side].t_selTemp[member].face2_data ~= nil then
						local applyFlag = pCfg.face2.applypal
						if applyFlag then
							start.p[side].t_selTemp[member].face2_data = start.loadPalettes(start.p[side].t_selTemp[member].face2_data, charRef, finalPal)
							animUpdate(start.p[side].t_selTemp[member].face2_data)
						end
					end
					selectState = 1
				end
			end
		--selection menu
		elseif selectState == 1 then
			if motif.select_info.paletteselect and motif.select_info.paletteselect > 0 then
				selectState = start.f_palMenu(side, cmd, player, member, selectState)
			else
				selectState = 3
			end
		--confirm selection
		elseif selectState == 3 then
			local valid = {1,2,3}
			local finalPal
			for _, v in ipairs(valid) do
				if motif.select_info.paletteselect == v then
					finalPal = start.p[side].t_selTemp[member].pal
					start.p[side].inPalMenu = false
					break
				end
			end
			finalPal = finalPal or start.f_selectPal(start.c[player].selRef, start.p[side].t_selTemp[member].pal)
			finalPal = resolvePalConflict(side, start.c[player].selRef, finalPal)
			applyPalette(start.p[side].t_selTemp[member], start.f_getCharData(start.c[player].selRef), finalPal)
			start.p[side].t_selected[member] = {
				ref = start.c[player].selRef,
				pal = finalPal,
				pn = start.f_getPlayerNo(side, member),
				cursor = {start.c[player].selX, start.c[player].selY},
			}
			main.f_preloadSetCharHighlight(player, nil)
			hook.run("start.f_selectMenu.selected", side, member, start.p[side].t_selected[member], start.p[side], player)
			if not gameOption('Options.Team.Duplicates') then
				t_reservedChars[side][start.c[player].selRef] = true
			end
			-- If we just confirmed a pick while hovering Random Select, force the random preview to advance for the next team member.
			-- This prevents mashing select from picking the same random character repeatedly.
			local cellIdx = (start.c[player].cell or -1) + 1
			if cellIdx > 0 then
				local g = start.f_selGrid(cellIdx)
				if g.char == 'randomselect' or g.hidden == 3 then
					start.c[player].randCnt = 0
					start.c[player].randRef = nil
				end
			end
			start.p[side].t_cursor[member] = {x = start.c[player].selX, y = start.c[player].selY}
			if main.f_tableLength(start.p[side].t_selected) == start.p[side].numChars then --if all characters have been chosen
				if side == 1 and main.cpuSide[2] and start.reset then --if player1 is allowed to select p2 characters
					if timerSelect == -1 then
						start.p[2].teamMode = start.p[1].teamMode
						start.p[2].numChars = start.p[1].numChars
						start.c[2].cell = start.c[1].cell
						start.c[2].selX = start.c[1].selX
						start.c[2].selY = start.c[1].selY
						start.p[2].teamEnd = false
					else
						start.p[2].teamEnd = false
					end
				end
				start.p[side].selEnd = true
			elseif not gameOption('Options.Team.Duplicates') and start.t_grid[start.c[player].selY + 1][start.c[player].selX + 1].char ~= 'randomselect' then
				local t_dirs = {'F', 'B', 'D', 'U'}
				if start.c[player].selY + 1 >= motif.select_info.rows then --next row not visible on the screen
					t_dirs = {'F', 'B', 'U', 'D'}
				end
				for _, v in ipairs(t_dirs) do
					local selX, selY = start.f_cellMovement(start.c[player].selX, start.c[player].selY, cmd, side, start.f_getCursorData(player).cursor.move.snd, v)
					if start.t_grid[selY + 1][selX + 1].char ~= nil and (selX ~= start.c[player].selX or selY ~= start.c[player].selY) then
						start.c[player].selX, start.c[player].selY = selX, selY
						break
					end
				end
			end
			if not start.p[1].teamEnd or not start.p[2].teamEnd or not start.p[1].selEnd or not start.p[2].selEnd then
				timerSelect = motif.select_info.timer.displaytime
			end
			if main.coop and (side == 1 or gameMode('versuscoop')) then --remaining members are controlled by different players
				selectState = 4
			elseif not start.p[side].selEnd then --next member controlled by this player should become selectable
				selectState = 0
			end
		end
	end
	return selectState, needUpdateDrawList
end

--;===========================================================
--; STAGE MENU
--;===========================================================
function start.f_stageMenu()
	local n = stageListNo
	local randomMode = {
		[0] = { init = 1, min = 1 }, -- disabled
		[1] = { init = 0, min = 0 }, -- default
		[2] = { init = 1, min = 0 }, -- random at the 'end'
	}

	local r = randomMode[motif.select_info.stage_randomselect] or randomMode[1]
	local stageListIdx = r.init
	local stageListMinIdx = r.min
	local stageListMaxIdx = #main.t_selectableStages

	if getInput(-1, motif.select_info.cell.left.key) then
		sndPlay(motif.Snd, motif.select_info.stage.move.snd[1], motif.select_info.stage.move.snd[2])
		stageListNo = stageListNo - 1
		if stageListNo < stageListMinIdx then stageListNo = stageListMaxIdx end
	elseif getInput(-1, motif.select_info.cell.right.key) then
		sndPlay(motif.Snd, motif.select_info.stage.move.snd[1], motif.select_info.stage.move.snd[2])
		stageListNo = stageListNo + 1
		if stageListNo > stageListMaxIdx then stageListNo = stageListMinIdx end
	elseif getInput(-1, motif.select_info.cell.up.key) then
		sndPlay(motif.Snd, motif.select_info.stage.move.snd[1], motif.select_info.stage.move.snd[2])
		for i = 1, 10 do
			stageListNo = stageListNo - 1
			if stageListNo < stageListMinIdx then stageListNo = stageListMaxIdx end
		end
	elseif getInput(-1, motif.select_info.cell.down.key) then
		sndPlay(motif.Snd, motif.select_info.stage.move.snd[1], motif.select_info.stage.move.snd[2])
		for i = 1, 10 do
			stageListNo = stageListNo + 1
			if stageListNo > stageListMaxIdx then stageListNo = stageListMinIdx end
		end
	end
	if n ~= stageListNo and stageListNo > 0 then
		animReset(main.t_selStages[main.t_selectableStages[stageListNo]].anim_data)
		animUpdate(main.t_selStages[main.t_selectableStages[stageListNo]].anim_data)
	end
	if stageListNo > 0 then
		main.f_preloadBoostStage(main.t_selectableStages[stageListNo])
	end
end

--;===========================================================
--; VERSUS SCREEN / ORDER SELECTION
--;===========================================================
-- Build params for loadStart()
function start.f_buildLoadStartParams(arg, doSelectMissing, t_orderRemap)
	local parts = {}
	local t = {}
	local musicParams = arg
	if type(arg) == "table" then
		t = arg
		musicParams = t.musicParams
	end
	if musicParams and musicParams ~= "" then
		parts[#parts + 1] = musicParams
	end
	local function addParam(k, v)
		if v == nil or v == "" then
			return
		end
		parts[#parts + 1] = k .. "=" .. tostring(v)
	end
	addParam("continue", t.continue)
	addParam("quickcontinue", t.quickcontinue)
	addParam("order", t.order)
	addParam("stage", t.stage)
	addParam("ai", t.ai)
	addParam("time", t.roundtime or t.time)
	addParam("vsscreen", t.vsscreen)
	addParam("victoryscreen", t.victoryscreen)
	addParam("winscreen", t.winscreen)
	addParam("lua", t.lua)
	addParam("charparam.ai", main.charparam.ai)
	addParam("charparam.arcadepath", main.charparam.arcadepath)
	addParam("charparam.music", main.charparam.music)
	addParam("charparam.rounds", main.charparam.rounds)
	addParam("charparam.single", main.charparam.single)
	addParam("charparam.stage", main.charparam.stage)
	addParam("charparam.time", main.charparam.time)
	addParam("p1.turnsoffset", start.p[1].turnsOffset or 0)
	addParam("p2.turnsoffset", start.p[2].turnsOffset or 0)
	addParam("persistlife", main.persistLife)
	addParam("persistmusic", main.persistMusic)
	addParam("persistrounds", main.persistRounds)
	addParam("rankingcondition", main.rankingCondition)
	return table.concat(parts, ", ")
end

-- Build per-member override params for selectChar().
function start.f_buildOverrideParams(side, member, v)
	local parts = {}
	local function addParam(field, val)
		if val == nil then return end
		parts[#parts + 1] = string.format("p%d.%d.%s=%s", side, member, field, tostring(val))
	end
	hook.run("start.f_selectLoading.member", v)
	addParam("life", v.life)
	addParam("lifemax", v.lifeMax)
	addParam("power", v.power)
	addParam("dizzypoints", v.dizzyPoints)
	addParam("guardpoints", v.guardPoints)
	addParam("existed", v.existed)
	if type(v.maps) == "table" then
		for mapName, mapValue in pairs(v.maps) do
			if type(mapName) == "string" and mapValue ~= nil then
				local key = mapName
				if key:sub(1, 4):lower() == "map." then
					key = key:sub(5)
				end
				if key ~= "" then
					addParam("map." .. key, mapValue)
				end
			end
		end
	end
	return table.concat(parts, ", ")
end

function start.f_selectVersus(active, t_orderSelect, loadStartArg)
	start.t_orderRemap = {{}, {}}
	start.bgLoadStarted = false
	for side = 1, 2 do
		-- populate order remap table with default values
		for i = 1, #start.p[side].t_selected do
			table.insert(start.t_orderRemap[side], i)
		end
		-- prevent order select if not enabled in screenpack or if team size = 1
		if t_orderSelect[side] then
			t_orderSelect[side] = motif.vs_screen.orderselect.enabled and #start.p[side].t_selected > 1
			-- In Turns Survival after any defeats, order selection would break the invariant that
			-- defeated members are a prefix of the roster, so should be disabled.
			if start.p[side].teamMode == 2 and (start.p[side].turnsOffset or 0) > 0 then
				t_orderSelect[side] = false
			end
		end
		-- reset order-confirm flags and selection flags
		for _, v in ipairs(start.p[side].t_selected) do
			main.f_preloadBoostChar(v.ref)
			v.loading = false
			v.selected = false
		end
	end
	-- skip versus screen if vs screen is disabled or p2 side char has vsscreen select.def flag set to 0
	for _, v in ipairs(start.p[2].t_selected) do
		if start.f_getCharData(v.ref).vsscreen == 0 then
			active = false
			break
		end
	end
	if not active then
		clearColor(motif.versusbgdef.bgclearcolor[1], motif.versusbgdef.bgclearcolor[2], motif.versusbgdef.bgclearcolor[3])
		return true
	end
	textImgReset(motif.vs_screen.match.TextSpriteData)
	textImgSetText(motif.vs_screen.match.TextSpriteData, string.format(motif.vs_screen.match.text, matchNo()))
	bgReset(motif.versusbgdef.BGDef)
	fadeInInit(motif.vs_screen.fadein.FadeData)
	local fadeOutStarted = false
	playBgm({source = "motif.vs"})
	start.f_resetTempData(motif.vs_screen, '')
	start.f_playWave(getStageNo(), 'stage', motif.vs_screen.stage.snd[1], motif.vs_screen.stage.snd[2])
	local counter = 0 - motif.vs_screen.fadein.time
	local bgLoading = gameOption('Config.VsScreenLoading')
	local done = (not t_orderSelect[1] and not t_orderSelect[2]) -- both sides having order disabled
	local timerActive = not done
	local timerCount = 0
	local escFlag = false
	local doneKeyReady = done
	local t_order = {{}, {}}
	local cpuOrderFinalized = {false, false}
	local t_icon = {false, false}
	local selStageNo = getStageNo()
	local loadStarted = false
	local netReady = false
	local readyToLeave = not bgLoading
	local wantSkip = false
	local wantDone = false

	-- Background loading: start async loader immediately.
	if bgLoading then
		local params = start.f_buildLoadStartParams(loadStartArg, false, start.t_orderRemap)
		if gameOption('Debug.DumpLuaTables') then main.f_printTable(params, "debug/loadStartParams.txt") end
		resetGameParams()
		loadStart(params)
		loadStarted = true
		start.bgLoadStarted = true
		-- Sides without order select: select everyone immediately so loading can begin.
		for side = 1, 2 do
			if not t_orderSelect[side] then
				for member, v in ipairs(start.p[side].t_selected) do
					if not v.selected then
						selectChar(side, v.ref, v.pal, start.f_buildOverrideParams(side, member, v))
						v.selected = true
					end
					v.loading = true
					t_order[side][#t_order[side] + 1] = member
				end
				t_icon[side] = nil
			end
		end
	end

	local function finishOrderSelection(side)
		for member, v in ipairs(start.p[side].t_selected) do
			if not v.loading then
				t_order[side][#t_order[side] + 1] = member
				local slot = #t_order[side]
				if not v.selected then
					if bgLoading then
						selectChar(side, v.ref, v.pal, start.f_buildOverrideParams(side, slot, v))
					else
						selectChar(side, v.ref, v.pal)
					end
					v.selected = true
				end
				v.loading = true
			end
		end
		if #start.p[side].t_selected == #t_order[side] then
			t_icon[side] = nil
		end
	end
	while true do
		main.f_preloadTick(4)
		local snd = false
		-- CPU order select: randomize first, then selectChar() using randomized slot order
		for side = 1, 2 do
			if main.cpuSide[side] and t_orderSelect[side] and not cpuOrderFinalized[side] then
				t_order[side] = {}
				for i = 1, #start.p[side].t_selected do
					t_order[side][i] = i
				end
				main.f_tableShuffle(t_order[side])
				for slot, idx in ipairs(t_order[side]) do
					local v = start.p[side].t_selected[idx]
					if bgLoading and not v.selected then
						selectChar(side, v.ref, v.pal, start.f_buildOverrideParams(side, slot, v))
						v.selected = true
					end
					v.loading = true
				end
				t_icon[side] = nil
				cpuOrderFinalized[side] = true
			end
		end
		-- for each team side member
		for side = 1, 2 do
			if not done and t_orderSelect[side] and not main.cpuSide[side] and getInput(side, motif.vs_screen.skip.key) then
				finishOrderSelection(side)
				if not snd then
					sndPlay(motif.Snd, motif.vs_screen['p' .. side].value.snd[1], motif.vs_screen['p' .. side].value.snd[2])
					snd = true
				end
			end
			for k, v in ipairs(start.p[side].t_selected) do
				local pn = 2 * (k - 1) + side
				local pCfg = f_getMotifP(motif.vs_screen, pn, side)
				-- until loading flag is set
				if not v.loading then
					-- Timeout: append all remaining members in default order and confirm them.
					if timerCount == -1 then
						for kk, vv in ipairs(start.p[side].t_selected) do
							if not vv.loading then
								t_order[side][#t_order[side] + 1] = kk
								local slot = #t_order[side]
								if not vv.selected then
									if bgLoading then
										selectChar(side, vv.ref, vv.pal, start.f_buildOverrideParams(side, slot, vv))
									else
										selectChar(side, vv.ref, vv.pal)
									end
									vv.selected = true
								end
								vv.loading = true
							end
						end
						t_icon[side] = nil
						if not snd then
							sndPlay(motif.Snd, motif.vs_screen['p' .. side].value.snd[1], motif.vs_screen['p' .. side].value.snd[2])
							snd = true
						end
					-- CPU / no-key / auto-confirm path
					elseif not t_orderSelect[side] or main.cpuSide[side] or (#pCfg.key == 0 and #t_order[side] == k - 1) then
						t_order[side][#t_order[side] + 1] = k
						local slot = #t_order[side]
						if bgLoading and not v.selected then
							selectChar(side, v.ref, v.pal, start.f_buildOverrideParams(side, slot, v))
							v.selected = true
						end
						v.loading = true
						if #start.p[side].t_selected == #t_order[side] then
							t_icon[side] = nil
						end
					elseif getInput(side, pCfg.key) or (#start.p[side].t_selected == #t_order[side] + 1) then
						t_order[side][#t_order[side] + 1] = k
						local slot = #t_order[side]
						if bgLoading and not v.selected then
							selectChar(side, v.ref, v.pal, start.f_buildOverrideParams(side, slot, v))
							v.selected = true
						end
						if not bgLoading and not v.selected then
							selectChar(side, v.ref, v.pal)
							v.selected = true
						end
						v.loading = true
						-- if it's the last unordered team member
						if #start.p[side].t_selected == #t_order[side] then
							t_icon[side] = nil
						end
						-- play sound only once in particular frame
						if not snd then
							sndPlay(motif.Snd, motif.vs_screen['p' .. side].value.snd[1], motif.vs_screen['p' .. side].value.snd[2])
							snd = true
						end
					end
				end
			end
		end
		-- do once if both sides confirmed order selection
		if not done and #start.p[1].t_selected == #t_order[1] and #start.p[2].t_selected == #t_order[2] then
			for side = 1, 2 do
				-- rearrange characters in selection order
				for k, v in ipairs(t_order[side]) do
					start.t_orderRemap[side][k] = v
				end
				-- update spr/anim data
				for member, v in ipairs(start.p[side].t_selected) do
					local pn = 2 * (member - 1) + side
					local pCfg = f_getMotifP(motif.vs_screen, pn, side)
					-- primary face "done" anim
					local done_anim = pCfg.done.anim
					if done_anim ~= -1 and start.p[side].t_selTemp[member].face_anim ~= done_anim then
						start.p[side].t_selTemp[member].face_data = start.f_animGet(v.ref, side, member, pCfg.done, pCfg, false, start.p[side].t_selTemp[member].face_data)
					end
					-- face2 "done" anim
					local done_anim2 = pCfg.face2.done.anim
					if done_anim2 ~= -1 and start.p[side].t_selTemp[member].face2_anim ~= done_anim2 then
						start.p[side].t_selTemp[member].face2_data = start.f_animGet(v.ref, side, member, pCfg.face2.done, pCfg.face2, false, start.p[side].t_selTemp[member].face2_data)
					end
				end
				if t_orderSelect[side] then
					t_icon[side] = true
				end
			end
			counter = motif.vs_screen.time - motif.vs_screen.done.time
			done = true
			doneKeyReady = false
		end
		counter = counter + 1
		--draw clearcolor
		clearColor(motif.versusbgdef.bgclearcolor[1], motif.versusbgdef.bgclearcolor[2], motif.versusbgdef.bgclearcolor[3])
		--draw layerno = 0 backgrounds
		bgDraw(motif.versusbgdef.BGDef, 0)
		--draw portraits and order icons
		for side = 1, 2 do
			start.f_drawPortraits(main.f_remapTable(start.p[side].t_selTemp, start.t_orderRemap[side]), side, motif.vs_screen, '', false, t_icon[side])
		end
		--draw order values
		for side = 1, 2 do
			if t_orderSelect[side] then
				for i = 1, math.min(#start.p[side].t_selected, motif.vs_screen['p' .. side].num) do
					local pn = 2 * (i - 1) + side
					local pCfg = f_getMotifP(motif.vs_screen, pn, side)
					if i > #t_order[side] and #start.p[side].t_selected > #t_order[side] then
						main.f_animPosDraw(
							pCfg.value.empty.icon.AnimData,
							(i - 1) * motif.vs_screen['p' .. side].value.icon.spacing[1],
							(i - 1) * motif.vs_screen['p' .. side].value.icon.spacing[2]
						)
					else
						main.f_animPosDraw(
							pCfg.value.icon.AnimData,
							(i - 1) * motif.vs_screen['p' .. side].value.icon.spacing[1],
							(i - 1) * motif.vs_screen['p' .. side].value.icon.spacing[2]
						)
					end
				end
			end
		end
		--draw names
		for side = 1, 2 do
			for k, v in ipairs(main.f_remapTable(start.p[side].t_selTemp, start.t_orderRemap[side])) do
				if k <= motif.vs_screen['p' .. side].name.num then
					textImgReset(motif.vs_screen['p' .. side].name.TextSpriteData)
					textImgAddPos(
						motif.vs_screen['p' .. side].name.TextSpriteData,
						(k - 1) * motif.vs_screen['p' .. side].name.spacing[1],
						(k - 1) * motif.vs_screen['p' .. side].name.spacing[2]
					)
					textImgSetText(motif.vs_screen['p' .. side].name.TextSpriteData, start.f_getName(v.ref, side))
					textImgDraw(motif.vs_screen['p' .. side].name.TextSpriteData)
				end
			end
		end
		--draw stage portrait
		if selStageNo then
			--draw stage portrait background
			main.f_animPosDraw(motif.vs_screen.stage.portrait.bg.AnimData)
			--draw stage portrait loaded from stage SFF
			local portrait = motif.vs_screen.stage.portrait
			local anim = main.t_selStages[selStageNo].vs_anim_data
			local loadingPortrait = false
			if getStagePreloadStatus(selStageNo) ~= 'ready' and hasPortraitAnim(portrait.loading) then
				portrait = portrait.loading
				anim = portrait.AnimData
				loadingPortrait = true
			end
			if anim then
				local x = motif.vs_screen.stage.pos[1]
				local y = motif.vs_screen.stage.pos[2]
				if not loadingPortrait then
					x = x + portrait.offset[1]
					y = y + portrait.offset[2]
				end
				main.f_animPosDraw(anim, x, y)
			end
		end
		--draw stage name
		if selStageNo and main.t_selStages[selStageNo] then
			textImgReset(motif.vs_screen.stage.TextSpriteData)
			local stage_text = main.f_formatBySpec(motif.vs_screen.stage.text, {i = selStageNo, s = main.t_selStages[selStageNo].name})
			textImgSetText(motif.vs_screen.stage.TextSpriteData, stage_text)
			textImgDraw(motif.vs_screen.stage.TextSpriteData)
		end
		--draw match counter
		if main.motif.vsmatchno then
			textImgDraw(motif.vs_screen.match.TextSpriteData)
		end
		--draw timer
		if not done and motif.vs_screen.timer.count ~= -1 and timerActive and counter >= 0 then
			timerCount, timerActive = main.f_drawTimer(timerCount, motif.vs_screen.timer)
		end
		-- Background loading status
		readyToLeave = not bgLoading
			if bgLoading and loadStarted then
				local localDone = not loading()
				if localDone and not netReady then
					netReady = netLoadingReady()
				end
				readyToLeave = localDone and netReady
				if not main.suppressFightLoadingDisplay then
					if not readyToLeave then
						main.f_animPosDraw(motif.vs_screen.loading.AnimData)
						textImgDraw(motif.vs_screen.loading.TextSpriteData)
					else
						main.f_animPosDraw(motif.vs_screen.loading.done.AnimData)
						textImgDraw(motif.vs_screen.loading.done.TextSpriteData)
					end
				end
			end
		--draw layerno = 1 backgrounds
		bgDraw(motif.versusbgdef.BGDef, 1)
		-- hook
		hook.run("start.f_selectVersus")
		-- done key
		if done and not doneKeyReady and not getInput(-1, motif.vs_screen.done.key) then
			doneKeyReady = true
		end
		--draw fadein / fadeout
		for side = 1, 2 do
			-- Latch skip/done while background loading is still in progress.
			if bgLoading and loadStarted and not readyToLeave then
				if not main.cpuSide[side] and getInput(side, motif.vs_screen.skip.key) then
					wantSkip = true
				end
				if done and doneKeyReady and getInput(side, motif.vs_screen.done.key) then
					wantDone = true
				end
			end
			if not fadeOutStarted and (
				(counter >= motif.vs_screen.time and (not (t_orderSelect[1] or t_orderSelect[2]) or done) and readyToLeave)
				or (readyToLeave and (not main.cpuSide[side] and (getInput(side, motif.vs_screen.skip.key) or wantSkip)))
				or (readyToLeave and (done and doneKeyReady and (getInput(side, motif.vs_screen.done.key) or wantDone)))) then
				fadeOutInit(motif.vs_screen.fadeout.FadeData)
				fadeOutStarted = true
				wantSkip = false
				wantDone = false
				break
			end
		end
		--frame transition
		if not escFlag and (esc() or getInput(-1, motif.vs_screen.cancel.key)) then
			esc(false)
			if bgLoading and loadStarted then
				loadCancel()
				clearSelected()
			end
			fadeOutInit(motif.vs_screen.fadeout.FadeData)
			fadeOutStarted = true
			escFlag = true
		end
		if not fadeActive() and (fadeOutStarted or start.escFlag) then
			clearColor(motif.versusbgdef.bgclearcolor[1], motif.versusbgdef.bgclearcolor[2], motif.versusbgdef.bgclearcolor[3])
			break --skip last frame rendering
		end
		refresh()
	end
	esc(escFlag) --force Esc detection
	return not escFlag
end

--loading loop called after versus screen is finished
function start.f_selectLoading(arg)
	clearAllSound()
	local t = {}
	if type(arg) == "table" then
		t = arg
	elseif type(arg) == "string" then
		t.musicParams = arg
	end
	local params = start.f_buildLoadStartParams(t, true, start.t_orderRemap)
	if gameOption('Debug.DumpLuaTables') then main.f_printTable(params, "debug/loadStartParams.txt") end
	resetGameParams()
	if not gameOption('Config.VsScreenLoading') then
		-- If background loading is disabled, first select all chars, then start the loader.
		for side = 1, 2 do
			local remap = start.t_orderRemap and start.t_orderRemap[side]
			for member = 1, #start.p[side].t_selected do
				local src = remap and remap[member] or member
				local v = start.p[side].t_selected[src]
				if not v.selected then
					selectChar(side, v.ref, v.pal, start.f_buildOverrideParams(side, member, v))
					v.selected = true
				end
			end
		end
		loadStart(params)
	else
		-- Background loading: start loader first, then feed selections.
		loadStart(params)
		for side = 1, 2 do
			local remap = start.t_orderRemap and start.t_orderRemap[side]
			for member = 1, #start.p[side].t_selected do
				local src = remap and remap[member] or member
				local v = start.p[side].t_selected[src]
				if not v.selected then
					selectChar(side, v.ref, v.pal, start.f_buildOverrideParams(side, member, v))
					v.selected = true
				end
			end
		end
			if not main.suppressFightLoadingDisplay and main.f_storyboard(motif.vs_screen.loading.storyboard) then
				loadCancel()
				clearSelected()
				return false
			end
		-- VS screen is normally responsible for showing loading progress.
		-- If it was skipped/disabled, keep a small Lua render loop alive here.
		local netReady = false
		while loading() or not netReady do
			if not loading() then
				netReady = netLoadingReady()
			end
				if not main.suppressFightLoadingDisplay then
					clearColor(0, 0, 0)
					main.f_animPosDraw(motif.vs_screen.loading.wait.AnimData)
					textImgDraw(motif.vs_screen.loading.wait.TextSpriteData)
				end
				refresh()
			end
	end
	return true
end

return start
