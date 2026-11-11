package game

import (
	"log"
	"sync"
	"time"
)

// Engine 管理所有牌桌与机器人
type Engine struct {
	mu     sync.RWMutex
	tables map[int]*Table
	bots   map[uint32]BotProfile

	OnChange func(*Table)
	OnHandEnd func(*Table, *Hand)
	OnLeave func(*Table, uint32, int64)
	OnHandStart func(*Table) bool
	OnRobotSeat func(*Table, SitRequest)

	stop chan struct{}
}

func NewEngine() *Engine {
	return &Engine{
		tables: map[int]*Table{},
		bots:   map[uint32]BotProfile{},
		stop:   make(chan struct{}),
	}
}

func (e *Engine) RegisterBot(p BotProfile) {
	e.mu.Lock()
	e.bots[p.UID] = p
	e.mu.Unlock()
}

func (e *Engine) Bot(uid uint32) (BotProfile, bool) {
	e.mu.RLock()
	p, ok := e.bots[uid]
	e.mu.RUnlock()
	return p, ok
}

func (e *Engine) Table(roomNo int) *Table {
	e.mu.RLock()
	t := e.tables[roomNo]
	e.mu.RUnlock()
	return t
}

// EnsureTable 按房间配置创建/更新桌台
func (e *Engine) EnsureTable(cfg RoomCfg) *Table {
	e.mu.Lock()
	t, ok := e.tables[cfg.RoomNo]
	if !ok {
		t = NewTable(cfg)
		t.OnChange = func(tb *Table) { e.broadcast(tb) }
		t.OnHandEnd = func(tb *Table, h *Hand) {
			if e.OnHandEnd != nil {
				e.OnHandEnd(tb, h)
			}
		}
		t.OnLeave = func(tb *Table, uid uint32, stack int64) {
			if e.OnLeave != nil {
				e.OnLeave(tb, uid, stack)
			}
		}
		t.OnHandStart = func(tb *Table) bool {
			if e.OnHandStart == nil {
				return true
			}
			return e.OnHandStart(tb)
		}
		t.OnLog = func(f string, a ...any) { log.Printf(f, a...) }
		e.tables[cfg.RoomNo] = t
		e.mu.Unlock()
		return t
	}
	e.mu.Unlock()
	return t
}

func (e *Engine) RemoveTable(roomNo int) {
	e.mu.Lock()
	t := e.tables[roomNo]
	delete(e.tables, roomNo)
	e.mu.Unlock()
	if t != nil {
		t.Close()
	}
}

func (e *Engine) Tables() []*Table {
	e.mu.RLock()
	out := make([]*Table, 0, len(e.tables))
	for _, t := range e.tables {
		out = append(out, t)
	}
	e.mu.RUnlock()
	return out
}

func (e *Engine) broadcast(t *Table) {
	if e.OnChange != nil {
		e.OnChange(t)
	}
}

// StartTick 全局定时器：120ms 精度，负责超时、机器人行动、续局
func (e *Engine) StartTick() {
	go func() {
		tk := time.NewTicker(120 * time.Millisecond)
		defer tk.Stop()
		for {
			select {
			case <-e.stop:
				return
			case now := <-tk.C:
				tables := e.Tables()
				for _, t := range tables {
					t.Tick(now, func(tb *Table, seat int) (string, int64) {
						s := tb.Seats[seat]
						if s == nil {
							return "fold", 0
						}
						p, ok := e.Bot(s.UID)
						if !ok {
							p = BotProfile{Style: BotStyles["tag"], Aggression: 50, Bluff: 15, VPIP: 30}
						}
						return RobotDecide(tb, p, seat)
					})
				}
			}
		}
	}()
}

func (e *Engine) Stop() {
	select {
	case <-e.stop:
	default:
		close(e.stop)
	}
}

// RobotSit 让机器人入座
func (e *Engine) RobotSit(t *Table, req SitRequest) (int, error) {
	seat, _, err := t.Sit(req)
	return seat, err
}
