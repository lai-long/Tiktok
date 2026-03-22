package db

import "log"

func (m *MySQLdb) InsertMsg(session, content, sender, receiver string) {
	sql := `insert into message(session_id,content,sender_id,receiver_id) values(?,?,?,?)`
	_, err := m.db.Exec(sql, session, content, sender, receiver)
	log.Println("数据库数据", session, content, sender, receiver)
	if err != nil {
		log.Println(err)
	}
}
