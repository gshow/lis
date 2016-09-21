package safemap

/*

package main
import(
  smap "github.com/ricolau/safemap"
 "fmt"
 "time"
)
func main(){
    b := smap.New()
    for i:=0;i<1000;i++{
        go func(i int){
            b.Set(i,i+1000)
        }(i)

        go func(i int){
            b.Exist(i)
        }(i)
        go func (i int){
            b.Get(i+10)
        }(i)
        go func (i int){
            b.Delete(i+10)
        }(i)
    }
    time.Sleep(time.Second * 2)
    fmt.Println(b.Size())
}

*/
