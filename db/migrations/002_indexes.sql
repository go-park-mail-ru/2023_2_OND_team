SET search_path TO pinspire;

CREATE INDEX IF NOT EXISTS pin_author_index
ON pin USING btree (author);

CREATE INDEX IF NOT EXISTS board_author_index
ON board USING btree (author);

CREATE INDEX IF NOT EXISTS comment_pin_author_index
ON comment USING btree (pin_id, author);

CREATE INDEX IF NOT EXISTS message_from_to_index
ON message USING btree (user_from, user_to);
