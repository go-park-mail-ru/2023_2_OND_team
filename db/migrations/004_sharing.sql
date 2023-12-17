CREATE TABLE IF NOT EXISTS link (
	id serial PRIMARY KEY,
	board_id int NOT NULL,
	role_id int NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	FOREIGN KEY (board_id) REFERENCES board (id) ON DELETE CASCADE,
	FOREIGN KEY (role_id) REFERENCES role (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS access_link (
	link_id int NOT NULL,
	user_id int NOT NULL,
	PRIMARY KEY (link_id, user_id),
	FOREIGN KEY (link_id) REFERENCES link (id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES profile (id) ON DELETE CASCADE
);
