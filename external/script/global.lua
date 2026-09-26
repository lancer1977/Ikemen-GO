--;===========================================================
--; DEBUG HOTKEYS
--;===========================================================
--key, ctrl, alt, shift, pause, debug key, function
addHotkey('c', true, false, false, true, false, 'toggleClsnDisplay()')
addHotkey('d', true, false, false, true, false, 'toggleDebugDisplay()')
addHotkey('d', false, false, true, true, false, 'toggleDebugDisplay(true)')
addHotkey('w', true, false, false, true, false, 'toggleWireframeDisplay()')
addHotkey('s', true, false, false, true, true, 'changeSpeed()')
addHotkey('KP_PLUS', true, false, false, true, true, 'changeSpeed(1)')
addHotkey('KP_MINUS', true, false, false, true, true, 'changeSpeed(-1)')
addHotkey('l', true, false, false, true, true, 'toggleLifebarDisplay()')
addHotkey('v', true, false, false, true, true, 'toggleVSync()')
addHotkey('1', true, false, false, true, true, 'toggleAI(1)')
addHotkey('1', true, true, false, true, true, 'togglePlayer(1)')
addHotkey('2', true, false, false, true, true, 'toggleAI(2)')
addHotkey('2', true, true, false, true, true, 'togglePlayer(2)')
addHotkey('3', true, false, false, true, true, 'toggleAI(3)')
addHotkey('3', true, true, false, true, true, 'togglePlayer(3)')
addHotkey('4', true, false, false, true, true, 'toggleAI(4)')
addHotkey('4', true, true, false, true, true, 'togglePlayer(4)')
addHotkey('5', true, false, false, true, true, 'toggleAI(5)')
addHotkey('5', true, true, false, true, true, 'togglePlayer(5)')
addHotkey('6', true, false, false, true, true, 'toggleAI(6)')
addHotkey('6', true, true, false, true, true, 'togglePlayer(6)')
addHotkey('7', true, false, false, true, true, 'toggleAI(7)')
addHotkey('7', true, true, false, true, true, 'togglePlayer(7)')
addHotkey('8', true, false, false, true, true, 'toggleAI(8)')
addHotkey('8', true, true, false, true, true, 'togglePlayer(8)')
addHotkey('9', true, true, false, true, true, 'togglePlayer(9)')
addHotkey('F1', false, false, false, false, true, 'kill(2); kill(4); kill(6); kill(8); debugFlag(1)')
addHotkey('F1', true, false, false, false, true, 'kill(1); kill(3); kill(5); kill(7); debugFlag(2)')
addHotkey('F2', false, false, false, false, true, 'kill(1,1); kill(2,1); kill(3,1); kill(4,1); kill(5,1); kill(6,1); kill(7,1); kill(8,1); debugFlag(1); debugFlag(2)')
addHotkey('F2', true, false, false, false, true, 'kill(1,1); kill(3,1); kill(5,1); kill(7,1); debugFlag(2)')
addHotkey('F2', false, false, true, false, true, 'kill(2,1); kill(4,1); kill(6,1); kill(8,1); debugFlag(1)')
addHotkey('F3', false, false, false, false, true, 'powMax(1); powMax(2); debugFlag(1); debugFlag(2)')
addHotkey('F3', true, false, true, false, true, 'toggleMaxPowerMode(); debugFlag(1); debugFlag(2)')
addHotkey('F4', false, false, false, false, true, 'roundReset(); closeMenu()')
addHotkey('F4', false, false, true, false, true, 'reload(); closeMenu()')
addHotkey('F5', false, false, false, false, true, 'setTime(0); debugFlag(1); debugFlag(2)')
addHotkey('F9', false, false, false, true, false, 'loadState()')
addHotkey('F10', false, false, false, true, false, 'saveState()')
addHotkey('SPACE', false, false, false, false, true, 'full(1); full(2); full(3); full(4); full(5); full(6); full(7); full(8); setTime(getRoundTime()); debugFlag(1); debugFlag(2); clearConsole()')
addHotkey('i', true, false, false, true, true, 'stand(1); stand(2); stand(3); stand(4); stand(5); stand(6); stand(7); stand(8)')
addHotkey('PAUSE', false, false, false, true, false, 'togglePause(); closeMenu()')
addHotkey('PAUSE', true, false, false, true, false, 'frameStep()')
addHotkey('SCROLLLOCK', false, false, false, true, false, 'frameStep()')

local speedMul = 1
local speedAdd = 0
function changeSpeed(add)
	if add ~= nil then
		speedAdd = speedAdd + add / 100
	elseif speedMul >= 4 then
		speedMul = 0.25
	else
		speedMul = speedMul * 2
	end
	setAccel(math.max(0.01, speedMul + speedAdd))
end

function toggleAI(p)
	local oldid = id()
	if player(p) then
		if ailevel() > 0 then
			setAILevel(0)
		else
			setAILevel(gameOption('Options.Difficulty'))
		end
		playerid(oldid)
	end
end

function kill(p, ...)
	local oldid = id()
	if player(p) then
		local n = ...
		if not n then n = 0 end
		setLife(n)
		setRedLife(0)
		playerid(oldid)
	end
end

function powMax(p)
	local oldid = id()
	if player(p) then
		setPower(powermax())
		setGuardPoints(guardpointsmax())
		setDizzyPoints(dizzypointsmax())
		playerid(oldid)
	end
end

function full(p)
	local oldid = id()
	if player(p) then
		setLife(lifemax())
		setPower(powermax())
		setGuardPoints(guardpointsmax())
		setDizzyPoints(dizzypointsmax())
		setRedLife(lifemax())
		removeDizzy()
		playerid(oldid)
	end
end

function stand(p)
	local oldid = id()
	if player(p) then
		selfState(0)
		playerid(oldid)
	end
