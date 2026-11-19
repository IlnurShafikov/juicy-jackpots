# Game algorithm: Juicy Jackpots, 5x3 classic fruit slot, 5 fixed lines, always-on Jackpot Coins (working title: Juicy Jackpots)

Config files live in `internal/configs/csv/`: `general/` (`paytable.csv`, `win_lines.csv`, `symbols_on_reel.csv`, `simple_cost_multiplier.csv`, `max_win_limit.csv`, `jp_pairs.csv`, `jp_types.csv`) and `<rtp>/` where `<rtp>` is 90, 92, 94 or 96 (`exp_rtp.csv` and, in `simple/base/`, `distributions.csv`, `jp_coin_count.csv`, `reels/group_N.csv`).

Symbol IDs start at 1 (ID 0 is the reserved empty symbol and is not used in this model).
1 High1, 2 High2, 3 Med1, 4 Med2, 5 Low1, 6 Low2, 7 Low3, 8 Low4, 9 Wild (reels 2-4 only). There is no Scatter and there are no free spins. Coins are not symbols on the strips: they are drawn separately in step 8.

1. **Bet.** One spin costs `simple_cost_multiplier.csv` = 5 units. Paytable payouts are in the same units (1 unit = line bet = 1/5 of the total bet; 5 lines x 1 unit = total bet). Jackpot prizes in `jp_pairs.csv` are already stored in units (prize in total bets x 5). Spin result in total bets = total win / 5.

2. **Reel set selection.** One draw per spin, with no draw per reel. A group is chosen by weights from `distributions.csv` (groups 1, 2 with weights 60, 40 in the RTP 96 version). All 5 reels use the strips of the chosen `group_N.csv` (column = reel 1..5, row = strip position).

3. **Reel stops.** For each reel independently, the stop position is chosen uniformly among all non-`-1` positions of that reel's column. Strips are circular (position 1 follows the last one). The visible window of reel _r_ is the 3 consecutive symbols (`symbols_on_reel.csv` = 3,3,3,3,3) at positions stop, stop+1, stop+2 (mod strip length); the symbol at the stop position is row 0 (top). All strips in the delivered groups have 40 positions, so no `-1` padding is used.

4. **Board modifications.** None.

5. **Win evaluation (lines, left to right).** Evaluated for each of the 5 lines in `win_lines.csv` (window row per reel, 0 = top: line 1 = 1,1,1,1,1; line 2 = 0,0,0,0,0; line 3 = 2,2,2,2,2; line 4 = 0,1,2,1,0; line 5 = 2,1,0,1,2).
   - The line symbol S is the symbol on reel 1 in the line's row. Wild (ID 9, `rules.WildID`) appears only on reels 2-4, so S is always a real symbol.
   - Wild substitutes for symbols 1-8 and has no payout of its own.
   - Let k be the number of consecutive reels counted from reel 1 whose cell on the line is S or a Wild. Counting stops at the first other cell.
   - The line pays `paytable.csv[S][k]` units. A payout of 0 means the count does not pay (symbol 8 Low4 pays from 2 reels, symbols 1-7 from 3 reels). Only the longest match on a line pays; shorter counts of the same line do not pay additionally.

6. **Win modifiers.** None.

7. **Summation.** The line win is the sum of the payouts of all 5 lines. Lines are independent: the same cell can belong to several lines.

8. **Jackpot Coins (on every spin, independent of the reels and lines).**
   - The coin count C (0..5) is drawn by weights from `jp_coin_count.csv` (`coins,weight`). One draw per spin.
   - Then C distinct reels are drawn uniformly without replacement, and for each coin a row 0..2 is drawn uniformly. This is a visual placement only and has no effect on any payout. Coins never replace symbols and never change line wins.
   - The prize depends on C only, from `jp_pairs.csv` (`coins,prize in units`): C = 2 MINI 50 units (10 total bets), C = 3 MINOR 100 (20), C = 4 MAJOR 500 (100), C = 5 GRAND 5000 (1000). `jp_types.csv` (`type id,prize in units`) maps the prize value to the jackpot type id 1..4 (MINI..GRAND) for the client event. C = 0 or 1 pays nothing. At most one prize per spin.

9. **Win cap.** Spin win = line win (step 7) + jackpot prize (step 8, converted to units). `max_win_limit.csv` = 1000 total bets = 5000 units, applied once to this sum. A spin with GRAND therefore always pays exactly 1000 total bets (the line win of that spin is lost to the cap). Without GRAND the maximum spin win is below the cap.

10. **Spin result.** Payout in total bets = spin win (after the cap) / `simple_cost_multiplier.csv` = / 5.

## RTP versions

RTP versions 90, 92, 94 and 96 use the same algorithm, strips, paytable and prize values. They differ only in `distributions.csv`, `jp_coin_count.csv` and `exp_rtp.csv` (see the folders `internal/configs/csv/90` ... `96`). In `jp_coin_count.csv` only the weight of 0 coins changes, so all coin prize frequencies scale together.

## Notes for the implementer

- Wild 9 is on reels 2-4 only. Reels 1 and 5 contain no Wild.
- A line always starts from reel 1; there is no right-to-left pay and no scatter pay.
- Coins and lines are independent random events. The only interaction is the win cap of step 9.
- Strips contain no two identical adjacent symbols (circularly).
