// Row Animation, inspiration: https://stonestoryrpg.com/

package ascii_generator

import (
	"math/rand"
	"pomossh/helper"
	"sync"
	"time"
)

type Row struct {
	cached                             string
	background                         string
	cur_animation, x, y, width, height int
	mu                                 sync.RWMutex
}

var (
	row_animation = []string{
		("                o |\n              _ \\ | __\n_ __ _ _ ___( ___\\| __ `)"),
		"                o _\n              _ \\  / _ \n_ __ ___ __ ( ___\\/ __ `)\n                 /_",
		"                o \n              _  V  _ \n_ __ _ ___ _( _ / ____ `)\n              _/ ~",
		"                 o  \n              _  (/ , _ \n_ __   _ ___( _  / ____ `)\n            __ ~/ ",
		"                 o / \n               _ </ , _ \n_  _ _ _ _ __( _ / ____ `)\n            ___ /",
	}
)

func (r *Row) NextAndString(percent int) string {
	return r.cached
}

func (r *Row) Width() int {
	return r.width
}

func (r *Row) Height() int {
	return r.height
}

func (r *Row) nextFrame() int {
	r.cur_animation++
	if r.cur_animation >= len(row_animation) {
		r.cur_animation = 0
	}

	return r.cur_animation
}

func (r *Row) backgroundCreator() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.background = ""
	for i := 0; i < r.height; i++ {
		for j := 0; j < r.width; j++ {
			if rand.Float32() > 0.9 {
				r.background += "_"
			} else {
				r.background += " "
			}
		}
		r.background += "\n"
	}
}

func GenerateRow(width, height int) *Row {
	r := &Row{
		x:      -30,
		y:      rand.Intn(height - 4),
		width:  width,
		height: height,
	}

	r.backgroundCreator()

	go func() {
		t := time.NewTicker(5 * time.Second)
		for {
			r.backgroundCreator()
			<-t.C
		}
	}()

	go func() {
		t := time.NewTicker(200 * time.Millisecond)
		for {
			r.mu.RLock()
			r.cached = helper.LayerString(r.background, row_animation[r.nextFrame()], r.x, r.y, true)
			r.mu.RUnlock()
			<-t.C
			if r.x > r.width {
				r.x = -30
				r.y = rand.Intn(height - 4)
			}
			r.x++
		}
	}()

	return r
}
