package game

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

const (
	StateWaiting  = "waiting"
	StatePlaying  = "playing"
	StateShowdown = "showdown"
	StateInsurance = "insurance"
	StateBidding   = "bidding"
	StateBetting   = "betting"
	StateClosed   = "closed"

	ActionTime     = 22 * time.Second
	NextHandDelay  = 4500 * time.Millisecond
	ResultHoldTime = 5200 * time.Millisecond
)

var (
	ErrSeatFull    = errors.New("座位已满")
	ErrNotSeated   = errors.New("未入座")
	ErrNotYourTurn = errors.New("未轮到行动")
	ErrBadAction   = errors.New("非法操作")
	ErrNoChips     = errors.New("筹码不足")
	ErrState       = errors.New("当前状态不可操作")
)

// ObserverInfo 围观者展示信息（报表「实时围观」用）
type ObserverInfo struct {
	UID    uint32 `json:"uid"`
	Nick   string `json:"nick"`
	Avatar string `json:"avatar"`
}

type Seat struct {
	Index      int
	UID        uint32
	Nick       string
	Avatar     string
	VIP        int
	IsRobot    bool
	Stack      int64
	Bet        int64
	Committed  int64
	Cards      []int
	Folded     bool
	AllIn      bool
	Acted      bool
	SittingOut bool
	WaitBB     bool
	WaitNext   bool
	LastAction string
	LastAmount int64
	Deadline   time.Time
	Insured    bool
	InsPremium int64
	InsPayout  int64
	Hands      int
	VPIPHands  int
	VPIPFlag   bool
	PFRHands   int
	PFRFlag    bool
	MuckBuyers map[uint32]bool
	Buyin      int64
	Net        int64
	Connected  bool
	DisconnAt  time.Time
	JoinedAt   time.Time
	Showdown   bool
	HandValue  HandValue
	Bid        int // 牛牛抢庄倍数（0=不抢）
	Stake      int // 牛牛闲家押注倍数（0=未押注，1/2/5/10）
}

type Hand struct {
	No        uint64
	DBID      int64
	Deck      []int
	Board     []int
	Street    int
	Pot       int64
	CurrentBet int64
	MinRaise  int64
	BankerBid  int // 抢庄成交倍数
	ToAct     int
	Dealer    int
	SBSeat    int
	BBSeat    int
	Seed      string
	SeedHash  string
	StartedAt time.Time
	Deadline  time.Time
	Actions   []HandAction
	Results   []HandResult
	Winners   []WinnerDTO
	Done      bool
	Showdown  bool
	AnteSum   int64
	Insurance *InsOffer
	InsUID     uint32
	InsPremium int64
	InsPayout  int64
	InsSettled bool
}

type HandAction struct {
	Seq      int
	Street   int
	Seat     int
	UID      uint32
	Action   string
	Amount   int64
	PotAfter int64
}

type HandResult struct {
	Seat     int    `json:"seat"`
	UID      uint32 `json:"uid"`
	Nick     string `json:"nick"`
	Cards    []string `json:"cards"`
	Best     []string `json:"best"`
	HandName string `json:"handName"`
	Win      int64  `json:"win"`
	Net      int64  `json:"net"`
	Insured  bool   `json:"insured"`
	InsPayout int64 `json:"insPayout"`
}

type InsOffer struct {
	Seat     int
	Equity   float64
	PotSize  int64
	Cover    int64
	Premium  int64
	Payout   int64
	Deadline time.Time
	Accepted bool
	Declined bool
	// 官方保单面板信息
	BoardCount int
	Outs       int
	OutCards   []int
	Odds       float64
	Plans      []InsPlan
	PlanKey    string
}

// InsPlan 一档投保方案（官方「选择投保方案」列表）
type InsPlan struct {
	Key     string `json:"key"`
	Label   string `json:"label"`
	Premium int64  `json:"premium"`
	Payout  int64  `json:"payout"`
}

// RoomCfg 房间配置（来自 DB）
type RoomCfg struct {
	RoomNo      int
	Name        string
	GameType    string
	SeatCount   int
	SB          int64
	BB          int64
	Ante        int64
	MinBuyin    int64
	MaxBuyin    int64
	BetLimit    string
	IsOfficial  bool
	ClubID      int64
	HostUID     uint32
	CardCost    int
	Insurance   bool
	DurationMin int
	Password    string
	Sort        int
	// 鱿鱼游戏模式：盈利提前离桌罚金 + 本轮赏金结算
	NiuBankerMode  string // 牛牛庄家模式：rob=抢庄 fixed=固定庄
	Squid          bool
	RequiredHands  int   // 盈利免罚所需的局数
	PenaltyPercent int   // 提前离桌按盈利百分比扣罚
	PenaltyFixed   int64 // 提前离桌固定扣罚（分）
}

func (r RoomCfg) Short() bool { return r.GameType == "short6" }

type Table struct {
	mu   sync.Mutex
	Cfg  RoomCfg
	Seats []*Seat
	Hand  *Hand

	State      string
	HandCounter uint64
	CreatedAt  time.Time
	ExpireAt   time.Time
	Version    int64
	cardPauseUntil time.Time
	Observers  map[uint32]bool
	ObsMeta    map[uint32]ObserverInfo
	NextHandAt time.Time
	Winners    []WinnerDTO
	SidePots   []SidePotDTO

	LastAction   *ActionAnim
	LastNewCards []CardAnim
	LastWin      *WinAnim

	OnChange func(*Table)
	OnHandEnd func(*Table, *Hand)
	OnLeave  func(*Table, uint32, int64)
	OnHandStart func(*Table) bool
	OnLog    func(format string, args ...any)

	rnd      *fastRand
	lastHandEnd time.Time
	lastHandID  int64 // 最近一手牌谱 id（持久化后回填，供回放入口使用）
	dealer   int
	closed   bool
	idleAt   time.Time
	startedAt time.Time
}

func NewTable(cfg RoomCfg) *Table {
	if cfg.SeatCount <= 0 || cfg.SeatCount > 9 {
		cfg.SeatCount = 6
	}
	t := &Table{
		Cfg:       cfg,
		Seats:     make([]*Seat, cfg.SeatCount),
		State:     StateWaiting,
		CreatedAt: time.Now(),
		Observers: map[uint32]bool{},
		ObsMeta:   map[uint32]ObserverInfo{},
		dealer:    -1,
		rnd:       newFastRand(uint64(time.Now().UnixNano()) ^ uint64(cfg.RoomNo)*2654435761),
		idleAt:    time.Now(),
	}
	if cfg.DurationMin > 0 {
		t.ExpireAt = t.CreatedAt.Add(time.Duration(cfg.DurationMin) * time.Minute)
	}
	return t
}

func (t *Table) Key() string { return fmt.Sprintf("room:%d", t.Cfg.RoomNo) }

func (t *Table) logf(f string, a ...any) {
	if t.OnLog != nil {
		t.OnLog(f, a...)
	}
}

// SetLastHandID 回填最近一手牌谱 id（持久化完成后调用）
func (t *Table) SetLastHandID(id int64) {
	t.mu.Lock()
	t.lastHandID = id
	t.mu.Unlock()
	t.notify()
}

// LastHandID 最近一手牌谱 id
func (t *Table) LastHandID() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastHandID
}

// refreshOfficialExpire 官方桌循环续期：倒计时走完自动续 1 小时
func (t *Table) refreshOfficialExpire(now time.Time) {
	if !t.Cfg.IsOfficial || t.ExpireAt.IsZero() || !now.After(t.ExpireAt) {
		return
	}
	t.ExpireAt = now.Add(time.Hour)
}

func (t *Table) notify() {
	t.Version++
	if t.OnChange != nil {
		t.OnChange(t)
	}
}

// ---------------- 入座 / 离座 ----------------

func (t *Table) SeatOf(uid uint32) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.seatOfLocked(uid)
}

func (t *Table) seatOfLocked(uid uint32) int {
	for i, s := range t.Seats {
		if s != nil && s.UID == uid {
			return i
		}
	}
	return -1
}

type SitRequest struct {
	UID     uint32
	Nick    string
	Avatar  string
	VIP     int
	IsRobot bool
	Buyin   int64
	Seat    int
}

func (t *Table) Sit(req SitRequest) (int, int64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return -1, 0, ErrState
	}
	if idx := t.seatOfLocked(req.UID); idx >= 0 {
		s := t.Seats[idx]
		// 已入座则补码
		add := req.Buyin
		if add <= 0 {
			return idx, 0, nil
		}
		if s.Stack+add > t.Cfg.MaxBuyin {
			add = t.Cfg.MaxBuyin - s.Stack
		}
		if add <= 0 {
			return idx, 0, ErrBadAction
		}
		s.Stack += add
		s.Buyin += add
		s.SittingOut = false
		t.idleAt = time.Now()
		t.notify()
		return idx, add, nil
	}
	if t.Cfg.SeatCount <= 0 {
		return -1, 0, ErrSeatFull
	}
	seat := -1
	if req.Seat >= 0 && req.Seat < len(t.Seats) && t.Seats[req.Seat] == nil {
		seat = req.Seat
	} else {
		for i, s := range t.Seats {
			if s == nil {
				seat = i
				break
			}
		}
	}
	if seat < 0 {
		return -1, 0, ErrSeatFull
	}
	buyin := req.Buyin
	if buyin < t.Cfg.MinBuyin {
		buyin = t.Cfg.MinBuyin
	}
	if buyin > t.Cfg.MaxBuyin {
		buyin = t.Cfg.MaxBuyin
	}
	t.Seats[seat] = &Seat{
		Index: seat, UID: req.UID, Nick: req.Nick, Avatar: req.Avatar, VIP: req.VIP,
		IsRobot: req.IsRobot, Stack: buyin, Buyin: buyin, Connected: true,
		JoinedAt: time.Now(), WaitBB: true,
		WaitNext: t.Hand != nil && !t.Hand.Done,
	}
	if t.State == StateWaiting && t.countActiveLocked() >= 2 {
		t.startHandLocked()
	}
	t.idleAt = time.Now()
	t.notify()
	return seat, buyin, nil
}

