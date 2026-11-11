package seed

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"u5poker/internal/auth"
	"u5poker/internal/store"
)

// coursesJSON 官方教学中心 28 门课程（分组/时长/热度/封面全部按官方抓取）
//
//go:embed courses.json
var coursesJSON []byte

// Run 幂等初始化：房间、机器人、课程、配置、后台账号
func Run(st *store.Store) error {
	if err := seedAdmin(st); err != nil {
		return err
	}
	if err := seedRooms(st); err != nil {
		return err
	}
	if err := seedRobots(st); err != nil {
		return err
	}
	if err := seedCourses(st); err != nil {
		return err
	}
	if err := seedJackpot(st); err != nil {
		return err
	}
	return seedSettings(st)
}

// seedJackpot 奖池公示历史（空表时灌入一批真实结构的样例，后台/前端都读真数据）
func seedJackpot(st *store.Store) error {
	var n int
	if err := st.DB.QueryRow("SELECT COUNT(1) FROM jackpot_winners").Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		rows := []struct {
			nick, label string
			pct         int
			reward      int64
			cards       string
			room        string
			roomNo      int
		}{
			{"德扑老王", "皇家同花顺", 10, 110370, "A♠K♠Q♠J♠10♠", "大师桌 5/10", 218},
			{"鱿鱼老六", "同花顺", 5, 55185, "9♥8♥7♥6♥5♥", "鱿鱼暴击场 0.5/1", 216},
			{"河牌刺客", "四条", 2, 22074, "K♦K♣K♠K♥9♠", "天天休闲桌 0.2/0.4", 215},
		}
		for _, r := range rows {
			_, _ = st.DB.Exec(`INSERT INTO jackpot_winners(uid,nickname,tier,hand_label,pct,reward_amount,pool_after,cards,room_no,room_name,hand_id)
				VALUES(?,?,?,?,?,?,?,?,?,?,0)`, 0, r.nick, r.label, r.label, r.pct, r.reward, 1077736, r.cards, r.roomNo, r.room)
		}
	}
	if err := st.DB.QueryRow("SELECT COUNT(1) FROM jackpot_penalties").Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		rows := []struct {
			nick, reason, room string
			amount             int64
		}{
			{"wx_8823", "伙牌双簧：多账号同桌互让底池，U5 卫士同 IP 关联查实", "天天休闲桌 0.2/0.4", 86420},
			{"tb_1907", "外挂脚本：客户端异常操作频率，牌局胜率显著偏离", "鱿鱼暴击场 0.5/1", 152300},
		}
		for _, r := range rows {
			_, _ = st.DB.Exec(`INSERT INTO jackpot_penalties(uid,nickname,reason,confiscated,room_name)
				VALUES(0,?,?,?,?)`, r.nick, r.reason, r.amount, r.room)
		}
	}
	return nil
}

func seedAdmin(st *store.Store) error {
	var n int
	_ = st.DB.QueryRow("SELECT COUNT(1) FROM admin_users").Scan(&n)
	if n > 0 {
		return nil
	}
	_, err := st.DB.Exec("INSERT INTO admin_users(username,password_hash,role) VALUES(?,?,'super')",
		"admin", auth.HashPassword("admin888"))
	if err == nil {
		log.Println("[seed] admin created: admin / admin888")
	}
	return err
}

type roomSeed struct {
	no     int
	name   string
	gt     string
	seats  int
	sb, bb int64
	minBuy int64
	maxBuy int64
	limit  string
	sort   int
	cost   int
}

