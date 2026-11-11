package game

// 手牌评估：52 张标准、36 张短牌(6+)、奥马哈(必须用 2 张手牌 + 3 张公共牌)
//
// 标准顺序: 高牌 < 一对 < 两对 < 三条 < 顺子 < 同花 < 葫芦 < 四条 < 同花顺
// 短牌(6+)顺序: 高牌 < 一对 < 两对 < 顺子 < 三条 < 葫芦 < 同花 < 四条 < 同花顺
//   —— 短牌里「三条 > 顺子」「同花 > 葫芦」，A6789 为最小顺子

const (
	CatHighCard = 0
	CatPair     = 1
	CatTwoPair  = 2
	CatTrips    = 3
	CatStraight = 4
	CatFlush    = 5
	CatFullHouse = 6
	CatQuads    = 7
	CatStraightFlush = 8
)

var catNames = [...]string{"高牌", "一对", "两对", "三条", "顺子", "同花", "葫芦", "四条", "同花顺"}

// 短牌类别重排映射：stdCat -> 短牌比较位
var shortCatRank = [9]int{0, 1, 2, 4, 3, 6, 5, 7, 8}

type HandValue struct {
	Score int64
	Cat   int
	CatName string
	Best  []int
}

func (h HandValue) Empty() bool { return h.Score == 0 }

// pack 把类别和 5 个踢脚压进 int64，便于直接比较大小
func pack(cat int, t [5]int) int64 {
	var v int64 = int64(cat)
	for i := 0; i < 5; i++ {
		v = v<<4 | int64(t[i]&0xF)
	}
	return v
}

func variantScore(cat int, t [5]int, short bool) int64 {
	if short {
		cat = shortCatRank[cat]
	}
	return pack(cat, t)
}

// eval5 返回标准类别与踢脚
func eval5(c [5]int, short bool) (int, [5]int) {
	var cnt [13]int
	var suitCnt [4]int
	for _, x := range c {
		cnt[Rank(x)]++
		suitCnt[Suit(x)]++
	}
	flush := false
	for s := 0; s < 4; s++ {
		if suitCnt[s] == 5 {
			flush = true
		}
	}

	// 顺子检测
	straightHigh := -1
	if !short {
		for hi := 12; hi >= 4; hi-- {
			if cnt[hi] > 0 && cnt[hi-1] > 0 && cnt[hi-2] > 0 && cnt[hi-3] > 0 && cnt[hi-4] > 0 {
				straightHigh = hi
				break
			}
		}
		// wheel A2345 -> 5 high
		if straightHigh < 0 && cnt[12] > 0 && cnt[0] > 0 && cnt[1] > 0 && cnt[2] > 0 && cnt[3] > 0 {
			straightHigh = 3
		}
	} else {
		for hi := 12; hi >= 8; hi-- {
			ok := true
			for k := 0; k < 5; k++ {
				if cnt[hi-k] == 0 {
					ok = false
					break
				}
			}
			if ok {
				straightHigh = hi
				break
			}
		}
		// 短牌 A6789 -> 视为比任何真实顺子都小
		if straightHigh < 0 && cnt[12] > 0 && cnt[4] > 0 && cnt[5] > 0 && cnt[6] > 0 && cnt[7] > 0 {
			straightHigh = 3
		}
	}

	// 四条
	for r := 12; r >= 0; r-- {
		if cnt[r] == 4 {
			k := 0
			for r2 := 12; r2 >= 0; r2-- {
				if r2 != r && cnt[r2] > 0 {
					k = r2
					break
				}
			}
			return CatQuads, [5]int{r, k, 0, 0, 0}
		}
	}
	// 葫芦
	t3, t2 := -1, -1
	for r := 12; r >= 0; r-- {
		if cnt[r] >= 3 && t3 < 0 {
			t3 = r
		}
	}
	if t3 >= 0 {
		for r := 12; r >= 0; r-- {
			if r != t3 && cnt[r] >= 2 {
				t2 = r
				break
			}
		}
	}
	if t3 >= 0 && t2 >= 0 {
		return CatFullHouse, [5]int{t3, t2, 0, 0, 0}
	}
	if flush {
		var t [5]int
		i := 0
		for r := 12; r >= 0 && i < 5; r-- {
			if cnt[r] > 0 {
				t[i] = r
				i++
			}
		}
		if straightHigh >= 0 {
			return CatStraightFlush, [5]int{straightHigh, 0, 0, 0, 0}
		}
		return CatFlush, t
	}
	if straightHigh >= 0 {
		return CatStraight, [5]int{straightHigh, 0, 0, 0, 0}
	}
	if t3 >= 0 {
		var t [5]int
		t[0] = t3
		i := 1
		for r := 12; r >= 0 && i < 5; r-- {
			if r != t3 && cnt[r] > 0 {
				t[i] = r
				i++
			}
		}
		return CatTrips, t
	}
	// 两对
	p1, p2 := -1, -1
	for r := 12; r >= 0; r-- {
		if cnt[r] >= 2 {
			if p1 < 0 {
				p1 = r
			} else if p2 < 0 {
				p2 = r
			}
		}
	}
	if p1 >= 0 && p2 >= 0 {
		k := 0
		for r := 12; r >= 0; r-- {
			if r != p1 && r != p2 && cnt[r] > 0 {
				k = r
				break
			}
		}
		return CatTwoPair, [5]int{p1, p2, k, 0, 0}
	}
	if p1 >= 0 {
		var t [5]int
		t[0] = p1
		i := 1
		for r := 12; r >= 0 && i < 5; r-- {
			if r != p1 && cnt[r] > 0 {
				t[i] = r
				i++
			}
		}
		return CatPair, t
	}
	var t [5]int
	i := 0
	for r := 12; r >= 0 && i < 5; r-- {
		if cnt[r] > 0 {
			t[i] = r
			i++
		}
	}
	return CatHighCard, t
}

