UPDATE problems
SET
sample_code = $EOF$list($player1_name, $player2_name) = fgetcsv(STDIN, 1024, ',', '"', '');
while ($row = fgetcsv(STDIN, 1024, ',', '"', '')) {
    if (count($row) === 0) {
        break;
    }
    $player1_hand = $row[0];
    $player2_hand = $row[1];
    switch ($player1_hand) {
        case "グー":
            printf("%sさん %s", $player1_name, match ($player1_hand) {
                "グー" => "\u{270A}",
                "チョキ" => "\u{270C}",
                "パー" => "\u{270B}",
            });
            echo " vs ";
            switch ($player2_hand) {
                case "グー":
                    printf("%sさん %s", $player2_name, match ($player2_hand) {
                        "グー" => "\u{270A}",
                        "チョキ" => "\u{270C}",
                        "パー" => "\u{270B}",
                    });
                    echo " => あいこ\n";
                    break;
                case "チョキ":
                    printf("%sさん %s", $player2_name, match ($player2_hand) {
                        "グー" => "\u{270A}",
                        "チョキ" => "\u{270C}",
                        "パー" => "\u{270B}",
                    });
                    echo " => {$player1_name}の勝ち\n";
                    break;
                case "パー":
                    printf("%sさん %s", $player2_name, match ($player2_hand) {
                        "グー" => "\u{270A}",
                        "チョキ" => "\u{270C}",
                        "パー" => "\u{270B}",
                    });
                    echo " => {$player2_name}の勝ち\n";
                    break;

            }
        break;
        case "チョキ":
            printf("%sさん %s", $player1_name, match ($player1_hand) {
                "グー" => "\u{270A}",
                "チョキ" => "\u{270C}",
                "パー" => "\u{270B}",
            });
            echo " vs ";
            switch ($player2_hand) {
                case "グー":
                    printf("%sさん %s", $player2_name, match ($player2_hand) {
                        "グー" => "\u{270A}",
                        "チョキ" => "\u{270C}",
                        "パー" => "\u{270B}",
                    });
                    echo " => {$player2_name}の勝ち\n";
                    break;
                case "チョキ":
                    printf("%sさん %s", $player2_name, match ($player2_hand) {
                        "グー" => "\u{270A}",
                        "チョキ" => "\u{270C}",
                        "パー" => "\u{270B}",
                    });
                    echo " => あいこ\n";
                    break;
                case "パー":
                    printf("%sさん %s", $player2_name, match ($player2_hand) {
                        "グー" => "\u{270A}",
                        "チョキ" => "\u{270C}",
                        "パー" => "\u{270B}",
                    });
                    echo " => {$player1_name}の勝ち\n";
                    break;

            }
        break;
        case "パー":
            printf("%sさん %s", $player1_name, match ($player1_hand) {
                "グー" => "\u{270A}",
                "チョキ" => "\u{270C}",
                "パー" => "\u{270B}",
            });
            echo " vs ";
            switch ($player2_hand) {
                case "グー":
                    printf("%sさん %s", $player2_name, match ($player2_hand) {
                        "グー" => "\u{270A}",
                        "チョキ" => "\u{270C}",
                        "パー" => "\u{270B}",
                    });
                    echo " => {$player1_name}の勝ち\n";
                    break;
                case "チョキ":
                    printf("%sさん %s", $player2_name, match ($player2_hand) {
                        "グー" => "\u{270A}",
                        "チョキ" => "\u{270C}",
                        "パー" => "\u{270B}",
                    });
                    echo " => {$player2_name}の勝ち\n";
                    break;
                case "パー":
                    printf("%sさん %s", $player2_name, match ($player2_hand) {
                        "グー" => "\u{270A}",
                        "チョキ" => "\u{270C}",
                        "パー" => "\u{270B}",
                    });
                    echo " => あいこ\n";
                    break;
            }
        break;
    }
}
$EOF$
WHERE title = 'じゃんけん';

UPDATE problems
SET
sample_code = $EOF$$aa = array_map(fn($s) => explode("\n", $s), [
  <<<'AA'
┌┐
││
└┘
AA,
  <<<'AA'
 ┐
 │
 ╵
AA,
  <<<'AA'
╶┐
┌┘
└╴
AA,
  <<<'AA'
┌┐
 ┤
└┘
AA,
  <<<'AA'
╷╷
└┤
 ╵
AA,
  <<<'AA'
┌╴
└┐
╶┘
AA,
  <<<'AA'
┌┐
├┐
└┘
AA,
  <<<'AA'
┌┐
 │
 ╵
AA,
  <<<'AA'
┌┐
├┤
└┘
AA,
  <<<'AA'
┌┐
└┤
└┘
AA,
]);

while ($line = fgets(STDIN)) {
    $digits = str_split(rtrim($line, "\n"));
    $output = '';
    foreach (range(0, 2) as $i) {
        foreach ($digits as $digit) {
            $output = $output . $aa[$digit][$i];
        }
        $output = $output . "\n";
    }
    echo $output;
}
$EOF$
WHERE title = 'でっかい数字';
