/**
* @program: send
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-09 17:15
**/

package main

type writer struct {
	total   int64
	current int64
	last    int64

	// onProgress func(p []byte, current int64, total int64)
	// rate       int64

	onData func([]byte)
}

func (w *writer) Write(p []byte) (int, error) {
	n := len(p)
	w.current += int64(n)

	w.onData(p)

	return n, nil
}
