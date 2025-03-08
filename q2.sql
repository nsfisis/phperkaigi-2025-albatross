UPDATE problems
SET
sample_code = $EOF$$line = fgets(STDIN);
$number = intval($line);
echo $number * $number, PHP_EOL;
$EOF$
WHERE title = '練習問題';

UPDATE problems
SET
sample_code = $EOF$fscanf(STDIN, "%d %d", $col, $row);

for ($i = 1; $i <= $row; $i++) {
    for ($j = 1; $j <= $col; $j++) {
        printf(str_repeat(" ", strlen($col * $row) - strlen($j * $i) + 1));
        printf("%d", $i * $j);
    }
    printf("\n");
}
$EOF$
WHERE title = '九九';

UPDATE problems
SET
sample_code = $EOF$$lines = [];
while ($line = fgets(STDIN)) {
    $lines[] = rtrim($line);
}

function to_iroha($chars) {
    return implode(',',array_map(fn ($c) => match($c) {
        'い' => 0, 'ろ' => 1, 'は' => 2, 'に' => 3, 'ほ' => 4, 'へ' => 5, 'と' => 6,
        'ち' => 7, 'り' => 8, 'ぬ' => 9, 'る' => 10, 'を' => 11,
    }, mb_str_split($chars)));
}

usort($lines, function (string $a, string $b): int {
    $as = to_iroha($a);
    $bs = to_iroha($b);

    if ($as === $bs) {
        return 0;
    } else {
        return ($as > $bs) ? 1 : -1;
    }
});

foreach ($lines as $line) echo $line, PHP_EOL;
$EOF$
WHERE title = 'いろは';

UPDATE problems
SET
sample_code = $EOF$$result = [];
while ([$a, $b] = fgetcsv(STDIN, 1024, ',', '"', '')) {
    if (!isset($result[$a])) {
        $result[$a] = [$b];
    } else {
        $result[$a][] = $b;
    }
}
ksort($result);
foreach ($result as $i => $v) {
    sort($v);
    fputcsv(STDOUT, array_merge([$i], $v), ',', '"', '');
}
$EOF$
WHERE title = 'グルーピング';

UPDATE problems
SET
sample_code = $EOF$$lines = [];
while ($line = fgets(STDIN)) {
    $lines[] = rtrim($line);
}

$current = array_shift($lines);
$words = $lines;
$lines = [];
echo $current, "\n";
do {
    $last = mb_str_split($current)[mb_strlen($current)-1];
    if ($last === 'ン') {
        echo "負けました\n";
        break;
    }
    foreach ($words as $i => $word) {
        if (str_starts_with($word, $last)) {
            echo $word, "\n";
            unset($words[$i]);
            $current = $word;
            continue 2;
        } 
    }
    exit("おかしいな\n");
} while(true);
$EOF$
WHERE title = 'しりとり';

UPDATE problems
SET
sample_code = $EOF$while($line = fgets(STDIN)) {
    $line = rtrim($line, "\n");
    $count = 0;
    foreach (str_split($line) as $char) {
        $bits = str_split(decbin(ord($char)));
        $count = $count + count(array_filter($bits, fn($c)=> $c == 1));
    }
    echo $count, PHP_EOL;
}
$EOF$
WHERE title = 'ポップカウント';
