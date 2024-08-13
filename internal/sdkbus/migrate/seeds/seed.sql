INSERT INTO users (user_id, name, email, roles, password_hash, enabled, email_confirm_token, date_created, date_updated) VALUES
	('97ee07e2-ebbb-4c69-a681-d5fe165c2cb9', 'Admin', 'admin@email.com', '{ADMIN}', '$2a$10$LYuQJ38ZzQ2vTah6LGVuROIwzovi4K5h3UN4mY6MtPVBb98/tBfom', true, null, '2024-03-24 01:02:03', '2024-03-24 04:05:06'),
	('272f05b5-b080-4e13-a976-153455926530', 'User', 'user@email.com', '{USER}', '$2a$10$LYuQJ38ZzQ2vTah6LGVuROIwzovi4K5h3UN4mY6MtPVBb98/tBfom', true, null, '2024-04-25 07:08:09', '2024-04-25 10:11:12')
ON CONFLICT DO NOTHING;
