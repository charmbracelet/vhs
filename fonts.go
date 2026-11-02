package main

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	gtfont "github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/fontscan"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// fontStyle selects one of the four faces of a font family. It is a bit set,
// so bold italic is styleBold|styleItalic.
type fontStyle uint8

const (
	styleRegular fontStyle = 0
	styleBold    fontStyle = 1
	styleItalic  fontStyle = 2
	styleCount             = 4
)

// boldWeight is the weight from which a face counts as bold.
const boldWeight = 600

// glyphFace is a face to draw a rune with, and whether bold has to be faked
// because the family has no bold face.
type glyphFace struct {
	face     font.Face
	fakeBold bool
}

// fontSet finds the system font to draw each rune with, like a browser does
// with a CSS font-family list: the first family in FontFamily that has the
// rune wins, the generic monospace family standing for the system's monospace
// fonts. Runes none of them have are drawn with the families listed after
// monospace, then the symbol fonts, then any installed font that has them.
//
// It is not safe for concurrent use.
type fontSet struct {
	size     float64
	maps     [styleCount]*fontscan.FontMap
	fallback []fontscan.Footprint
	// listed is how many fonts at the start of fallback are of a family
	// listed in FontFamily or a symbol font.
	listed int
	runes  [styleCount]map[rune]*glyphFace
	faces  map[fontscan.Location]font.Face
	fonts  map[string]*opentype.Collection
}

// fontFamilies splits a CSS font-family list into the families to draw with,
// ending with the generic monospace family, and the families only used for
// runes those lack: the ones listed after monospace and the symbol fonts.
// fontscan prefers any installed family listed over the generic one, so
// these must be kept out of its query, or letters would be drawn with them.
func fontFamilies(family string) (families, fallback []string) {
	isMonospace := func(f string) bool { return strings.EqualFold(f, fontscan.Monospace) }
	for f := range strings.SplitSeq(family, fontsSeparator) {
		f = strings.Trim(strings.TrimSpace(f), `"'`)
		switch {
		case f == "":
		case slices.ContainsFunc(families, isMonospace) || slices.Contains(symbolsFallback, f):
			fallback = append(fallback, f)
		default:
			families = append(families, f)
		}
	}
	if !slices.ContainsFunc(families, isMonospace) {
		families = append(families, fontscan.Monospace)
	}
	return families, append(fallback, symbolsFallback...)
}

// fontCacheDir returns where the index of the system fonts is kept.
func fontCacheDir() string {
	dir, err := os.UserCacheDir()
	if err == nil {
		dir = filepath.Join(dir, "vhs", "fonts")
		err = os.MkdirAll(dir, 0o700)
	}
	if err != nil {
		return filepath.Join(os.TempDir(), "vhs-fonts")
	}
	return dir
}

// newFontSet indexes the system fonts and prepares lookups for the given
// font-family list and size in pixels.
func newFontSet(family string, size int) (*fontSet, error) {
	cacheDir := fontCacheDir()
	logger := log.New(io.Discard, "", 0)
	system, err := fontscan.SystemFonts(logger, cacheDir)
	if err != nil {
		return nil, fmt.Errorf("could not load system fonts: %w", err)
	}

	families, fallback := fontFamilies(family)
	fs := &fontSet{
		size:  float64(size),
		faces: map[fontscan.Location]font.Face{},
		fonts: map[string]*opentype.Collection{},
	}
	fs.fallback, fs.listed = fallbackFonts(system, fallback)
	for s := range fontStyle(styleCount) {
		fm := fontscan.NewFontMap(logger)
		if err := fm.UseSystemFonts(cacheDir); err != nil {
			return nil, fmt.Errorf("could not load system fonts: %w", err)
		}
		aspect := gtfont.Aspect{Style: gtfont.StyleNormal, Weight: gtfont.WeightNormal}
		if s&styleBold != 0 {
			aspect.Weight = gtfont.WeightBold
		}
		if s&styleItalic != 0 {
			aspect.Style = gtfont.StyleItalic
		}
		aspect.SetDefaults()
		fm.SetQuery(fontscan.Query{Families: families, Aspect: aspect})
		fs.maps[s] = fm
		fs.runes[s] = map[rune]*glyphFace{}
	}

	if fs.face('W', styleRegular) == nil {
		return nil, fmt.Errorf("no fonts found for %q", family)
	}
	return fs, nil
}

