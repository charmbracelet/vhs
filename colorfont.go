package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"

	"golang.org/x/image/font/sfnt"
)

// maxColorGlyphSize bounds how much is read for one color glyph, in case a
// font is damaged.
const maxColorGlyphSize = 4 << 20

// colorGlyphs reads the color bitmap glyphs of an emoji font: PNG images in
// the sbix table (Apple Color Emoji) or in the CBLC and CBDT tables (Noto
// Color Emoji). Emoji fonts are large (Apple's is about 180MB), so only the
// image of each glyph drawn is read from the file.
type colorGlyphs struct {
	r     io.ReaderAt
	sbix  []bitmapStrike
	cbdt  int64
	sizes []bitmapStrike
}

// bitmapStrike is a set of glyph images drawn for one size, in pixels per em.
type bitmapStrike struct {
	ppem   int
	offset int64
	// For CBLC sizes: the index subtables, and the glyphs they cover.
	subtables  int
	start, end sfnt.GlyphIndex
}

// tableRecord is where a font table is in the file.
type tableRecord struct {
	offset, length int64
}

func readBytes(r io.ReaderAt, off int64, n int) ([]byte, error) {
	if n < 0 || n > maxColorGlyphSize {
		return nil, fmt.Errorf("invalid size %d", n)
	}
	b := make([]byte, n)
	if _, err := r.ReadAt(b, off); err != nil {
		return nil, err //nolint:wrapcheck
	}
	return b, nil
}

func be16(b []byte) int { return int(binary.BigEndian.Uint16(b)) }
func be32(b []byte) int { return int(binary.BigEndian.Uint32(b)) }

// fontTables reads the table directory of the font at index in a font file
// or collection.
func fontTables(r io.ReaderAt, index int) (map[string]tableRecord, error) {
	header, err := readBytes(r, 0, 12)
	if err != nil {
		return nil, err
	}
	off := int64(0)
	if string(header[:4]) == "ttcf" {
		if index >= be32(header[8:]) {
			return nil, errors.New("no such font in the collection")
		}
		b, err := readBytes(r, 12+4*int64(index), 4)
		if err != nil {
			return nil, err
		}
		off = int64(be32(b))
		if header, err = readBytes(r, off, 12); err != nil {
			return nil, err
		}
	}
	n := be16(header[4:])
	records, err := readBytes(r, off+12, 16*n)
	if err != nil {
		return nil, err
	}
	tables := make(map[string]tableRecord, n)
	for i := range n {
		rec := records[16*i:]
		tables[string(rec[:4])] = tableRecord{offset: int64(be32(rec[8:])), length: int64(be32(rec[12:]))}
	}
	return tables, nil
}

// readColorGlyphs returns the color glyphs of the font at index, or nil if it
// has none.
func readColorGlyphs(r io.ReaderAt, index int) (*colorGlyphs, error) {
	tables, err := fontTables(r, index)
	if err != nil {
		return nil, err
	}
	c := &colorGlyphs{r: r}
	if t, ok := tables["sbix"]; ok {
		header, err := readBytes(r, t.offset, 8)
		if err != nil {
			return nil, err
		}
		offsets, err := readBytes(r, t.offset+8, 4*be32(header[4:]))
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(offsets); i += 4 {
			strike := t.offset + int64(be32(offsets[i:]))
			b, err := readBytes(r, strike, 2)
			if err != nil {
				return nil, err
			}
			c.sbix = append(c.sbix, bitmapStrike{ppem: be16(b), offset: strike})
		}
		return c, nil
	}

	cblc, ok := tables["CBLC"]
	cbdt, ok2 := tables["CBDT"]
	if !ok || !ok2 {
		return nil, nil
	}
	c.cbdt = cbdt.offset
	header, err := readBytes(r, cblc.offset, 8)
	if err != nil {
		return nil, err
	}
	sizes, err := readBytes(r, cblc.offset+8, 48*be32(header[4:]))
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(sizes); i += 48 {
		s := sizes[i:]
		c.sizes = append(c.sizes, bitmapStrike{
			offset:    cblc.offset + int64(be32(s)),
			subtables: be32(s[8:]),
			start:     gidAt(s[40:]),
			end:       gidAt(s[42:]),
			ppem:      int(s[45]),
		})
	}
	return c, nil
}

// pickStrike returns the smallest strike at least ppem pixels per em, or the
// largest one, among those ok accepts.
func pickStrike(strikes []bitmapStrike, ppem int, ok func(bitmapStrike) bool) (bitmapStrike, bool) {
	var best bitmapStrike
	found := false
	for _, s := range strikes {
		if !ok(s) {
			continue
		}
		switch {
		case !found,
			best.ppem < ppem && s.ppem > best.ppem,
			s.ppem >= ppem && s.ppem < best.ppem:
			best, found = s, true
		}
	}
	return best, found
}

// glyph returns the image of the glyph, drawn for about ppem pixels per em,
// and the pixels per em it was actually drawn for.
func (c *colorGlyphs) glyph(gid sfnt.GlyphIndex, ppem int) (image.Image, int, error) {
	if len(c.sbix) > 0 {
		return c.sbixGlyph(gid, ppem)
	}
	return c.cbdtGlyph(gid, ppem)
}

