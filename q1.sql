INSERT INTO problems (title, description, sample_code)
VALUES
    ('練習問題',       '', ''),
    ('九九',           '', ''),
    ('いろは',         '', ''),
    ('グルーピング',   '', ''),
    ('しりとり',       '', ''),
    ('ポップカウント', '', '');

INSERT INTO games (game_type, is_public, display_name, duration_seconds, problem_id)
VALUES
    ('multiplayer', true, '練習問題',        900, 21),
    ('multiplayer', true, '予選ラウンド1',   900, 22),
    ('multiplayer', true, '予選ラウンド2',   900, 23),
    ('multiplayer', true, '予選ラウンド3',   900, 24),
    ('1v1',         true, 'エキシビション1', 900, 25),
    ('1v1',         true, 'エキシビション2', 900, 26);