// fallbackFonts orders the system fonts to search for runes the chosen fonts
// lack: the given families first, then regular faces, then the rest. It also
// returns how many fonts are of the given families.
func fallbackFonts(system []fontscan.Footprint, families []string) ([]fontscan.Footprint, int) {
	rank := func(fp *fontscan.Footprint) int {
		for i, f := range families {
			if fp.Family == gtfont.NormalizeFamily(f) {
				return i
			}
		}
		if fp.Aspect.Weight == gtfont.WeightNormal && fp.Aspect.Style == gtfont.StyleNormal {
			return len(families)
		}
		return len(families) + 1
	}
	fonts := slices.Clone(system)
	slices.SortStableFunc(fonts, func(a, b fontscan.Footprint) int {
		return cmp.Compare(rank(&a), rank(&b))
	})
	listed := 0
	for listed < len(fonts) && rank(&fonts[listed]) < len(families) {
		listed++
	}
	return fonts, listed
}

// face returns the face to draw r with in the given style, or nil if no
// installed font can draw it.
func (fs *fontSet) face(r rune, style fontStyle) *glyphFace {
	if f, ok := fs.runes[style][r]; ok {
		return f
	}
	f := fs.resolve(r, style)
	switch {
	case f != nil:
	case style != styleRegular:
		// Draw with the regular face, faking bold if needed.
		if f = fs.face(r, styleRegular); f != nil && style&styleBold != 0 {
			f = &glyphFace{face: f.face, fakeBold: true}
		}
	default:
		f = fs.fallbackFace(r)
	}
	fs.runes[style][r] = f
	return f
}

func (fs *fontSet) resolve(r rune, style fontStyle) *glyphFace {
	fm := fs.maps[style]
	gf := fm.ResolveFace(r)
	if gf == nil {
		return nil
	}
	loc := fm.FontLocation(gf.Font)
	face := fs.load(loc)
	if !hasGlyph(face, r) {
		return nil
	}
	f := &glyphFace{face: face}
	if style&styleBold != 0 {
		_, aspect := fm.FontMetadata(gf.Font)
		// The rasterizer draws variable fonts at their default instance, so
		// a bold instance of one still needs faking.
		f.fakeBold = aspect.Weight < boldWeight || loc.Instance != 0
	}
	return f
}

func (fs *fontSet) fallbackFace(r rune) *glyphFace {
	fonts := fs.fallback
	if unicode.Is(unicode.Co, r) {
		// Private use characters, like Nerd Font icons, mean something
		// different in every font, so only listed fonts may draw them.
		fonts = fonts[:fs.listed]
	}
	for i := range fonts {
		fp := &fonts[i]
		if !fp.Runes.Contains(r) {
			continue
		}
		if face := fs.load(fp.Location); hasGlyph(face, r) {
			return &glyphFace{face: face}
		}
	}
	return nil
}

// hasGlyph reports whether face can draw r. Bitmap-only glyphs, like those of
// color emoji fonts, can't be drawn.
func hasGlyph(face font.Face, r rune) bool {
	if face == nil {
		return false
	}
	_, _, ok := face.GlyphBounds(r)
	return ok
}

// load returns the face of a font, or nil for formats the rasterizer can't
// draw.
func (fs *fontSet) load(loc fontscan.Location) font.Face {
	// The rasterizer doesn't support font variations, so every instance of a
	// variable font is the same face.
	loc.Instance = 0
	if face, ok := fs.faces[loc]; ok {
		return face
	}
	var face font.Face
	if f := fs.collectionFont(loc); f != nil {
		face, _ = opentype.NewFace(f, &opentype.FaceOptions{
			Size:    fs.size,
			DPI:     72, //nolint:mnd // one point per pixel, like CSS px
			Hinting: font.HintingNone,
		})
	}
	fs.faces[loc] = face
	return face
}

func (fs *fontSet) collectionFont(loc fontscan.Location) *opentype.Font {
	coll, ok := fs.fonts[loc.File]
	if !ok {
		data, err := os.ReadFile(loc.File)
		if err == nil {
			coll, err = opentype.ParseCollection(data)
		}
		if err != nil {
			coll = nil
		}
		fs.fonts[loc.File] = coll
	}
	if coll == nil {
		return nil
	}
	f, err := coll.Font(int(loc.Index))
	if err != nil {
		return nil
	}
	return f
}
