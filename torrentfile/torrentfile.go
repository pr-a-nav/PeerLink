// Package torrentfile parses torrent files or magnet links.
package torrentfile

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/zeebo/bencode"
)

// Parse and return a TorrentFile that is a useful shape for actual downloading

// TorrentFile represents the contents of a .torrent file, reformatted for ease
// of use in the download process.
//
// The InfoHash is generated via SHA-1 from the entire Info field of the file.
//
// The 20-byte SHA1 hashes are formatted into a slice of 20-byte arrays for easy
// comparison with pieces downloaded from a peer.
type TorrentFile struct {
	TrackerURLs []string   // tracker URLs (from announce-list, announce or magnet link tr)
	InfoHash    [20]byte   // SHA-1 hash of the entire info field, uniquely identifies the torrent
	PieceHashes [][20]byte // SHA-1 hashes of each file piece
	PieceLength int        // number of bytes per piece
	Files       []File     // in the 1 file case, this will only have one element
	TotalLength int        // calculated as the sum of all files
	DisplayName string     // human readable display name (.torrent filename or magnet link dn)
}

// File contains metadata about the final downloaded files, namely their path
// and file length.
type File struct {
	Length   int    // length in bytes
	FullPath string // download path
	SHA1Hash string // optional for final validation
	MD5Hash  string // optional for final validation
}

// New returns a new TorrentFile.
//
// If the source is a .torrent file, it will be ready for use.
//
// If the source is a magnet link, metadata will need to be acquired from peers
// already in the swarm, then added using TorrentFile.AppendMetadata()
func New(source string) (TorrentFile, error) {
	if strings.HasSuffix(source, ".torrent") {
		return parseTorrentFile(source)
	}
	if strings.HasPrefix(source, "magnet") {
		return parseMagnetLink(source)
	}
	return TorrentFile{}, fmt.Errorf("invalid source (torrent file and magnet links supported)")
}

// parseTorrentFile parses a raw .torrent file.
func parseTorrentFile(filename string) (TorrentFile, error) {
	filename = os.ExpandEnv(filename)
	f, err := os.Open(filename)
	if err != nil {
		return TorrentFile{}, err
	}

	var btor bencodeTorrent
	err = bencode.NewDecoder(f).Decode(&btor)
	if err != nil {
		return TorrentFile{}, fmt.Errorf("unmarshalling file: %w", err)
	}

	var trackerURLs []string
	for _, list := range btor.AnnounceList {
		trackerURLs = append(trackerURLs, list...)
	}
	// BEP0012, only use `announce` if `announce-list` is not present
	if len(trackerURLs) == 0 {
		trackerURLs = append(trackerURLs, btor.Announce)
	}
	tf := TorrentFile{
		TrackerURLs: trackerURLs,
		DisplayName: filename,
	}

	err = tf.AppendMetadata(btor.Info)
	if err != nil {
		return TorrentFile{}, fmt.Errorf("parsing metadata: %w", err)
	}

	return tf, nil
}

// parseMagnetLink parses a magnet URI and returns a TorrentFile with as much info as possible.
func parseMagnetLink(uri string) (TorrentFile, error) {
	if !strings.HasPrefix(uri, "magnet:") {
		return TorrentFile{}, fmt.Errorf("not a magnet link")
	}
	// Remove the magnet:? prefix
	params := strings.TrimPrefix(uri, "magnet:?")
	parts := strings.Split(params, "&")

	var tf TorrentFile
	for _, part := range parts {
		if strings.HasPrefix(part, "xt=") {
			// xt=urn:btih:<infohash>
			xtVal := strings.TrimPrefix(part, "xt=")
			if strings.HasPrefix(xtVal, "urn:btih:") {
				hash := strings.TrimPrefix(xtVal, "urn:btih:")
				// infohash can be hex (40 chars) or base32 (32 chars)
				if len(hash) == 40 {
					// hex
					decoded, err := decodeHexString(hash)
					if err != nil {
						return TorrentFile{}, fmt.Errorf("invalid hex infohash: %w", err)
					}
					copy(tf.InfoHash[:], decoded)
				} else if len(hash) == 32 {
					// base32
					decoded, err := decodeBase32String(hash)
					if err != nil {
						return TorrentFile{}, fmt.Errorf("invalid base32 infohash: %w", err)
					}
					copy(tf.InfoHash[:], decoded)
				} else {
					return TorrentFile{}, fmt.Errorf("unexpected infohash length in magnet link")
				}
			}
		} else if strings.HasPrefix(part, "tr=") {
			tracker := strings.TrimPrefix(part, "tr=")
			tracker, _ = url.QueryUnescape(tracker)
			tf.TrackerURLs = append(tf.TrackerURLs, tracker)
		} else if strings.HasPrefix(part, "dn=") {
			dn := strings.TrimPrefix(part, "dn=")
			dn, _ = url.QueryUnescape(dn)
			tf.DisplayName = dn
		}
	}
	if tf.DisplayName == "" {
		tf.DisplayName = uri
	}
	if tf.InfoHash == [20]byte{} {
		return TorrentFile{}, fmt.Errorf("magnet link missing infohash")
	}
	return tf, nil
}