// ReduceChips 桌上筹码提前退回钱包：保留 keepMin 以防空座
func (t *Table) ReduceChips(uid uint32, amount int64) (int64, int64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	idx := t.seatOfLocked(uid)
	if idx < 0 {
		return 0, 0, ErrNotSeated
	}
	s := t.Seats[idx]
	if s == nil {
		return 0, 0, ErrNotSeated
	}
	keep := t.Cfg.BB
	if keep < 1 {
		keep = 1
	}
	if t.Hand != nil && !t.Hand.Done && s.Committed > 0 {
		return 0, 0, errors.New("本手进行中，结算后再退码")
	}
	if amount <= 0 {
		return 0, 0, errors.New("退码金额需大于 0")
	}
	if s.Stack-amount < keep {
		return 0, 0, fmt.Errorf("在座期间需至少保留 %d 筹码", keep)
	}
	s.Stack -= amount
	t.notify()
	return amount, s.Stack, nil
}

func (t *Table) Stand(uid uint32) (int, int64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	idx := t.seatOfLocked(uid)
	if idx < 0 {
		return -1, 0, ErrNotSeated
	}
	s := t.Seats[idx]
	cash := s.Stack
	if t.Hand != nil && !t.Hand.Done && t.Hand.Street >= 0 && s.Committed > 0 && t.State == StatePlaying {
		// 手中还有牌，先标记离座，结算后再退
		s.SittingOut = true
		s.Connected = false
		t.notify()
		return idx, 0, nil
	}
	// 鱿鱼模式：离桌前先结清盈利罚金 / 本轮赏金，罚款与赏金都从退回筹码里扣
	if t.Cfg.Squid {
		cash, _, _ = t.settleSquidExitLocked(idx, cash)
	}
	t.Seats[idx] = nil
	t.idleAt = time.Now()
	if t.countActiveLocked() < 2 && t.State == StatePlaying {
		t.finishHandLocked()
	}
	t.notify()
	return idx, cash, nil
}

// settleSquidExitLocked 结清鱿鱼模式离桌费用：返回实退筹码、罚金、赏金
// 罚金直接从筹码扣除（销毁），赏金按比例补给本轮仍盈利的在座玩家（不给机器人以外的空位）
func (t *Table) settleSquidExitLocked(idx int, cash int64) (int64, int64, int64) {
	s := t.Seats[idx]
	if s == nil {
		return cash, 0, 0
	}
	var penalty, bounty int64

	req := t.Cfg.RequiredHands
	if req <= 0 {
		req = 20
	}
	// 盈利提前离桌罚金
	if s.Net > 0 && s.Hands < req {
		pct := t.Cfg.PenaltyPercent
		if pct <= 0 {
			pct = 5
		}
		penalty = s.Net*int64(pct)/100 + t.Cfg.PenaltyFixed
		if penalty > cash {
			penalty = cash
		}
		if penalty < 0 {
			penalty = 0
		}
		cash -= penalty
	}
	// 尚未脱险：把亏损额作为赏金赔付给本轮盈利的其它在座玩家
	if s.Net < 0 {
		beneficiaries := []*Seat{}
		for i, x := range t.Seats {
			if i == idx || x == nil || x.Net <= 0 {
				continue
			}
			beneficiaries = append(beneficiaries, x)
		}
		if len(beneficiaries) > 0 {
			want := -s.Net
			if want > cash {
				want = cash // 按可用筹码赔付，不产生欠款
			}
			if want > 0 {
				each := want / int64(len(beneficiaries))
				rest := want - each*int64(len(beneficiaries))
				for k, b := range beneficiaries {
					add := each
					if int64(k) < rest {
						add++
					}
					b.Stack += add
					b.Net += add
				}
				bounty = want
				cash -= want
			}
		}
	}
	if cash < 0 {
		cash = 0
	}
	if penalty > 0 || bounty > 0 {
		t.logf("room %d uid=%d 鱿鱼离桌结算 罚金=%d 赏金=%d 实退=%d", t.Cfg.RoomNo, s.UID, penalty, bounty, cash)
	}
	return cash, penalty, bounty
}

func (t *Table) countActiveLocked() int {
	n := 0
	for _, s := range t.Seats {
		if s != nil && !s.SittingOut && s.Stack > 0 {
			n++
		}
	}
	return n
}

func (t *Table) SetConnected(uid uint32, on bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if i := t.seatOfLocked(uid); i >= 0 {
		t.Seats[i].Connected = on
		if on {
			t.Seats[i].DisconnAt = time.Time{}
		} else {
			t.Seats[i].DisconnAt = time.Now()
		}
		t.notify()
	}
}

// ---------------- 开牌 ----------------

func (t *Table) TryStartHand() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.State != StateWaiting {
		return
	}
	if t.countActiveLocked() >= 2 {
		t.startHandLocked()
	}
}

func (t *Table) startHandLocked() {
	if t.Cfg.GameType == "niuniu" {
		t.startNiuNiuLocked()
		return
	}
	active := []int{}
	for i, s := range t.Seats {
		if s != nil && !s.SittingOut && s.Stack > 0 {
			active = append(active, i)
		}
	}
	if len(active) < 2 {
		return
	}
	if t.Cfg.CardCost > 0 && !t.Cfg.IsOfficial && t.OnHandStart != nil {
		// 房卡不足时冷却 30 秒再试，避免每 tick 都去打库与刷日志
		if time.Now().Before(t.cardPauseUntil) {
			return
		}
		if !t.OnHandStart(t) {
			t.cardPauseUntil = time.Now().Add(30 * time.Second)
			t.logf("room %d 房卡不足，暂停开局（30 秒后重试）", t.Cfg.RoomNo)
			return
		}
	}
	t.HandCounter++
	hand := &Hand{
		No: t.HandCounter, Street: 0, Dealer: t.dealer,
		StartedAt: time.Now(), MinRaise: t.Cfg.BB,
	}
	seedEntropy := fmt.Sprintf("%d|%d|%d|%d", t.Cfg.RoomNo, t.HandCounter, time.Now().UnixNano(), t.rnd.next())
	hand.Seed = NewSeed(seedEntropy)
	hand.SeedHash = SeedHash(hand.Seed)
	hand.Deck = Shuffle(NewDeck(t.Cfg.Short()), hand.Seed, t.HandCounter)

	// 移动按钮
	t.dealer = t.nextSeatLocked(hand.Dealer, func(s *Seat) bool {
		return s != nil && !s.SittingOut && s.Stack > 0
	})
	hand.Dealer = t.dealer

	// 清空座位手牌状态
	for _, s := range t.Seats {
		if s == nil {
			continue
		}
		s.Cards = nil
		s.Bet = 0
		s.Committed = 0
		s.Folded = false
		s.AllIn = false
		s.Acted = false
		s.LastAction = ""
		s.LastAmount = 0
		s.Deadline = time.Time{}
		s.Showdown = false
		s.Insured = false
		s.InsPremium = 0
		s.InsPayout = 0
		s.HandValue = HandValue{}
	}

	order := t.seatsFromLocked(hand.Dealer, true)
	activeOrder := []int{}
	for _, i := range order {
		s := t.Seats[i]
		if s != nil && !s.SittingOut && s.Stack > 0 {
			s.WaitNext = false
			activeOrder = append(activeOrder, i)
		}
	}
	if len(activeOrder) == 2 {
		hand.SBSeat = activeOrder[0]
		hand.BBSeat = activeOrder[1]
	} else {
		hand.SBSeat = activeOrder[0]
		hand.BBSeat = activeOrder[1]
	}

	// 前注
	if t.Cfg.Ante > 0 {
		for _, i := range activeOrder {
			s := t.Seats[i]
			pay := t.Cfg.Ante
			if pay > s.Stack {
				pay = s.Stack
			}
			s.Stack -= pay
			s.Committed += pay
			hand.Pot += pay
			hand.AnteSum += pay
			if s.Stack == 0 {
				s.AllIn = true
				s.Acted = true
			}
		}
	}

	// 盲注
	post := func(i int, amt int64, label string) {
		s := t.Seats[i]
		pay := amt
		if pay > s.Stack {
			pay = s.Stack
		}
		s.Stack -= pay
		s.Bet += pay
		s.Committed += pay
		hand.Pot += pay
		if s.Stack == 0 {
			s.AllIn = true
		}
		s.LastAction = label
		s.LastAmount = pay
		hand.Actions = append(hand.Actions, HandAction{Seq: len(hand.Actions) + 1, Street: 0, Seat: i, UID: s.UID, Action: label, Amount: pay, PotAfter: hand.Pot})
	}
	post(hand.SBSeat, t.Cfg.SB, "sb")
	post(hand.BBSeat, t.Cfg.BB, "bb")
	hand.CurrentBet = t.Cfg.BB
	if t.Cfg.SB > hand.CurrentBet {
		hand.CurrentBet = t.Cfg.SB
	}
	hand.MinRaise = t.Cfg.BB

	// 发底牌
	holeCount := 2
	if t.Cfg.GameType == "omaha" {
		holeCount = 4
	}
	newCards := []CardAnim{}
	for k := 0; k < holeCount; k++ {
		for _, i := range activeOrder {
			s := t.Seats[i]
			c := NextCard(&hand.Deck)
			s.Cards = append(s.Cards, c)
			newCards = append(newCards, CardAnim{Seat: i, Card: CardStr(c), Index: len(s.Cards) - 1, Street: 0})
		}
	}

	// 谁先行动
	if len(activeOrder) == 2 {
		hand.ToAct = hand.SBSeat // 单挑时按钮位(小盲)先行动
	} else {
		hand.ToAct = t.nextSeatLocked(hand.BBSeat, func(s *Seat) bool {
			return s != nil && !s.Folded && !s.AllIn && !s.SittingOut && s.Stack > 0 && len(s.Cards) > 0
		})
		if hand.ToAct < 0 {
			// 没有任何可行动座位（极端情况）：直接按现有牌力结算，避免整桌卡死
			t.Hand = hand
			t.State = StatePlaying
			t.finishHandLocked()
			return
		}
	}
	for _, s := range t.Seats {
		if s != nil && s.Stack > 0 {
			s.Acted = false
			s.WaitBB = false
		}
	}
	hand.Deadline = t.actorDeadlineLocked(hand.ToAct)
	if s := t.Seats[hand.ToAct]; s != nil {
		s.Deadline = hand.Deadline
	}
	sort.Slice(newCards, func(i, j int) bool { return newCards[i].Index < newCards[j].Index })
	t.Hand = hand
	t.State = StatePlaying
	t.Winners = nil
	t.SidePots = nil
	t.LastWin = nil
	t.LastAction = nil
	t.LastNewCards = newCards
	t.startedAt = time.Now()
	t.logf("room %d hand %d start seats=%v", t.Cfg.RoomNo, hand.No, activeOrder)
	t.notify()
}


