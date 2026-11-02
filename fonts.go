package main

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	gtfont "github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/fontscan"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
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

// monospaceFamilies are the families the generic monospace family stands for,
// in order of preference: the macOS terminal font, then the list fontconfig
// prefers, which covers Linux, then the Windows ones.
var monospaceFamilies = []string{
	"Menlo", "Noto Sans Mono", "DejaVu Sans Mono", "Liberation Mono", "Ubuntu Mono",
	"Inconsolata", "Andale Mono", "Consolas", "Cascadia Mono", "Courier New",
	"Nimbus Mono PS", "Courier", "Droid Sans Mono",
}

// fontFace is a face of a font file, at the font size.
type fontFace struct {
	font *opentype.Font
	face font.Face
	// color holds the color glyphs of emoji fonts, which are drawn instead
	// of outlines.
	color *colorGlyphs
	buf   sfnt.Buffer
}

// glyphIndex returns the glyph for r, or 0 if the font has none.
func (f *fontFace) glyphIndex(r rune) sfnt.GlyphIndex {
	gid, err := f.font.GlyphIndex(&f.buf, r)
	if err != nil {
		return 0
	}
	return gid
}

// has reports whether the face can draw r.
func (f *fontFace) has(r rune) bool {
	if f.color != nil {
		return f.glyphIndex(r) != 0
	}
	_, _, ok := f.face.GlyphBounds(r)
	return ok
}

// glyphFace is a face to draw a rune with, and how to fake a style the
// family has no face for.
type glyphFace struct {
	*fontFace
	fakeBold, fakeItalic bool
}

// family is a font family from FontFamily, with the face that best matches
// each style.
type family struct {
	best [styleCount]*fontscan.Footprint
}

// fontSet finds the system font to draw each rune with, like a browser does
// with a CSS font-family list: the first family in FontFamily that has the
// rune wins, the generic monospace family standing for the system's monospace
// fonts. Runes none of them have are drawn with the families listed after
// monospace, then the symbol fonts, then any installed font that has them.
//
// Fonts are picked from the index of the system fonts, which lists the
// characters each one has, and only the fonts drawn with are opened.
//
// It is not safe for concurrent use.
type fontSet struct {
	size     float64
	families []family
	fallback []*fontscan.Footprint
	// listed is how many fonts at the start of fallback are of a family
	// listed in FontFamily or a symbol font.
	listed int
	runes  [styleCount]map[rune]*glyphFace
	emoji  map[rune]*glyphFace
	faces  map[fontscan.Location]*fontFace
	files  map[string]*os.File
}

// fontFamilies splits a CSS font-family list into the families to draw with,
// ending with the generic monospace family, and the families only used for
// runes those lack: the ones listed after monospace and the symbol fonts.
// Families listed after monospace would never be used for letters, which
// every monospace font has, so they are only fallbacks.
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
func newFontSet(familyList string, size int) (*fontSet, error) {
	system, err := fontscan.SystemFonts(log.New(io.Discard, "", 0), fontCacheDir())
	if err != nil {
		return nil, fmt.Errorf("could not load system fonts: %w", err)
	}
	byFamily := map[string][]*fontscan.Footprint{}
	for i := range system {
		fp := &system[i]
		byFamily[fp.Family] = append(byFamily[fp.Family], fp)
	}

	fs := &fontSet{
		size:  float64(size),
		faces: map[fontscan.Location]*fontFace{},
		files: map[string]*os.File{},
	}
	names, fallback := fontFamilies(familyList)
	for _, name := range names {
		expanded := []string{name}
		if strings.EqualFold(name, fontscan.Monospace) || strings.EqualFold(name, "ui-monospace") {
			expanded = monospaceFamilies
		}
		for _, name := range expanded {
			faces := byFamily[gtfont.NormalizeFamily(name)]
			if len(faces) == 0 {
				continue
			}
			var f family
			for s := range fontStyle(styleCount) {
				f.best[s] = slices.MinFunc(faces, func(a, b *fontscan.Footprint) int {
					return cmp.Compare(aspectDistance(a.Aspect, s), aspectDistance(b.Aspect, s))
				})
			}
			fs.families = append(fs.families, f)
		}
	}
	fs.fallback, fs.listed = fallbackFonts(system, fallback)
	for s := range fs.runes {
		fs.runes[s] = map[rune]*glyphFace{}
	}

	if f := fs.face('W', styleRegular); f == nil || f.color != nil {
		_ = fs.Close()
		return nil, fmt.Errorf("no fonts found for %q", familyList)
	}
	return fs, nil
}

// aspectDistance tells how far a face is from the style asked for, lower
// being better, following CSS font matching: width first, then slant, then
// weight, where a normal weight prefers lighter faces to bolder ones and a
// bold weight bolder ones to lighter ones.
func aspectDistance(a gtfont.Aspect, style fontStyle) float64 {
	stretch := math.Abs(float64(a.Stretch) - 1)
	if a.Stretch == 0 {
		stretch = 0
	}
	slant := 0.0
	if (a.Style == gtfont.StyleItalic) != (style&styleItalic != 0) {
		slant = 1
	}
	w, weight := float64(a.Weight), 0.0
	if w == 0 {
		w = 400
	}
	if style&styleBold != 0 {
		if w >= 700 {
			weight = w - 700
		} else {
			weight = 1000 + 700 - w
		}
	} else {
		switch {
		case w >= 400 && w <= 500:
			weight = w - 400
		case w < 400:
			weight = 100 + 400 - w
		default:
			weight = 1000 + w - 500
		}
	}
	return stretch*1e7 + slant*1e5 + weight
}

