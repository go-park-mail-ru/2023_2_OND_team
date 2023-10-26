SET search_path TO pinspire;

CREATE EXTENSION IF NOT EXISTS moddatetime;

CREATE OR REPLACE TRIGGER modify_profile_updated_at
	BEFORE UPDATE
	ON profile
	FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_pin_updated_at
	BEFORE UPDATE
	ON pin
	FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_board_updated_at
	BEFORE UPDATE
	ON board
	FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_contributor_updated_at
	BEFORE UPDATE
	ON contributor
	FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_comment_updated_at
	BEFORE UPDATE
	ON comment
	FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_message_updated_at
	BEFORE UPDATE
	ON message
	FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);
