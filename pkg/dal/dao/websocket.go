package dao

func (m *MySQLdb) InsertMsg(session, content string) error {
	sql := `insert into message(session_id,content) values(?,?)`
	_, err := m.db.Exec(sql, session, content)
	if err != nil {
		return err
	}
	return nil
}
func (m *MySQLdb) GetWebsocketHistory(session1, session2 string, pageNum, pageSize int) ([]string, error) {
	sql := `select content from message where session_id=? or session_id=? LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	var messages []string
	err := m.db.Select(&messages, sql, session1, session2, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