func (c *colorGlyphs) sbixGlyph(gid sfnt.GlyphIndex, ppem int) (image.Image, int, error) {
	strike, ok := pickStrike(c.sbix, ppem, func(bitmapStrike) bool { return true })
	if !ok {
		return nil, 0, errors.New("no strikes")
	}
	// A glyph may be a duplicate of another: follow a few of those.
	for range 4 {
		b, err := readBytes(c.r, strike.offset+4+4*int64(gid), 8)
		if err != nil {
			return nil, 0, err
		}
		start, end := be32(b), be32(b[4:])
		if end <= start+8 {
			return nil, 0, errors.New("no image for the glyph")
		}
		data, err := readBytes(c.r, strike.offset+int64(start), end-start)
		if err != nil {
			return nil, 0, err
		}
		switch string(data[4:8]) {
		case "png ":
			img, err := png.Decode(bytes.NewReader(data[8:]))
			return img, strike.ppem, err //nolint:wrapcheck
		case "dupe":
			gid = gidAt(data[8:])
		default:
			return nil, 0, fmt.Errorf("unsupported image type %q", data[4:8])
		}
	}
	return nil, 0, errors.New("too many duplicate glyphs")
}

func (c *colorGlyphs) cbdtGlyph(gid sfnt.GlyphIndex, ppem int) (image.Image, int, error) {
	size, ok := pickStrike(c.sizes, ppem, func(s bitmapStrike) bool {
		return gid >= s.start && gid <= s.end
	})
	if !ok {
		return nil, 0, errors.New("no image for the glyph")
	}

	// Find the index subtable for the glyph.
	array, err := readBytes(c.r, size.offset, 8*size.subtables)
	if err != nil {
		return nil, 0, err
	}
	sub, first := int64(-1), sfnt.GlyphIndex(0)
	for i := 0; i < len(array); i += 8 {
		if f, l := gidAt(array[i:]), gidAt(array[i+2:]); gid >= f && gid <= l {
			sub, first = size.offset+int64(be32(array[i+4:])), f
			break
		}
	}
	if sub < 0 {
		return nil, 0, errors.New("no image for the glyph")
	}
	header, err := readBytes(c.r, sub, 12)
	if err != nil {
		return nil, 0, err
	}
	indexFormat, imageFormat := be16(header), be16(header[2:])
	imageData := c.cbdt + int64(be32(header[4:]))
	i := int64(gid - first)

	var off int64
	var n int
	switch indexFormat {
	case 1: // 32-bit offsets
		b, err := readBytes(c.r, sub+8+4*i, 8)
		if err != nil {
			return nil, 0, err
		}
		off, n = imageData+int64(be32(b)), be32(b[4:])-be32(b)
	case 3: // 16-bit offsets
		b, err := readBytes(c.r, sub+8+2*i, 4)
		if err != nil {
			return nil, 0, err
		}
		off, n = imageData+int64(be16(b)), be16(b[2:])-be16(b)
	case 2: // images of the same size
		n = be32(header[8:])
		off = imageData + i*int64(n)
	case 4: // sparse glyphs with offsets
		count := be32(header[8:])
		pairs, err := readBytes(c.r, sub+12, 4*(count+1))
		if err != nil {
			return nil, 0, err
		}
		n = -1
		for k := 0; k < count; k++ {
			if gidAt(pairs[4*k:]) == gid {
				off, n = imageData+int64(be16(pairs[4*k+2:])), be16(pairs[4*k+6:])-be16(pairs[4*k+2:])
				break
			}
		}
	case 5: // sparse glyphs of the same size
		n = be32(header[8:])
		b, err := readBytes(c.r, sub+20, 4)
		if err != nil {
			return nil, 0, err
		}
		ids, err := readBytes(c.r, sub+24, 2*be32(b))
		if err != nil {
			return nil, 0, err
		}
		k := -1
		for j := 0; j < len(ids); j += 2 {
			if gidAt(ids[j:]) == gid {
				k = j / 2
				break
			}
		}
		if k < 0 {
			n = -1
		}
		off = imageData + int64(k)*int64(n)
	default:
		return nil, 0, fmt.Errorf("unsupported index format %d", indexFormat)
	}
	if n <= 0 {
		return nil, 0, errors.New("no image for the glyph")
	}

	data, err := readBytes(c.r, off, n)
	if err != nil {
		return nil, 0, err
	}
	// Skip the glyph metrics before the PNG data.
	var skip int
	switch imageFormat {
	case 17:
		skip = 5 + 4
	case 18:
		skip = 8 + 4
	case 19:
		skip = 4
	default:
		return nil, 0, fmt.Errorf("unsupported image format %d", imageFormat)
	}
	if len(data) < skip {
		return nil, 0, errors.New("truncated glyph image")
	}
	img, err := png.Decode(bytes.NewReader(data[skip:]))
	return img, size.ppem, err //nolint:wrapcheck
}

func gidAt(b []byte) sfnt.GlyphIndex { return sfnt.GlyphIndex(binary.BigEndian.Uint16(b)) }
