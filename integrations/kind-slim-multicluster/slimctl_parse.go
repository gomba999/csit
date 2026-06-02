// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package kindslimmulticluster

import (
	"regexp"
	"strings"
)

var splitCellsRE = regexp.MustCompile(`\s{2,}`)

func splitTableCells(line string) []string {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	parts := splitCellsRE.Split(line, -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}

var nodeConnectedRE = regexp.MustCompile(`(?m)^Node ID: (\S+) status: Connected`)

// ParseConnectedNodeIDs returns node IDs reported as Connected in `slimctl controller node list` output.
func ParseConnectedNodeIDs(output string) []string {
	var ids []string
	for _, m := range nodeConnectedRE.FindAllStringSubmatch(output, -1) {
		ids = append(ids, m[1])
	}
	return ids
}

// ParseAppliedLinkID finds a link row with STATUS APPLIED and matching source/dest node IDs
// in `slimctl controller link outline` output (table from agntcy/slim slimctl).
func ParseAppliedLinkID(output, sourceNodeID, destNodeID string) (linkID string, ok bool) {
	for _, raw := range strings.Split(output, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "---") {
			continue
		}
		if strings.HasPrefix(line, "LINK_ID") || strings.HasPrefix(line, "Outline links") || strings.HasPrefix(line, "Number of links") {
			continue
		}
		if !strings.Contains(line, "APPLIED") {
			continue
		}
		cells := splitTableCells(line)
		if len(cells) < 8 {
			continue
		}
		// LINK_ID SOURCE DEST_NODE DEST_ENDPOINT STATUS STATUS_MSG DELETED LAST_UPDATED
		lid := cells[0]
		src := cells[1]
		dst := cells[2]
		status := cells[4]
		if status != "APPLIED" {
			continue
		}
		if src == sourceNodeID && dst == destNodeID {
			return lid, true
		}
	}
	return "", false
}