func seedRooms(st *store.Store) error {
	var n int
	_ = st.DB.QueryRow("SELECT COUNT(1) FROM rooms").Scan(&n)
	if n > 0 {
		return nil
	}
	rooms := []roomSeed{
		{217, "长牌高额 1/2 · 9人", "nlh", 9, 100, 200, 10000, 40000, "nl", 9, 0},
		{216, "长牌进阶 0.5/1 · 6人", "nlh", 6, 50, 100, 5000, 20000, "nl", 8, 0},
		{215, "长牌休闲 0.2/0.4 · 6人", "nlh", 6, 20, 40, 2000, 8000, "nl", 7, 0},
		{214, "长牌微盲 0.1/0.2 · 6人", "nlh", 6, 10, 20, 1000, 4000, "nl", 6, 0},
		{223, "奥马哈休闲 0.2/0.4 · 6人", "omaha", 6, 20, 40, 2000, 8000, "pl", 5, 0},
		{224, "奥马哈微盲 0.1/0.2 · 6人", "omaha", 6, 10, 20, 1000, 4000, "pl", 4, 0},
		{219, "短牌休闲 0.2/0.4 · 6人", "short6", 6, 20, 40, 2000, 8000, "nl", 4, 0},
		{218, "短牌微盲 0.1/0.2 · 6人", "short6", 6, 10, 20, 1000, 4000, "nl", 3, 0},
		{220, "短牌微盲 0.1/0.2 · 8人", "short6", 8, 10, 20, 1000, 4000, "nl", 2, 0},
		{221, "短牌进阶 0.5/1 · 6人", "short6", 6, 50, 100, 5000, 20000, "nl", 1, 0},
		{230, "牛牛微盲 0.1/0.2 · 6人", "niuniu", 6, 10, 20, 1000, 4000, "nl", 1, 0},
		{231, "牛牛休闲 0.2/0.4 · 6人", "niuniu", 6, 20, 40, 2000, 8000, "nl", 1, 0},
	}
	for _, r := range rooms {
		_, err := st.DB.Exec(`INSERT INTO rooms(room_no,name,game_type,seat_count,sb,bb,ante,min_buyin,max_buyin,bet_limit,
			is_official,club_id,host_uid,card_cost,insurance,duration_min,status,sort)
			VALUES(?,?,?,?,?,?,0,?,?,?,1,0,0,?,1,120,1,?)`,
			r.no, r.name, r.gt, r.seats, r.sb, r.bb, r.minBuy, r.maxBuy, r.limit, r.cost, r.sort)
		if err != nil {
			return fmt.Errorf("seed room %d: %w", r.no, err)
		}
	}
	log.Printf("[seed] %d rooms created", len(rooms))
	return nil
}

func seedRobots(st *store.Store) error {
	var n int
	_ = st.DB.QueryRow("SELECT COUNT(1) FROM robots").Scan(&n)
	if n > 0 {
		return nil
	}
	type bot struct {
		nick  string
		style string
		aggr  int
		bluff int
		vpip  int
	}
	bots := []bot{
		{"逃离天花板", "tag", 60, 18, 28},
		{"甜酷布丁", "station", 30, 6, 52},
		{"老郑", "nit", 42, 5, 16},
		{"阿南", "lag", 72, 32, 44},
		{"底池收割机", "maniac", 86, 46, 66},
		{"银座古董", "tag", 58, 14, 26},
		{"曼谷佛牌阁", "fish", 34, 10, 48},
		{"浮生", "tag", 62, 20, 30},
		{"子轩", "lag", 70, 28, 40},
		{"Kimi", "station", 26, 8, 55},
		{"罐头终结者", "lag", 75, 35, 42},
		{"荒野旅行家", "fish", 32, 12, 50},
		{"Fox", "tag", 55, 16, 27},
		{"Dean", "nit", 40, 6, 18},
		{"闷声发大财", "tag", 60, 12, 25},
		{"孤岛沉思", "lag", 68, 26, 38},
		{"怡保白咖啡", "station", 28, 10, 58},
		{"盛哥", "maniac", 90, 52, 70},
		{"追逐黄昏", "tag", 57, 15, 29},
		{"闪电", "lag", 74, 30, 45},
		{"清迈香薰坊", "fish", 30, 8, 47},
		{"Dave", "tag", 63, 20, 31},
		{"零点回撤", "nit", 45, 8, 20},
		{"深夜长街", "station", 25, 5, 53},
		{"半山腰", "lag", 71, 27, 41},
		{"空港旅人", "tag", 56, 13, 24},
		{"小满", "fish", 33, 11, 49},
		{"北纬三十度", "tag", 61, 17, 28},
		{"温莎公爵", "lag", 69, 24, 37},
		{"浪里白条", "maniac", 84, 44, 64},
	}
	for i, b := range bots {
		uid, err := nextUID(st)
		if err != nil {
			return err
		}
		av := fmt.Sprintf("/img/avatars/avatar_%03d.jpg", i%43+1)
		_, err = st.DB.Exec(`INSERT INTO users(uid,username,nickname,avatar,balance,room_cards,is_robot,status)
			VALUES(?,?,?,?,?,0,1,1)`, uid, fmt.Sprintf("bot_%03d", i+1), b.nick, av, 20000000)
		if err != nil {
			return fmt.Errorf("seed robot %s: %w", b.nick, err)
		}
		res, err := st.DB.Exec(`INSERT INTO robot_profiles(name,avatar,style,aggression,vpip,bluff) VALUES(?,?,?,?,?,?)`,
			b.nick, av, b.style, b.aggr, b.vpip, b.bluff)
		if err != nil {
			return err
		}
		pid, _ := res.LastInsertId()
		_, err = st.DB.Exec(`INSERT INTO robots(uid,profile_id,room_no,status) VALUES(?,?,0,1)`, uid, pid)
		if err != nil {
			return err
		}
	}
	log.Printf("[seed] %d robots created", len(bots))
	return nil
}

