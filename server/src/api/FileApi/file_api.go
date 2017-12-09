/* ====================================================================*/
/* Copyright (c) 2017.  All rights reserved.                           */
/* Author:     libinbin_1014@sina.com                                  */
/* Date :      2017/12/09                                              */
/* ====================================================================*/
package FileApi

import (
	"os"
)

type FileApi struct {
}

func ChkExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
