/**
* @program: send
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-09 16:36
**/

package main

func help() string {
	return `
Usage: send server
  -- start send server

Usage: send client [-f | --file] [path name]
  -- send file

Usage: send client str1
  -- send string
`
}