end

local function kfmFactionAIBuff()
	if start == nil or start.f_getActiveCharRef == nil or start.f_getCharFaction == nil then
		return
	end
	if roundstate() ~= 2 or main.pauseMenu or paused() then
		return
	end
	local oldid = id()
	for side = 1, 2 do
		if player(side) then
			local ref = start.f_getActiveCharRef(side)
			local faction = ref ~= nil and start.f_getCharFaction(ref) or nil
			if tostring(faction or ''):lower() == 'kfm' and ailevel() > 0 then
				if ailevel() < 8 then
					setAILevel(8)
				end
				if gameTime() % 20 == 0 then
					setPower(math.min(powermax(), power() + 35))
					setGuardPoints(math.min(guardpointsmax(), guardpoints() + 8))
				end
				if gameTime() % 90 == 0 and life() > 0 and life() < lifemax() then
					setLife(math.min(lifemax(), life() + 3))
				end
			end
		end
	end
	playerid(oldid)
end

local function specialCharacterAIBuff()
	if start == nil or start.f_getActiveCharRef == nil or start.f_getCharData == nil then
		return
	end
	if roundstate() ~= 2 or main.pauseMenu or paused() then
		return
	end
	local oldid = id()
	for side = 1, 2 do
		if player(side) then
			local ref = start.f_getActiveCharRef(side)
			local data = ref ~= nil and start.f_getCharData(ref) or nil
			local key = ''
			if data ~= nil then
				key = table.concat({
					tostring(data.char or ''),
					tostring(data.def or ''),
					tostring(data.name or ''),
					tostring(data.displayname or ''),
					tostring(data.recordKey or ''),
				}, ' '):lower()
			end
			local target = key:find('mizuchi%-type%-m', 1, false) ~= nil
				or key:find('whiteness_of_g_kfm', 1, false) ~= nil
				or key:find('whiteness of g kfm', 1, false) ~= nil
				or key:find('whiteness of g kung fu man', 1, false) ~= nil
			if target and ailevel() > 0 then
				if ailevel() < 8 then
					setAILevel(8)
				end
				if gameTime() % 12 == 0 then
					setPower(math.min(powermax(), power() + 60))
					setGuardPoints(math.min(guardpointsmax(), guardpoints() + 15))
				end
				if gameTime() % 60 == 0 and life() > 0 and life() < lifemax() then
					setLife(math.min(lifemax(), life() + 5))
				end
				if gameTime() % 90 == 0 then
					removeDizzy()
				end
			end
		end
	end
	playerid(oldid)
end

local smokeActiveFrameCount = 0

local function smokeEndAfterActiveFrames()
	if main == nil or main.smokeEndAfterActiveFrames == nil then
		smokeActiveFrameCount = 0
		return
	end
	if roundstate() == 2 and not paused() and not main.pauseMenu then
		smokeActiveFrameCount = smokeActiveFrameCount + 1
		if smokeActiveFrameCount >= main.smokeEndAfterActiveFrames then
			endMatch()
		end
	else
		smokeActiveFrameCount = 0
	end
end

local smokeMotionState = nil

local function smokeMotionProbe()
	if getCommandLineValue("-smokemotion") == nil then
		smokeMotionState = nil
		return false
	end
	local targetFrames = tonumber(getCommandLineValue("-smokemotionframes")) or 3600
	local active = roundstate() == 2 and not paused()
	if not active then
		return false
	end
	local snap = afkPlayerSnapshot(1)
	if snap == nil then
		return false
	end
	if smokeMotionState == nil then
		smokeMotionState = {
			frames = 0,
			signature = afkMovementSignature(snap),
			signatureChanges = 0,
			stateChanges = 0,
			animChanges = 0,
			moveTypeChanges = 0,
			ctrlFrames = 0,
			attackFrames = 0,
			lastState = snap.state,
			lastAnim = snap.anim,
			lastMoveType = snap.movetype,
		}
	end
	local sig = afkMovementSignature(snap)
	if sig ~= smokeMotionState.signature then
		smokeMotionState.signatureChanges = smokeMotionState.signatureChanges + 1
		smokeMotionState.signature = sig
	end
	if snap.state ~= smokeMotionState.lastState then
		smokeMotionState.stateChanges = smokeMotionState.stateChanges + 1
		smokeMotionState.lastState = snap.state
	end
	if snap.anim ~= smokeMotionState.lastAnim then
		smokeMotionState.animChanges = smokeMotionState.animChanges + 1
		smokeMotionState.lastAnim = snap.anim
	end
	if snap.movetype ~= smokeMotionState.lastMoveType then
		smokeMotionState.moveTypeChanges = smokeMotionState.moveTypeChanges + 1
		smokeMotionState.lastMoveType = snap.movetype
	end
	if snap.ctrl then
		smokeMotionState.ctrlFrames = smokeMotionState.ctrlFrames + 1
	end
	if snap.movetype == 'A' then
		smokeMotionState.attackFrames = smokeMotionState.attackFrames + 1
	end
	smokeMotionState.frames = smokeMotionState.frames + 1
	if smokeMotionState.frames >= targetFrames then
		printConsole(string.format(
			'smokemotion: frames=%d signatureChanges=%d stateChanges=%d animChanges=%d moveTypeChanges=%d ctrlFrames=%d attackFrames=%d finalState=%s finalAnim=%s finalMoveType=%s',
			smokeMotionState.frames,
			smokeMotionState.signatureChanges,
			smokeMotionState.stateChanges,
			smokeMotionState.animChanges,
			smokeMotionState.moveTypeChanges,
			smokeMotionState.ctrlFrames,
			smokeMotionState.attackFrames,
			tostring(snap.state),
			tostring(snap.anim),
			tostring(snap.movetype)
		))
		endMatch()
		os.exit()
		return true
	end
	return false
