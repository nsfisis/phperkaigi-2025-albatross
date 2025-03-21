UPDATE problems
SET
description = $EOF$標準入力から行区切りでテキストが入力されます。テキストに含まれる整数(n>0)を序数表現(1→1st, 2→2nd)のように置換して出力してください。

入力:
1 penguin
2 chance
3 time's the charm

出力:
1st penguin
2nd chance
3rd time's the charm
$EOF$
WHERE title = '序数';

UPDATE problems
SET
description = $EOF$標準入力から整数widthが入力されます。入力した数字から例のような図形を出力してください。ただし、入力は4以上の整数のみです。

入力:
6

出力:
******
 *
  *
   *
  *
 *
******
$EOF$
WHERE title = 'Σ';

UPDATE problems
SET
description = $EOF$標準入力から行区切りでテキストが入力されます。テキストに含まれる数字を漢数字に置換して出力してください。置換は一桁ずつおこない、数字が連続している場合でも「十」や「百」にはしないでください。

入力:
1234567890
第8回
ペチパー会議2025
気温は摂氏7度
総勢123人
税込98700円

出力:
一二三四五六七八九〇
第八回
ペチパー会議二〇二五
気温は摂氏七度
総勢一二三人
税込九八七〇〇円
$EOF$
WHERE title = '漢数字';

UPDATE problems
SET
description = $EOF$標準入力から「西暦年-月」の形式で 2025-02 のような行が1行入力されます。入力された月に合わせて「[Y年M月]」と出力してから、入力された月のカレンダーを日曜始まり形式で整列して出力してください。日付はすべて「3桁空白右寄せ」で表示します。

入力:
2025-02

出力
[2025年2月]
                    1
  2  3  4  5  6  7  8
  9 10 11 12 13 14 15
 16 17 18 19 20 21 22
 23 24 25 26 27 28
$EOF$
WHERE title = 'カレンダー';

UPDATE problems
SET
description = $EOF$標準入力の最初の行に出力の最大行数 $max が、次以降の行に置換リストが 数,名前 の形式で改行区切りで入力されます。置換リストに含まれる数は2以上の素数で、昇順に入力されます。1 から $max までの連続した整数を改行区切りで出力してください。その際、出力しようとしている数が置換リストの数の倍数なら対応する名前に置き換え、その数が置換リストに含まれる複数の数の公倍数ならリスト内の数が小さい順に連結して置き換えてください。

入力:
10
2,Dizz
3,Fizz
5,Buzz

出力:
1
Dizz
Fizz
Dizz
Buzz
DizzFizz
7
Dizz
Fizz
DizzBuzz
$EOF$
WHERE title = 'ダイナミックFizzBuzz';
