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
	if main.credits ~= -1 and getKey(motif.attract_mode.credits_key) then
		sndPlay(motif.files.snd_data, motif.attract_mode.credits_snd[1], motif.attract_mode.credits_snd[2])
		main.credits = main.credits + 1
		resetKey()
	end
	--music
	start.f_stageMusic()
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
			return
		--result screen
		elseif start.f_result() then
			return
		--continue screen
		elseif start.f_continue() then
			return
		end
		clearColor(motif.selectbgdef.bgclearcolor[1], motif.selectbgdef.bgclearcolor[2], motif.selectbgdef.bgclearcolor[3])
		togglePostMatch(false)
		end
		hook.run("loop#" .. gamemode())
		if start ~= nil and not indialogue() and not postmatch() and not main.pauseMenu and roundstate() > 0 and roundstate() < 4 then
			if start.txt_fightRecord == nil then
				start.txt_fightRecord = {
					text:create({
						font = 8,
						bank = 0,
						align = 1,
						text = '',
						x = 578,
						y = 24,
						scaleX = 1,
						scaleY = 1,
						r = 255,
						g = 255,
						b = 255,
						a = 255,
						height = -1,
						xshear = 0,
						angle = 0,
						window = nil,
						defsc = false,
					}),
					text:create({
						font = 8,
						bank = 0,
						align = -1,
						text = '',
						x = 704,
						y = 24,
						scaleX = 1,
						scaleY = 1,
						r = 255,
						g = 255,
						b = 255,
						a = 255,
						height = -1,
						xshear = 0,
						angle = 0,
						window = nil,
						defsc = false,
					}),
				}
			end
			for side = 1, 2 do
				local ref = start.f_getActiveCharRef(side)
				if ref ~= nil then
					local record = start.f_getCharRecord(ref)
					local textValue = string.format('W:%d L:%d %s', record.wins, record.losses, start.f_getRecordTierText(record))
					local x = side == 1 and 578 or 704
					local align = side == 1 and 1 or -1
					start.txt_fightRecord[side]:update({text = textValue, x = x, y = 24, align = align})
					start.txt_fightRecord[side]:draw()
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