// decodeHexString decodes a 40-character hex string to 20 bytes.
func decodeHexString(s string) ([]byte, error) {
	if len(s) != 40 {
		return nil, fmt.Errorf("hex string must be 40 chars")
	}
	b := make([]byte, 20)
	for i := 0; i < 20; i++ {
		var v byte
		_, err := fmt.Sscanf(s[i*2:i*2+2], "%02x", &v)
		if err != nil {
			return nil, err
		}
		b[i] = v
	}
	return b, nil
}

// decodeBase32String decodes a 32-character base32 string to 20 bytes.
func decodeBase32String(s string) ([]byte, error) {
	if len(s) != 32 {
		return nil, fmt.Errorf("base32 string must be 32 chars")
	}
	// Use standard base32 decoder, upper-case only
	enc := "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	var out [20]byte
	var bits, bitsLeft, count int
	for _, c := range s {
		idx := strings.IndexRune(enc, unicode.ToUpper(c))
		if idx == -1 {
			return nil, fmt.Errorf("invalid base32 char: %c", c)
		}
		bits = (bits << 5) | idx
		bitsLeft += 5
		if bitsLeft >= 8 {
			bitsLeft -= 8
			if count < 20 {
				out[count] = byte(bits >> bitsLeft)
				count++
			}
			bits &= (1 << bitsLeft) - 1
		}
	}
	if count != 20 {
		return nil, fmt.Errorf("base32 decode did not yield 20 bytes")
	}
	return out[:], nil
}

// serialization struct the represents the structure of a .torrent file
// it is not immediately usable, so it can be converted to a TorrentFile struct
type bencodeTorrent struct {
	// URL of tracker server to get peers from
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	// Info is parsed as a RawMessage to ensure that the final info_hash is
	// correct even in the case of the info dictionary being an unexpected shape
	Info bencode.RawMessage `bencode:"info"`
}

// Only Length OR Files will be present per BEP0003
// spec: http://bittorrent.org/beps/bep_0003.html#info-dictionary
type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`       // binary blob of all SHA1 hash of each piece
	PieceLength int    `bencode:"piece length"` // length in bytes of each piece
	Name        string `bencode:"name"`         // Name of file (or folder if there are multiple files)
	Length      int    `bencode:"length"`       // total length of file (in single file case)
	Files       []struct {
		Length   int      `bencode:"length"` // length of this file
		Path     []string `bencode:"path"`   // list of subdirectories, last element is file name
		SHA1Hash string   `bencode:"sha1"`   // optional, to validate this file
		MD5Hash  string   `bencode:"md5"`    // optional, to validate this file
	} `bencode:"files"`
}

// AppendMetadata adds the metadata (aka the info dictionary of a torrent file).
// It must be called after torrentfile.New() is invoked with a magnet link
// source with the metadata acquired from a peer in the swarm.
func (t *TorrentFile) AppendMetadata(metadata []byte) error {
	var info bencodeInfo
	err := bencode.DecodeBytes(metadata, &info)
	if err != nil {
		return fmt.Errorf("unmarshalling info dict: %w", err)
	}

	// SHA-1 hash the entire info dictionary to get the info_hash
	t.InfoHash = sha1.Sum(metadata)

	// split the Pieces blob into the 20-byte SHA-1 hashes for comparison later
	const hashLen = 20 // length of a SHA-1 hash
	if len(info.Pieces)%hashLen != 0 {
		return errors.New("invalid length for info pieces")
	}
	t.PieceHashes = make([][20]byte, len(info.Pieces)/hashLen)
	for i := 0; i < len(t.PieceHashes); i++ {
		piece := info.Pieces[i*hashLen : (i+1)*hashLen]
		copy(t.PieceHashes[i][:], piece)
	}

	t.PieceLength = info.PieceLength

	// either Length OR Files field must be present (but not both)
	if info.Length == 0 && len(info.Files) == 0 {
		return fmt.Errorf("invalid torrent file info dict: no length OR files")
	}

	if info.Length != 0 {
		t.Files = append(t.Files, File{
			Length:   info.Length,
			FullPath: info.Name,
		})
		t.TotalLength = info.Length
	} else {
		for _, f := range info.Files {
			subPaths := append([]string{info.Name}, f.Path...)
			t.Files = append(t.Files, File{
				Length:   f.Length,
				FullPath: filepath.Join(subPaths...),
				SHA1Hash: f.SHA1Hash,
				MD5Hash:  f.MD5Hash,
			})
			t.TotalLength += f.Length
		}
	}

	return nil
}
