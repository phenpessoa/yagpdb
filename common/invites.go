package common

import (
	"bytes"
	"regexp"
	"sync"
)

type InviteSource struct {
	Name  string
	Regex *regexp.Regexp
}

var DiscordInviteSource = &InviteSource{
	Name:  "Discord",
	Regex: regexp.MustCompile(`(?i)(discord\.gg|discordapp\.com\/invite|discord\.com\/invite)(?:\/#)?\/([a-zA-Z0-9-]+)`),
}

var ThirdpartyDiscordSites = []*InviteSource{
	{Name: "discord.me", Regex: regexp.MustCompile(`(?i)discord\.me\/.+`)},
	{Name: "invite.gg", Regex: regexp.MustCompile(`(?i)invite\.gg\/.+`)},
	{Name: "discord.io", Regex: regexp.MustCompile(`(?i)discord\.io\/.+`)},
	{Name: "discord.li", Regex: regexp.MustCompile(`(?i)discord\.li\/.+`)},
	{Name: "disboard.org", Regex: regexp.MustCompile(`(?i)disboard\.org\/server\/join\/.+`)},
	{Name: "discordy.com", Regex: regexp.MustCompile(`(?i)discordy\.com\/server\.php`)},

	// regexp.MustCompile(`disco\.gg\/.+`), Youc can't actually link to specific servers here can you, so not needed for now?
}

var AllInviteSources = append([]*InviteSource{DiscordInviteSource}, ThirdpartyDiscordSites...)

func ReplaceServerInvites(msg string, guildID int64, replacement string) string {

	for _, s := range AllInviteSources {
		msg = s.Regex.ReplaceAllString(msg, replacement)
	}

	return msg
}

func ContainsInvite(s string, checkDiscordSource, checkThirdPartySources bool) *InviteSource {
	for _, source := range AllInviteSources {
		if source == DiscordInviteSource && !checkDiscordSource {
			continue
		} else if source != DiscordInviteSource && !checkThirdPartySources {
			continue
		}

		if source.Regex.MatchString(s) {
			return source
		}
	}

	return nil
}

const (
	discordggInvite     = "discord.gg/"
	discordappcomInvite = "discordapp.com/invite/"
	discordcomInvite    = "discord.com/invite/"

	matchesPreAlloc = 20
)

var (
	discordggBuf     = make([]byte, 0, len(discordggInvite)+100)
	discordappcomBuf = make([]byte, 0, len(discordappcomInvite)+100)
	discordcomBuf    = make([]byte, 0, len(discordcomInvite)+100)

	prefixToTrim = []byte{'#', '/'}

	mu sync.Mutex
)

func InvitesInString(str string) []string {
	mu.Lock()
	defer func() {
		discordggBuf = discordggBuf[:0]
		discordappcomBuf = discordappcomBuf[:0]
		discordcomBuf = discordcomBuf[:0]

		mu.Unlock()
	}()

	matches := make([]string, 0, matchesPreAlloc)

	for i := range str {
		cur := str[i]

		var match string

		match, discordggBuf = invitesInner(cur, discordggBuf, discordggInvite)
		if match != "" {
			matches = append(matches, match)
		}

		match, discordappcomBuf = invitesInner(cur, discordappcomBuf, discordappcomInvite)
		if match != "" {
			matches = append(matches, match)
		}

		match, discordcomBuf = invitesInner(cur, discordcomBuf, discordcomInvite)
		if match != "" {
			matches = append(matches, match)
		}
	}

	if len(discordggBuf) > len(discordggInvite) {
		trimmed := discordggBuf[len(discordggInvite):]
		trimmed = bytes.TrimPrefix(trimmed, prefixToTrim)
		if str := string(trimmed); str != "" {
			matches = append(matches, string(trimmed))
		}
	}

	if len(discordappcomBuf) > len(discordappcomInvite) {
		trimmed := discordappcomBuf[len(discordappcomInvite):]
		trimmed = bytes.TrimPrefix(trimmed, prefixToTrim)
		if str := string(trimmed); str != "" {
			matches = append(matches, string(trimmed))
		}
	}

	if len(discordcomBuf) > len(discordcomInvite) {
		trimmed := discordcomBuf[len(discordcomInvite):]
		trimmed = bytes.TrimPrefix(trimmed, prefixToTrim)
		if str := string(trimmed); str != "" {
			matches = append(matches, string(trimmed))
		}
	}

	return matches
}

func invitesInner(cur byte, buf []byte, base string) (string, []byte) {
	if len(buf) >= len(base) {
		if len(buf) == len(base) && cur == '#' {
			buf = append(buf, cur)
			return "", buf
		}

		if len(buf) == len(base)+1 {
			if cur == '/' && buf[len(base)] == '#' {
				buf = append(buf, cur)
				return "", buf
			}

			if cur != '/' && buf[len(base)] == '#' {
				buf = buf[:0]
				return "", buf
			}
		}

		if !isValidInviteByte(cur) {
			if len(buf) == len(base) {
				buf = buf[:0]
				return "", buf
			}

			trimmed := buf[len(base):]
			trimmed = bytes.TrimPrefix(trimmed, prefixToTrim)
			match := string(trimmed)
			buf = buf[:0]
			return match, buf
		}

		buf = append(buf, cur)
		return "", buf
	}

	if cur == base[len(buf)] {
		buf = append(buf, cur)
		return "", buf
	}

	buf = buf[:0]
	return "", buf
}

func isValidInviteByte(b byte) bool {
	return b == '-' ||
		(b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z')
}
