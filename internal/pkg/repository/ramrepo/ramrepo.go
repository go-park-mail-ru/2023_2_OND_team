package ramrepo

import (
	"database/sql"
	"fmt"

	_ "github.com/proullon/ramsql/driver"
)

func OpenDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("ramsql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = createUsersTable(db)
	if err != nil {
		return nil, err
	}

	err = createPinTable(db)
	if err != nil {
		return nil, err
	}

	err = createSessionTable(db)
	if err != nil {
		return nil, err
	}

	err = fillPinTableRows(db)
	if err != nil {
		return nil, err
	}

	err = fillUsersTableRows(db)
	if err != nil {
		return nil, err
	}

	err = fillSessionTableRows(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createUsersTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE users(
		id bigserial PRIMARY KEY,
		username varchar(30) UNIQUE,
		password varchar(50),
		email varchar(50) UNIQUE,
		avatar varchar(50) DEFAULT 'https://pinspire.online:8081/upload/avatars/default-avatar.png'
	);`)
	if err != nil {
		return fmt.Errorf("create table users: %w", err)
	}
	return nil
}

func createPinTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE pin(
		id bigserial PRIMARY KEY,
		author int,
		picture varchar(50)
	);`)
	if err != nil {
		return fmt.Errorf("create table pin: %w", err)
	}
	return nil
}

func createSessionTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE session(
		session_key varchar(30) PRIMARY KEY,
		user_id int,
		expire timestamp
	);`)
	if err != nil {
		return fmt.Errorf("create table session: %w", err)
	}
	return nil
}

func fillUsersTableRows(db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO users (username, password, email) VALUES 
	("dogsLover", "bf62b19f2f755d892f0ee1efb591795c198bcd0eecb69e58b153064e7ca11f384bf2e2746d91bf36", "dogslove@gmail.com"),
	("professional_player", "2f45a4f97b2d849448ac28cf95d4a55ddbc146f607e158e78b25a1906b469fe9ebde41b8127dd50e", "fortheplayers@yandex.ru"),
	("goodJobBer", "ade1af872d23126858c289e0c1bfc8b57502f7f0237e35fc64d08fab2d6b667358f04ac4174b736b", "jobjobjob@mail.ru");`)
	if err != nil {
		return fmt.Errorf("fill users table: %w", err)
	}
	return nil
}

