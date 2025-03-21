INSERT INTO problems (title, description, sample_code)
VALUES
    ('序数',                 '', ''),
    ('Σ',                    '', ''),
    ('漢数字',               '', ''),
    ('カレンダー',           '', ''),
    ('ダイナミックFizzBuzz', '', '');

INSERT INTO games (game_type, is_public, display_name, duration_seconds, problem_id)
VALUES
    ('1v1', false, '準々決勝マッチ1', 15*60, (SELECT problem_id FROM problems WHERE title = '序数')),
    ('1v1', false, '準々決勝マッチ2', 15*60, (SELECT problem_id FROM problems WHERE title = 'Σ')),
    ('1v1', false, '準決勝マッチ1',   15*60, (SELECT problem_id FROM problems WHERE title = '漢数字')),
    ('1v1', false, '準決勝マッチ2',   15*60, (SELECT problem_id FROM problems WHERE title = 'カレンダー')),
    ('1v1', false, '決勝',            15*60, (SELECT problem_id FROM problems WHERE title = 'ダイナミックFizzBuzz'));