end

local function saltyBetStandAnimation(p)
	local oldid = id()
	if player(p) then
		if stateNo() ~= 0 then
			selfState(0)
		end
		if anim() ~= 0 and animExist(0) then
			changeAnim(0)
		end
		playerid(oldid)
	end
end

local function saltyBetStandAll()
	for p = 1, 8 do
		saltyBetStandAnimation(p)
	end
end

function debugFlag(side)
	if start ~= nil and start.t_savedData.debugFlag ~= nil then
		start.t_savedData.debugflag[side] = true
	end
end

function closeMenu()
	main.pauseMenu = false
end

local afkStuckRoundResetSeconds = 20
local afkStuckRoundResetState = {}
local afkNoDamageSince = nil

local function afkStuckRoundResetActive()
	-- Disabled: the external no_damage_watchdog.py (480s) is the only reset timer.
	if true then return false end
	return roundstate() == 2 and not network() and not gamemode('training') and not main.pauseMenu and not paused()
end

local function afkStunnedState(st)
	return st == 5500 or (st >= 6565300 and st <= 6565303)
end

local function afkOptionalNumber(names)
	for _, name in ipairs(names) do
		local fn = _G[name]
		if type(fn) == 'function' then
			local ok, ret = pcall(fn)
			if ok and type(ret) == 'number' then
				return ret
			end
		end
	end
	return 0
end

local function afkWatchdogClock()
	local fps = afkOptionalNumber({'ticksPerSecond', 'tickspersecond'})
	if fps <= 0 then
		fps = 60
	end
	for _, name in ipairs({'gameTime', 'roundtime'}) do
		local fn = _G[name]
		if type(fn) == 'function' then
			local ok, ret = pcall(fn)
			if ok and type(ret) == 'number' then
				return ret, afkStuckRoundResetSeconds * fps, fps
			end
		end
	end
	return os.time(), afkStuckRoundResetSeconds, 1
end

local function afkQuantize(value, scale)
	return math.floor((value or 0) * scale + 0.5)
end

local function afkPlayerSnapshot(p)
	local oldid = id()
	playerid(p)
	local ret = nil
	if player(p) then
		ret = {
			life = life(),
			redlife = redlife(),
			state = stateno(),
			anim = anim(),
			ctrl = ctrl(),
			statetype = statetype(),
			movetype = movetype(),
			physics = physics(),
			posX = afkOptionalNumber({'posX'}),
			posY = afkOptionalNumber({'posY'}),
			posZ = afkOptionalNumber({'posZ'}),
			velX = afkOptionalNumber({'velX'}),
			velY = afkOptionalNumber({'velY'}),
			velZ = afkOptionalNumber({'velZ'}),
			stageBackEdgeDist = afkOptionalNumber({'stageBackEdgeDist', 'backEdgeBodyDist', 'backEdgeDist'}),
			stageFrontEdgeDist = afkOptionalNumber({'stageFrontEdgeDist', 'frontEdgeBodyDist', 'frontEdgeDist'}),
			topBoundDist = afkOptionalNumber({'topBoundDist', 'topBoundBodyDist'}),
			botBoundDist = afkOptionalNumber({'botBoundDist', 'botBoundBodyDist'}),
		}
	end
	playerid(oldid)
	return ret
end

local function afkMovementSignature(snap)
	return table.concat({
		afkQuantize(snap.posX, 2),
		afkQuantize(snap.posY, 2),
		afkQuantize(snap.posZ, 2),
		afkQuantize(snap.velX, 100),
		afkQuantize(snap.velY, 100),
		afkQuantize(snap.velZ, 100),
		snap.state,
		snap.anim,
		boolToInt(snap.ctrl),
		snap.statetype,
		snap.movetype,
		snap.physics,
	}, ':')
end

local function afkOutOfBounds(snap)
	return snap.stageBackEdgeDist < -160
		or snap.stageFrontEdgeDist < -160
		or snap.topBoundDist < -240
		or snap.botBoundDist < -240
		or snap.posY < -720
		or snap.posY > 720
end

local function afkSetPlayerPosition(x, y, z)
	if type(setPos) == 'function' then
		setPos(x, y, z or 0)
		return true
	end
	return false
end

local function afkSetPlayerVelocity(x, y, z)
	if type(setVel) == 'function' then
		setVel(x, y, z or 0)
		return true
	end
	return false
end

local function afkRescueOutOfBoundsPlayer(p, data, now, fps)
	if data.rescueCooldownUntil ~= nil and now < data.rescueCooldownUntil then
		return false
	end
	local oldid = id()
	local rescued = false
	playerid(p)
	if player(p) then
		local x = -80
		if p % 2 == 0 then
			x = 80
		end
		rescued = afkSetPlayerPosition(x, 0, 0)
		if rescued then
			afkSetPlayerVelocity(0, 0, 0)
			selfState(0)
			if type(animExist) == 'function' and animExist(0) then
				changeAnim(0)
			end
			data.rescueCooldownUntil = now + math.max(30, math.floor((fps or 60) / 2))
			data.outOfBoundsSince = nil
			data.signature = nil
			data.stillSince = now
			printConsole('antistuck: moved player ' .. tostring(p) .. ' back into stage')
		end
	end
	playerid(oldid)
	return rescued
end

local function afkRoundLifeSignature()
	local ret = {}
	for p = 1, 8 do
		local snap = afkPlayerSnapshot(p)
		if snap then
			table.insert(ret, p .. ':' .. snap.life .. ':' .. snap.redlife)
		end
	end
	return table.concat(ret, '|')
end

