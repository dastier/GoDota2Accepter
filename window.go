package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)



func findIds(id string) {
	// find the process id by the process name
	fmt.Println("we are in findIds")
	sum := 1
	for sum < 10 {
		sum += sum

	
		fpid, err := robotgo.FindIds(id)

		fmt.Println("fpid: ", fpid)
		if err == nil {
			fmt.Println("pids...", fpid)
			if len(fpid) > 0 {
				if robotgo.ActivePID(fpid[0]) != nil {
					fmt.Println("Could not find dota window", err)
				}


				tl := robotgo.GetTitle(fpid[0])
				fmt.Println("pid[0] title is: ", tl)

				x, y, w, h := robotgo.GetBounds(fpid[0])
				fmt.Println("GetBounds is: ", x, y, w, h)
				robotgo.MaxWindow(fpid[0])
				// robotgo.CloseWindow(fpid[0])
				robotgo.Sleep(2)
				robotgo.KeyTap("enter")

				// robotgo.Kill(fpid[0])
				break
			}
		} else {
			fmt.Println("Found errors while robotgo.FindIds:", err)
		}

}
}