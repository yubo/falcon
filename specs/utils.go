package specs

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

/*
func Checksum(endpoint string, metric string, tags map[string]string) string {
	if tags == nil || len(tags) == 0 {
		return md5sum(fmt.Sprintf("%s/%s", endpoint, metric))
	}
	return md5sum(fmt.Sprintf("%s/%s/%s", endpoint, metric, sortedTags(tags)))
}
*/

func md5sum(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func fmtTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

/*
func sortedTags(tags map[string]string) string {
	if tags == nil {
		return ""
	}

	size := len(tags)

	if size == 0 {
		return ""
	}

	if size == 1 {
		for k, v := range tags {
			return fmt.Sprintf("%s=%s", k, v)
		}
	}

	keys := make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	ret := make([]string, size)
	for j, key := range keys {
		ret[j] = fmt.Sprintf("%s=%s", key, tags[key])
	}

	return strings.Join(ret, ",")
}
func readableFloat(raw float64) string {
	val := strconv.FormatFloat(raw, 'f', 5, 64)
	if strings.Contains(val, ".") {
		val = strings.TrimRight(val, "0")
		val = strings.TrimRight(val, ".")
	}

	return val
}
*/