// actorDeadlineLocked 机器人思考时间更短，人类 22 秒
func (t *Table) actorDeadlineLocked(seat int) time.Time {
	if s := t.Seats[seat]; s != nil && s.IsRobot {
		d := 900 + t.rnd.next()%1700
		return time.Now().Add(time.Duration(d) * time.Millisecond)
	}
	return time.Now().Add(ActionTime)
}

func (t *Table) seatsFromLocked(start int, inclusive bool) []int {
	n := len(t.Seats)
	out := make([]int, 0, n)
	if n == 0 {
		return out
	}
	i := start
	if !inclusive {
		i = (start + 1) % n
	}
	for k := 0; k < n; k++ {
		out = append(out, i%n)
		i++
	}
	return out
}

func (t *Table) nextSeatLocked(from int, ok func(*Seat) bool) int {
	n := len(t.Seats)
	for k := 1; k <= n; k++ {
		i := ((from + k) % n + n) % n
		if ok(t.Seats[i]) {
			return i
		}
	}
	return -1
}

func (t *Table) prevSeatLocked(from int, ok func(*Seat) bool) int {
	n := len(t.Seats)
	for k := 1; k <= n; k++ {
		i := ((from-k)%n + n) % n
		if ok(t.Seats[i]) {
			return i
		}
	}
	return -1
}

// ---------------- 行动 ----------------

func (t *Table) Action(uid uint32, action string, amount int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.actionLocked(uid, action, amount)
}

func (t *Table) actionLocked(uid uint32, action string, amount int64) error {
	if t.State != StatePlaying || t.Hand == nil || t.Hand.Done {
		return ErrState
	}
	idx := t.seatOfLocked(uid)
	if idx < 0 {
		return ErrNotSeated
	}
	if t.Hand.ToAct != idx {
		return ErrNotYourTurn
	}
	s := t.Seats[idx]
	if s.Folded || s.AllIn || s.WaitNext || len(s.Cards) == 0 {
		return ErrBadAction
	}
	h := t.Hand
	toCall := h.CurrentBet - s.Bet
	if toCall < 0 {
		toCall = 0
	}
	effective := toCall
	if effective > s.Stack {
		effective = s.Stack
	}

	switch action {
	case "fold":
		s.Folded = true
		s.Acted = true
		s.LastAction = "fold"
		s.LastAmount = 0
		h.Actions = append(h.Actions, HandAction{Seq: len(h.Actions) + 1, Street: h.Street, Seat: idx, UID: uid, Action: "fold", PotAfter: h.Pot})
		t.LastAction = &ActionAnim{Seat: idx, Action: "fold", Street: h.Street, TS: time.Now().UnixMilli()}
	case "check":
		if toCall > 0 {
			return ErrBadAction
		}
		s.Acted = true
		s.LastAction = "check"
		s.LastAmount = 0
		h.Actions = append(h.Actions, HandAction{Seq: len(h.Actions) + 1, Street: h.Street, Seat: idx, UID: uid, Action: "check", PotAfter: h.Pot})
		t.LastAction = &ActionAnim{Seat: idx, Action: "check", Street: h.Street, TS: time.Now().UnixMilli()}
	case "call":
		if toCall <= 0 {
			return ErrBadAction
		}
		if h.Street == 0 {
			s.VPIPFlag = true
		}
		pay := effective
		s.Stack -= pay
		s.Bet += pay
		s.Committed += pay
		h.Pot += pay
		if s.Stack == 0 {
			s.AllIn = true
		}
		s.Acted = true
		s.LastAction = "call"
		s.LastAmount = pay
		h.Actions = append(h.Actions, HandAction{Seq: len(h.Actions) + 1, Street: h.Street, Seat: idx, UID: uid, Action: "call", Amount: pay, PotAfter: h.Pot})
		t.LastAction = &ActionAnim{Seat: idx, Action: "call", Amount: pay, Street: h.Street, TS: time.Now().UnixMilli()}
	case "bet", "raise", "allin":
		if h.Street == 0 {
			s.VPIPFlag = true
			s.PFRFlag = true
		}
		var target int64
		if action == "allin" {
			target = s.Bet + s.Stack
		} else {
			target = amount
		}
		if target > s.Bet+s.Stack {
			target = s.Bet + s.Stack
		}
		if target <= h.CurrentBet && target < s.Bet+s.Stack {
			return ErrBadAction
		}
		minTo := h.CurrentBet + h.MinRaise
		isAllIn := target >= s.Bet+s.Stack
		if !isAllIn && target < minTo {
			return ErrBadAction
		}
		if t.Cfg.BetLimit == "pl" {
			maxTo := h.Pot + toCall + s.Bet
			if target > maxTo && !isAllIn {
				target = maxTo
			}
		}
		pay := target - s.Bet
		if pay > s.Stack {
			pay = s.Stack
		}
		s.Stack -= pay
		s.Bet += pay
		s.Committed += pay
		h.Pot += pay
		raiseSize := s.Bet - h.CurrentBet
		if s.Bet > h.CurrentBet {
			if raiseSize >= h.MinRaise {
				h.MinRaise = raiseSize
				for _, o := range t.Seats {
					if o != nil && o.Index != idx {
						o.Acted = false
					}
				}
			}
			h.CurrentBet = s.Bet
		}
		if s.Stack == 0 {
			s.AllIn = true
		}
		s.Acted = true
		label := action
		if s.Bet <= h.CurrentBet {
			label = "call"
		} else if action != "allin" {
			if toCall == 0 {
				label = "bet"
			} else {
				label = "raise"
			}
		}
		s.LastAction = label
		s.LastAmount = pay
		h.Actions = append(h.Actions, HandAction{Seq: len(h.Actions) + 1, Street: h.Street, Seat: idx, UID: uid, Action: label, Amount: pay, PotAfter: h.Pot})
		t.LastAction = &ActionAnim{Seat: idx, Action: label, Amount: pay, Street: h.Street, TS: time.Now().UnixMilli()}
	default:
		return ErrBadAction
	}

	s.Deadline = time.Time{}
	h.ToAct = -1

	// 只剩一人
	remain := []int{}
	for _, x := range t.Seats {
		if x != nil && !x.Folded && len(x.Cards) > 0 {
			remain = append(remain, x.Index)
		}
	}
	if len(remain) <= 1 {
		t.finishHandLocked()
		t.notify()
		return nil
	}

	// 本轮是否结束
	roundDone := true
	for _, x := range t.Seats {
		if x == nil || x.Folded || x.AllIn || x.WaitNext || x.SittingOut || len(x.Cards) == 0 {
			continue
		}
		if !x.Acted || x.Bet < h.CurrentBet {
			roundDone = false
			break
		}
	}
	if roundDone {
		t.advanceStreetLocked()
	} else {
		next := t.nextSeatLocked(idx, func(x *Seat) bool {
			return x != nil && !x.Folded && !x.AllIn && !x.WaitNext && len(x.Cards) > 0
		})
		if next < 0 {
			t.advanceStreetLocked()
		} else {
			h.ToAct = next
			h.Deadline = t.actorDeadlineLocked(next)
			t.Seats[next].Deadline = h.Deadline
		}
	}
	t.notify()
	return nil
}

func (t *Table) advanceStreetLocked() {
	h := t.Hand
	// 收注
	for _, s := range t.Seats {
		if s == nil {
			continue
		}
		s.Bet = 0
		s.Acted = false
		s.Deadline = time.Time{}
	}
	h.CurrentBet = 0
	h.MinRaise = t.Cfg.BB
	newCards := []CardAnim{}

	canAct := 0
	for _, s := range t.Seats {
		if s == nil || s.Folded || s.WaitNext || len(s.Cards) == 0 {
			continue
		}
		if !s.AllIn {
			canAct++
		}
	}
	if canAct <= 1 {
		// 所有人全下：先摊牌（亮牌）→ 可选保险 → 跑完公共牌
		if t.pauseForInsuranceLocked() {
			return
		}
		t.runoutLocked(&newCards)
		t.LastNewCards = newCards
		t.showdownLocked()
		return
	}
	if h.Street >= 3 {
		t.showdownLocked()
		return
	}
	h.Street++
	n := 1
	if h.Street == 1 {
		n = 3
	}
	for k := 0; k < n; k++ {
		c := NextCard(&h.Deck)
		h.Board = append(h.Board, c)
		newCards = append(newCards, CardAnim{Seat: -1, Card: CardStr(c), Index: len(h.Board) - 1, Street: h.Street})
	}
	t.LastNewCards = newCards

	next := t.nextSeatLocked(h.Dealer, func(s *Seat) bool {
		return s != nil && !s.Folded && !s.AllIn && !s.WaitNext && !s.SittingOut && len(s.Cards) > 0
	})
	if next < 0 {
		t.showdownLocked()
		return
	}
	h.ToAct = next
	h.Deadline = t.actorDeadlineLocked(next)
	t.Seats[next].Deadline = h.Deadline
}

