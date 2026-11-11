package game

// 牛牛（经典牛牛）：5 张牌，任意 3 张点数和为 10 的倍数即为「有牛」，
// 余下 2 张点数和 mod 10 即牛几（0 为牛牛）。
// 牌型优先级：五小牛 > 炸弹 > 五花牛 > 牛牛 > 牛九 … 牛一 > 无牛
// 倍数：牛1-6 / 无牛 ×1；牛7 ×2；牛8 ×3；牛9 ×4；牛牛 ×5；五花牛 ×5；炸弹 ×6；五小牛 ×8

const (
	NiuNone      = 0
	NiuOne       = 1
	NiuTwo       = 2
	NiuThree     = 3
	NiuFour      = 4
	NiuFive      = 5
	NiuSix       = 6
	NiuSeven     = 7
	NiuEight     = 8
	NiuNine      = 9
	NiuTen       = 10 // 牛牛
	NiuFiveFlower = 11
	NiuBomb      = 12
	NiuFiveSmall = 13
)

var niuNames = [...]string{
	"无牛", "牛一", "牛二", "牛三", "牛四", "牛五", "牛六", "牛七", "牛八", "牛九", "牛牛",
	"五花牛", "炸弹", "五小牛",
}

var niuMults = [...]int{1, 1, 1, 1, 1, 1, 1, 2, 3, 4, 5, 5, 6, 8}

type NiuResult struct {
	Cat    int
	Name   string
	Mult   int
	Score  int64
	NiuNum int // 有牛时的 1..10，0 表示无牛
}

// niuPoint A=1，2-10 按面值，J/Q/K=10
func niuPoint(c int) int {
	r := Rank(c)
	if r == 12 {
		return 1
	}
	if r >= 8 {
		return 10
	}
	return r + 2
}

func NiuValue(cards []int) NiuResult {
	if len(cards) < 5 {
		return NiuResult{}
	}
	cards = cards[:5]
	pts := [5]int{}
	sum := 0
	maxRank := 0
	allFlower := true
	allSmall := true
	sameRank := map[int]int{}
	for i, c := range cards {
		pts[i] = niuPoint(c)
		sum += pts[i]
		if r := Rank(c); r > maxRank {
			maxRank = r
		}
		if r := Rank(c); r < 9 { // J,Q,K 之外
			allFlower = false
		}
		if pts[i] > 4 {
			allSmall = false
		}
		sameRank[Rank(c)]++
	}
	bomb := false
	for _, n := range sameRank {
		if n >= 4 {
			bomb = true
		}
	}
	cat := NiuNone
	niuNum := 0
	switch {
	case allSmall && sum <= 10:
		cat = NiuFiveSmall
	case bomb:
		cat = NiuBomb
	case allFlower:
		cat = NiuFiveFlower
	default:
		// 找三张和为 10 的倍数
		found := false
		for a := 0; a < 5 && !found; a++ {
			for b := a + 1; b < 5 && !found; b++ {
				for c := b + 1; c < 5 && !found; c++ {
					if (pts[a]+pts[b]+pts[c])%10 == 0 {
						rest := 0
						for k := 0; k < 5; k++ {
							if k != a && k != b && k != c {
								rest += pts[k]
							}
						}
						n := rest % 10
						if n == 0 {
							cat = NiuTen
							niuNum = 10
						} else {
							cat = n
							niuNum = n
						}
						found = true
					}
				}
			}
		}
		if !found {
			cat = NiuNone
		}
	}
	if cat >= NiuTen {
		niuNum = 10
	}
	return NiuResult{
		Cat: cat, Name: niuNames[cat], Mult: niuMults[cat],
		Score:  int64(cat)<<20 | int64(maxRank)<<4 | int64(Suit(cards[0])),
		NiuNum: niuNum,
	}
}
