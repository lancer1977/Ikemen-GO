package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const defaultLiveOverlayPath = "save/live_overlays.json"
const liveOverlaySmokeCaptureEnv = "IKEMEN_OVERLAY_SMOKE_CAPTURE"

type LiveOverlayFile struct {
	Schema         string              `json:"schema"`
	SourceSequence int64               `json:"sourceSequence,omitempty"`
	Overlays       []LiveOverlayEffect `json:"overlays"`
}

type LiveOverlayEffect struct {
	ID         string     `json:"id"`
	Kind       string     `json:"kind"`
	Text       string     `json:"text,omitempty"`
	AssetPath  string     `json:"assetPath,omitempty"`
	X          float32    `json:"x,omitempty"`
	Y          float32    `json:"y,omitempty"`
	Align      int32      `json:"align,omitempty"`
	Layer      int16      `json:"layer,omitempty"`
	DurationMs int32      `json:"durationMs,omitempty"`
	Color      [4]int32   `json:"color,omitempty"`
	Scale      [2]float32 `json:"scale,omitempty"`
}

type liveOverlayRuntime struct {
	effect         LiveOverlayEffect
	expiresAtFrame int32
	textSprite     *TextSprite
	sprite         *Sprite
}

func (s *System) liveOverlayPath() string {
	if path := strings.TrimSpace(s.cmdFlags["-overlayfile"]); path != "" {
		return path
	}
	return defaultLiveOverlayPath
}

func (s *System) maybeRefreshLiveOverlays() {
	path := s.liveOverlayPath()
	if path == "" {
		s.liveOverlayRuntime = nil
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			LogMessage("live overlay read failed: %v", err)
		}
		s.liveOverlayRuntime = nil
		return
	}

	file, err := parseLiveOverlayFile(data)
	if err != nil {
		LogMessage("live overlay parse failed: %v", err)
		return
	}

	runtime := make([]*liveOverlayRuntime, 0, len(file.Overlays))
	for _, effect := range file.Overlays {
		if entry := s.buildLiveOverlayRuntime(effect); entry != nil {
			runtime = append(runtime, entry)
		}
	}
	s.liveOverlayRuntime = runtime
	if os.Getenv(liveOverlaySmokeCaptureEnv) != "" {
		LogMessage("live overlay refresh: loaded=%d capturePending=%v done=%v", len(runtime), s.liveOverlayCapturePending, s.liveOverlayCaptureDone)
	}
	if len(runtime) > 0 && os.Getenv(liveOverlaySmokeCaptureEnv) != "" && !s.liveOverlayCaptureDone {
		s.liveOverlayCapturePending = true
	}
}

func parseLiveOverlayFile(data []byte) (*LiveOverlayFile, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return &LiveOverlayFile{Schema: "live-lancero/overlay-state/v1"}, nil
	}

	var file LiveOverlayFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	if file.Schema == "" {
		file.Schema = "live-lancero/overlay-state/v1"
	}
	return &file, nil
}

func (s *System) buildLiveOverlayRuntime(effect LiveOverlayEffect) *liveOverlayRuntime {
	kind := strings.ToLower(strings.TrimSpace(effect.Kind))
	if kind == "" {
		return nil
	}

	durationFrames := effectDurationFrames(effect.DurationMs)
	entry := &liveOverlayRuntime{
		effect:         effect,
		expiresAtFrame: s.frameCounter + durationFrames,
	}

	switch kind {
	case "text", "emoji":
		if s.debugFont == nil || s.debugFont.fnt == nil {
			return nil
		}
		ts := s.debugFont.Copy()
		ts.Reset()
		if effect.Align != 0 {
			ts.align = effect.Align
		}
		if effect.Layer != 0 {
			ts.layerno = effect.Layer
		} else {
			ts.layerno = 2
		}
		ts.SetPos(effect.X, effect.Y)
		if effect.Scale != [2]float32{} {
			ts.SetScale(effect.Scale[0], effect.Scale[1])
		}
		ts.text = effect.Text
		ts.textInit = effect.Text
		if effect.Color != [4]int32{} {
			ts.SetColor(effect.Color[0], effect.Color[1], effect.Color[2], effect.Color[3])
		}
		ts.SetWindow([4]float32{0, 0, float32(sys.scrrect[2]), float32(sys.scrrect[3])})
		entry.textSprite = ts
	case "image", "emote":
		sprite, err := loadOverlaySprite(effect.AssetPath)
		if err != nil {
			LogMessage("live overlay image load failed: %v", err)
			return nil
		}
		entry.sprite = sprite
	default:
		return nil
	}

	return entry
}

