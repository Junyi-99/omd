# Input: OMD_TMUX_MIN (major.minor). `tmux -V` prints "tmux 3.5a" or "tmux next-3.9";
# only major.minor are compared.
command -v tmux >/dev/null && tmux -V | awk -v min="$OMD_TMUX_MIN" '{
    sub(/^[^0-9]*/, "", $2); split($2, v, "."); split(min, m, ".")
    exit !(v[1]+0 > m[1]+0 || (v[1]+0 == m[1]+0 && v[2]+0 >= m[2]+0))
}'