func (t *Table) runoutLocked(newCards *[]CardAnim) {
	h := t.Hand
	for h.Street < 3 {
		h.Street++
		n := 1
		if h.Street == 1 {
			n = 3
		}
		for k := 0; k < n; k++ {
			c := NextCard(&h.Deck)
			h.Board = append(h.Board, c)
			*newCards = append(*newCards, CardAnim{Seat: -1, Card: CardStr(c), Index: len(h.Board) - 1, Street: h.Street})
		}
	}
}

// pauseForInsuranceLocked 全下摊牌时暂停 12 秒提供保险；返回 true 表示已暂停
func (t *Table) pauseForInsuranceLocked() bool {
	h := t.Hand
	if h == nil || !t.Cfg.Insurance || h.Insurance != nil {
		return false
	}
	players := []*Seat{}
	for _, s := range t.Seats {
		if s != nil && !s.Folded && len(s.Cards) > 0 {
			players = append(players, s)
		}
	}
	if len(players) != 2 {
		return false
	}
	a, b := players[0], players[1]
	if !(a.AllIn && b.AllIn) {
		return false
	}
	eq := Equity(t.Cfg.GameType, a.Cards, h.Board, 1, 900, t.Cfg.Short())
	if eq < 0.05 || eq > 0.95 {
		return false
	}
	target := a
	other := b
	behind := eq
	if eq > 0.5 {
		target = b
		other = a
		behind = 1 - eq
	}
	cover := h.Pot / 2
	if cover < t.Cfg.BB*2 {
		return false
	}
	premium := int64(float64(cover) * (1 - behind) * 1.06)
	if premium <= 0 || premium >= target.Stack+cover {
		return false
	}
	// 亮牌
	for _, s := range players {
		s.HandValue = Evaluate(t.Cfg.GameType, s.Cards, h.Board)
		s.Showdown = true
	}
	outs, outCards := insuranceOuts(t.Cfg.GameType, target, other, h.Board, h.Deck, t.Cfg.Short())
	h.Insurance = &InsOffer{
		Seat: target.Index, Equity: behind, PotSize: h.Pot,
		Cover: cover, Premium: premium, Payout: cover,
		Deadline: time.Now().Add(12 * time.Second),
		BoardCount: len(h.Board), Outs: outs, OutCards: outCards,
		Odds:  oddsRatio(cover, premium),
		Plans: buildInsPlans(cover, behind),
		PlanKey: defaultInsPlan(cover, behind),
	}
	t.State = StateInsurance
	t.notify()
	return true
}

// insuranceOuts 统计「下一张公牌就能反超」的牌的张数与牌面（官方「对手可能逆转的反超牌」）
func insuranceOuts(gameType string, target, other *Seat, board []int, deck []int, short bool) (int, []int) {
	if target == nil || other == nil || len(board) == 0 {
		return 0, nil
	}
	out := make([]int, 0, 12)
	for _, c := range deck {
		nb := make([]int, len(board)+1)
		copy(nb, board)
		nb[len(board)] = c
		mine := Evaluate(gameType, target.Cards, nb)
		theirs := Evaluate(gameType, other.Cards, nb)
		if theirs.Score > mine.Score {
			out = append(out, c)
		}
	}
	return len(out), out
}

// oddsRatio 参考赔率（1 : payout/premium）
func oddsRatio(cover, premium int64) float64 {
	if premium <= 0 {
		return 0
	}
	return float64(cover) / float64(premium)
}

// buildInsPlans 三档投保方案：半池 / 七成 / 满池，保费按承保额等比例缩放
func buildInsPlans(cover int64, behind float64) []InsPlan {
	if cover <= 0 {
		return nil
	}
	tiers := []struct {
		key   string
		label string
		ratio float64
	}{
		{"half", "半池保险 · 保守", 0.5},
		{"most", "七成保险 · 平衡", 0.7},
		{"full", "满池保险 · 全额", 1.0},
	}
	out := make([]InsPlan, 0, len(tiers))
	for _, tr := range tiers {
		payout := int64(float64(cover) * tr.ratio)
		prem := int64(float64(payout) * (1 - behind) * 1.06)
		if prem < 1 {
			prem = 1
		}
		out = append(out, InsPlan{Key: tr.key, Label: tr.label, Premium: prem, Payout: payout})
	}
	return out
}

func defaultInsPlan(cover int64, behind float64) string {
	for _, p := range buildInsPlans(cover, behind) {
		if p.Key == "most" {
			return p.Key
		}
	}
	return "full"
}

// QAOfferInsurance 联调专用：为指定座位伪造一份保险报价
func (t *Table) QAOfferInsurance(seat int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.Hand == nil {
		return
	}
	pot := t.Hand.Pot
	if pot < t.Cfg.BB*4 {
		pot = t.Cfg.BB * 20
	}
	cover := pot / 2
	premium := int64(float64(cover) * 0.42)
	var target, other *Seat
	if seat >= 0 && seat < len(t.Seats) {
		target = t.Seats[seat]
	}
	for _, s := range t.Seats {
		if s != nil && s.Index != seat && len(s.Cards) > 0 {
			other = s
			break
		}
	}
	board := t.Hand.Board
	if len(board) == 0 {
		board = []int{NextCard(&t.Hand.Deck), NextCard(&t.Hand.Deck), NextCard(&t.Hand.Deck)}
	}
	outs, outCards := 0, []int(nil)
	if target != nil && other != nil && len(target.Cards) > 0 {
		outs, outCards = insuranceOuts(t.Cfg.GameType, target, other, board, t.Hand.Deck, t.Cfg.Short())
	}
	t.Hand.Insurance = &InsOffer{
		Seat: seat, Equity: 0.58, PotSize: pot, Cover: cover,
		Premium: premium, Payout: cover, Deadline: time.Now().Add(12 * time.Second),
		BoardCount: len(board), Outs: outs, OutCards: outCards,
		Odds:     oddsRatio(cover, premium),
		Plans:    buildInsPlans(cover, 0.58),
		PlanKey:  defaultInsPlan(cover, 0.58),
	}
	t.State = StateInsurance
	t.notify()
}

func (t *Table) resumeFromInsuranceLocked() {
	var newCards []CardAnim
	t.runoutLocked(&newCards)
	t.LastNewCards = newCards
	t.showdownLocked()
}

// ---------------- 结算 ----------------

func (t *Table) finishHandLocked() {
	h := t.Hand
	if h == nil || h.Done {
		return
	}
	h.Done = true
	h.ToAct = -1
	for _, s := range t.Seats {
		if s != nil {
			s.Deadline = time.Time{}
		}
	}
	// 只剩一人，直接收池
	remain := []int{}
	for _, s := range t.Seats {
		if s != nil && !s.Folded && s.Cards != nil {
			remain = append(remain, s.Index)
		}
	}
	if len(remain) == 1 {
		s := t.Seats[remain[0]]
		win := h.Pot
		s.Stack += win
		s.Net += win - s.Committed
		s.Hands++
		if s.VPIPFlag {
			s.VPIPHands++
			s.VPIPFlag = false
		}
		if s.PFRFlag {
			s.PFRHands++
			s.PFRFlag = false
		}
		h.Winners = []WinnerDTO{{Seat: s.Index, UID: s.UID, Nick: s.Nick, Amount: win, HandName: "", Cards: cardsToStr(s.Cards)}}
		t.Winners = h.Winners
		t.LastWin = &WinAnim{Seat: s.Index, Amount: win, Pot: win}
		for _, x := range t.Seats {
			if x != nil && !x.Folded {
				x.Showdown = true
			}
		}
	} else {
		t.showdownLocked()
		return
	}
	t.finishCommonLocked()
}

func (t *Table) showdownLocked() {
	h := t.Hand
	if h == nil {
		return
	}
	for _, s := range t.Seats {
		if s == nil || s.Folded || len(s.Cards) == 0 {
			continue
		}
		s.HandValue = Evaluate(t.Cfg.GameType, s.Cards, h.Board)
		s.Showdown = true
	}
	// 边池
	pots := t.buildPotsLocked()
	t.SidePots = nil
	winBySeat := map[int]int64{}
	for _, p := range pots {
		t.SidePots = append(t.SidePots, SidePotDTO{Amount: p.amount, Seats: p.seats})
		if len(p.seats) == 0 {
			// 全部弃牌，退回最近的下注者
			continue
		}
		var best int64 = -1
		winners := []int{}
		for _, seat := range p.seats {
			s := t.Seats[seat]
			if s == nil || s.Folded {
				continue
			}
			if s.HandValue.Score > best {
				best = s.HandValue.Score
				winners = []int{seat}
			} else if s.HandValue.Score == best {
				winners = append(winners, seat)
			}
		}
		if len(winners) == 0 || best < 0 {
			continue
		}
		share := p.amount / int64(len(winners))
		rem := p.amount - share*int64(len(winners))
		for i, seat := range winners {
			add := share
			if i == 0 {
				add += rem
			}
			t.Seats[seat].Stack += add
			winBySeat[seat] += add
		}
	}
	h.Winners = nil
	for _, s := range t.Seats {
		if s == nil {
			continue
		}
		if w, ok := winBySeat[s.Index]; ok {
			h.Winners = append(h.Winners, WinnerDTO{
				Seat: s.Index, UID: s.UID, Nick: s.Nick, Amount: w,
				HandName: s.HandValue.CatName, Cards: cardsToStr(s.Cards),
				Best: cardsToStr(s.HandValue.Best),
			})
			if w > 0 {
				t.LastWin = &WinAnim{Seat: s.Index, Amount: w, Pot: h.Pot}
			}
		}
	}
	sort.Slice(h.Winners, func(i, j int) bool { return h.Winners[i].Amount > h.Winners[j].Amount })
	t.Winners = h.Winners
	t.finishCommonLocked()
}

