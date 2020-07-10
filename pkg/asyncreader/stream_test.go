/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

package asyncreader

import (
	"bytes"
	"crypto/rand"
	"io"
	"sync"
	"testing"
)

func TestStream(t *testing.T) {
	w, rs := NewWithMaxsize(5, 10)
	savedBuf := new(bytes.Buffer)
	Treader := io.TeeReader(rand.Reader, savedBuf)
	wg := &sync.WaitGroup{}
	wg.Add(6)
	go func() {
		defer wg.Done()
		if copied, err := io.CopyN(w, Treader, 32); err != nil {
			t.Errorf("unable to write data(%d): %s", copied, err)
		}
		w.Close()
	}()
	results := make([]*bytes.Buffer, 5)
	for i := 0; i < 5; i++ {
		results[i] = new(bytes.Buffer)
		go func(b *bytes.Buffer, r io.Reader, indx int) {
			defer wg.Done()
			if copied, err := io.Copy(b, r); err != nil {
				t.Errorf("failed to read %d's data(%d): %s", indx, copied, err)
			}
		}(results[i], rs[i], i)
	}
	wg.Wait()
	if t.Failed() {
		t.FailNow()
	}
	for i := 0; i < 5; i++ {
		if !bytes.Equal(savedBuf.Bytes(), results[i].Bytes()) {
			t.Errorf("incorrectly copied values: %d, expected: %v, resulted: %v", i, savedBuf.Bytes(), results[i].Bytes())
		}
	}
}