// fallbackFonts orders the system fonts to search for runes the chosen fonts
// lack: the given families first, then regular faces, then the rest. It also
// returns how many fonts are of the given families.
func fallbackFonts(system []fontscan.Footprint, families []string) ([]*fontscan.Footprint, int) {
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
	fonts := make([]*fontscan.Footprint, len(system))
	for i := range system {
		fonts[i] = &system[i]
	}
	slices.SortStableFunc(fonts, func(a, b *fontscan.Footprint) int {
		return cmp.Compare(rank(a), rank(b))
	})
	listed := 0
	for listed < len(fonts) && rank(fonts[listed]) < len(families) {
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
		// Draw with the regular face, faking the style.
		if f = fs.face(r, styleRegular); f != nil {
			f = &glyphFace{
				fontFace:   f.fontFace,
				fakeBold:   style&styleBold != 0,
				fakeItalic: style&styleItalic != 0,
			}
		}
	default:
		f = fs.fallbackFace(r)
	}
	fs.runes[style][r] = f
	return f
}

// resolve returns the face of the first family in FontFamily that can draw r
// in the given style.
func (fs *fontSet) resolve(r rune, style fontStyle) *glyphFace {
	for _, fam := range fs.families {
		fp := fam.best[style]
		if !fp.Runes.Contains(r) {
			continue
		}
		face := fs.load(fp.Location)
		if face == nil || !face.has(r) {
			continue
		}
		return &glyphFace{
			fontFace: face,
			// The rasterizer draws variable fonts at their default instance,
			// so a bold instance of one still needs faking.
			fakeBold:   style&styleBold != 0 && (fp.Aspect.Weight < boldWeight || fp.Location.Instance != 0),
			fakeItalic: style&styleItalic != 0 && fp.Aspect.Style != gtfont.StyleItalic,
		}
	}
	return nil
}

func (fs *fontSet) fallbackFace(r rune) *glyphFace {
	return fs.search(r, func(*fontFace) bool { return true })
}

// emojiFace returns a color font that has r, or nil if there is none.
func (fs *fontSet) emojiFace(r rune) *glyphFace {
	if f, ok := fs.emoji[r]; ok {
		return f
	}
	f := fs.search(r, func(f *fontFace) bool { return f.color != nil })
	if fs.emoji == nil {
		fs.emoji = map[rune]*glyphFace{}
	}
	fs.emoji[r] = f
	return f
}

// search returns the first of the fallback fonts that has r and that ok
// accepts.
func (fs *fontSet) search(r rune, ok func(*fontFace) bool) *glyphFace {
	fonts := fs.fallback
	if unicode.Is(unicode.Co, r) {
		// Private use characters, like Nerd Font icons, mean something
		// different in every font, so only listed fonts may draw them.
		fonts = fonts[:fs.listed]
	}
	for _, fp := range fonts {
		if !fp.Runes.Contains(r) {
			continue
		}
		if face := fs.load(fp.Location); face != nil && ok(face) && face.has(r) {
			return &glyphFace{fontFace: face}
		}
	}
	return nil
}

// load returns the face of a font, or nil if it can't be drawn. The font file
// is read as glyphs are needed rather than loaded into memory.
func (fs *fontSet) load(loc fontscan.Location) *fontFace {
	// The rasterizer doesn't support font variations, so every instance of a
	// variable font is the same face.
	loc.Instance = 0
	if f, ok := fs.faces[loc]; ok {
		return f
	}
	f := fs.open(loc)
	fs.faces[loc] = f
	return f
}

func (fs *fontSet) open(loc fontscan.Location) *fontFace {
	file, ok := fs.files[loc.File]
	if !ok {
		file, _ = os.Open(loc.File)
		fs.files[loc.File] = file
	}
	if file == nil {
		return nil
	}
	coll, err := opentype.ParseCollectionReaderAt(file)
	if err != nil {
		return nil
	}
	ft, err := coll.Font(int(loc.Index))
	if err != nil {
		return nil
	}
	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    fs.size,
		DPI:     72, //nolint:mnd // one point per pixel, like CSS px
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil
	}
	color, err := readColorGlyphs(file, int(loc.Index))
	if err != nil {
		color = nil
	}
	return &fontFace{font: ft, face: face, color: color}
}

// Close closes the font files.
func (fs *fontSet) Close() error {
	var errs []error
	for _, f := range fs.files {
		if f != nil {
			errs = append(errs, f.Close())
		}
	}
	clear(fs.files)
	clear(fs.faces)
	for _, runes := range fs.runes {
		clear(runes)
	}
	clear(fs.emoji)
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("could not close fonts: %w", err)
	}
	return nil
}
