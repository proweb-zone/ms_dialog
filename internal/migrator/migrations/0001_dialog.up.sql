CREATE TABLE dialog (
    id SERIAL,
    user_id_sender INTEGER,
	  user_id_recipient INTEGER,
	  msg TEXT,
    state BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    primary key (id, user_id_sender, user_id_recipient)
);
