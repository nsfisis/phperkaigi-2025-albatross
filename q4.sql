INSERT INTO testcases (problem_id, stdin, stdout)
VALUES
    ((SELECT problem_id FROM problems WHERE title = '序数'),
    $$1 penguin
2 chance
3 time's the charm
$$,
    $$1st penguin
2nd chance
3rd time's the charm
$$),
    ((SELECT problem_id FROM problems WHERE title = '序数'),
    $$aaa 10 11 12 13 14 15 16 17 18 19 20 bbb
ccc 21 22 23 24 25 26 27 28 29 30 ddd
aaa 110 111 112 113 114 bbb
$$,
    $$aaa 10th 11th 12th 13th 14th 15th 16th 17th 18th 19th 20th bbb
ccc 21st 22nd 23rd 24th 25th 26th 27th 28th 29th 30th ddd
aaa 110th 111th 112th 113th 114th bbb
$$);

INSERT INTO testcases (problem_id, stdin, stdout)
VALUES
    ((SELECT problem_id FROM problems WHERE title = 'Σ'),
    $$6
$$,
    $$******
 *
  *
   *
  *
 *
******
$$),
    ((SELECT problem_id FROM problems WHERE title = 'Σ'),
    $$10
$$,
    $$**********
 *
  *
   *
    *
     *
    *
   *
  *
 *
**********
$$);

INSERT INTO testcases (problem_id, stdin, stdout)
VALUES
    ((SELECT problem_id FROM problems WHERE title = '漢数字'),
    $$1234567890
第8回
ペチパー会議2025
気温は摂氏7度
総勢123人
税込98700円
$$,
    $$一二三四五六七八九〇
第八回
ペチパー会議二〇二五
気温は摂氏七度
総勢一二三人
税込九八七〇〇円
$$),
    ((SELECT problem_id FROM problems WHERE title = '漢数字'),
    $$-123
$$,
    $$-一二三
$$);

INSERT INTO testcases (problem_id, stdin, stdout)
VALUES
    ((SELECT problem_id FROM problems WHERE title = 'カレンダー'),
    $$2025-02
$$,
    $$[2025年2月]
                    1
  2  3  4  5  6  7  8
  9 10 11 12 13 14 15
 16 17 18 19 20 21 22
 23 24 25 26 27 28
$$),
    ((SELECT problem_id FROM problems WHERE title = 'カレンダー'),
    $$2000-02
$$,
    $$[2000年2月]
        1  2  3  4  5
  6  7  8  9 10 11 12
 13 14 15 16 17 18 19
 20 21 22 23 24 25 26
 27 28 29
$$),
    ((SELECT problem_id FROM problems WHERE title = 'カレンダー'),
    $$2025-03
$$,
    $$[2025年3月]
                    1
  2  3  4  5  6  7  8
  9 10 11 12 13 14 15
 16 17 18 19 20 21 22
 23 24 25 26 27 28 29
 30 31
$$),
    ((SELECT problem_id FROM problems WHERE title = 'カレンダー'),
    $$2024-02
$$,
    $$[2024年2月]
              1  2  3
  4  5  6  7  8  9 10
 11 12 13 14 15 16 17
 18 19 20 21 22 23 24
 25 26 27 28 29
$$);

INSERT INTO testcases (problem_id, stdin, stdout)
VALUES
    ((SELECT problem_id FROM problems WHERE title = 'ダイナミックFizzBuzz'),
    $$10
2,Dizz
3,Fizz
5,Buzz
$$,
    $$1
Dizz
Fizz
Dizz
Buzz
DizzFizz
7
Dizz
Fizz
DizzBuzz
$$),
    ((SELECT problem_id FROM problems WHERE title = 'ダイナミックFizzBuzz'),
    $$21
3,Fizz
7,Pazz
$$,
    $$1
2
Fizz
4
5
Fizz
Pazz
8
Fizz
10
11
Fizz
13
Pazz
Fizz
16
17
Fizz
19
20
FizzPazz
$$);