local function afkRoundMovementSignature()
	local ret = {}
	for p = 1, 8 do
		local snap = afkPlayerSnapshot(p)
		if snap and snap.life > 0 then
			table.insert(ret, p .. ':' .. afkMovementSignature(snap))
		end
	end
	return table.concat(ret, '|')
end

local function afkResetRound(oldid)
	afkStuckRoundResetState = {}
	afkNoDamageSince = nil
	playerid(oldid)
	roundReset()
	closeMenu()
end

hook.add("loop", "afkStuckRoundReset", function()
	if not afkStuckRoundResetActive() then
		afkStuckRoundResetState = {}
		afkNoDamageSince = nil
		return
	end
	local oldid = id()
	local now, resetLimit, fps = afkWatchdogClock()
	local roundLifeSignature = afkRoundLifeSignature()
	local roundMovementSignature = afkRoundMovementSignature()
	if afkNoDamageSince == nil then
		afkNoDamageSince = now
	elseif afkStuckRoundResetState.roundLifeSignature ~= roundLifeSignature
		or afkStuckRoundResetState.roundMovementSignature ~= roundMovementSignature then
		afkNoDamageSince = now
	end
	for p = 1, 8 do
		local snap = afkPlayerSnapshot(p)
		if snap and snap.life > 0 then
			local data = afkStuckRoundResetState[p] or {}
			local signature = afkMovementSignature(snap)
			if data.signature ~= signature then
				data.signature = signature
				data.stillSince = now
			elseif data.stillSince ~= nil and now - data.stillSince >= resetLimit then
				afkResetRound(oldid)
				return
			end
				if (afkStunnedState(snap.state) or (snap.anim == 5300 and snap.movetype == 'I')) and not snap.ctrl then
					data.stunnedSince = data.stunnedSince or now
					if now - data.stunnedSince >= resetLimit then
						afkResetRound(oldid)
						return
					end
				else
					data.stunnedSince = nil
				end
				if afkOutOfBounds(snap) then
					if not afkRescueOutOfBoundsPlayer(p, data, now, fps) then
						data.outOfBoundsSince = data.outOfBoundsSince or now
						if now - data.outOfBoundsSince >= resetLimit then
							afkResetRound(oldid)
							return
						end
					end
				else
					data.outOfBoundsSince = nil
				end
				afkStuckRoundResetState[p] = data
		else
			afkStuckRoundResetState[p] = nil
		end
	end
	if afkNoDamageSince ~= nil and now - afkNoDamageSince >= resetLimit then
		afkResetRound(oldid)
		return
	end
	afkStuckRoundResetState.roundLifeSignature = roundLifeSignature
	afkStuckRoundResetState.roundMovementSignature = roundMovementSignature
	playerid(oldid)
end)

--;===========================================================
--; MCONSOLE EQUIVALENTS
--;===========================================================
function toggleDebugPause()
	togglePause()
	closeMenu()
end

function toggleMaxPowerModeAll() -- maxpowermode
	toggleMaxPowerMode()
	debugFlag(1)
	debugFlag(2)
end

function matchReload() -- matchreset
	reload()
	closeMenu()
end

function powMaxAll()
	powMax(1)
	powMax(2)
	debugFlag(1)
	debugFlag(2)
end

function roundResetNow()
	roundReset()
	closeMenu()
end

function fullAll()
	full(1)
	full(2)
	full(3)
	full(4)
	full(5)
	full(6)
	full(7)
	full(8)
	setTime(getRoundTime())
	debugFlag(1)
	debugFlag(2)
	clearConsole()
end

function standAll()
	stand(1)
	stand(2)
	stand(3)
	stand(4)
	stand(5)
	stand(6)
	stand(7)
	stand(8)
end

--;===========================================================
--; DEBUG STATUS INFO
--;===========================================================
function statusInfo(p)
	local oldid = id()
	if not player(p) then return false end
	local ret = string.format(
		'P%d: %d; LIF:%4d; POW:%4d; ATK:%4d; DEF:%4d; RED:%4d; GRD:%4d; STN:%4d',
		playerno(), id(), life(), power(), attack(), defence(), redlife(), guardpoints(), dizzypoints()
	)
	playerid(oldid)
	return ret
end

loadDebugStatus('statusInfo')

--;===========================================================
--; DEBUG PLAYER/HELPER INFO
--;===========================================================
function customState()
	if not incustomstate() then return "" end
	return " (in " .. stateownername() .. " " .. stateownerid() .. "'s state)"
end

function boolToInt(bool)
	if bool then return 1 end
	return 0
end

function engineInfo()
	return string.format('Frames: %d, VSync: %d; Speed: %d/%d%%; FPS: %.3f', roundtime(), gameOption('Video.VSync'), tickspersecond(), gamespeed(), gamefps())
end

function playerInfo()
	return string.format('%s %d%s', name(), id(), customState())
end

function actionInfo()
	return string.format(
		'ActionID: %d (P%d); SPR: %d,%d; ElemNo: %d/%d; Time: %d/%d (%d/%d)',
		anim(), animplayerno(), animelemvar("group"), animelemvar("image"), animelemno(0), animelemcount(), animelemtime(animelemno(0)), animelemvar("time"), animtimesum(), animlength()
	)
end

function stateInfo()
	return string.format(
		'State No: %d (P%d); CTRL: %s; Type: %s; MoveType: %s; Physics: %s; Time: %d',
		stateno(), stateownerplayerno(), boolToInt(ctrl()), statetype(), movetype(), physics(), time()-1
	)
end

loadDebugInfo({'engineInfo', 'playerInfo', 'actionInfo', 'stateInfo'})

--;===========================================================
--; MATCH LOOP
--;===========================================================
local endFlag = false

