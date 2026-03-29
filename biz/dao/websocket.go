package dao

import "log"

func (m *MySQLdb) InsertMsg(session, content string) {
	sql := `insert into message(session_id,content) values(?,?)`
	_, err := m.db.Exec(sql, session, content)
	if err != nil {
		log.Println(err)
	}
}
func (m *MySQLdb) GetWebsocketHistory(session1, session2 string) []string {
	sql := `select content from message where session_id=? or session_id=?`
	var messages []string
	err := m.db.Select(&messages, sql, session1, session2)
	if err != nil {
		log.Println(err)
		return nil
	}
	return messages
}
