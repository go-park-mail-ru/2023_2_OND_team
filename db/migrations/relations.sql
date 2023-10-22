CREATE SCHEMA IF NOT EXISTS pinspire;

SET search_path TO pinspire;

CREATE TABLE IF NOT EXISTS profile (
	id serial PRIMARY KEY,
	email text NOT NULL,
	avatar text NOT NULL DEFAULT 'default-avatar.png',
	name text,
	surname text,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	CONSTRAINT profile_email_uniq UNIQUE (email)
);

ALTER TABLE profile ALTER COLUMN avatar SET DEFAULT 'avatar.jpg'; 

CREATE TABLE IF NOT EXISTS auth (
	id serial PRIMARY KEY,
	username text NOT NULL,
	password text NOT NULL,
	profile_id int NOT NULL,
	CONSTRAINT auth_username_uniq UNIQUE (username),
	CONSTRAINT auth_profile_id_uniq UNIQUE (profile_id),
	FOREIGN KEY (profile_id) REFERENCES profile (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tag (
	id serial PRIMARY KEY,
	title text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT tag_title_uniq UNIQUE (title)
);

CREATE TABLE IF NOT EXISTS pin (
	id serial PRIMARY KEY,
	author int NOT NULL,
	title text,
	description text,
	picture text NOT NULL,
	public bool NOT NULL DEFAULT TRUE,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	FOREIGN KEY (author) REFERENCES profile (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pin_tag (
	pin_id int NOT NULL,
	tag_id int NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (pin_id, tag_id),
	FOREIGN KEY (pin_id) REFERENCES pin (id) ON DELETE CASCADE,
	FOREIGN KEY (tag_id) REFERENCES tag (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS like_pin (
	user_id int NOT NULL,
	pin_id int NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (pin_id, user_id),
	FOREIGN KEY (user_id) REFERENCES profile (id) ON DELETE CASCADE,
	FOREIGN KEY (pin_id) REFERENCES pin (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS board (
	id serial PRIMARY KEY,
	author int NOT NULL,
	title text,
	description text,
	public bool NOT NULL DEFAULT TRUE,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	FOREIGN KEY (author) REFERENCES profile (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS board_tag (
	board_id int NOT NULL,
	tag_id int NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (board_id, tag_id),
	FOREIGN KEY (board_id) REFERENCES board (id) ON DELETE CASCADE,
	FOREIGN KEY (tag_id) REFERENCES tag (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscription_board (
	user_id int NOT NULL,
	board_id int NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (board_id, user_id),
	FOREIGN KEY (user_id) REFERENCES profile (id) ON DELETE CASCADE,
	FOREIGN KEY (board_id) REFERENCES board (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS membership (
	pin_id int NOT NULL,
	board_id int NOT NULL,
	added_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (board_id, pin_id),
	FOREIGN KEY (pin_id) REFERENCES pin (id) ON DELETE CASCADE,
	FOREIGN KEY (board_id) REFERENCES board (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS role (
	id serial PRIMARY KEY,
	name text NOT NULL,
	CONSTRAINT role_name_uniq UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS contributor (
	board_id int NOT NULL,
	user_id int NOT NULL,
	role_id int NOT NULL,
	added_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (user_id, board_id),
	FOREIGN KEY (board_id) REFERENCES board (id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES profile (id) ON DELETE CASCADE,
	FOREIGN KEY (role_id) REFERENCES role (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscription_user (
	who int NOT NULL,
	whom int NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (whom, who),
	FOREIGN KEY (who) REFERENCES profile (id) ON DELETE CASCADE,
	FOREIGN KEY (whom) REFERENCES profile (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment (
	id serial PRIMARY KEY,
	author int NOT NULL,
	pin_id int NOT NULL,
	content text,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	FOREIGN KEY (author) REFERENCES profile (id) ON DELETE CASCADE,
	FOREIGN KEY (pin_id) REFERENCES pin (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS like_comment (
	user_id int NOT NULL,
	comment_id int NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (comment_id, user_id),
	FOREIGN KEY (user_id) REFERENCES profile (id) ON DELETE CASCADE,
	FOREIGN KEY (comment_id) REFERENCES comment (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS message (
	id serial PRIMARY KEY,
	user_from int NOT NULL,
	user_to int NOT NULL,
	content text,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	FOREIGN KEY (user_from) REFERENCES profile (id) ON DELETE CASCADE,
	FOREIGN KEY (user_to) REFERENCES profile (id) ON DELETE CASCADE
);
