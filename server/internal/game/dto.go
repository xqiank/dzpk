package game

// DTO：下发给客户端的数据结构（严格按座位视角脱敏）

type SeatDTO struct {
	Seat       int      `json:"seat"`
	UID        uint32   `json:"uid"`
	Nick       string   `json:"nick"`
	Avatar     string   `json:"avatar"`
	VIP        int      `json:"vip"`
	Robot      bool     `json:"robot"`
	Stack      int64    `json:"stack"`
	Bet        int64    `json:"bet"`
	Committed  int64    `json:"committed"`
	Folded     bool     `json:"folded"`
	AllIn      bool     `json:"allin"`
	Acted      bool     `json:"acted"`
	LastAction string   `json:"lastAction"`
	LastAmount int64    `json:"lastAmount"`
	Deadline   int64    `json:"deadline"`
	Cards      []string `json:"cards"`
	CardCount  int      `json:"cardCount"`
	VPIP       int      `json:"vpip"`
	PFR        int      `json:"pfr"`
	MuckShown  bool     `json:"muckShown"`
	MuckFee    int64    `json:"muckFee"`
	Show       bool     `json:"show"`
	HandName   string   `json:"handName,omitempty"`
	Win        int64    `json:"win"`
	Insured    bool     `json:"insured"`
	Mult       int      `json:"mult"`
	SittingOut bool     `json:"sittingOut"`
	Bid        int      `json:"bid"`
	Stake      int      `json:"stake"`
	Hands      int      `json:"hands"`
	Net        int64    `json:"net"`
	Buyin      int64    `json:"buyin"`
	Connected  bool     `json:"connected"`
	IsMe       bool     `json:"isMe"`
}

type ActionOptions struct {
	CanCheck  bool  `json:"canCheck"`
	CanFold   bool  `json:"canFold"`
	CanCall   bool  `json:"canCall"`
	CallAmount int64 `json:"callAmount"`
	CanRaise  bool  `json:"canRaise"`
	MinRaiseTo int64 `json:"minRaiseTo"`
	MaxRaiseTo int64 `json:"maxRaiseTo"`
	CanAllIn  bool  `json:"canAllIn"`
	IsRaiseOfBet bool `json:"isRaiseOfBet"`
}

type WinnerDTO struct {
	Seat     int      `json:"seat"`
	UID      uint32   `json:"uid"`
	Nick     string   `json:"nick"`
	Amount   int64    `json:"amount"`
	HandName string   `json:"handName"`
	Cards    []string `json:"cards"`
	Best     []string `json:"best"`
}

type SidePotDTO struct {
	Amount int64 `json:"amount"`
	Seats  []int `json:"seats"`
}

type InsuranceOfferDTO struct {
	Seat    int   `json:"seat"`
	EquityP float64 `json:"equityP"` // 百分比
	Premium int64 `json:"premium"`
	Payout  int64 `json:"payout"`
	Deadline int64 `json:"deadline"`
	PotSize int64 `json:"potSize"`
	Cover   int64 `json:"cover"`
	BoardCount int `json:"boardCount"`
	Outs    int    `json:"outs"`
	OutCards []string `json:"outCards"`
	Odds    float64 `json:"odds"`
	Plans   []InsPlan `json:"plans"`
	PlanKey string  `json:"planKey"`
	ClaimFeeBps int `json:"claimFeeBps"`
}

type TableStateDTO struct {
	RoomNo     int          `json:"roomNo"`
	Name       string       `json:"name"`
	GameType   string       `json:"gameType"`
	SeatCount  int          `json:"seatCount"`
	SB         int64        `json:"sb"`
	BB         int64        `json:"bb"`
	Ante       int64        `json:"ante"`
	BetLimit   string       `json:"betLimit"`
	Pot        int64        `json:"pot"`
	Board      []string     `json:"board"`
	Street     int          `json:"street"`
	HandNo     uint64       `json:"handNo"`
	HandID     int64        `json:"handId"`
	Dealer     int          `json:"dealer"`
	Banker     int          `json:"banker"`
	ToAct      int          `json:"toAct"`
	State      string       `json:"state"`
	Deadline   int64        `json:"deadline"`
	ExpireAt   int64        `json:"expireAt"`
	MySeat     int          `json:"mySeat"`
	MyCards    []string     `json:"myCards"`
	Actions    *ActionOptions `json:"actions"`
	Seats      []SeatDTO    `json:"seats"`
	Winners    []WinnerDTO  `json:"winners"`
	SidePots   []SidePotDTO `json:"sidePots"`
	Insurance  *InsuranceOfferDTO `json:"insurance"`
	ServerTime int64        `json:"serverTime"`
	DeckLeft   int          `json:"deckLeft"`
	InsuranceEnabled bool   `json:"insuranceEnabled"`
	MinBuyin   int64        `json:"minBuyin"`
	MaxBuyin   int64        `json:"maxBuyin"`
	IsOfficial bool         `json:"isOfficial"`
	ClubID     int64        `json:"clubId"`
	HostUID    uint32       `json:"hostUid"`
	Observers  int          `json:"observers"`
	SeedHash   string       `json:"seedHash"`
	Seed       string       `json:"seed"`
	NewCards   []CardAnim   `json:"newCards"`
	ActionEv   *ActionAnim  `json:"actionEv"`
	WinEv      *WinAnim     `json:"winEv"`
	Version    int64        `json:"version"`
	NextHandAt int64        `json:"nextHandAt"`
	ResultHoldMs int64      `json:"resultHoldMs"`
	BankerBid   int        `json:"bankerBid"`
	LastHandAt int64        `json:"lastHandAt"`
	LastHandID int64        `json:"lastHandId"`
}

type CardAnim struct {
	Seat   int    `json:"seat"`   // -1 公共牌
	Card   string `json:"card"`
	Index  int    `json:"index"`
	Street int    `json:"street"`
}

type ActionAnim struct {
	Seat   int    `json:"seat"`
	Action string `json:"action"`
	Amount int64  `json:"amount"`
	Street int    `json:"street"`
	TS     int64  `json:"ts"`
}

type WinAnim struct {
	Seat   int   `json:"seat"`
	Amount int64 `json:"amount"`
	Pot    int64 `json:"pot"`
}

// LobbyRoomDTO 大厅房间卡片
type LobbyRoomDTO struct {
	RoomNo    int    `json:"roomNo"`
	Name      string `json:"name"`
	GameType  string `json:"gameType"`
	GameName  string `json:"gameName"`
	SeatCount int    `json:"seatCount"`
	Online    int    `json:"online"`
	SB        int64  `json:"sb"`
	BB        int64  `json:"bb"`
	MinBuyin  int64  `json:"minBuyin"`
	MaxBuyin  int64  `json:"maxBuyin"`
	VPIP      int    `json:"vpip"`
	ExpireAt  int64  `json:"expireAt"`
	IsOfficial bool  `json:"isOfficial"`
	ClubID    int64  `json:"clubId"`
	HostUID   uint32 `json:"hostUid"`
	CardCost  int    `json:"cardCost"`
	Insurance bool   `json:"insurance"`
	Status    string `json:"status"`
	Sort      int    `json:"sort"`
	BetLimit  string `json:"betLimit"`
	LastHandAt int64 `json:"lastHandAt"`
	HandNo     uint64 `json:"handNo"`
	Pot        int64  `json:"pot"`
	AvgPot     int64  `json:"avgPot"`
	NeedsCode  bool   `json:"needsCode"`
}
