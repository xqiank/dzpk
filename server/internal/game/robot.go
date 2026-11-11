package game

import "math"

// 机器人风格
type BotStyle struct {
	Name       string
	VPIP       float64 // 入池率
	Aggression float64 // 翻后加注倾向 0..1
	Bluff      float64 // 诈唬频率
	CallDown   float64 // 抓鸡倾向
}

var BotStyles = map[string]BotStyle{
	"tag":     {Name: "紧凶", VPIP: 0.28, Aggression: 0.62, Bluff: 0.18, CallDown: 0.42},
	"lag":     {Name: "松凶", VPIP: 0.44, Aggression: 0.74, Bluff: 0.34, CallDown: 0.5},
	"nit":     {Name: "岩石", VPIP: 0.16, Aggression: 0.4, Bluff: 0.05, CallDown: 0.25},
	"station": {Name: "跟注站", VPIP: 0.52, Aggression: 0.18, Bluff: 0.04, CallDown: 0.72},
	"maniac":  {Name: "疯子", VPIP: 0.72, Aggression: 0.88, Bluff: 0.5, CallDown: 0.6},
	"fish":    {Name: "娱乐玩家", VPIP: 0.55, Aggression: 0.3, Bluff: 0.12, CallDown: 0.55},
}

type BotProfile struct {
	UID        uint32
	Style      BotStyle
	Aggression int // 0..100 微调
	Bluff      int // 0..100
	VPIP       int
}

// RobotDecide 机器人决策（在桌锁内调用）。返回 action 与金额（raise to 语义）。
func RobotDecide(t *Table, profile BotProfile, seat int) (string, int64) {
	h := t.Hand
	if h == nil {
		return "fold", 0
	}
	s := t.Seats[seat]
	if s == nil || len(s.Cards) == 0 {
		return "fold", 0
	}
	toCall := h.CurrentBet - s.Bet
	if toCall < 0 {
		toCall = 0
	}
	opp := 0
	for _, x := range t.Seats {
		if x != nil && !x.Folded && x.Index != seat {
			opp++
		}
	}
	if opp < 1 {
		return "check", 0
	}
	short := t.Cfg.Short()
	iters := 420
	if h.Street >= 2 {
		iters = 700
	}
	eq := Equity(t.Cfg.GameType, s.Cards, h.Board, opp, iters, short)
	pot := float64(h.Pot)
	callF := float64(toCall)
	stack := float64(s.Stack)
	potOdds := 1.0
	if callF > 0 {
		potOdds = callF / (pot + callF)
	}
	st := profile.Style
	aggr := st.Aggression + float64(profile.Aggression-50)/260
	if aggr < 0.05 {
		aggr = 0.05
	}
	bluff := st.Bluff + float64(profile.Bluff-50)/400
	if bluff < 0.01 {
		bluff = 0.01
	}
	noise := float64(t.rnd.next()%1000)/1000.0 - 0.5
	eq += noise * 0.05

	// 底池比
	potRatio := 0.0
	if pot > 0 {
		potRatio = callF / pot
	}
	_ = stack

	// 翻前收窄范围：按 VPIP 过滤垃圾手
	preflop := h.Street == 0
	if preflop {
		openEQ := 0.5 / float64(opp+1) * 1.22
		if eq < openEQ*(1-st.VPIP) {
			if toCall == 0 {
				return "check", 0
			}
			// 小注时可以跟一下
			if potRatio < 0.12 && float64(t.rnd.next()%100) < st.VPIP*40 {
				return "call", 0
			}
			return "fold", 0
		}
	}

	// 价值加注
	strong := 0.52/float64(opp)*1.0 + 0.16
	if eq > strong && toCall == 0 {
		size := int64(pot * (0.55 + 0.25*aggr))
		if size < t.Cfg.BB {
			size = t.Cfg.BB
		}
		return "bet", clampTo(s, h, size)
	}
	if eq > strong && toCall > 0 {
		if float64(t.rnd.next()%100) < (35 + aggr*55) {
			target := h.CurrentBet + max64(h.MinRaise, int64((pot+callF)*(0.6+0.3*aggr)))
			return "raise", clampTo(s, h, target)
		}
		return "call", 0
	}

	// 中等牌力：跟注为主，少量加注
	if eq > potOdds+0.06 {
		if toCall == 0 {
			if eq > 0.42 && float64(t.rnd.next()%100) < aggr*45 {
				return "bet", clampTo(s, h, int64(pot*0.5))
			}
			return "check", 0
		}
		if eq > 0.45 && float64(t.rnd.next()%100) < aggr*22 {
			target := h.CurrentBet + max64(h.MinRaise, int64((pot+callF)*0.7))
			return "raise", clampTo(s, h, target)
		}
		if potRatio < 0.45 || float64(t.rnd.next()%100) < st.CallDown*100 {
			return "call", 0
		}
		return "fold", 0
	}

	// 诈唬
	if toCall == 0 {
		if float64(t.rnd.next()%1000) < bluff*220 {
			return "bet", clampTo(s, h, int64(pot*0.6))
		}
		return "check", 0
	}
	if potRatio < 0.25 && float64(t.rnd.next()%1000) < bluff*180 {
		if t.rnd.next()%100 < 25 {
			target := h.CurrentBet + max64(h.MinRaise, int64((pot+callF)*0.8))
			return "raise", clampTo(s, h, target)
		}
		return "call", 0
	}
	return "fold", 0
}

func clampTo(s *Seat, h *Hand, target int64) int64 {
	maxTo := s.Bet + s.Stack
	if target > maxTo {
		target = maxTo
	}
	minTo := h.CurrentBet + h.MinRaise
	if target < minTo && maxTo >= minTo {
		target = minTo
	}
	return target
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// 供教学计算器使用的精确胜率（迭代更多）
func EquityPrecise(gameType string, hole, board []int, opponents int) float64 {
	short := gameType == "short6"
	iters := 6000
	if len(board) >= 3 {
		iters = 9000
	}
	return math.Round(Equity(gameType, hole, board, opponents, iters, short)*10000) / 10000
}