func nextUID(st *store.Store) (uint32, error) {
	for i := 0; i < 100; i++ {
		var b [4]byte
		for j := range b {
			b[j] = auth.RandomDigits(1)[0]
		}
		n := uint32(100000) + uint32(b[0]-'0')*100000 + uint32(b[1]-'0')*10000 +
			uint32(b[2]-'0')*1000 + uint32(b[3]-'0')*137
		if n > 999999 {
			n = n % 900000
			n += 100000
		}
		var c int
		_ = st.DB.QueryRow("SELECT COUNT(1) FROM users WHERE uid=?", n).Scan(&c)
		if c == 0 {
			return n, nil
		}
	}
	return 0, fmt.Errorf("uid allocation failed")
}

type courseSeed struct {
	Group           int    `json:"group"`
	Lesson          int    `json:"lesson"`
	Sort            int    `json:"sort"`
	Title           string `json:"title"`
	Summary         string `json:"summary"`
	Cover           string `json:"cover"`
	Video           string `json:"video"`
	VideoFallback   string `json:"videoFallback"`
	DurationSeconds int    `json:"durationSeconds"`
	Heat            int    `json:"heat"`
	ViewCount       int    `json:"viewCount"`
	SourceName      string `json:"sourceName"`
	Publish         string `json:"publish"`
}

func mmss(sec int) string {
	if sec <= 0 {
		return ""
	}
	return fmt.Sprintf("%02d:%02d", sec/60, sec%60)
}

func seedCourses(st *store.Store) error {
	var n int
	_ = st.DB.QueryRow("SELECT COUNT(1) FROM courses").Scan(&n)
	if n > 0 {
		return nil
	}
	var list []courseSeed
	if err := json.Unmarshal(coursesJSON, &list); err != nil {
		return err
	}
	for _, c := range list {
		publish := strings.TrimSpace(c.Publish)
		if publish == "" {
			publish = "2026-08-05"
		}
		_, err := st.DB.Exec(`INSERT INTO courses(lesson_no,course_group,sort_no,title,summary,cover,video_url,video_fallback,
				duration,duration_seconds,heat,view_count,source_name,publish_at,status)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,1)`,
			c.Lesson, c.Group, c.Sort, c.Title, c.Summary, c.Cover, c.Video, c.VideoFallback,
			mmss(c.DurationSeconds), c.DurationSeconds, c.Heat, c.ViewCount, c.SourceName, publish)
		if err != nil {
			return err
		}
	}
	log.Printf("[seed] %d courses created", len(list))
	return nil
}

func seedSettings(st *store.Store) error {
	defaults := map[string]string{
		"announcement":  "安全卫士已启动，同环境下IP两人在房间安全卫士实时追踪提醒全房间用户，禁止任何形式的伙牌与串通。",
		"site_name":     "U5 德扑",
		"withdraw_fee":  "0",
		"robot_fill":    "1",
		"min_buyin_rate": "40",
		"mtt_title":      "每晚 21:00 免费分奖金争霸赛",
		"mtt_schedule":   "每日 21:00 (60分钟)",
		"mtt_points":     "3000",
		"mtt_prize":      "18800",
		"jackpot_pool":   "1077736",
		"jackpot_base":   "888800",
		"jackpot_fee":    "100",
		"jackpot_hours":  "24",
		"jackpot_royal_pct": "10",
		"jackpot_sf_pct":    "5",
		"jackpot_quads_pct": "2",
	}
	for k, v := range defaults {
		_ = st.SetSetting(k, v)
	}
	return nil
}
