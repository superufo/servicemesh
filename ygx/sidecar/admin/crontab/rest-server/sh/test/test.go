package main
import(
		"fmt"
		"os/exec"
		"syscall"
      )
func main(){
           c := exec.Command("C://Windows//system32//cmd.exe","/c","start", "D://gopro//src//github.com//go-chassis//ygx//sidecar//admin//crontab//rest-server//sh//test//test.bat"  )
	   c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	   if err := c.Run(); err != nil {
		   fmt.Println("Error: ", err)
	   }
}
