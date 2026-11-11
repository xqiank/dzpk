package game

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"strings"
)

// 牌 ID: rank*4 + suit
// rank: 0=2 1=3 2=4 3=5 4=6 5=7 6=8 7=9 8=T 9=J 10=Q 11=K 12=A
// suit: 0=club 1=diamond 2=heart 3=spade
const (
	SuitClub = 0
	SuitDiam = 1
	SuitHeart = 2
	SuitSpade = 3
)

var rankChars = [...]byte{'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A'}
var suitChars = [...]byte{'c', 'd', 'h', 's'}

// 短牌 6+ 使用 36 张：A 6 7 8 9 T J Q K
var shortDeckRanks = map[int]bool{12: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true, 10: true, 11: true}

func Rank(c int) int { return c / 4 }
func Suit(c int) int { return c % 4 }

func CardStr(c int) string {
	if c < 0 || c > 51 {
		return ""
	}
	return string([]byte{rankChars[c/4], suitChars[c%4]})
}

func CardsStr(cs []int) string {
	var b strings.Builder
	for _, c := range cs {
		b.WriteString(CardStr(c))
	}
	return b.String()
}

func NewDeck(short bool) []int {
	d := make([]int, 0, 52)
	for r := 0; r < 13; r++ {
		if short && !shortDeckRanks[r] {
			continue
		}
		for s := 0; s < 4; s++ {
			d = append(d, r*4+s)
		}
	}
	return d
}

func ParseCards(s string) []int {
	out := []int{}
	s = strings.TrimSpace(s)
	for i := 0; i+1 < len(s); i += 2 {
		rc := strings.IndexByte(string(rankChars[:]), upper(s[i]))
		sc := strings.IndexByte(string(suitChars[:]), lower(s[i+1]))
		if rc >= 0 && sc >= 0 {
			out = append(out, rc*4+sc)
		}
	}
	return out
}

func upper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 32
	}
	return b
}
func lower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}

// NewSeed 生成一手牌的随机种子
func NewSeed(entropy string) string {
	h := sha256.Sum256([]byte(entropy))
	return hex.EncodeToString(h[:])
}

func SeedHash(seed string) string {
	h := sha256.Sum256([]byte("u5-fair|" + seed))
	return hex.EncodeToString(h[:])
}

// Shuffle 基于 HMAC-SHA256 的可验证洗牌：同一 seed 必得同一副牌。
// 客户端可在结算后用公开的 seed 复算，实现「公平哈希存证」。
func Shuffle(deck []int, seed string, handNo uint64) []int {
	stream := prngStream(seed, handNo)
	n := len(deck)
	out := make([]int, n)
	copy(out, deck)
	for i := n - 1; i > 0; i-- {
		j := int(stream() % uint64(i+1))
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func prngStream(seed string, handNo uint64) func() uint64 {
	key := []byte(seed)
	var counter uint64
	var buf [8]byte
	var pool []byte
	pos := 0
	return func() uint64 {
		for {
			if pos+8 <= len(pool) {
				v := binary.BigEndian.Uint64(pool[pos : pos+8])
				pos += 8
				return v
			}
			binary.BigEndian.PutUint64(buf[:], counter)
			counter++
			m := hmac.New(sha256.New, key)
			m.Write(buf[:])
			binary.BigEndian.PutUint64(buf[:], handNo)
			m.Write(buf[:])
			pool = m.Sum(nil)
			pos = 0
		}
	}
}

// NextCard 从牌堆顶取一张（牌堆用切片模拟，末尾为顶）
func NextCard(deck *[]int) int {
	d := *deck
	if len(d) == 0 {
		return -1
	}
	c := d[len(d)-1]
	*deck = d[:len(d)-1]
	return c
}