// EvaluateHoldem 7 张选 5（也支持 5/6 张）
func EvaluateHoldem(cards []int, short bool) HandValue {
	n := len(cards)
	if n < 5 {
		return HandValue{}
	}
	var bestScore int64 = -1
	var bestCat int
	var bestCards []int
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	var rec func(start int, picked []int)
	rec = func(start int, picked []int) {
		if len(picked) == 5 {
			var cc [5]int
			for i, p := range picked {
				cc[i] = cards[p]
			}
			cat, t := eval5(cc, short)
			s := variantScore(cat, t, short)
			if s > bestScore {
				bestScore = s
				bestCat = cat
				bestCards = append([]int{}, cc[:]...)
			}
			return
		}
		for i := start; i < n; i++ {
			rec(i+1, append(picked, i))
		}
	}
	rec(0, []int{})
	return HandValue{Score: bestScore, Cat: bestCat, CatName: catNames[bestCat], Best: bestCards}
}

// EvaluateOmaha 必须恰好 2 张手牌 + 3 张公共牌
func EvaluateOmaha(hole []int, board []int, short bool) HandValue {
	if len(hole) < 4 || len(board) < 3 {
		return HandValue{}
	}
	var bestScore int64 = -1
	var bestCat int
	var best []int
	hole = hole[:4]
	if len(board) > 5 {
		board = board[:5]
	}
	for a := 0; a < 4; a++ {
		for b := a + 1; b < 4; b++ {
			for x := 0; x < len(board); x++ {
				for y := x + 1; y < len(board); y++ {
					for z := y + 1; z < len(board); z++ {
						cs := [5]int{hole[a], hole[b], board[x], board[y], board[z]}
						cat, t := eval5(cs, short)
						s := variantScore(cat, t, short)
						if s > bestScore {
							bestScore = s
							bestCat = cat
							best = cs[:]
						}
					}
				}
			}
		}
	}
	return HandValue{Score: bestScore, Cat: bestCat, CatName: catNames[bestCat], Best: best}
}

// Evaluate 统一入口
func Evaluate(gameType string, hole []int, board []int) HandValue {
	switch gameType {
	case "omaha":
		return EvaluateOmaha(hole, board, false)
	case "short6":
		return EvaluateHoldem(append(append([]int{}, hole...), board...), true)
	default:
		return EvaluateHoldem(append(append([]int{}, hole...), board...), false)
	}
}

