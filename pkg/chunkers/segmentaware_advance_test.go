package chunkers

import (
	"bytes"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSegmentAwareAdvance(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	// Mostly small segments with some big ones that are chunked on their own
	var segments []uint64
	var size uint64
	for range 300 {
		s := uint64(1000 + rng.Intn(60000))
		if rng.Intn(10) == 0 {
			s = uint64(300000 + rng.Intn(300000))
		}
		segments = append(segments, s)
		size += s
	}
	data := make([]byte, size)
	rng.Read(data)
	dir := t.TempDir()

	for _, merge := range []bool{true, false} {
		t.Run("merge="+strconv.FormatBool(merge), func(t *testing.T) {
			options := ""
			if !merge {
				options = ",mergeSmallSegments=false"
			}
			// Returns a chunker that starts at offset, as if the input began there
			newChunker := func(offset uint64) Chunker {
				name := filepath.Join(dir, strconv.FormatUint(offset, 10)+".csv")
				require.NoError(t, os.WriteFile(name, segmentsCSVFrom(segments, offset), 0644))
				params := DefaultChunkerParams()
				params.Options = "segmentsSource=" + name + options
				c, err := NewChunker("segmentaware", bytes.NewReader(data[offset:]), params)
				require.NoError(t, err)
				return c
			}
			breakpoint := 2 * DefaultChunkSizeAvg
			if !merge {
				breakpoint = 0
			}

			c := newChunker(0)
			var pos uint64
			var expected []byte // first chunk of a fresh chunker after an advance
			for pos < size {
				if rng.Intn(3) == 0 {
					n := min(uint64(1+rng.Intn(400000)), size-pos)
					require.NoError(t, c.Advance(int(n)))
					pos += n
					expected = nil
					if !insideBigSegment(segments, pos, breakpoint) && pos < size {
						_, b, err := newChunker(pos).Next()
						require.NoError(t, err)
						expected = append([]byte(nil), b...)
					}
					continue
				}
				start, b, err := c.Next()
				require.NoError(t, err)
				require.NotEmpty(t, b, "unexpected end at %d of %d", pos, size)
				require.Equal(t, pos, start)
				require.Equal(t, data[pos:pos+uint64(len(b))], b)
				if expected != nil {
					require.Equal(t, expected, b, "chunk after advancing to %d differs from a fresh chunker", pos)
					expected = nil
				}
				pos += uint64(len(b))
			}
			_, b, err := c.Next()
			require.NoError(t, err)
			require.Empty(t, b)

			// Can't advance past the end
			require.Error(t, newChunker(0).Advance(int(size)+1))
		})
	}
}

// Returns true if offset is inside, not at the start of, a segment that is
// chunked on its own.
func insideBigSegment(segments []uint64, offset, breakpoint uint64) bool {
	for _, s := range segments {
		if offset < s {
			return offset > 0 && s > breakpoint
		}
		offset -= s
	}
	return false
}

// Returns the segments after offset as CSV, shortening the one offset falls into.
func segmentsCSVFrom(segments []uint64, offset uint64) []byte {
	var csv strings.Builder
	for _, s := range segments {
		if offset >= s {
			offset -= s
			continue
		}
		csv.WriteString(strconv.FormatUint(s-offset, 10) + "\n")
		offset = 0
	}
	return []byte(csv.String())
}