func (t *Table) finishCommonLocked() {
	h := t.Hand
	h.Done = true
	h.Showdown = true
	t.State = StateShowdown
	t.NextHandAt = time.Now().Add(ResultHoldTime)
	t.lastHandEnd = time.Now()
	if h.Insurance != nil && h.Insurance.Accepted {
		if s := t.Seats[h.Insurance.Seat]; s != nil {
			s.Insured = true
			s.InsPremium = h.Insurance.Premium
			s.InsPayout = h.Insurance.Payout
		}
	}
	// 结算结果
	results := []HandResult{}
	for _, s := range t.Seats {
		if s == nil || len(s.Cards) == 0 {
			continue
		}
		win := int64(0)
		for _, w := range h.Winners {
			if w.Seat == s.Index {
				win = w.Amount
			}
		}
		s.Net += win - s.Committed
		s.Hands++
		if s.VPIPFlag {
			s.VPIPHands++
			s.VPIPFlag = false
		}
		if s.PFRFlag {
			s.PFRHands++
			s.PFRFlag = false
		}
		if s.Insured && win == 0 && s.InsPayout > 0 {
			s.Stack += s.InsPayout
		}
		res := HandResult{
			Seat: s.Index, UID: s.UID, Nick: s.Nick,
			Cards: cardsToStr(s.Cards), HandName: s.HandValue.CatName,
			Win: win, Net: win - s.Committed, Insured: s.Insured, InsPayout: s.InsPayout,
		}
		if len(s.HandValue.Best) > 0 {
			res.Best = cardsToStr(s.HandValue.Best)
		}
		results = append(results, res)
	}
	h.Results = results
	t.LastNewCards = nil
	if h.Insurance != nil && h.Insurance.Accepted {
		s := t.Seats[h.Insurance.Seat]
		if s != nil {
			won := false
			for _, w := range h.Winners {
				if w.Seat == s.Index {
					won = true
				}
			}
			h.InsUID = s.UID
			h.InsPremium = h.Insurance.Premium
			if !won {
				h.InsPayout = h.Insurance.Payout
			}
		}
	}
	h.InsSettled = true
	if t.OnHandEnd != nil {
		hh := h
		tb := t
		go func() { tb.OnHandEnd(tb, hh) }()
	}
	t.logf("room %d hand %d end pot=%d winners=%v", t.Cfg.RoomNo, h.No, h.Pot, h.Winners)
}

type potLayer struct {
	amount int64
	seats  []int
}

func (t *Table) buildPotsLocked() []potLayer {
	levels := []int64{}
	seen := map[int64]bool{}
	var total int64
	for _, s := range t.Seats {
		if s == nil {
			continue
		}
		total += s.Committed
		if s.Committed > 0 && !seen[s.Committed] {
			seen[s.Committed] = true
			levels = append(levels, s.Committed)
		}
	}
	sort.Slice(levels, func(i, j int) bool { return levels[i] < levels[j] })
	out := []potLayer{}
	var prev int64
	for _, l := range levels {
		layer := (l - prev) * int64(t.countCommittedAtLeastLocked(l))
		if layer > 0 {
			seats := []int{}
			for _, s := range t.Seats {
				if s != nil && s.Committed >= l && !s.Folded {
					seats = append(seats, s.Index)
				}
			}
			out = append(out, potLayer{amount: layer, seats: seats})
		}
		prev = l
	}
	// 保险等额补差（理论上已覆盖）
	var sum int64
	for _, p := range out {
		sum += p.amount
	}
	if diff := total - sum; diff > 0 {
		seats := []int{}
		for _, s := range t.Seats {
			if s != nil && !s.Folded {
				seats = append(seats, s.Index)
			}
		}
		out = append(out, potLayer{amount: diff, seats: seats})
	}
	return out
}

func (t *Table) countCommittedAtLeastLocked(v int64) int {
	n := 0
	for _, s := range t.Seats {
		if s != nil && s.Committed >= v {
			n++
		}
	}
	return n
}

func cardsToStr(cs []int) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, CardStr(c))
	}
	return out
}

// ---------------- Tick：超时、机器人、续局 ----------------

func (t *Table) Tick(now time.Time, robotThink func(*Table, int) (string, int64)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return
	}
	t.vacateIdleLocked(now)
	t.refreshOfficialExpire(now)
	h := t.Hand
	switch t.State {
	case StateWaiting:
		// 卡池/房卡恢复后自动续局：Waiting 状态下满 2 人即重新开局
		if t.countActiveLocked() >= 2 {
			if t.Hand == nil || t.Hand.Done || now.After(t.NextHandAt) {
				t.startHandLocked()
			}
		}
	case StateBidding:
		if h == nil || h.Done {
			return
		}
		if now.After(h.Deadline) {
			// 到点还没喊完：机器人按牌力代喊，真人未操作视为不抢
			for _, s := range t.Seats {
				if s == nil || len(s.Cards) < 5 {
					continue
				}
				if s.Bid == 0 && s.IsRobot {
					v := NiuValue(s.Cards)
					bid := 1
					if v.Cat >= 2 {
						bid = v.Cat + 1
					}
					if bid > 4 {
						bid = 4
					}
					s.Bid = bid
				}
			}
			t.finishBiddingLocked()
		}
		return
	case StateBetting:
		if h == nil || h.Done {
			return
		}
		if now.After(h.Deadline) {
			// 到点未押注：机器人与真人一律按底注 1 倍
			t.startNiuRevealLocked()
		}
		return
	case StatePlaying:
		if h == nil || h.Done {
			// 手牌对象丢失/已结算但状态仍停留在 playing：收回等待态，交给下一 tick 正常开局
			t.Hand = nil
			t.State = StateWaiting
			t.NextHandAt = now
			t.notify()
			return
		}
		if t.Cfg.GameType == "niuniu" {
			if now.After(h.Deadline) {
				t.settleNiuNiuLocked()
			}
			return
		}
		if h.ToAct < 0 {
			return
		}
		s := t.Seats[h.ToAct]
		if s == nil || s.Folded || s.AllIn || s.WaitNext || s.SittingOut || len(s.Cards) == 0 {
			// 行动座位异常（已离桌/无牌/暂离/等待下一局）：把行动权推进到下一个可行动座位，无人可行动就直接结算
			t.resyncTurnLocked()
			return
		}
		deadline := s.Deadline
		if deadline.IsZero() {
			deadline = h.Deadline
		}
		// 断线玩家立即代为过牌/弃牌，避免牌局被 22 秒超时拖长
		if !s.IsRobot && !s.Connected {
			uid := s.UID
			if h.CurrentBet-s.Bet <= 0 {
				if err := t.actionLocked(uid, "check", 0); err == nil {
					return
				}
			}
			_ = t.actionLocked(uid, "fold", 0)
			return
		}
		if now.Before(deadline) {
			return
		}
		if s.IsRobot && robotThink != nil {
			act, amt := robotThink(t, s.Index)
			uid := s.UID
			// 在锁内直接执行，避免竞态
			err := t.actionLocked(uid, act, amt)
			if err != nil {
				_ = t.actionLocked(uid, "fold", 0)
			}
			return
		}
		// 人类超时：能过牌则过牌，否则弃牌
		uid := s.UID
		if h.CurrentBet-s.Bet <= 0 {
			if err := t.actionLocked(uid, "check", 0); err == nil {
				return
			}
		}
		_ = t.actionLocked(uid, "fold", 0)
	case StateInsurance:
		if h == nil {
			return
		}
		if h.Insurance == nil {
			t.resumeFromInsuranceLocked()
			return
		}
		if h.Insurance.Accepted || h.Insurance.Declined || now.After(h.Insurance.Deadline) {
			if !h.Insurance.Accepted {
				h.Insurance.Declined = true
			}
			t.resumeFromInsuranceLocked()
		}
	case StateShowdown:
		if now.After(t.NextHandAt) {
			t.State = StateWaiting
			if t.countActiveLocked() >= 2 {
				t.startHandLocked()
			} else {
				t.notify()
			}
		}
	}
}

// handBankerBid 当前手的庄家抢庄倍数（没有手牌时为 0）
func (t *Table) handBankerBid() int {
	if t.Hand == nil {
		return 0
	}
	return t.Hand.BankerBid
}

// ---------------- 退桌结算预览（官方 SquidExitDialog） ----------------

// ExitPreview 退桌/起身前的实时结算预览，全部金额单位「分」
type ExitPreview struct {
	Squid              bool  `json:"squid"`
	Standing           bool  `json:"standing"`        // 本轮尚未脱险（赏金待结算）
	InActiveHand       bool  `json:"in_active_hand"`  // 本手仍在牌局中且未弃牌
	InActiveHandNote   string `json:"in_active_hand_note,omitempty"`
	WinExitRestricted  bool  `json:"win_exit_restricted"`
	SessionHands       int   `json:"session_hands"`
	RequiredHands      int   `json:"required_hands"`
	RemainingHands     int   `json:"remaining_hands"`
	SessionBuyIn       int64 `json:"session_buy_in"`
	SessionProfit      int64 `json:"session_profit"`
	PenaltyPercent     int   `json:"win_exit_penalty_percent"`
	PenaltyFixed       int64 `json:"win_exit_penalty_fixed"`
	WinExitPenalty     int64 `json:"win_exit_penalty"`
	Recipients         int   `json:"recipients"`
	Expected           int64 `json:"expected"`
	Paid               int64 `json:"paid"`
	Refund             int64 `json:"refund"`
	Leaving            bool  `json:"leaving"`
}

