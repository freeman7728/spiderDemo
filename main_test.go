/*
 * @Description:
 * @author: freeman7728
 * @Date: 2024-08-29 19:28:02
 * @LastEditTime: 2024-08-30 14:19:23
 * @LastEditors: freeman7728
 */

package main

import "testing"

func TestInitDB(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"test1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitDB()
		})
	}
}
