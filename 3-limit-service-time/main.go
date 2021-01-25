//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import "time"

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

var elapsed_time map[int]float64

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {

	var start <-chan time.Time
	if _, ok := elapsed_time[u.ID]; !ok {
		elapsed_time[u.ID] = 0.0
	}
	if elapsed_time[u.ID] >= 10.0 {
		return false
	} else if u.IsPremium {
		process()
		return true
	} else {
		start = time.Tick((time.Duration)((10.0 - elapsed_time[u.ID]) * 1000000000))
		time_start := time.Now()
		process()
		for {
			select {
			case <-start:
				{
					elapsed_time[u.ID] = 0.0
					return false
				}
			default:
				{
					elapsed_time[u.ID] += float64(time.Now().Sub(time_start) / 1000000000.0)
					return true
				}
			}
		}
		return true
	}
}

func main() {
	elapsed_time = make(map[int]float64)
	RunMockServer()
}