// ExitPreview 计算某个座位的退桌预览；未入座返回 nil
func (t *Table) ExitPreview(uid uint32, leaving bool) *ExitPreview {
	t.mu.Lock()
	defer t.mu.Unlock()
	idx := t.seatOfLocked(uid)
	if idx < 0 {
		return nil
	}
	s := t.Seats[idx]
	h := t.Hand
	p := &ExitPreview{
		Squid:         t.Cfg.Squid,
		Leaving:       leaving,
		SessionHands:  s.Hands,
		SessionBuyIn:  s.Buyin,
		SessionProfit: s.Net,
	}
	req := t.Cfg.RequiredHands
	if req <= 0 {
		req = 20
	}
	p.RequiredHands = req
	if rem := req - s.Hands; rem > 0 {
		p.RemainingHands = rem
	}

	// 是否仍在有效手牌中（有牌且未弃牌）
	if h != nil && !h.Done && !s.Folded && len(s.Cards) > 0 && (t.State == StatePlaying || t.State == StateInsurance) {
		p.InActiveHand = true
		p.InActiveHandNote = "您正在参与本局对战且未弃牌。退桌将立即为您执行【弃牌】，本手已下注筹码不予退还，剩余桌上筹码将全额退回钱包。"
	}

	// 盈利提前离桌罚金（鱿鱼模式生效）
	if t.Cfg.Squid {
		p.Standing = s.Net < 0
		if s.Net > 0 && s.Hands < req {
			p.WinExitRestricted = true
			pct := t.Cfg.PenaltyPercent
			if pct <= 0 {
				pct = 5
			}
			p.PenaltyPercent = pct
			p.PenaltyFixed = t.Cfg.PenaltyFixed
			p.WinExitPenalty = s.Net*int64(pct)/100 + t.Cfg.PenaltyFixed
			if p.WinExitPenalty > s.Stack {
				p.WinExitPenalty = s.Stack
			}
		}
		// 尚未脱险：先结算本轮赏金，赔付给本轮盈利的其它玩家
		if s.Net < 0 {
			n := 0
			for _, x := range t.Seats {
				if x != nil && x.UID != uid && x.Net > 0 {
					n++
				}
			}
			p.Recipients = n
			if n > 0 {
				p.Expected = -s.Net
				p.Paid = p.Expected
			}
		}
	}

	// 预计退回钱包：桌上筹码 - 本手已下注 - 罚金 - 赏金（不为负）
	refund := s.Stack
	if p.InActiveHand && s.Committed > 0 {
		refund -= s.Committed
	}
	refund -= p.WinExitPenalty
	if p.Paid > 0 {
		refund -= p.Paid
	}
	if refund < 0 {
		if p.Paid > 0 {
			p.Paid += refund // 筹码不足时按可用筹码赔付，不产生欠款
			if p.Paid < 0 {
				p.Paid = 0
			}
		}
		refund = 0
	}
	p.Refund = refund
	return p
}

// ---------------- 快照 ----------------

func (t *Table) Snapshot(uid uint32) *TableStateDTO {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.snapshotLocked(uid)
}

func (t *Table) snapshotLocked(uid uint32) *TableStateDTO {
	h := t.Hand
	mySeat := t.seatOfLocked(uid)
	dto := &TableStateDTO{
		RoomNo: t.Cfg.RoomNo, Name: t.Cfg.Name, GameType: t.Cfg.GameType,
		SeatCount: t.Cfg.SeatCount, SB: t.Cfg.SB, BB: t.Cfg.BB, Ante: t.Cfg.Ante,
		BetLimit: t.Cfg.BetLimit, State: t.State, MySeat: mySeat,
		ServerTime: time.Now().UnixMilli(), ExpireAt: t.ExpireAt.UnixMilli(),
		InsuranceEnabled: t.Cfg.Insurance, MinBuyin: t.Cfg.MinBuyin, MaxBuyin: t.Cfg.MaxBuyin,
		IsOfficial: t.Cfg.IsOfficial, ClubID: t.Cfg.ClubID, HostUID: t.Cfg.HostUID, Observers: len(t.Observers),
		Version: t.Version, Winners: t.Winners, SidePots: t.SidePots,
		BankerBid: t.handBankerBid(),
		LastHandAt: t.lastHandEnd.UnixMilli(),
		LastHandID: t.lastHandID,
		NextHandAt: t.NextHandAt.UnixMilli(),
		ResultHoldMs: ResultHoldTime.Milliseconds(),
		ActionEv: t.LastAction, NewCards: t.LastNewCards, WinEv: t.LastWin,
	}
	if dto.Winners == nil {
		dto.Winners = []WinnerDTO{}
	}
	if dto.SidePots == nil {
		dto.SidePots = []SidePotDTO{}
	}
	if dto.NewCards == nil {
		dto.NewCards = []CardAnim{}
	}
	if dto.Seats == nil {
		dto.Seats = []SeatDTO{}
	}
	if dto.Board == nil {
		dto.Board = []string{}
	}
	if h != nil {
		dto.HandNo = h.No
		dto.HandID = h.DBID
		if dto.HandID == 0 {
			dto.HandID = t.lastHandID
		}
		dto.Board = cardsToStr(h.Board)
		dto.Street = h.Street
		dto.Pot = h.Pot
		dto.Dealer = h.Dealer
		dto.ToAct = h.ToAct
		dto.Deadline = h.Deadline.UnixMilli()
		dto.DeckLeft = len(h.Deck)
		dto.SeedHash = h.SeedHash
		if h.Done {
			dto.Seed = h.Seed
		}
		if h.Insurance != nil && !h.Insurance.Accepted && !h.Insurance.Declined {
			outs := make([]string, 0, len(h.Insurance.OutCards))
			for _, c := range h.Insurance.OutCards {
				outs = append(outs, CardStr(c))
			}
			dto.Insurance = &InsuranceOfferDTO{
				Seat: h.Insurance.Seat, EquityP: h.Insurance.Equity * 100,
				Premium: h.Insurance.Premium, Payout: h.Insurance.Payout,
				Deadline: h.Insurance.Deadline.UnixMilli(), PotSize: h.Insurance.PotSize,
				Cover: h.Insurance.Cover, BoardCount: h.Insurance.BoardCount,
				Outs: h.Insurance.Outs, OutCards: outs, Odds: h.Insurance.Odds,
				Plans: h.Insurance.Plans, PlanKey: h.Insurance.PlanKey,
			}
		}
	}
	dto.Seats = make([]SeatDTO, 0, len(t.Seats))
	for i, s := range t.Seats {
		if s == nil {
			continue
		}
		sd := SeatDTO{
			Seat: i, UID: s.UID, Nick: s.Nick, Avatar: s.Avatar, VIP: s.VIP, Robot: s.IsRobot,
			Stack: s.Stack, Bet: s.Bet, Committed: s.Committed, Folded: s.Folded, AllIn: s.AllIn,
			Acted: s.Acted, LastAction: s.LastAction, LastAmount: s.LastAmount,
			CardCount: len(s.Cards), Show: s.Showdown, Win: winFor(t.Winners, i),
			Insured: s.Insured, SittingOut: s.SittingOut, Hands: s.Hands, Net: s.Net,
			Buyin: s.Buyin, Connected: s.Connected, IsMe: s.UID == uid,
			VPIP: pctOf(s.VPIPHands, s.Hands), PFR: pctOf(s.PFRHands, s.Hands),
			MuckShown: s.MuckBuyers[uid], MuckFee: t.muckFeeLocked(),
		}
		if !s.Deadline.IsZero() && h != nil && h.ToAct == i && !h.Done {
			sd.Deadline = s.Deadline.UnixMilli()
		}
		// 手牌：自己的牌永远可见；摊牌阶段其它人可见；买过盖牌窥探的买家可见
		sd.Bid = s.Bid
		sd.Stake = s.Stake
		if len(s.Cards) > 0 {
			if s.UID == uid || s.Showdown || (s.Folded && s.MuckBuyers[uid]) {
				sd.Cards = cardsToStr(s.Cards)
			}
			if s.Showdown {
				sd.HandName = s.HandValue.CatName
				if t.Cfg.GameType == "niuniu" {
					sd.Mult = niuMults[s.HandValue.Cat&0xF]
				}
			}
		}
		dto.Seats = append(dto.Seats, sd)
	}
	if mySeat >= 0 && h != nil {
		if len(t.Seats[mySeat].Cards) > 0 {
			dto.MyCards = cardsToStr(t.Seats[mySeat].Cards)
		}
		if t.State == StatePlaying && h.ToAct == mySeat && !h.Done {
			dto.Actions = t.actionOptionsLocked(mySeat)
		}
	}
	return dto
}

// pctOf 百分比（0~100，无样本返回 0）
func pctOf(part, whole int) int {
	if whole <= 0 {
		return 0
	}
	return int(math.Round(float64(part) / float64(whole) * 100))
}

func winFor(ws []WinnerDTO, seat int) int64 {
	for _, w := range ws {
		if w.Seat == seat {
			return w.Amount
		}
	}
	return 0
}

func (t *Table) actionOptionsLocked(seat int) *ActionOptions {
	h := t.Hand
	s := t.Seats[seat]
	if h == nil || s == nil {
		return nil
	}
	toCall := h.CurrentBet - s.Bet
	if toCall < 0 {
		toCall = 0
	}
	opts := &ActionOptions{CanFold: true, CanCheck: toCall == 0, CanCall: toCall > 0, CallAmount: toCall}
	if opts.CallAmount > s.Stack {
		opts.CallAmount = s.Stack
	}
	minTo := h.CurrentBet + h.MinRaise
	maxTo := s.Bet + s.Stack
	if minTo > maxTo {
		minTo = maxTo
	}
	opts.MinRaiseTo = minTo
	opts.MaxRaiseTo = maxTo
	opts.CanRaise = s.Stack > toCall
	opts.CanAllIn = s.Stack > 0
	opts.IsRaiseOfBet = h.CurrentBet == 0
	return opts
}

func (t *Table) ClearAnim() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.LastAction = nil
	t.LastNewCards = nil
	t.LastWin = nil
}

func (t *Table) IsClosed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.closed
}

func (t *Table) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
	t.State = StateClosed
}

func (t *Table) CountSeated() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	n := 0
	for _, s := range t.Seats {
		if s != nil {
			n++
		}
	}
	return n
}

func (t *Table) ObserversCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.Observers)
}

// SetObserverMeta 记录围观者昵称与头像，供房间报表展示
func (t *Table) SetObserverMeta(uid uint32, nick, avatar string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ObsMeta == nil {
		t.ObsMeta = map[uint32]ObserverInfo{}
	}
	t.ObsMeta[uid] = ObserverInfo{UID: uid, Nick: nick, Avatar: avatar}
}

