/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	IndentSize    = 4
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	crc64_table = crc64.MakeTable(crc64.ECMA)
	crc32_table = crc32.MakeTable(0xD5828281)
	_randSrc    = rand.NewSource(time.Now().UnixNano())

	f_network = regexp.MustCompile(`^(tcp)|(unix)+:`)
	F_ip      = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`)
	F_addr    = regexp.MustCompile(`^([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)?:[0-9]+`)
	F_num     = regexp.MustCompile(`^0x[0-9a-fA-F]+|^[0-9]+`)
	F_keyword = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_]+`)
	F_text    = regexp.MustCompile(`(^"[^"]+")|(^[^"\n \t;]+)`)
	F_env     = regexp.MustCompile(`\$\{[a-zA-Z][0-9a-zA-Z_]+\}`)
)

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, _randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = _randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func Md5sum(raw []byte) string {
	h := md5.New()
	h.Write(raw)
	//io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sum64(raw []byte) uint64 {
	h := crc64.New(crc64_table)
	h.Write(raw)
	//io.WriteString(h, raw)
	return h.Sum64()
}

func Sum32(raw []byte) uint32 {
	return crc32.Checksum(raw, crc32_table)
}

func FmtTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func ReadFileInt(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	if i, err := strconv.Atoi(strings.TrimSpace(string(data))); err != nil {
		return 0, err
	} else {
		return i, nil
	}
}

func GetType(v interface{}) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		//return "*" + t.Elem().Name()
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func PreByte(data []byte, pos int) (int, byte) {
	for i := pos; i >= 0; i-- {
		c := data[i]
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}
		return i, c
	}
	return -1, ' '
}

func IndentLines(i int, lines string) (ret string) {
	ls := strings.Split(strings.Trim(lines, "\n"), "\n")
	indent := strings.Repeat(" ", i*IndentSize)
	for _, l := range ls {
		ret += fmt.Sprintf("%s%s\n", indent, l)
	}
	return string([]byte(ret)[:len(ret)-1])
}

func AddrIsDisable(addr string) bool {
	if addr == "" || addr == "disable" {
		return true
	}
	return false
}

func Dialer(addr string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{Timeout: timeout}
	return d.Dial(ParseAddr(addr))
}

func IsFile(file string) bool {
	f, e := os.Stat(file)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func CleanSockFile(net, addr string) (string, string) {
	if net == "unix" && IsFile(addr) {
		os.Remove(addr)
	}
	return net, addr
}

func ParseAddr(url string) (net, addr string) {
	if f := f_network.Find([]byte(url)); f != nil {
		return url[:len(f)-1], url[len(f):]
	}
	return "tcp", url
}

func sortTags(s []byte) []byte {
	str := strings.Replace(string(s), " ", "", -1)
	if str == "" {
		return []byte{}
	}

	tags := strings.Split(str, ",")
	sort.Strings(tags)
	return []byte(strings.Join(tags, ","))
}

func Override(target, over interface{}) error {
	ov := reflect.ValueOf(over).Elem()
	ot := ov.Type()

	tv := reflect.ValueOf(target).Elem()
	tt := tv.Type()

	if !tv.CanSet() {
		return errors.New("target can't set")
	}

	for i := 0; i < ov.NumField(); i++ {
		fname := ot.Field(i).Name

		if _, ok := tt.FieldByName(fname); !ok {
			continue
		}

		of := ov.Field(i)
		tf := tv.FieldByName(fname)

		if !tf.CanSet() {
			fmt.Printf("fname %s\n", fname)
			continue
		}

		if of.Type().Kind() != tf.Type().Kind() {
			continue
		}

		switch tf.Type().Kind() {
		case reflect.Struct, reflect.Map:
			// skip
		default:
			tf.Set(of)
		}
	}
	return nil
}

func NewOrm(name, dsn string, maxIdleConns, maxOpenConns int) (o orm.Ormer, db *sql.DB, err error) {
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	if err = db.Ping(); err != nil {
		return
	}

	o, err = orm.NewOrmWithDB("mysql", name, db)
	return
}

// ############################################################################3

type RegisterServer func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// newGateway returns a new gateway server which translates HTTP into gRPC.
func newGateway(registerServer RegisterServer,
	ctx context.Context, address string, opts ...runtime.ServeMuxOption) (http.Handler, error) {

	mux := runtime.NewServeMux(opts...)

	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDialer(Dialer),
		grpc.WithBlock(),
	}

	err := registerServer(ctx, mux, address, dialOpts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

func Gateway(registerServer RegisterServer, ctx context.Context, mux *http.ServeMux, upstream string,
	opts ...runtime.ServeMuxOption) error {

	mux.HandleFunc("/swagger/", serveSwagger)

	gw, err := newGateway(registerServer, ctx, upstream, opts...)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)

	return nil
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
		glog.Errorf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	glog.Infof("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join("service", p)
	http.ServeFile(w, r, p)
}

// http client
func GetJson(url string, resp interface{}, timeout time.Duration) error {
	var cli *http.Client

	if strings.HasPrefix(url, "https://") {
		cli = &http.Client{
			Timeout:   timeout,
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		}

	} else {
		cli = &http.Client{Timeout: timeout}
	}

	r, err := cli.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	glog.V(4).Infof("getJson %s \n-> %s", url, string(b))
	return json.Unmarshal(b, resp)
}

func jsonStr(i interface{}) string {
	if ret, err := json.Marshal(i); err != nil {
		return ""
	} else {
		return string(ret)
	}
}

func PostJson(url string, param interface{}, resp interface{}) error {
	cli := &http.Client{Timeout: 60 * time.Second}
	r, err := cli.Post(
		url,
		"application/json",
		bytes.NewBuffer([]byte(jsonStr(param))),
	)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(resp)
}

func GetByteTls(url string, timeout time.Duration) ([]byte, error) {
	cli := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	r, err := cli.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

func KeyAttr(key []byte) (string, string, string, string, error) {
	var err error
	s := strings.Split(string(key), "/")
	if len(s) != 4 {
		err = EINVAL
	}

	return s[0], s[1], s[2], s[3], err
}

func AttrKey(endpoint, metric, tags, typ string) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s/%s", endpoint, metric, tags, typ))
}