// Equity 蒙特卡洛胜率（用于机器人决策与教学计算器）
func Equity(gameType string, hole []int, board []int, opponents int, iters int, short bool) float64 {
	if iters <= 0 {
		iters = 2000
	}
	if opponents < 1 {
		opponents = 1
	}
	base := NewDeck(short)
	used := map[int]bool{}
	for _, c := range append(append([]int{}, hole...), board...) {
		used[c] = true
	}
	rest := make([]int, 0, len(base))
	for _, c := range base {
		if !used[c] {
			rest = append(rest, c)
		}
	}
	wins := 0.0
	rnd := newFastRand(uint64(len(rest))*7919 + uint64(len(hole))*104729)
	buf := make([]int, len(rest))
	copy(buf, rest)
	needBoard := 5 - len(board)
	for it := 0; it < iters; it++ {
		// Fisher-Yates 局部洗牌
		for i := 0; i < 8+needBoard+4*opponents && i < len(buf); i++ {
			j := int(rnd.next() % uint64(len(buf)))
			buf[i], buf[j] = buf[j], buf[i]
		}
		idx := 0
		newBoard := append([]int{}, board...)
		for k := 0; k < needBoard; k++ {
			newBoard = append(newBoard, buf[idx])
			idx++
		}
		my := Evaluate(gameType, hole, newBoard)
		best := true
		tie := false
		for o := 0; o < opponents; o++ {
			var oh []int
			if gameType == "omaha" {
				oh = buf[idx : idx+4]
				idx += 4
			} else {
				oh = buf[idx : idx+2]
				idx += 2
			}
			ov := Evaluate(gameType, oh, newBoard)
			if ov.Score > my.Score {
				best = false
				break
			}
			if ov.Score == my.Score {
				tie = true
			}
		}
		if best {
			if tie {
				wins += 0.5
			} else {
				wins += 1
			}
		}
	}
	return wins / float64(iters)
}

// 轻量 xorshift，避免 math/rand 全局锁
type fastRand struct{ s uint64 }

func newFastRand(seed uint64) *fastRand {
	if seed == 0 {
		seed = 88172645463325252
	}
	return &fastRand{s: seed}
}
func (f *fastRand) next() uint64 {
	x := f.s
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	f.s = x
	return x
}

// EquityArgs 指定对手确切手牌的胜率计算
type EquityArgs struct {
	GameType string
	Hole     []int
	Board    []int
	OppHole  []int
	Iters    int
}

// EquityVs 对手手牌已知时的胜率（蒙特卡洛剩余公共牌）
func EquityVs(a EquityArgs) float64 {
	iters := a.Iters
	if iters <= 0 {
		iters = 5000
	}
	short := a.GameType == "short6"
	deck := NewDeck(short)
	used := map[int]bool{}
	for _, c := range a.Hole {
		used[c] = true
	}
	for _, c := range a.Board {
		used[c] = true
	}
	for _, c := range a.OppHole {
		used[c] = true
	}
	rest := make([]int, 0, len(deck))
	for _, c := range deck {
		if !used[c] {
			rest = append(rest, c)
		}
	}
	need := 5 - len(a.Board)
	if need < 0 {
		need = 0
	}
	wins := 0.0
	rnd := newFastRand(uint64(len(rest))*7919 + 33)
	for it := 0; it < iters; it++ {
		for i := 0; i < need+3 && i < len(rest); i++ {
			j := int(rnd.next() % uint64(len(rest)))
			rest[i], rest[j] = rest[j], rest[i]
		}
		board := append([]int{}, a.Board...)
		board = append(board, rest[:need]...)
		me := Evaluate(a.GameType, a.Hole, board)
		op := Evaluate(a.GameType, a.OppHole, board)
		if me.Score > op.Score {
			wins++
		} else if me.Score == op.Score {
			wins += 0.5
		}
	}
	return wins / float64(iters)
}