// Observers 列表（按加入顺序稳定输出）
func (t *Table) spectatorList() []ObserverInfo {
	out := make([]ObserverInfo, 0, len(t.ObsMeta))
	for uid, m := range t.ObsMeta {
		if !t.Observers[uid] {
			continue
		}
		if t.seatOfLocked(uid) >= 0 {
			continue
		}
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UID < out[j].UID })
	return out
}

func (t *Table) AddObserver(uid uint32, on bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if on {
		t.Observers[uid] = true
	} else {
		delete(t.Observers, uid)
	}
}

func (t *Table) EmptyTime() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.idleAt
}

// AcceptInsurance 玩家接受保险（保费由钱包侧记账，不动筹码）
func (t *Table) AcceptInsurance(uid uint32, accept bool) (bool, int64, error) {
	return t.AcceptInsurancePlan(uid, accept, "")
}

// AcceptInsurancePlan 带方案键的投保（官方「选择投保方案」）
func (t *Table) AcceptInsurancePlan(uid uint32, accept bool, planKey string) (bool, int64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	h := t.Hand
	if h == nil || h.Insurance == nil || h.Insurance.Declined || h.Insurance.Accepted {
		return false, 0, ErrState
	}
	idx := t.seatOfLocked(uid)
	if idx != h.Insurance.Seat {
		return false, 0, ErrBadAction
	}
	if !accept {
		h.Insurance.Declined = true
		t.notify()
		return false, 0, nil
	}
	s := t.Seats[idx]
	if s == nil {
		return false, 0, ErrBadAction
	}
	if planKey != "" {
		for _, p := range h.Insurance.Plans {
			if p.Key == planKey {
				h.Insurance.Premium = p.Premium
				h.Insurance.Payout = p.Payout
				h.Insurance.PlanKey = p.Key
				break
			}
		}
	}
	h.Insurance.Accepted = true
	s.Insured = true
	s.InsPremium = h.Insurance.Premium
	s.InsPayout = h.Insurance.Payout
	t.notify()
	return true, h.Insurance.Premium, nil
}

// muckFeeLocked 盖牌窥探价格（2 个大盲，最低 1 筹码）
func (t *Table) muckFeeLocked() int64 {
	fee := t.Cfg.BB * 2
	if fee < 100 {
		fee = 100
	}
	return fee
}

// MuckFee 对外暴露窥探价格
func (t *Table) MuckFee() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.muckFeeLocked()
}

// RevealMuck 付费查看某位已弃牌玩家的底牌（官方「查看对手盖牌」），返回牌面与扣费
func (t *Table) RevealMuck(uid uint32, seat int) ([]string, int64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.Hand == nil {
		return nil, 0, ErrState
	}
	if seat < 0 || seat >= len(t.Seats) {
		return nil, 0, ErrBadAction
	}
	s := t.Seats[seat]
	if s == nil || !s.Folded || len(s.Cards) == 0 {
		return nil, 0, ErrBadAction
	}
	if s.UID == uid {
		return nil, 0, ErrBadAction
	}
	if t.seatOfLocked(uid) < 0 {
		return nil, 0, ErrNotSeated
	}
	if s.MuckBuyers == nil {
		s.MuckBuyers = map[uint32]bool{}
	}
	fee := int64(0)
	if !s.MuckBuyers[uid] {
		fee = t.muckFeeLocked()
		s.MuckBuyers[uid] = true
		t.notify()
	}
	return cardsToStr(s.Cards), fee, nil
}

// InsuranceResult 结算后返回保险赔付信息，供钱包记账
func (t *Table) InsuranceResult() (uid uint32, premium, payout int64, ok bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	h := t.Hand
	if h == nil || h.Insurance == nil || !h.Insurance.Accepted {
		return 0, 0, 0, false
	}
	s := t.Seats[h.Insurance.Seat]
	if s == nil {
		return 0, 0, 0, false
	}
	won := false
	for _, w := range h.Winners {
		if w.Seat == s.Index {
			won = true
		}
	}
	if won {
		return s.UID, h.Insurance.Premium, 0, true
	}
	return s.UID, h.Insurance.Premium, h.Insurance.Payout, true
}

// vacateIdleLocked 断线人类玩家 120 秒未回来自动离桌结算
func (t *Table) vacateIdleLocked(now time.Time) {
	for i, s := range t.Seats {
		if s == nil || s.IsRobot || s.Connected || s.DisconnAt.IsZero() {
			continue
		}
		if now.Sub(s.DisconnAt) < 120*time.Second {
			continue
		}
		if t.Hand != nil && !t.Hand.Done && len(s.Cards) > 0 && !s.Folded {
			continue
		}
		stack := s.Stack
		uid := s.UID
		// 鱿鱼模式：断线自动离座同样要结清罚金 / 赏金
		if t.Cfg.Squid && !s.SittingOut {
			stack, _, _ = t.settleSquidExitLocked(i, stack)
		}
		t.Seats[i] = nil
		t.idleAt = now
		t.logf("room %d auto vacate uid=%d stack=%d", t.Cfg.RoomNo, uid, stack)
		if t.OnLeave != nil && stack > 0 {
			go t.OnLeave(t, uid, stack)
		}
		if t.countActiveLocked() < 2 && t.State == StatePlaying {
			t.finishHandLocked()
		}
		t.notify()
	}
}

// startNiuNiuLocked 牛牛一手：发 5 张、轮庄、庄闲对赌结算
func (t *Table) startNiuNiuLocked() {
	active := []int{}
	for i, s := range t.Seats {
		if s != nil && !s.SittingOut && s.Stack > 0 {
			active = append(active, i)
		}
	}
	if len(active) < 2 {
		return
	}
	t.HandCounter++
	h := &Hand{No: t.HandCounter, Street: 0, Dealer: t.dealer, StartedAt: time.Now()}
	h.Seed = NewSeed(fmt.Sprintf("%d|%d|%d|%d", t.Cfg.RoomNo, t.HandCounter, time.Now().UnixNano(), t.rnd.next()))
	h.SeedHash = SeedHash(h.Seed)
	h.Deck = Shuffle(NewDeck(false), h.Seed, t.HandCounter)
	t.dealer = t.nextSeatLocked(h.Dealer, func(s *Seat) bool { return s != nil && !s.SittingOut && s.Stack > 0 })
	h.Dealer = t.dealer

	for _, s := range t.Seats {
		if s == nil {
			continue
		}
		s.Cards = nil
		s.Bet = 0
		s.Committed = 0
		s.Folded = false
		s.AllIn = false
		s.Acted = false
		s.LastAction = ""
		s.LastAmount = 0
		s.Showdown = false
		s.HandValue = HandValue{}
		s.Insured = false
		s.WaitNext = false
		s.Deadline = time.Time{}
	}
	newCards := []CardAnim{}
	for k := 0; k < 5; k++ {
		for _, i := range active {
			s := t.Seats[i]
			c := NextCard(&h.Deck)
			s.Cards = append(s.Cards, c)
			newCards = append(newCards, CardAnim{Seat: i, Card: CardStr(c), Index: len(s.Cards) - 1, Street: 0})
		}
	}
	for _, s := range t.Seats {
		if s != nil {
			s.Bid = 0
		}
	}
	t.Hand = h
	t.Winners = nil
	t.SidePots = nil
	t.LastWin = nil
	t.LastAction = nil
	t.LastNewCards = newCards
	h.ToAct = -1
	if t.Cfg.NiuBankerMode == "fixed" {
		// 固定庄：按钮位坐庄，接着走闲家押注
		h.BankerBid = 1
		for _, s := range t.Seats {
			if s != nil {
				s.Stake = 0
			}
		}
		h.Deadline = time.Now().Add(12 * time.Second)
		t.State = StateBetting
	} else {
		// 抢庄：先给 6 秒喊庄，再由最高倍数坐庄
		h.Deadline = time.Now().Add(8 * time.Second)
		t.State = StateBidding
	}
	t.logf("room %d niuniu hand %d banker=%d mode=%s seats=%v", t.Cfg.RoomNo, h.No, h.Dealer, t.Cfg.NiuBankerMode, active)
	t.notify()
}

// settleNiuNiuLocked 开牌并结算，Tick 到点调用
// NiuBid 玩家抢庄：mult 为 1–4 倍，0 表示不抢
func (t *Table) NiuBid(uid uint32, mult int) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.State != StateBidding || t.Hand == nil || t.Hand.Done {
		return ErrState
	}
	idx := t.seatOfLocked(uid)
	if idx < 0 {
		return ErrNotSeated
	}
	if mult < 0 || mult > 4 {
		return ErrBadAction
	}
	s := t.Seats[idx]
	if s == nil || len(s.Cards) < 5 {
		return ErrState
	}
	s.Bid = mult
	s.LastAction = "bid"
	s.LastAmount = int64(mult)
	t.LastAction = &ActionAnim{Seat: idx, Action: "bid", Amount: int64(mult), Street: 0, TS: time.Now().UnixMilli()}
	// 全员喊完立即定庄
	pending := false
	for _, x := range t.Seats {
		if x == nil || len(x.Cards) < 5 {
			continue
		}
		if x.Bid == 0 && !x.IsRobot {
			pending = true
			break
		}
	}
	if !pending {
		t.finishBiddingLocked()
		return nil
	}
	t.notify()
	return nil
}

// NiuStake 牛牛闲家押注：mult 取 1/2/5/10（底注倍数），庄家不押注
func (t *Table) NiuStake(uid uint32, mult int) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.State != StateBetting || t.Hand == nil || t.Hand.Done {
		return ErrState
	}
	idx := t.seatOfLocked(uid)
	if idx < 0 {
		return ErrNotSeated
	}
	if idx == t.Hand.Dealer {
		return errors.New("庄家无需押注")
	}
	switch mult {
	case 1, 2, 5, 10:
	default:
		return ErrBadAction
	}
	s := t.Seats[idx]
	if s == nil || len(s.Cards) < 5 {
		return ErrState
	}
	s.Stake = mult
	s.LastAction = "stake"
	s.LastAmount = int64(mult)
	t.LastAction = &ActionAnim{Seat: idx, Action: "stake", Amount: int64(mult), Street: 0, TS: time.Now().UnixMilli()}
	// 闲家全部押完立即开牌
	for _, x := range t.Seats {
		if x == nil || len(x.Cards) < 5 || x.Index == t.Hand.Dealer {
			continue
		}
		if x.Stake == 0 && !x.IsRobot {
			t.notify()
			return nil
		}
	}
	t.startNiuRevealLocked()
	return nil
}

