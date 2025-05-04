use v5.12;
use utf8;
use open ':utf8';
use open ':std', ':utf8';

binmode(STDOUT,":encoding(utf8)");
 
my @subtitle_number;
say "# 目录";
while (<>) {
    # next if /^```/ ... /^```/;
    if (/^#(#+)\s*(.*?)\s*$/) {
        my ($level, $title) = (length($1), $2);
 
        my $indent = "  " x $level;
 
        my $id = $title;
        $id =~ s/[^_[:^punct:]]//g;
        $id =~ s/[[:space:]]/-/g;
        $id = lc $id;
 
        @subtitle_number = splice @subtitle_number, 0, $level;
        $subtitle_number[$level - 1] += 1;
        my $subtitle_number = join ".", @subtitle_number;
 
        say "$indent+ $subtitle_number [$title](#$id)";
    }
}
say "";