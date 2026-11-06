package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/bits"
	"sort"
	"strings"

	"github.com/go-text/typesetting/di"
	gtfont "github.com/go-text/typesetting/font"
	ot "github.com/go-text/typesetting/font/opentype"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// shapingTables are the font tables the shaper reads. Emoji fonts are large
// because of their images, which are not among them: Apple's has 191MB of
// images and under 500KB of everything the shaper needs.
var shapingTables = []string{
	"cmap", "head", "hhea", "hmtx", "maxp", "OS/2", "GDEF", "GSUB", "GPOS",
	"morx", "feat", "kerx", "kern", "trak",
}

// slimFont builds a font file with only the tables used for shaping from the
// font at index of a font file or collection. Glyph numbers are the same as in
// the original font, so shaped glyphs index its images.
func slimFont(r io.ReaderAt, index int) ([]byte, error) {
	tables, err := fontTables(r, index)
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, tag := range shapingTables {
		if _, ok := tables[tag]; ok {
			tags = append(tags, tag)
		}
	}
	sort.Strings(tags)

	n := len(tags)
	entrySelector := bits.Len(uint(n)) - 1
	searchRange := 16 << entrySelector //nolint:mnd
	header := make([]byte, 12+16*n)    //nolint:mnd
	binary.BigEndian.PutUint32(header, 0x00010000)
	binary.BigEndian.PutUint16(header[4:], uint16(n))                 //nolint:gosec,mnd
	binary.BigEndian.PutUint16(header[6:], uint16(searchRange))       //nolint:gosec,mnd
	binary.BigEndian.PutUint16(header[8:], uint16(entrySelector))     //nolint:gosec,mnd
	binary.BigEndian.PutUint16(header[10:], uint16(16*n-searchRange)) //nolint:gosec,mnd

	var data bytes.Buffer
	for i, tag := range tags {
		t := tables[tag]
		b, err := readBytes(r, t.offset, int(t.length))
		if err != nil {
			return nil, fmt.Errorf("could not read table %s: %w", tag, err)
		}
		rec := header[12+16*i:]
		copy(rec, tag)
		binary.BigEndian.PutUint32(rec[8:], uint32(len(header)+data.Len())) //nolint:gosec,mnd
		binary.BigEndian.PutUint32(rec[12:], uint32(len(b)))                //nolint:gosec,mnd
		data.Write(b)
		for data.Len()%4 != 0 {
			data.WriteByte(0)
		}
	}
	font := make([]byte, 0, len(header)+data.Len())
	font = append(font, header...)
	return append(font, data.Bytes()...), nil
}

// clusterShaper shapes grapheme clusters with the substitution tables of an
// emoji font, which join sequences (skin tones, ZWJ sequences, flags) into the
// single glyph that has their image.
type clusterShaper struct {
	face   *gtfont.Face
	shaper shaping.HarfbuzzShaper
	cache  map[string]sfnt.GlyphIndex
}

func newClusterShaper(r io.ReaderAt, index int) (*clusterShaper, error) {
	data, err := slimFont(r, index)
	if err != nil {
		return nil, err
	}
	loader, err := ot.NewLoader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("could not read font: %w", err)
	}
	font, err := gtfont.NewFont(loader)
	if err != nil {
		return nil, fmt.Errorf("could not read font: %w", err)
	}
	return &clusterShaper{face: gtfont.NewFace(font), cache: map[string]sfnt.GlyphIndex{}}, nil
}

// glyph returns the glyph the font draws the cluster with, or 0 if it takes
// more than one, which means it has no image for the whole sequence.
func (c *clusterShaper) glyph(cluster string) sfnt.GlyphIndex {
	if g, ok := c.cache[cluster]; ok {
		return g
	}
	g := c.shape(cluster)
	if g == 0 && strings.ContainsRune(cluster, emojiPresentation) {
		// Fonts without variation selector mappings keep the selector as a
		// glyph of its own, but have the sequence without it.
		g = c.shape(strings.ReplaceAll(cluster, string(emojiPresentation), ""))
	}
	c.cache[cluster] = g
	return g
}

// shape returns the single glyph the cluster shapes to, or 0.
func (c *clusterShaper) shape(cluster string) sfnt.GlyphIndex {
	text := []rune(cluster)
	out := c.shaper.ShapeNoExtents(shaping.Input{
		Text:      text,
		RunStart:  0,
		RunEnd:    len(text),
		Direction: di.DirectionLTR,
		Face:      c.face,
		Size:      fixed.I(16), //nolint:mnd
		Script:    language.Common,
		Language:  "en",
	})
	if len(out.Glyphs) != 1 {
		return 0
	}
	return sfnt.GlyphIndex(out.Glyphs[0].GlyphID) //nolint:gosec
}
