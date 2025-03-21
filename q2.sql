UPDATE problems
SET
sample_code = $EOF$while($line = fgets(STDIN)) {
    echo preg_replace_callback('/(?P<num>\b[1-9][0-9]*)\b/', function ($m) {
        $num = $m['num'];
        return toOrd($num);
    }, $line);
}

function toOrd(int $num): string {
    $_2 = $num % 100;
    if ($_2 === 11 || $_2 === 12 || $_2 === 13) {
        return $num . 'th';
    }

    $_1 = $num % 10;
    if ($_1 === 1) {
        return $num . 'st';
    } elseif ($_1 === 2) {
        return $num . 'nd';
    } elseif ($_1 === 3) {
        return $num . 'rd';
    } else {
        return $num . 'th';
    }
}
$EOF$
WHERE title = '序数';

UPDATE problems
SET
sample_code = $EOF$fscanf(STDIN, "%d", $width);

foreach (range(1, $width) as $i) {
    echo '*';
}
echo PHP_EOL;

foreach (range(1, intdiv($width, 2)) as $i) {
    foreach (range(0, $i) as $j) {
        if ($i === $j) {
            echo '*';
        } else {
            echo ' ';
        }
    }
    echo PHP_EOL;
}
foreach (range(intdiv($width, 2)-1, 1) as $i) {
    foreach (range(0, $i) as $j) {
        if ($i === $j) {
            echo '*';
        } else {
            echo ' ';
        }
    }
    echo PHP_EOL;
}
foreach (range(1, $width) as $i) {
    echo '*';
}
$EOF$
WHERE title = 'Σ';

UPDATE problems
SET
sample_code = $EOF$while($line = fgets(STDIN)) {
    echo preg_replace_callback('/(?P<num>\b[1-9]?[0-9]+)\b/', function ($m) {
        $num = $m['num'];
        return to漢数字($num);
    }, $line);
}

function to漢数字(int $n): string
{
    $kansuji = [
        0 => '〇',
        1 => '一',
        2 => '二',
        3 => '三',
        4 => '四',
        5 => '五',
        6 => '六',
        7 => '七',
        8 => '八',
        9 => '九',
    ];

    do {
        $_ = (int)$n % 10;
        $digits[] = $kansuji[$_];
        $n = $n / 10;
    } while($n >= 1);

    return implode(array_reverse($digits));
}
$EOF$
WHERE title = '漢数字';

UPDATE problems
SET
sample_code = $EOF$[$year, $month] = explode("-", trim(fgets(STDIN)));
$first_day = strtotime("first day of {$year}-{$month}");
$last_day = strtotime("last day of {$year}-{$month}");

$calendar = [];
$this_week = array_fill(0, idate('w', $first_day), '');
$days = range(1, idate('d', $last_day));
while (count($days) > 0) {
    $this_week[] = array_shift($days);
    if (count($this_week) === 7) {
        $calendar[] = $this_week;
        $this_week = [];
    }
}
$calendar[] = $this_week;

echo "[", idate('Y', $first_day), "年", idate("n", $first_day), "月]\n";

foreach($calendar as $week){
    foreach ($week as $day){
        printf("%3s", $day);
    }
    echo "\n";
}
$EOF$
WHERE title = 'カレンダー';

UPDATE problems
SET
sample_code = $EOF$$max = intval(fgets(STDIN));
$nums = [];
while ($row = fgetcsv(STDIN, escape: ',')) {
    [$n, $name] = $row;
    $nums[$n] = $name;
}

foreach (range(1, $max) as $i) {
    $s = '';
    foreach ($nums as $n => $name) {
        if ($i % $n === 0) {
            $s = $s . $name;
        }
    }
    if ($s === '') {
        echo $i, "\n";
    } else {
        echo $s, "\n";
    }
}
$EOF$
WHERE title = 'ダイナミックFizzBuzz';