--function called during match via config.ini CommonLua
function loop()
	hook.run("loop")
	if start == nil then --match started via command line without -loadmotif flag
		if smokeMotionProbe() then
			return
		end
		if esc() then
			endMatch()
			os.exit()
		end
		if indialogue() then
			dialogueReset()
		end
		togglePostMatch(false)
		toggleDialogueBars(false)
		return
	end
	--credits
	if main.credits ~= -1 and type(motif.attract_mode.credits_snd) == 'table' and getKey(motif.attract_mode.credits_key) then
		sndPlay(motif.files.snd_data, motif.attract_mode.credits_snd[1], motif.attract_mode.credits_snd[2])
		main.credits = main.credits + 1
		resetKey()
	end
	--music
	if type(start.f_stageMusic) == 'function' then
		start.f_stageMusic()
	end
	--match start
	if roundstart() then
		setLifebarElements({bars = main.lifebar.bars})
		if roundno() == 1 then
			speedMul = 1
			speedAdd = 0
			start.charRecordSaved = false
			start.victoryInit = false
			start.resultInit = false
			start.continueInit = false
			start.hiscoreInit = false
			endFlag = false
			if indialogue() then
				dialogueReset()
			end
			if gamemode('training') then
				menu.f_trainingReset()
			end
		end
		start.turnsRecoveryInit = false
		start.roundRecordSaved = false
		start.dialogueInit = false
	end
	if winnerteam() ~= -1 and player(winnerteam()) and roundstate() == 4 and isasserted("over") then
		if start.p[1] ~= nil and start.p[2] ~= nil and (start.p[1].teamMode ~= 0 or start.p[2].teamMode ~= 0) then
			start.f_updateRoundCharRecords(winnerteam())
		end
		--turns or solo-vs-team life recovery
		start.f_turnsRecovery()
	end
	--dialogue
	if indialogue() then
		start.f_dialogue()
	--match end
	elseif postmatch() then
		if not endFlag then
			start.f_recordCharMatchResult(winnerteam())
			resetMatchData(false)
			endFlag = true
			end
			--victory screen
			if start.f_victory() then
				if type(start.f_drawTierChangeOverlay) == 'function' then
					start.f_drawTierChangeOverlay()
			end
				return
			--result screen
			elseif start.f_result() then
				if type(start.f_drawTierChangeOverlay) == 'function' then
					start.f_drawTierChangeOverlay()
			end
				return
			--continue screen
			elseif start.f_continue() then
				if type(start.f_drawTierChangeOverlay) == 'function' then
					start.f_drawTierChangeOverlay()
			end
				return
			end
			clearColor(motif.selectbgdef.bgclearcolor[1], motif.selectbgdef.bgclearcolor[2], motif.selectbgdef.bgclearcolor[3])
			if type(start.f_drawTierChangeOverlay) == 'function' then
				start.f_drawTierChangeOverlay()
			end
			togglePostMatch(false)
			end
		hook.run("loop#" .. gamemode())
					kfmFactionAIBuff()
					specialCharacterAIBuff()
					smokeEndAfterActiveFrames()
					if start ~= nil and not indialogue() and not postmatch() and not main.pauseMenu and roundstate() > 0 and (roundstate() < 4 or (main.currentTournamentFinalMatch and roundstate() == 4)) then
						if start.txt_fightRecord == nil or start.txt_fightTierUnderName == nil or start.txt_fightFactionUnderTier == nil or start.txt_fightFactionRankUnderFaction == nil or start.txt_fightRecordUnderTier == nil or start.txt_tournamentStatus == nil or start.txt_tournamentFighters == nil or start.txt_tournamentConcluded == nil or start.txt_countdownMatchInfo == nil or start.txt_matchesUntilTournament == nil then
					local function makeFightRecordText(x, y, align, r, g, b, scale)
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
					start.txt_fightRecord = {
						shadow = {
							makeFightRecordText(114, 28, 1, 0, 0, 0, 1.35),
							makeFightRecordText(1166, 28, -1, 0, 0, 0, 1.35),
						},
						value = {
							makeFightRecordText(112, 26, 1, 255, 255, 255, 1.35),
							makeFightRecordText(1168, 26, -1, 255, 255, 255, 1.35),
						},
					}
						start.txt_fightTierUnderName = {
							shadow = {
								makeFightRecordText(114, 44, 1, 0, 0, 0, 1.45),
								makeFightRecordText(1166, 44, -1, 0, 0, 0, 1.45),
							},
							value = {
								makeFightRecordText(112, 42, 1, 255, 255, 255, 1.45),
								makeFightRecordText(1168, 42, -1, 255, 255, 255, 1.45),
							},
						}
						start.txt_fightFactionUnderTier = {
							shadow = {
								makeFightRecordText(114, 60, 1, 0, 0, 0, 1.05),
								makeFightRecordText(1166, 60, -1, 0, 0, 0, 1.05),
							},
							value = {
								makeFightRecordText(112, 58, 1, 215, 235, 255, 1.05),
								makeFightRecordText(1168, 58, -1, 215, 235, 255, 1.05),
							},
						}
						start.txt_fightFactionRankUnderFaction = {
							shadow = {
								makeFightRecordText(114, 76, 1, 0, 0, 0, 1.0),
								makeFightRecordText(1166, 76, -1, 0, 0, 0, 1.0),
							},
							value = {
								makeFightRecordText(112, 74, 1, 255, 210, 130, 1.0),
								makeFightRecordText(1168, 74, -1, 255, 210, 130, 1.0),
							},
						}
						start.txt_fightRecordUnderTier = {
							winsShadow = {
								makeFightRecordText(114, 92, 1, 0, 0, 0, 1.2),
								makeFightRecordText(1166, 92, -1, 0, 0, 0, 1.2),
							},
							winsValue = {
								makeFightRecordText(112, 90, 1, 80, 255, 120, 1.2),
								makeFightRecordText(1168, 90, -1, 80, 255, 120, 1.2),
							},
							lossesShadow = {
								makeFightRecordText(114, 108, 1, 0, 0, 0, 1.2),
								makeFightRecordText(1166, 108, -1, 0, 0, 0, 1.2),
							},
							lossesValue = {
								makeFightRecordText(112, 106, 1, 255, 70, 70, 1.2),
								makeFightRecordText(1168, 106, -1, 255, 70, 70, 1.2),
							},
						}
						start.txt_saltyBetCountdown = {
							shadow = makeFightRecordText(642, 324, 0, 0, 0, 0, 2.75),
							value = makeFightRecordText(640, 322, 0, 255, 255, 255, 2.75),
						}
								start.txt_tournamentStatus = {
									shadow = makeFightRecordText(642, 14, 0, 0, 0, 0, 1.2),
									value = makeFightRecordText(640, 12, 0, 255, 230, 80, 1.2),
								}
									start.txt_countdownMatchInfo = {
										shadow = makeFightRecordText(642, 204, 0, 0, 0, 0, 1.45),
										value = makeFightRecordText(640, 202, 0, 255, 255, 255, 1.45),
									}
										start.txt_matchesUntilTournament = {
											shadow = makeFightRecordText(642, 304, 0, 0, 0, 0, 1.55),
											value = makeFightRecordText(640, 302, 0, 255, 230, 40, 1.55),
										}
									start.txt_tournamentFighters = {
									shadow = {},
									value = {},
								}
							for i = 1, 16 do
								start.txt_tournamentFighters.shadow[i] = makeFightRecordText(0, 0, 0, 0, 0, 0, 0.72)
								start.txt_tournamentFighters.value[i] = makeFightRecordText(0, 0, 0, 230, 230, 230, 0.72)
							end
							start.txt_tournamentConcluded = {
								shadow = makeFightRecordText(0, 252, 0, 0, 0, 0, 2.2),
								value = makeFightRecordText(0, 250, 0, 255, 230, 80, 2.2),
							}
						end
					for side = 1, 2 do
						local ref = start.f_getActiveCharRef(side)
						if ref ~= nil then
						local record = start.f_getCharRecord(ref)
						local tierText = start.f_getRecordTierDisplayText(record)
							local align = side == 1 and 1 or -1
							local r, g, b = start.f_getRecordTierColor(record)
								local tierX = side == 1 and 112 or 1168
								local tierShadowX = side == 1 and 114 or 1166
								local tierY = 42
							start.txt_fightTierUnderName.shadow[side]:update({text = tierText, x = tierShadowX, y = tierY + 2, align = align, r = 0, g = 0, b = 0})
								start.txt_fightTierUnderName.shadow[side]:draw()
								start.txt_fightTierUnderName.value[side]:update({text = tierText, x = tierX, y = tierY, align = align, r = r, g = g, b = b})
									start.txt_fightTierUnderName.value[side]:draw()
										local factionText = nil
										if start.f_getActiveCharFaction ~= nil then
											factionText = start.f_getActiveCharFaction(side)
						elseif start.f_getCharFaction ~= nil then
							factionText = start.f_getCharFaction(ref)
						end
						if factionText == nil or tostring(factionText) == '' then
							local liveData = main.t_selChars ~= nil and main.t_selChars[tonumber(ref) + 1] or nil
							if liveData ~= nil and liveData.faction ~= nil and tostring(liveData.faction) ~= '' then
								factionText = liveData.faction
							end
							local identity = tostring(start.f_getCharData(ref) and (start.f_getCharData(ref).name or start.f_getCharData(ref).char or '') or ''):lower()
							if identity:find('knuckles', 1, true) then
								factionText = 'Lancero'
							elseif identity:find('riki', 1, true) then
								factionText = 'DBC'
							end
						end
										if main.debugFactionTraceOnce ~= true then
											main.debugFactionTraceOnce = true
											printConsole('DEBUGFACTION: side=' .. tostring(side) .. ' ref=' .. tostring(ref) .. ' char=' .. tostring(start.f_getCharData(ref) and start.f_getCharData(ref).char) .. ' factionText=' .. tostring(factionText) .. ' hasActiveFn=' .. tostring(start.f_getActiveCharFaction ~= nil) .. ' hasFactionFn=' .. tostring(start.f_getCharFaction ~= nil))
										end
										if factionText ~= nil and tostring(factionText) ~= '' then
											local factionDisplay = tostring(factionText):gsub('%s+[Ff]action%s*$', '') .. ' Faction'
											local fr, fg, fb = 215, 235, 255
											if start.f_getFactionColor ~= nil then
												fr, fg, fb = start.f_getFactionColor(factionText)
											end
											start.txt_fightFactionUnderTier.shadow[side]:update({text = factionDisplay, x = tierShadowX, y = 62, align = align, r = 0, g = 0, b = 0})
											start.txt_fightFactionUnderTier.shadow[side]:draw()
											start.txt_fightFactionUnderTier.value[side]:update({text = factionDisplay, x = tierX, y = 60, align = align, r = fr, g = fg, b = fb})
											start.txt_fightFactionUnderTier.value[side]:draw()
											if start.f_getCharFactionRank ~= nil then
												local factionRank, factionTotal = start.f_getCharFactionRank(ref)
												if factionRank ~= nil and factionTotal ~= nil then
													local rankDisplay = string.format('Rank %d/%d', factionRank, factionTotal)
													start.txt_fightFactionRankUnderFaction.shadow[side]:update({text = rankDisplay, x = tierShadowX, y = 78, align = align, r = 0, g = 0, b = 0})
													start.txt_fightFactionRankUnderFaction.shadow[side]:draw()
													start.txt_fightFactionRankUnderFaction.value[side]:update({text = rankDisplay, x = tierX, y = 76, align = align, r = 255, g = 210, b = 130})
													start.txt_fightFactionRankUnderFaction.value[side]:draw()
												end
											end
										end
										local winsText = string.format('Wins: %d', tonumber(record.wins) or 0)
										local lossesText = string.format('Losses: %d', tonumber(record.losses) or 0)
									local recordX = side == 1 and 112 or 1168
									local recordShadowX = side == 1 and 114 or 1166
								start.txt_fightRecordUnderTier.winsShadow[side]:update({text = winsText, x = recordShadowX, y = 92, align = align, r = 0, g = 0, b = 0})
								start.txt_fightRecordUnderTier.winsShadow[side]:draw()
								start.txt_fightRecordUnderTier.winsValue[side]:update({text = winsText, x = recordX, y = 90, align = align, r = 80, g = 255, b = 120})
								start.txt_fightRecordUnderTier.winsValue[side]:draw()
								start.txt_fightRecordUnderTier.lossesShadow[side]:update({text = lossesText, x = recordShadowX, y = 108, align = align, r = 0, g = 0, b = 0})
								start.txt_fightRecordUnderTier.lossesShadow[side]:draw()
								start.txt_fightRecordUnderTier.lossesValue[side]:update({text = lossesText, x = recordX, y = 106, align = align, r = 255, g = 70, b = 70})
								start.txt_fightRecordUnderTier.lossesValue[side]:draw()
								end
					end
					if main.currentTournamentFinalMatch and roundstate() == 4 and winnerteam() ~= -1 and start.txt_tournamentConcluded ~= nil then
						local scrollText = 'Tournament Concluded.'
						local scrollX = 1450 - ((gameTime() * 5) % 1900)
						start.txt_tournamentConcluded.shadow:update({text = scrollText, x = scrollX + 3, y = 252, align = 0, r = 0, g = 0, b = 0})
						start.txt_tournamentConcluded.shadow:draw()
						start.txt_tournamentConcluded.value:update({text = scrollText, x = scrollX, y = 250, align = 0, r = 255, g = 230, b = 80})
						start.txt_tournamentConcluded.value:draw()
					end
							if roundno() == 1 and roundstate() == 2 and not start.saltyBetRoundResetDone then
							if start.saltyBetHoldStartGameTime == nil then
								start.saltyBetHoldStartGameTime = gameTime()
							end
								local remaining = 30 - math.floor((gameTime() - start.saltyBetHoldStartGameTime) / 60)
										if remaining > 0 then
											saltyBetStandAll()
											local countdownInfoLines = {}
											local matchupLines = {}
											local fightsUntilTournamentText = nil
										local function countdownFighterLine(side)
											if start.f_getActiveCharRef == nil or start.f_getCharRecord == nil then
												return nil
											end
											local ref = start.f_getActiveCharRef(side)
											if ref == nil then
												return nil
											end
											local data = start.f_getCharData ~= nil and start.f_getCharData(ref) or nil
											local name = data and (data.displayname or data.name or data.char) or tostring(ref)
											name = tostring(name or ''):gsub('^%s+', ''):gsub('%s+$', '')
											if #name > 24 then
												name = name:sub(1, 23) .. '.'
											end
											local record = start.f_getCharRecord(ref)
							local tierText = start.f_getRecordTierDisplayText ~= nil and start.f_getRecordTierDisplayText(record) or tostring(record.tier or 'U')
							local factionText = start.f_getCharFaction ~= nil and start.f_getCharFaction(ref) or nil
							local factionRank, factionTotal = nil, nil
							if start.f_getCharFactionRank ~= nil then
								factionRank, factionTotal = start.f_getCharFactionRank(ref)
							end
							local factionRankText = ''
							if factionRank ~= nil and factionTotal ~= nil then
								factionRankText = string.format('  Rank %d/%d', tonumber(factionRank) or 0, tonumber(factionTotal) or 0)
							end
							if factionText ~= nil and tostring(factionText) ~= '' then
								return string.format('P%d %s  %s Tier  %s  Faction %s%s  W:%d  L:%d', side, name, tierText, tostring(factionText), factionRankText, tonumber(record.wins) or 0, tonumber(record.losses) or 0)
							end
											return string.format('P%d %s  %s Tier  W:%d  L:%d', side, name, tierText, tonumber(record.wins) or 0, tonumber(record.losses) or 0)
										end
											if start.txt_countdownMatchInfo ~= nil then
											for side = 1, 2 do
												local line = countdownFighterLine(side)
												if line ~= nil then
													table.insert(matchupLines, line)
												end
											end
										end
										local tournamentLines = {}
												if gameMode() == 'utiergraduation' and main.uTierGraduationStatus ~= nil then
													table.insert(tournamentLines, tostring(main.uTierGraduationStatus))
												end
												if main.endlessRandomTournamentLoading then
													table.insert(tournamentLines, 'LOADING RANDOM TOURNAMENT')
												end
											if main.currentTournamentName ~= nil and tostring(main.currentTournamentName) ~= '' then
												table.insert(tournamentLines, 'TOURNAMENT: ' .. tostring(main.currentTournamentName))
										end
												if main.endlessRandomActive and not main.endlessRandomTournamentLoading and main.currentTournamentName == nil and main.endlessRandomFightsUntilTournament ~= nil and start.txt_matchesUntilTournament ~= nil then
													local fightsLeft = tonumber(main.endlessRandomFightsUntilTournament) or 0
												if fightsLeft >= 0 then
													local matchesText = fightsLeft == 0 and 'Tournament is Approaching!' or string.format('%d fights until next tournament', fightsLeft)
													fightsUntilTournamentText = matchesText
													start.txt_matchesUntilTournament.shadow:update({text = matchesText, x = 642, y = 304, align = 0, r = 0, g = 0, b = 0})
													start.txt_matchesUntilTournament.shadow:draw()
													start.txt_matchesUntilTournament.value:update({text = matchesText, x = 640, y = 302, align = 0, r = 255, g = 230, b = 40})
													start.txt_matchesUntilTournament.value:draw()
											end
										end
										if #tournamentLines > 0 then
											for _, line in ipairs(tournamentLines) do
												table.insert(countdownInfoLines, line)
											end
										end
										if #matchupLines > 0 then
											for _, line in ipairs(matchupLines) do
												table.insert(countdownInfoLines, line)
											end
										end
										if #countdownInfoLines > 0 and start.txt_countdownMatchInfo ~= nil then
											local countdownInfoText = table.concat(countdownInfoLines, '\n')
											start.txt_countdownMatchInfo.shadow:update({text = countdownInfoText, x = 642, y = 204, align = 0, r = 0, g = 0, b = 0})
											start.txt_countdownMatchInfo.shadow:draw()
											start.txt_countdownMatchInfo.value:update({text = countdownInfoText, x = 640, y = 202, align = 0, r = 255, g = 255, b = 255})
											start.txt_countdownMatchInfo.value:draw()
										end
										if #tournamentLines > 0 then
										local tournamentText = table.concat(tournamentLines, '\n')
										start.txt_tournamentStatus.shadow:update({text = tournamentText, x = 642, y = 14, align = 0, r = 0, g = 0, b = 0})
										start.txt_tournamentStatus.shadow:draw()
										start.txt_tournamentStatus.value:update({text = tournamentText, x = 640, y = 12, align = 0, r = 255, g = 230, b = 80})
										start.txt_tournamentStatus.value:draw()
									end
										if type(main.currentTournamentFighters) == 'table' and start.txt_tournamentFighters ~= nil then
										for i, fighter in ipairs(main.currentTournamentFighters) do
											if i <= 16 and fighter ~= nil then
													local col = (i - 1) % 4
													local row = math.floor((i - 1) / 4)
													local x = 205 + col * 290
													local y = 112 + row * 30
													local rankText = string.format('#%02d', tonumber(fighter.rank) or i)
													if fighter.eliminated then
														rankText = rankText .. '  E'
													end
													local wins = 0
													local losses = 0
													if fighter.ref ~= nil and start.f_getCharRecord ~= nil then
														local record = start.f_getCharRecord(fighter.ref)
														wins = tonumber(record.wins) or 0
														losses = tonumber(record.losses) or 0
													end
													local fighterText = tostring(fighter.name or '') .. '\n' .. string.format('%s  W:%d L:%d', rankText, wins, losses)
												start.txt_tournamentFighters.shadow[i]:update({text = fighterText, x = x + 2, y = y + 2, align = 0, r = 0, g = 0, b = 0})
												start.txt_tournamentFighters.shadow[i]:draw()
												start.txt_tournamentFighters.value[i]:update({text = fighterText, x = x, y = y, align = 0, r = fighter.eliminated and 180 or 230, g = fighter.eliminated and 80 or 230, b = fighter.eliminated and 80 or 230})
												start.txt_tournamentFighters.value[i]:draw()
											end
										end
									end
									local countdownText = string.format('MATCH STARTS IN\n%d', remaining)
									if fightsUntilTournamentText ~= nil then
										countdownText = countdownText .. '\n' .. fightsUntilTournamentText
									end
									start.txt_saltyBetCountdown.shadow:update({text = countdownText, x = 642, y = 324, align = 0, r = 0, g = 0, b = 0})
									start.txt_saltyBetCountdown.shadow:draw()
									start.txt_saltyBetCountdown.value:update({text = countdownText, x = 640, y = 322, align = 0, r = 255, g = 230, b = 40})
								start.txt_saltyBetCountdown.value:draw()
							else
								start.saltyBetRoundResetDone = true
								start.saltyBetHoldStartGameTime = nil
								roundReset()
							end
						else
							if not start.saltyBetRoundResetDone then
								start.saltyBetHoldStartGameTime = nil
							end
				end
			end
			if start ~= nil and type(start.f_drawPlacementGrid) == 'function' and not postmatch() then
			start.f_drawPlacementGrid()
		end
	--pause menu
	if main.pauseMenu then
		playerBufReset()
		menu.f_run()
	else
		main.f_cmdInput()
		--esc / m
		if (esc() or (main.f_input(main.t_players, {'m'}) and not network())) and not start.challengerInit then
			if network() or gamemode('demo') or (not gameOption('Config.EscOpensMenu') and esc()) then
			endMatch()
			else
				menu.f_init()
			end
		--demo mode
		elseif gamemode('demo') and ((motif.attract_mode.enabled == 1 and main.credits > 0 and not sndPlaying(motif.files.snd_data, motif.attract_mode.credits_snd[1], motif.attract_mode.credits_snd[2])) or (motif.attract_mode.enabled == 0 and main.f_input(main.t_players, {'pal'})) or fighttime() >= motif.demo_mode.fight_endtime) then
			endMatch()
		--challenger
		elseif motif.challenger_info.enabled ~= 0 and gamemode('arcade') then
			if start.challenger > 0 then
				start.f_challenger()
			else
				--TODO: detecting players that are part of P1 team
				--[[for i = 1, #main.t_cmd do
					if commandGetState(main.t_cmd[i], '/s') then
						print(i)
			end
			end]]
				if main.f_input(main.t_players, {'s'}) and main.playerInput ~= 1 and (motif.attract_mode.enabled == 0 or main.credits ~= 0) then
					start.challenger = main.playerInput
			end
			end
		end
	end
end
