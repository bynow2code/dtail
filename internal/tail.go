package internal

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	ErrStop = fmt.Errorf("tail should now stop")
)

type fileChange struct {
	new       chan string
	modified  chan bool
	truncated chan bool
	deleted   chan bool
}

func newFileChange() *fileChange {
	return &fileChange{
		new:       make(chan string),
		modified:  make(chan bool),
		truncated: make(chan bool),
		deleted:   make(chan bool),
	}
}

func (fc *fileChange) notifyNew(filename string) {
	select {
	case fc.new <- filename:
	default:
	}
}

func (fc *fileChange) notifyModified() {
	sendOnlyIfEmpty(fc.modified)
}

func (fc *fileChange) notifyTruncated() {
	sendOnlyIfEmpty(fc.truncated)
}

func (fc *fileChange) notifyDeleted() {
	sendOnlyIfEmpty(fc.deleted)
}

func (fc *fileChange) close() {
	close(fc.modified)
	close(fc.truncated)
	close(fc.deleted)
}

func sendOnlyIfEmpty(ch chan bool) {
	select {
	case ch <- true:
	default:
	}
}

type TailDir struct {
	dirname    string
	ticker     *time.Ticker
	tail       *Tail
	fileChange *fileChange
	Lines      chan *line
}

func NewTailDir(dirname string) {
	td := &TailDir{
		dirname:    dirname,
		ticker:     time.NewTicker(1 * time.Second),
		fileChange: newFileChange(),
		Lines:      make(chan *line),
	}

	go td.watchNewFile()

	newFile := <-td.fileChange.new

	color.NoColor = false
	color.Green("打开文件：%s", newFile)
	time.Sleep(1 * time.Second)

	td.tail = newTail(newFile)
	td.tail.file.Seek(0, io.SeekEnd)

	go td.sync()

	for l := range td.Lines {
		newText := l.text
		color.NoColor = false

		// Blue
		blueReplacement := color.BlueString(`$1$2`)
		newText = regexp.MustCompile(`("level":"info")|("func":".*?")`).ReplaceAllString(l.text, blueReplacement)

		// Red
		redReplacement := color.RedString(`$1`)
		newText = regexp.MustCompile(`("level":"error")`).ReplaceAllString(newText, redReplacement)

		// Yellow
		yellowReplacement := color.YellowString(`$1`)
		newText = regexp.MustCompile(`("level":"warning")`).ReplaceAllString(newText, yellowReplacement)

		// Green
		greenReplacement := color.GreenString(`$1`)
		newText = regexp.MustCompile(`("time":".*?")`).ReplaceAllString(newText, greenReplacement)

		// Cyan
		cyanReplacement := color.CyanString(`$1`)
		newText = regexp.MustCompile(`("message":".*?")`).ReplaceAllString(newText, cyanReplacement)

		fmt.Println(newText)
	}
}

func (td *TailDir) watchNewFile() {
	for range td.ticker.C {
		open, err := os.Open(td.dirname)
		if err != nil {
			log.Fatalln(err)
		}

		dirEntry, err := open.ReadDir(-1)
		if err != nil {
			return
		}

		var files []string
		for _, entry := range dirEntry {
			if entry.IsDir() {
				continue
			} else {
				files = append(files, entry.Name())
			}
		}

		if len(files) > 0 {
			sort.Sort(sort.Reverse(sort.StringSlice(files)))
			filename := filepath.Join(td.dirname, files[0])
			td.fileChange.notifyNew(filename)
		}
	}
}

func (td *TailDir) waitForChanges() error {
	if td.tail.size == 0 {
		var err error
		td.tail.size, err = td.tail.file.Seek(0, io.SeekCurrent)
		if err != nil {
			log.Fatalln(err)
		}
		td.events()
	}

	select {
	case <-td.fileChange.modified:
		return nil
	case <-td.fileChange.truncated:
		td.tail.reopen()
		td.tail.openReader()
		return nil
	case <-td.fileChange.deleted:
		return ErrStop
	case filename := <-td.fileChange.new:
		if filename == td.tail.filename {
			break
		}
		err := td.tail.watcher.Remove(td.tail.filename)
		if err != nil {
			log.Fatalln(err)
		}

		td.tail.filename = filename
		td.tail.reopen()
		td.tail.openReader()
		td.tail.size = 0
		td.tail.watcher.Add(filename)

		color.NoColor = false
		color.Green("打开文件：%s", filename)
		time.Sleep(1 * time.Second)

		return nil
	}
	return nil
}

func (t *Tail) reopen() {
	t.closeFile()

	var err error
	t.file, err = os.Open(t.filename)
	if err != nil {
		log.Fatalln(err)
	}

	t.openReader()
}

type Tail struct {
	filename string
	file     *os.File
	reader   *bufio.Reader
	watcher  *fsnotify.Watcher
	size     int64
}

type line struct {
	now  time.Time
	text string
}

func newTail(filename string) *Tail {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}

	t := &Tail{
		filename: filename,
		file:     file,
		watcher:  watcher,
	}

	return t
}

func (t *Tail) openReader() {
	t.reader = bufio.NewReader(t.file)
}

func (t *Tail) readLine() (text string, err error) {
	text, err = t.reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")
	return text, err
}

func (td *TailDir) sendLine(text string) {
	td.Lines <- &line{
		now:  time.Now(),
		text: text,
	}
}

func (td *TailDir) sync() {
	defer td.tail.close()

	td.tail.openReader()
	for {
		text, err := td.tail.readLine()
		if err == nil {
			td.sendLine(text)
		} else {
			if err == io.EOF {
				if text != "" {
					td.sendLine(text)
				}

				err := td.waitForChanges()
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				log.Fatalln(fmt.Sprintf("non-eof err:"), err)
			}
		}
	}
}

func (td *TailDir) events() {
	go func() {
		defer td.tail.watcher.Close()

		for {
			prevSize := td.tail.size

			var event fsnotify.Event
			var ok bool

			select {
			case event, ok = <-td.tail.watcher.Events:
				if !ok {
					return
				}
			}
			switch {
			case event.Op&fsnotify.Rename == fsnotify.Rename:
				fallthrough
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				td.fileChange.notifyDeleted()
				return
			case event.Op&fsnotify.Write == fsnotify.Write:
				stat, err := os.Stat(td.tail.filename)
				if err != nil {
					td.fileChange.notifyDeleted()
					return
				}

				td.tail.size = stat.Size()

				if prevSize > 0 && prevSize > td.tail.size {
					td.fileChange.notifyTruncated()
				} else {
					td.fileChange.notifyModified()
				}

				prevSize = td.tail.size
			}
		}
	}()

	err := td.tail.watcher.Add(td.tail.filename)
	if err != nil {
		if errors.Is(err, fsnotify.ErrClosed) {
			log.Println(err)
		} else {
			log.Fatalln(err)
		}
	}
}

func (t *Tail) reset(filename string) {
	t.watcher.Remove(t.filename)
	t.filename = filename
	t.watcher.Add(t.filename)
	t.reopen()
}

func (t *Tail) close() {
	t.closeFile()
}

func (t *Tail) closeFile() {
	t.file.Close()
	t.file = nil
}
