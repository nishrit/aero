package cache

import (
	"github.com/thejackrabbit/aero/conf"
	"github.com/thejackrabbit/aero/panik"
	"os"
	"time"
)

type debugCache struct {
	SimpleKeyFormat
	base Cacher
	w    *os.File
	r    *os.File
}

func NewDebug(w string, r string, inner Cacher) Cacher {
	var err error
	var rf, wf *os.File

	if r != "" {
		rf, err = os.OpenFile(r, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		panik.On(err)
	}
	if w != "" {
		wf, err = os.OpenFile(w, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		panik.On(err)
	}
	return debugCache{
		base: inner,
		w:    wf,
		r:    rf,
	}
}

// debug
// - inner
//   - type
//   - ..
// - read-file
// - write-file
func DebugFromConfig(container string) Cacher {
	return NewDebug(
		conf.String("", container, "write-file"),
		conf.String("", container, "read-file"),
		FromConfig(container+".inner"),
	)
}

func (c debugCache) Set(key string, data []byte, expireIn time.Duration) {
	k := c.Format(key)
	c.base.Set(key, data, expireIn)
	c.writeSetLog(k, data)
}

func (c debugCache) Get(key string) ([]byte, error) {
	data, err := c.base.Get(key)
	k := c.Format(key)

	if err != nil {
		c.writeGetLog(k, []byte("<not-found>"))
		return nil, err
	} else {
		c.writeGetLog(k, data)
		return data, nil
	}
}

func (c debugCache) Delete(key string) error {
	k := c.Format(key)
	c.writeDeleteLog(k, []byte(""))
	return c.base.Delete(key)
}

func (c debugCache) Close() {
	c.base.Close()
}

var newLine string = "\r\n"

func (c debugCache) writeGetLog(k string, b []byte) {
	if c.r == nil {
		return
	}

	c.r.WriteString(newLine + "~ ~ ~")
	c.r.WriteString(newLine + "GET: " + k + " @ " + time.Now().Format("2006-01-02_15:04:05"))

	if b != nil {
		c.r.WriteString(newLine + string(b))
	} else {
		c.r.WriteString(newLine + "<not-found>")
	}
}

func (c debugCache) writeSetLog(k string, b []byte) {
	if c.w == nil {
		return
	}

	c.w.WriteString(newLine + "~ ~ ~")
	c.w.WriteString(newLine + "SET: " + k + " @ " + time.Now().Format("2006-01-02_15:04:05"))

	if b != nil {
		c.w.WriteString(newLine + string(b))
	} else {
		c.w.WriteString(newLine + "<nil>")
	}
}

func (c debugCache) writeDeleteLog(k string, b []byte) {
	c.w.WriteString(newLine + "~ ~ ~")
	c.w.WriteString(newLine + "DELETE: " + k + " @ " + time.Now().Format("2006-01-02_15:04:05"))
}