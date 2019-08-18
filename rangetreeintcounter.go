package rangecounter

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type rangeTreeIntCounter struct {
	backend     Backend
	heightLimit int
	bitLength   uint
}

func (rtic *rangeTreeIntCounter) determineSumKeys(ctx context.Context, from, to int64) ([]string) {
	frompath := rtic.getTreePath(uint64(from))
	frompathKeys := rtic.getTreePathKeys(frompath)

	if from == to {
		return []string{frompathKeys[len(frompathKeys)-1]}
	}

	topath := rtic.getTreePath(uint64(to))
	topathKeys := rtic.getTreePathKeys(topath)
	maxPath := uint64(1 << rtic.bitLength)

	nonCommonIdx := 0
	for i := range frompath {
		if frompath[i] != topath[i] {
			nonCommonIdx = i
			break
		}
	}

	parentKey := ""
	if nonCommonIdx != 0 {
		parentKey = frompathKeys[nonCommonIdx-1]
	}

	inBetweenKeys := []string{}
	for i := frompath[nonCommonIdx] + 1; i < topath[nonCommonIdx]; i++ {
		inBetweenKeys = append(inBetweenKeys, rtic.appendKey(parentKey, i))
	}

	fromSubtreeKey := make([]string, 0, rtic.heightLimit*int(maxPath))
	for i := nonCommonIdx+1;i<len(frompath);i++ {
		parentKey := ""
		if i != 0 {
			parentKey = frompathKeys[i-1]
		}

		for nextPath := frompath[i] + 1; nextPath < maxPath; nextPath++ {
			fromSubtreeKey = append(fromSubtreeKey, rtic.appendKey(parentKey, nextPath))
		}
	}
	fromSubtreeKey = append(fromSubtreeKey, frompathKeys[len(frompathKeys)-1])

	toSubtreeKey := make([]string, 0, rtic.heightLimit*int(maxPath))
	for i := nonCommonIdx+1;i<len(topath);i++ {
		parentKey := ""
		if i != 0 {
			parentKey = topathKeys[i-1]
		}

		for beforePath := uint64(0); beforePath < topath[i]; beforePath++ {
			toSubtreeKey = append(toSubtreeKey, rtic.appendKey(parentKey, beforePath))
		}
	}
	toSubtreeKey = append(toSubtreeKey, topathKeys[len(topathKeys)-1])

	keys := make([]string, 0, len(fromSubtreeKey) + len(inBetweenKeys) + len(toSubtreeKey))
	keys = append(keys, fromSubtreeKey...)
	keys = append(keys, inBetweenKeys...)
	keys = append(keys, toSubtreeKey...)
	return keys
}

func (rtic *rangeTreeIntCounter) QuerySum(ctx context.Context, from, to int64) (int64, error) {
	keys := rtic.determineSumKeys(ctx, from, to)

	backendResult, err := rtic.backend.Query(ctx, keys)
	if err != nil {
		return 0, err
	}

	sum := int64(0)
	for _, it := range backendResult {
		sum = sum + it
	}
	return sum, nil
}

func (rtic *rangeTreeIntCounter) Increment(ctx context.Context, at int64, by int64) error {
	treepath := rtic.getTreePath(uint64(at))
	treepathKeys := rtic.getTreePathKeys(treepath)

	increments := []int64{}
	for _ = range treepath {
		increments = append(increments, by)
	}

	return rtic.backend.Increment(ctx, treepathKeys, increments)
}

func (rtic *rangeTreeIntCounter) getTreePathKeys(paths []uint64) []string {
	builder := strings.Builder{}
	keys := []string{}
	for _, path := range paths {
		builder.WriteRune(':')
		builder.WriteString(strconv.FormatUint(path, 10))
		keys = append(keys, builder.String())
	}
	return keys
}

func (rtic *rangeTreeIntCounter) appendKey(parent string, idx uint64) string {
	return parent + ":" + fmt.Sprint(idx)
}

func (rtic *rangeTreeIntCounter) getTreePath(idx uint64) []uint64 {
	lowerMask := uint64((1 << (rtic.bitLength)) - 1)
	treePath := []uint64{}

	for i := 0; i < (rtic.heightLimit)-1; i++ {
		cPath := idx & lowerMask
		treePath = append(treePath, cPath)
		idx = idx >> rtic.bitLength
	}

	treePath = append(treePath, idx)

	reversed := make([]uint64, len(treePath))
	for i := 0;i<len(treePath);i++ {
		reversed[i] = treePath[len(treePath)-i-1]
	}

	return reversed
}

func NewRangeTreeIntCounter(backend Backend, heightLimit int, bitLength uint) IntRangeCounter {
	if heightLimit < 0 {
		panic("heightLimit must be nonzero")
	}
	if bitLength < 0 {
		panic("bitLength must be nonzero")
	}
	return &rangeTreeIntCounter{
		backend:     backend,
		heightLimit: heightLimit,
		bitLength:   bitLength,
	}
}
