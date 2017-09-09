package main

import (
	"fmt"
	"os"
	"strings"
)

const ellipsis = "\u2026"

type pathSegment struct {
	path     string
	home     bool
	root     bool
	ellipsis bool
}

func cwdToPathSegments(cwd string) []pathSegment {
	pathSegments := make([]pathSegment, 0)

	home, _ := os.LookupEnv("HOME")
	if strings.HasPrefix(cwd, home) {
		pathSegments = append(pathSegments, pathSegment{
			path: "~",
			home: true,
		})
		cwd = cwd[len(home):]
	} else if cwd == "/" {
		pathSegments = append(pathSegments, pathSegment{
			path: "/",
			root: true,
		})
	}

	cwd = strings.Trim(cwd, "/")
	names := strings.Split(cwd, "/")
	if names[0] == "" {
		names = names[1:]
	}

	for _, name := range names {
		pathSegments = append(pathSegments, pathSegment{
			path: name,
		})
	}

	return pathSegments
}

func maybeShortenName(p *powerline, pathSegment string) string {
	if *p.args.CwdMaxDirSize > 0 && len(pathSegment) > *p.args.CwdMaxDirSize {
		return pathSegment[:*p.args.CwdMaxDirSize]
	} else {
		return pathSegment
	}
}

func escapeVariables(p *powerline, pathSegment string) string {
	pathSegment = strings.Replace(pathSegment, `\`, p.shellInfo.escapedBackslash, -1)
	pathSegment = strings.Replace(pathSegment, "`", p.shellInfo.escapedBacktick, -1)
	pathSegment = strings.Replace(pathSegment, `$`, p.shellInfo.escapedDollar, -1)
	return pathSegment
}

func getColor(p *powerline, pathSegment pathSegment, isLastDir bool) (uint8, uint8) {
	if pathSegment.home && p.theme.HomeSpecialDisplay {
		return p.theme.HomeFg, p.theme.HomeBg
	} else if isLastDir {
		return p.theme.CwdFg, p.theme.PathBg
	} else {
		return p.theme.PathFg, p.theme.PathBg
	}
}

func segmentCwd(p *powerline) {
	cwd := p.cwd
	if cwd == "" {
		cwd, _ = os.LookupEnv("PWD")
	}

	if *p.args.CwdMode == "plain" {
		home, _ := os.LookupEnv("HOME")
		if strings.HasPrefix(cwd, home) {
			cwd = "~" + cwd[len(home):]
		}

		p.appendSegment("cwd", segment{
			content:    fmt.Sprintf(" %s ", cwd),
			foreground: p.theme.CwdFg,
			background: p.theme.PathBg,
		})
	} else {
		pathSegments := cwdToPathSegments(cwd)

		if *p.args.CwdMode == "dironly" {
			pathSegments = pathSegments[len(pathSegments)-1:]
		} else {
			maxDepth := *p.args.CwdMaxDepth
			if maxDepth <= 0 {
				warn("Ignoring -cwd-max-depth argument since it's smaller than or equal to 0")
			} else if len(pathSegments) > maxDepth {
				var firstPart = make([]pathSegment, 0)
				secondPart := pathSegments[len(pathSegments)-maxDepth:]
				pathSegments = append(append(firstPart, pathSegment{
					path:     ellipsis,
					ellipsis: true,
				}), secondPart...)
			}

			for idx, pathSegment := range pathSegments {
				isLastDir := idx == len(pathSegments)-1
				foreground, background := getColor(p, pathSegment, isLastDir)

				segment := segment{
					content:    fmt.Sprintf(" %s ", escapeVariables(p, maybeShortenName(p, pathSegment.path))),
					foreground: foreground,
					background: background,
				}

				if !(pathSegment.home && p.theme.HomeSpecialDisplay) && !isLastDir {
					segment.separator = p.symbolTemplates.SeparatorThin
					segment.separatorForeground = p.theme.SeparatorFg
				}

				p.appendSegment("cwd", segment)
			}
		}
	}
}
