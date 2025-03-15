INSERT INTO problems (title, description, sample_code)
VALUES
    ('じゃんけん',   '', ''),
    ('でっかい数字', '', '');

INSERT INTO games (game_type, is_public, display_name, duration_seconds, problem_id)
VALUES
    ('multiplayer', false, 'オンライン予選1', 86400*5 + 3600*11, (SELECT problem_id FROM problems WHERE title = 'じゃんけん')),
    ('multiplayer', false, 'オンライン予選2', 86400*5 + 3600*11, (SELECT problem_id FROM problems WHERE title = 'でっかい数字'));