func fillSessionTableRows(db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO session (session_key, user_id, expire) VALUES
	("461afabf38b3147c", 1, 2024-10-03 10:52:09.243860007 +0000 UTC),
	("f4280a941b664d02", 3434, 2024-10-03 10:52:09.243860007 +0000 UTC);`)
	if err != nil {
		return fmt.Errorf("fill session table: %w", err)
	}
	return nil
}

func fillPinTableRows(db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO pin (picture) VALUES
	('https://pinspire.online:8081/upload/pins/d7dc22616d514788b514fc2edb60920b.png'),
	('https://pinspire.online:8081/upload/pins/ec66fd27b7f74524894740c5c830327e.png'),
	('https://pinspire.online:8081/upload/pins/ec66fd27b7f74524894740c5c830327e.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/9921115ac96c4223ab36a9bb77fd1e3a.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/7b92db4d365d4c409e18b7869df1e448.png'),
	('https://pinspire.online:8081/upload/pins/aed4d0f218564a5a819ac2348fa97d27.png'),
	('https://pinspire.online:8081/upload/pins/aa9128758dfd442382f41ef5bf3c55d3.png'),
	('https://pinspire.online:8081/upload/pins/369cb5bac5ce4496be9d96c2517d45f9.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/a5a7a7871ead45b9a0626b2a1e1f3a44.png'),
	('https://pinspire.online:8081/upload/pins/5d3c0b4625fc4b73b60351be053e46af.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/9921115ac96c4223ab36a9bb77fd1e3a.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/7b92db4d365d4c409e18b7869df1e448.png'),
	('https://pinspire.online:8081/upload/pins/aed4d0f218564a5a819ac2348fa97d27.png'),
	('https://pinspire.online:8081/upload/pins/aa9128758dfd442382f41ef5bf3c55d3.png'),
	('https://pinspire.online:8081/upload/pins/369cb5bac5ce4496be9d96c2517d45f9.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/a5a7a7871ead45b9a0626b2a1e1f3a44.png'),
	('https://pinspire.online:8081/upload/pins/5d3c0b4625fc4b73b60351be053e46af.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/9921115ac96c4223ab36a9bb77fd1e3a.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/7b92db4d365d4c409e18b7869df1e448.png'),
	('https://pinspire.online:8081/upload/pins/aed4d0f218564a5a819ac2348fa97d27.png'),
	('https://pinspire.online:8081/upload/pins/aa9128758dfd442382f41ef5bf3c55d3.png'),
	('https://pinspire.online:8081/upload/pins/369cb5bac5ce4496be9d96c2517d45f9.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/a5a7a7871ead45b9a0626b2a1e1f3a44.png'),
	('https://pinspire.online:8081/upload/pins/5d3c0b4625fc4b73b60351be053e46af.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/9921115ac96c4223ab36a9bb77fd1e3a.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/7b92db4d365d4c409e18b7869df1e448.png'),
	('https://pinspire.online:8081/upload/pins/aed4d0f218564a5a819ac2348fa97d27.png'),
	('https://pinspire.online:8081/upload/pins/aa9128758dfd442382f41ef5bf3c55d3.png'),
	('https://pinspire.online:8081/upload/pins/369cb5bac5ce4496be9d96c2517d45f9.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/a5a7a7871ead45b9a0626b2a1e1f3a44.png'),
	('https://pinspire.online:8081/upload/pins/5d3c0b4625fc4b73b60351be053e46af.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/9921115ac96c4223ab36a9bb77fd1e3a.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/7b92db4d365d4c409e18b7869df1e448.png'),
	('https://pinspire.online:8081/upload/pins/aed4d0f218564a5a819ac2348fa97d27.png'),
	('https://pinspire.online:8081/upload/pins/aa9128758dfd442382f41ef5bf3c55d3.png'),
	('https://pinspire.online:8081/upload/pins/369cb5bac5ce4496be9d96c2517d45f9.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/a5a7a7871ead45b9a0626b2a1e1f3a44.png'),
	('https://pinspire.online:8081/upload/pins/5d3c0b4625fc4b73b60351be053e46af.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/9921115ac96c4223ab36a9bb77fd1e3a.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/7b92db4d365d4c409e18b7869df1e448.png'),
	('https://pinspire.online:8081/upload/pins/aed4d0f218564a5a819ac2348fa97d27.png'),
	('https://pinspire.online:8081/upload/pins/aa9128758dfd442382f41ef5bf3c55d3.png'),
	('https://pinspire.online:8081/upload/pins/369cb5bac5ce4496be9d96c2517d45f9.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/a5a7a7871ead45b9a0626b2a1e1f3a44.png'),
	('https://pinspire.online:8081/upload/pins/5d3c0b4625fc4b73b60351be053e46af.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/9921115ac96c4223ab36a9bb77fd1e3a.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/7b92db4d365d4c409e18b7869df1e448.png'),
	('https://pinspire.online:8081/upload/pins/aed4d0f218564a5a819ac2348fa97d27.png'),
	('https://pinspire.online:8081/upload/pins/aa9128758dfd442382f41ef5bf3c55d3.png'),
	('https://pinspire.online:8081/upload/pins/369cb5bac5ce4496be9d96c2517d45f9.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/a5a7a7871ead45b9a0626b2a1e1f3a44.png'),
	('https://pinspire.online:8081/upload/pins/5d3c0b4625fc4b73b60351be053e46af.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/9921115ac96c4223ab36a9bb77fd1e3a.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/7b92db4d365d4c409e18b7869df1e448.png'),
	('https://pinspire.online:8081/upload/pins/aed4d0f218564a5a819ac2348fa97d27.png'),
	('https://pinspire.online:8081/upload/pins/aa9128758dfd442382f41ef5bf3c55d3.png'),
	('https://pinspire.online:8081/upload/pins/369cb5bac5ce4496be9d96c2517d45f9.png'),
	('https://pinspire.online:8081/upload/pins/df4f0038efe24e86a471444f94bc863e.png'),
	('https://pinspire.online:8081/upload/pins/2ac2ad104cdd4ca0981bac73777ad368.png'),
	('https://pinspire.online:8081/upload/pins/a5a7a7871ead45b9a0626b2a1e1f3a44.png'),
	('https://pinspire.online:8081/upload/pins/5d3c0b4625fc4b73b60351be053e46af.png'),
	('https://pinspire.online:8081/upload/pins/43b8b6602f9d404ca3510f28fb712026.png'),
	('https://pinspire.online:8081/upload/pins/60ce511226c14573a0151f5e5893b15d.png'),
	('https://pinspire.online:8081/upload/pins/a5a7a7871ead45b9a0626b2a1e1f3a44.png'),
	('https://pinspire.online:8081/upload/pins/aa9128758dfd442382f41ef5bf3c55d3.png'),
	('https://pinspire.online:8081/upload/pins/369cb5bac5ce4496be9d96c2517d45f9.png');`)
	if err != nil {
		return fmt.Errorf("fill pin table: %w", err)
	}
	return nil
}
