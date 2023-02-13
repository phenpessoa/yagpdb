package common

import (
	"sync"
	"testing"
)

func TestDiscordInviteRegex(t *testing.T) {
	testcases := []struct {
		input    string
		inviteID string
	}{
		{input: "https://discordapp.com/invite/FPxNX2", inviteID: "FPxNX2"},
		{input: "discordapp.com/developers/docs/reference#message-formatting", inviteID: ""},
		{input: "https://discord.gg/FPxNX2", inviteID: "FPxNX2"},
		{input: "https://discord.gg/landfall", inviteID: "landfall"},
		{input: "HElllo there", inviteID: ""},
		{input: "Jajajaj", inviteID: ""},
	}

	for _, v := range testcases {
		t.Run("Case "+v.input, func(t *testing.T) {
			matches := DiscordInviteSource.Regex.FindAllStringSubmatch(v.input, -1)
			if len(matches) < 1 && v.inviteID != "" {
				t.Error("No matches")
			}

			if len(matches) < 1 {
				return
			}

			if len(matches[0]) < 2 {
				if v.inviteID != "" {
					t.Error("ID Not found!")
				}

				return
			}

			id := matches[0][2]
			if id != v.inviteID {
				t.Errorf("Found ID %q, expected %q", id, v.inviteID)
			}
		})
	}
}

const (
	benchmarkStr = `\0abc␀discord.gg/123discord\n
	discord.gg/123as
	discord.gg/#/123asZWD
	asd.gg/sa123\n
discord.com/invite/#/123
https://discord.gg/123456
http://discord.gg/123456
discordapp.com/invite/0def世界您好something
https://discordapp.com/invite/FPxNX2
discordapp.com/developers/docs/reference#message-formatting
https://discord.gg/FPxNX2
https://discord.gg/landfall
HElllo there
Jajajaj
discord.com/invite/#123
discord.com/invite/##
discord.com/invite/##//
discord.com/invite/##132
discord.com/invite/##//123
discord.com/invite//123
discord.com/invite/1#23
discord.com/invite/1/23
discord.com/invite/`
)

var (
	testMu sync.Mutex
	sink   []string
)

// logic to retrieve matches from regex
func rgx(str string) []string {
	matches := DiscordInviteSource.Regex.FindAllStringSubmatch(str, -1)
	if len(matches) < 1 {
		return nil
	}

	_matches := make([]string, 0, matchesPreAlloc)

	for _, v := range matches {
		if len(v) < 3 {
			continue
		}

		_matches = append(_matches, v[2])
	}

	return _matches
}

func BenchmarkInvitesInString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = InvitesInString(benchmarkStr)
	}
}

func BenchmarkInvitesInString_Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			testMu.Lock()
			sink = InvitesInString(benchmarkStr)
			testMu.Unlock()
		}
	})
}

func TestSameBehaviour(t *testing.T) {
	for _, tc := range []string{
		"",
		benchmarkStr,
	} {
		rgxMatches := rgx(tc)
		newMatches := InvitesInString(tc)

		if len(rgxMatches) != len(newMatches) {
			t.Errorf(
				"\nUnequal results from old and new invite matching logic:\nOld len: %d\nNew len: %d\nFrom string: %s",
				len(rgxMatches), len(newMatches), tc,
			)
			continue
		}

		for idx, match := range rgxMatches {
			if match != newMatches[idx] {
				t.Errorf(
					"\nUnequal results from old and new invite matching logic:\nOld match #%d: %s\nNew match #%d: %s",
					idx, match, idx, newMatches[idx],
				)
			}
		}
	}
}