// startNiuRevealLocked 押注结束→开牌（留 2.6 秒给客户端翻牌动画）
func (t *Table) startNiuRevealLocked() {
	h := t.Hand
	if h == nil {
		return
	}
	for _, s := range t.Seats {
		if s != nil && len(s.Cards) == 5 && s.Index != h.Dealer && s.Stake == 0 {
			s.Stake = 1 // 未押注默认底注
		}
	}
	h.Deadline = time.Now().Add(2600 * time.Millisecond)
	t.State = StatePlaying
	t.notify()
}

// finishBiddingLocked 定庄：最高倍数坐庄（同为最高时按座位顺序取最近的），全不抢则由按钮位坐庄
func (t *Table) finishBiddingLocked() {
	h := t.Hand
	if h == nil {
		return
	}
	bestSeat, bestBid := -1, 0
	start := h.Dealer
	for k := 0; k < len(t.Seats); k++ {
		i := (start + k) % len(t.Seats)
		s := t.Seats[i]
		if s == nil || len(s.Cards) < 5 {
			continue
		}
		if s.Bid > bestBid {
			bestBid = s.Bid
			bestSeat = i
		}
	}
	if bestSeat < 0 {
		bestSeat = start
		bestBid = 1
		if s := t.Seats[start]; s == nil || len(s.Cards) < 5 {
			for i, s2 := range t.Seats {
				if s2 != nil && len(s2.Cards) == 5 {
					bestSeat = i
					break
				}
			}
		}
	}
	h.Dealer = bestSeat
	h.BankerBid = bestBid
	for _, s := range t.Seats {
		if s != nil {
			s.Stake = 0
		}
	}
	h.Deadline = time.Now().Add(12 * time.Second)
	t.State = StateBetting
	t.LastAction = &ActionAnim{Seat: bestSeat, Action: "banker", Amount: int64(bestBid), Street: 0, TS: time.Now().UnixMilli()}
	t.logf("room %d niuniu hand %d 定庄 seat=%d bid=%d", t.Cfg.RoomNo, h.No, bestSeat, bestBid)
	t.notify()
}

func (t *Table) settleNiuNiuLocked() {
	h := t.Hand
	if h == nil || h.Done {
		return
	}
	banker := t.Seats[h.Dealer]
	if banker == nil {
		t.finishNiuNiuLocked()
		return
	}
	base := t.Cfg.BB
	if base <= 0 {
		base = 10
	}
	type stake struct {
		seat int
		niu  NiuResult
	}
	list := []stake{}
	for _, s := range t.Seats {
		if s == nil || len(s.Cards) < 5 {
			continue
		}
		v := NiuValue(s.Cards)
		s.HandValue = HandValue{Score: v.Score, Cat: v.Cat, CatName: v.Name, Best: s.Cards}
		s.Showdown = true
		list = append(list, stake{seat: s.Index, niu: v})
	}
	var bankerNiu NiuResult
	for _, it := range list {
		if it.seat == h.Dealer {
			bankerNiu = it.niu
		}
	}
	net := map[int]int64{}
	var transferTotal int64
	for _, it := range list {
		if it.seat == h.Dealer || it.niu.Cat < 0 {
			continue
		}
		player := t.Seats[it.seat]
		isPlayerWin := it.niu.Score > bankerNiu.Score
		m := it.niu.Mult
		if !isPlayerWin {
			m = bankerNiu.Mult
		}
		if m < 1 {
			m = 1
		}
		// 庄家抢庄倍数 × 闲家押注倍数，一起放大本手输赢
		if h.BankerBid > 1 {
			m *= h.BankerBid
		}
		stake := player.Stake
		if stake <= 0 {
			stake = 1
		}
		amount := base * int64(m) * int64(stake)
		if isPlayerWin {
			if banker.Stack < amount {
				amount = banker.Stack
			}
			banker.Stack -= amount
			player.Stack += amount
			net[banker.Index] -= amount
			net[player.Index] += amount
		} else {
			if player.Stack < amount {
				amount = player.Stack
			}
			player.Stack -= amount
			banker.Stack += amount
			net[player.Index] -= amount
			net[banker.Index] += amount
		}
		transferTotal += amount
	}
	h.Pot = transferTotal
	h.Done = true
	h.Showdown = true
	for _, s := range t.Seats {
		if s == nil || len(s.Cards) < 5 {
			continue
		}
		s.Hands++
		s.Net += net[s.Index]
		s.Committed = 0
		s.Stake = 0
	}
	// 汇总赢家
	h.Winners = []WinnerDTO{}
	for _, s := range t.Seats {
		if s == nil {
			continue
		}
		if n := net[s.Index]; n > 0 {
			h.Winners = append(h.Winners, WinnerDTO{
				Seat: s.Index, UID: s.UID, Nick: s.Nick, Amount: n, HandName: s.HandValue.CatName, Cards: cardsToStr(s.Cards),
			})
		}
	}
	sort.Slice(h.Winners, func(i, j int) bool { return h.Winners[i].Amount > h.Winners[j].Amount })
	t.Winners = h.Winners
	if len(h.Winners) > 0 {
		t.LastWin = &WinAnim{Seat: h.Winners[0].Seat, Amount: h.Winners[0].Amount, Pot: transferTotal}
	}
	// 结算明细
	results := []HandResult{}
	for _, s := range t.Seats {
		if s == nil || len(s.Cards) < 5 {
			continue
		}
		res := HandResult{
			Seat: s.Index, UID: s.UID, Nick: s.Nick, Cards: cardsToStr(s.Cards),
			HandName: s.HandValue.CatName, Win: net[s.Index], Net: net[s.Index],
			Best: cardsToStr(s.Cards),
		}
		results = append(results, res)
	}
	h.Results = results
	t.State = StateShowdown
	t.NextHandAt = time.Now().Add(ResultHoldTime + 2*time.Second)
	if t.OnHandEnd != nil {
		hh := h
		tb := t
		go func() { tb.OnHandEnd(tb, hh) }()
	}
	t.notify()
}

func (t *Table) finishNiuNiuLocked() {
	t.State = StateShowdown
	t.NextHandAt = time.Now().Add(ResultHoldTime)
	t.notify()
}

// PeersOf 返回同桌其它玩家 uid 列表
func (t *Table) PeersOf(uid uint32) []uint32 {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := []uint32{}
	for _, s := range t.Seats {
		if s != nil && s.UID != uid {
			out = append(out, s.UID)
		}
	}
	return out
}

// SetSitOut 暂离/回座；牌局进行中暂离视为自动弃牌，避免整桌等一个不在的人
func (t *Table) SetSitOut(uid uint32, out bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	i := t.seatOfLocked(uid)
	if i < 0 {
		return
	}
	s := t.Seats[i]
	s.SittingOut = out
	if out && t.Hand != nil && !t.Hand.Done && t.State == StatePlaying && !s.Folded && len(s.Cards) > 0 {
		_ = t.actionLocked(uid, "fold", 0)
		return
	}
	t.notify()
}

// resyncTurnLocked 行动座位不可行动时推进行动权；全桌都不可行动则直接结算，避免牌桌永久卡死
func (t *Table) resyncTurnLocked() {
	h := t.Hand
	if h == nil || h.Done {
		return
	}
	start := h.ToAct
	if start < 0 || start >= len(t.Seats) {
		start = h.Dealer
	}
	next := t.nextSeatLocked(start, func(s *Seat) bool {
		return s != nil && !s.Folded && !s.AllIn && !s.WaitNext && !s.SittingOut && len(s.Cards) > 0
	})
	if next < 0 {
		t.logf("room %d hand %d 无可行动座位，直接结算（自愈）", t.Cfg.RoomNo, h.No)
		t.finishHandLocked()
		t.notify()
		return
	}
	h.ToAct = next
	h.Deadline = t.actorDeadlineLocked(next)
	t.Seats[next].Deadline = h.Deadline
	t.logf("room %d hand %d 行动座位自愈 -> seat %d", t.Cfg.RoomNo, h.No, next)
	t.notify()
}

// Report 房间报表：总带入、手数、战绩排行
func (t *Table) Report() map[string]any {
	t.mu.Lock()
	defer t.mu.Unlock()
	rows := []map[string]any{}
	for _, s := range t.Seats {
		if s == nil {
			continue
		}
		vpip := pctOf(s.VPIPHands, s.Hands)
		rows = append(rows, map[string]any{
			"seat": s.Index, "uid": s.UID, "nick": s.Nick, "avatar": s.Avatar,
			"buyin": s.Buyin, "hands": s.Hands, "vpip": vpip, "net": s.Net, "stack": s.Stack,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i]["net"].(int64) > rows[j]["net"].(int64)
	})
	return map[string]any{
		"roomNo": t.Cfg.RoomNo, "name": t.Cfg.Name, "players": rows,
		"handNo": func() uint64 {
			if t.Hand != nil {
				return t.Hand.No
			}
			return 0
		}(),
		"expireAt": t.ExpireAt.UnixMilli(),
		"spectators": t.spectatorList(),
		"insuranceHarvest": t.insuranceHarvest(),
	}
}

// insuranceHarvest 本桌保险理赔累计（报表「保险理赔」页签用）
func (t *Table) insuranceHarvest() []map[string]any {
	out := []map[string]any{}
	for _, s := range t.Seats {
		if s == nil || s.InsPayout == 0 {
			continue
		}
		out = append(out, map[string]any{
			"seat": s.Index, "uid": s.UID, "nick": s.Nick,
			"premium": s.InsPremium, "payout": s.InsPayout,
		})
	}
	return out
}
