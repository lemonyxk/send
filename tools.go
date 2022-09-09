/**
* @program: send
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-09 15:44
**/

package main

func HasArgs(flag string, args []string) bool {
	for i := 0; i < len(args); i++ {
		if args[i] == flag {
			return true
		}
	}
	return false
}

func GetArgs(flag []string, args []string) string {
	for i := 0; i < len(args); i++ {
		for j := 0; j < len(flag); j++ {
			if args[i] == flag[j] {
				if i+1 < len(args) {
					return args[i+1]
				}
			}
		}
	}
	return ""
}