func loadOverlaySprite(assetPath string) (*Sprite, error) {
	assetPath = strings.TrimSpace(assetPath)
	if assetPath == "" {
		return nil, fmt.Errorf("missing assetPath")
	}

	path := assetPath
	if !filepath.IsAbs(path) {
		path = filepath.Clean(path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	rect := img.Bounds()
	rgba, ok := img.(*image.RGBA)
	if !ok {
		rgba = image.NewRGBA(rect)
		draw.Draw(rgba, rect, img, rect.Min, draw.Src)
	}

	sprite := newSprite()
	sprite.Size = [2]uint16{uint16(rect.Dx()), uint16(rect.Dy())}
	sprite.Offset = [2]int16{0, 0}
	sprite.SetRaw(rgba.Pix, int32(rect.Dx()), int32(rect.Dy()), 32)
	return sprite, nil
}

func (s *System) drawLiveOverlays() {
	if s.frameSkip || len(s.liveOverlayRuntime) == 0 {
		return
	}

	kept := s.liveOverlayRuntime[:0]
	for _, entry := range s.liveOverlayRuntime {
		if entry == nil || s.frameCounter > entry.expiresAtFrame {
			continue
		}

		switch strings.ToLower(strings.TrimSpace(entry.effect.Kind)) {
		case "text", "emoji":
			if entry.textSprite != nil {
				if os.Getenv(liveOverlaySmokeCaptureEnv) != "" {
					LogMessage("live overlay draw text id=%s layer=%d frame=%d", entry.effect.ID, entry.textSprite.layerno, s.frameCounter)
				}
				entry.textSprite.Draw(entry.textSprite.layerno)
			}
		case "image", "emote":
			if entry.sprite != nil && !entry.sprite.isBlank() {
				scaleX, scaleY := float32(1), float32(1)
				if entry.effect.Scale != [2]float32{} {
					scaleX = entry.effect.Scale[0]
					scaleY = entry.effect.Scale[1]
				}
				if os.Getenv(liveOverlaySmokeCaptureEnv) != "" {
					LogMessage("live overlay draw image id=%s x=%.1f y=%.1f scale=%.2f,%.2f frame=%d", entry.effect.ID, entry.effect.X, entry.effect.Y, scaleX, scaleY, s.frameCounter)
				}
				entry.sprite.Draw(entry.effect.X, entry.effect.Y, scaleX, scaleY, 0, Rotation{}, 0, 0, nil, &sys.scrrect)
			}
		}

		kept = append(kept, entry)
	}

	s.liveOverlayRuntime = kept
	if len(kept) > 0 && os.Getenv(liveOverlaySmokeCaptureEnv) != "" && s.liveOverlayCapturePending && !s.liveOverlayCaptureDone {
		s.isTakingScreenshot = true
		s.liveOverlayCapturePending = false
		s.liveOverlayCaptureDone = true
	}
}

func effectDurationFrames(durationMs int32) int32 {
	if durationMs <= 0 {
		return 72
	}
	fps := sys.cfg.Video.Framerate
	if fps <= 0 {
		fps = 60
	}
	frames := int32(float32(durationMs) * float32(fps) / 1000.0)
	if frames < 1 {
		return 1
	}
	return frames
}
