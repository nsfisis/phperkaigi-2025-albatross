UPDATE problems
SET
description = $EOF$標準入力から数字が一つ入力されます。二乗して標準出力へ出力してください。
$EOF$
WHERE title = '練習問題';

UPDATE problems
SET
description = $EOF$標準入力にスペース区切りの整数col rowが与えられます。「かける数」は1〜col、「かけられる数」は1〜rowまでの連続した整数です。全ての組み合わせのかけ算の結果が右揃えに整列するように出力してください。空白の個数は最も横幅が長いマスに合わせてください。
$EOF$
WHERE title = '九九';

UPDATE problems
SET
description = $EOF$「いろはにほへとちりぬるを」の12種類の文字から構成される文字列が改行区切りで入力されます。標準入力の内容を全て読み取り、文字列を「いろは」順に並べて出力してください。

入力:
るりいろ
りぬる
いちにち
ほへとちに
いち
ほとはろ
りにへ
いい
いいい
とりほい
いろ
いろいろ

出力:
いい
いいい
いろ
いろいろ
いち
いちにち
るりいろ
ほへとちに
ほとはろ
とりほい
りにへ
りぬる
$EOF$
WHERE title = 'いろは';

UPDATE problems
SET
description = $EOF$標準入力からRFC 4180準拠の行区切りのデータが与えられます。行は name,value のような構造になっています。nameごとにvalueを集約し、name,value1,value2,...のようなCSVとして出力してください。ただし、行はnameのアルファベット昇順、列はnameを先頭にして、valueをアルファベット昇順でそれぞれソートして出力すること。

入力:
Foo,zzz
Foo,xxx
Bar,ccc
Bar,bbb
Foo,yyy
Bar,"a, aa, aaa"
Bar,aaa
Foo,"hello, world"

出力:
Bar,"a, aa, aaa",aaa,bbb,ccc
Foo,"hello, world",xxx,yyy,zzz
$EOF$
WHERE title = 'グルーピング';

UPDATE problems
SET
description = $EOF$標準入力から改行区切りでカタカナ文字列の単語が入力されます。全ての行を読み取り、最初に入力された行から開始して「しりとり」になるように単語を改行区切りで順番に出力してください。一度使った文字を再利用することはできません。出力した単語の最後の文字が「ン」で終わったときは"負けました\n"と出力して、プログラムを終了してください。最初の文字と最後の文字は重複しないように入力されます。

入力:
リンゴ
ライオン
ゴリラ

出力:
リンゴ
ゴリラ
ライオン
負けました

入力:
リンゴ
ライオン
ゴリラ

出力:
リンゴ
ゴリラ
ライオン
負けました
$EOF$
WHERE title = 'しりとり';

UPDATE problems
SET
description = $EOF$標準入力から行区切りのデータが入力されます。行末の改行コードを取り除いた入力文字列について、2進数で表現した際にビットで表現した時の1の個数を行ごとに合計した数を改行区切りで出力してください。たとえば文字列のXは数値で表すと88ですが、2進数で表示すると01011000なので、1のビットの個数は3になります。

入力:
ABC
Hello, world

出力:
7
47
$EOF$
WHERE title = 'ポップカウント';
