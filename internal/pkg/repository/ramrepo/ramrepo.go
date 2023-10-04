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
		avatar varchar(50) DEFAULT 'https://cdn-icons-png.flaticon.com/512/149/149071.png'
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
	('https://i.pinimg.com/564x/e2/43/10/e24310fe1909ec1f1de347fedc6318b0.jpg'),
	('https://i.pinimg.com/564x/91/39/51/913951d97d3cc3ac5a4ecb58da2ffdf5.jpg'),
	('https://i.pinimg.com/564x/91/39/51/913951d97d3cc3ac5a4ecb58da2ffdf5.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/49/23/a9/4923a9a174fc87ab806121e79fda51e4.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/57/21/90/5721907848655c918c667d84defb99f8.jpg'),
	('https://i.pinimg.com/564x/f8/bd/0a/f8bd0aeae74e94e12eb57b6ae3280d6c.jpg'),
	('https://i.pinimg.com/564x/ff/03/1f/ff031f62ad3e9e3733ed78216064978c.jpg'),
	('https://i.pinimg.com/564x/b0/17/fe/b017fea78ff90de1187b857166f12af8.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/30/da/d2/30dad2f5d5923e7a7715fe25ea590d35.jpg'),
	('https://i.pinimg.com/564x/bc/07/62/bc07626808f2f1385e6d38765ff115cc.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/49/23/a9/4923a9a174fc87ab806121e79fda51e4.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/57/21/90/5721907848655c918c667d84defb99f8.jpg'),
	('https://i.pinimg.com/564x/f8/bd/0a/f8bd0aeae74e94e12eb57b6ae3280d6c.jpg'),
	('https://i.pinimg.com/564x/ff/03/1f/ff031f62ad3e9e3733ed78216064978c.jpg'),
	('https://i.pinimg.com/564x/b0/17/fe/b017fea78ff90de1187b857166f12af8.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/30/da/d2/30dad2f5d5923e7a7715fe25ea590d35.jpg'),
	('https://i.pinimg.com/564x/bc/07/62/bc07626808f2f1385e6d38765ff115cc.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/49/23/a9/4923a9a174fc87ab806121e79fda51e4.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/57/21/90/5721907848655c918c667d84defb99f8.jpg'),
	('https://i.pinimg.com/564x/f8/bd/0a/f8bd0aeae74e94e12eb57b6ae3280d6c.jpg'),
	('https://i.pinimg.com/564x/ff/03/1f/ff031f62ad3e9e3733ed78216064978c.jpg'),
	('https://i.pinimg.com/564x/b0/17/fe/b017fea78ff90de1187b857166f12af8.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/30/da/d2/30dad2f5d5923e7a7715fe25ea590d35.jpg'),
	('https://i.pinimg.com/564x/bc/07/62/bc07626808f2f1385e6d38765ff115cc.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/49/23/a9/4923a9a174fc87ab806121e79fda51e4.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/57/21/90/5721907848655c918c667d84defb99f8.jpg'),
	('https://i.pinimg.com/564x/f8/bd/0a/f8bd0aeae74e94e12eb57b6ae3280d6c.jpg'),
	('https://i.pinimg.com/564x/ff/03/1f/ff031f62ad3e9e3733ed78216064978c.jpg'),
	('https://i.pinimg.com/564x/b0/17/fe/b017fea78ff90de1187b857166f12af8.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/30/da/d2/30dad2f5d5923e7a7715fe25ea590d35.jpg'),
	('https://i.pinimg.com/564x/bc/07/62/bc07626808f2f1385e6d38765ff115cc.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/49/23/a9/4923a9a174fc87ab806121e79fda51e4.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/57/21/90/5721907848655c918c667d84defb99f8.jpg'),
	('https://i.pinimg.com/564x/f8/bd/0a/f8bd0aeae74e94e12eb57b6ae3280d6c.jpg'),
	('https://i.pinimg.com/564x/ff/03/1f/ff031f62ad3e9e3733ed78216064978c.jpg'),
	('https://i.pinimg.com/564x/b0/17/fe/b017fea78ff90de1187b857166f12af8.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/30/da/d2/30dad2f5d5923e7a7715fe25ea590d35.jpg'),
	('https://i.pinimg.com/564x/bc/07/62/bc07626808f2f1385e6d38765ff115cc.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/49/23/a9/4923a9a174fc87ab806121e79fda51e4.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/57/21/90/5721907848655c918c667d84defb99f8.jpg'),
	('https://i.pinimg.com/564x/f8/bd/0a/f8bd0aeae74e94e12eb57b6ae3280d6c.jpg'),
	('https://i.pinimg.com/564x/ff/03/1f/ff031f62ad3e9e3733ed78216064978c.jpg'),
	('https://i.pinimg.com/564x/b0/17/fe/b017fea78ff90de1187b857166f12af8.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/30/da/d2/30dad2f5d5923e7a7715fe25ea590d35.jpg'),
	('https://i.pinimg.com/564x/bc/07/62/bc07626808f2f1385e6d38765ff115cc.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/49/23/a9/4923a9a174fc87ab806121e79fda51e4.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/57/21/90/5721907848655c918c667d84defb99f8.jpg'),
	('https://i.pinimg.com/564x/f8/bd/0a/f8bd0aeae74e94e12eb57b6ae3280d6c.jpg'),
	('https://i.pinimg.com/564x/ff/03/1f/ff031f62ad3e9e3733ed78216064978c.jpg'),
	('https://i.pinimg.com/564x/b0/17/fe/b017fea78ff90de1187b857166f12af8.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/30/da/d2/30dad2f5d5923e7a7715fe25ea590d35.jpg'),
	('https://i.pinimg.com/564x/bc/07/62/bc07626808f2f1385e6d38765ff115cc.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/49/23/a9/4923a9a174fc87ab806121e79fda51e4.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/57/21/90/5721907848655c918c667d84defb99f8.jpg'),
	('https://i.pinimg.com/564x/f8/bd/0a/f8bd0aeae74e94e12eb57b6ae3280d6c.jpg'),
	('https://i.pinimg.com/564x/ff/03/1f/ff031f62ad3e9e3733ed78216064978c.jpg'),
	('https://i.pinimg.com/564x/b0/17/fe/b017fea78ff90de1187b857166f12af8.jpg'),
	('https://i.pinimg.com/564x/32/80/5e/32805ec1935f0e4d2e4544d328512e03.jpg'),
	('https://i.pinimg.com/564x/f7/f8/d4/f7f8d4200cb60af122be89a39fd45c57.jpg'),
	('https://i.pinimg.com/564x/30/da/d2/30dad2f5d5923e7a7715fe25ea590d35.jpg'),
	('https://i.pinimg.com/564x/bc/07/62/bc07626808f2f1385e6d38765ff115cc.jpg'),
	('https://i.pinimg.com/564x/ec/b9/ca/ecb9cae2e1f174aca65d5d369f9a71d9.jpg'),
	('https://i.pinimg.com/564x/43/67/15/4367152cd5654e8e74afab54823732ef.jpg'),
	('https://i.pinimg.com/564x/30/da/d2/30dad2f5d5923e7a7715fe25ea590d35.jpg'),
	('https://i.pinimg.com/564x/ff/03/1f/ff031f62ad3e9e3733ed78216064978c.jpg'),
	('https://i.pinimg.com/564x/b0/17/fe/b017fea78ff90de1187b857166f12af8.jpg');`)
	if err != nil {
		return fmt.Errorf("fill pin table: %w", err)
	}
	return nil
}
